package ImageTools

import (
	"bytes"
	"fmt"
	"github.com/opensatelliteproject/SatHelperApp/ImageProcessor/ImageData"
	"github.com/opensatelliteproject/SatHelperApp/ImageProcessor/MapCutter"
	"github.com/opensatelliteproject/SatHelperApp/ImageProcessor/MapDrawer"
	"github.com/opensatelliteproject/SatHelperApp/ImageProcessor/Projector"
	"github.com/opensatelliteproject/SatHelperApp/ImageProcessor/Structs"
	"github.com/opensatelliteproject/SatHelperApp/Logger"
	"github.com/opensatelliteproject/SatHelperApp/Tools"
	"github.com/opensatelliteproject/SatHelperApp/XRIT"
	"github.com/opensatelliteproject/SatHelperApp/XRIT/Geo"
	"github.com/opensatelliteproject/SatHelperApp/XRIT/PacketData"
	"image"
	"image/draw"
	"image/jpeg"
	"image/png"
	"io"
	"io/ioutil"
	"os"
	"path"
	"strings"
)

var saveNoMap = false
var saveNoProj = false

func SetSaveNoMap(s bool) {
	saveNoMap = s
}

func SetSaveNoProj(s bool) {
	saveNoProj = s
}

func DrawGray8At(data []byte, px, py int, image *image.Gray) {
	b := image.Bounds()
	p := b.Dx()*py + px

	copy(image.Pix[p:], data)
}

func MultiSegmentAssemble(msi *Structs.MultiSegmentImage) (error, image.Image) {
	width := int(msi.FirstSegmentHeader.SegmentIdentificationHeader.MaxColumns)
	height := int(msi.FirstSegmentHeader.SegmentIdentificationHeader.MaxRows)

	if width == 0 || height == 0 {
		// COMS-1 for example doesn't have so we sum up the sizes
		w := int(msi.FirstSegmentHeader.ImageStructureHeader.Columns)
		h := int(msi.FirstSegmentHeader.ImageStructureHeader.Lines) * int(msi.FirstSegmentHeader.SegmentIdentificationHeader.MaxSegments)

		if width == 0 {
			width = w
		}
		if height == 0 {
			height = h
		}
	}

	img := image.NewGray(image.Rect(0, 0, width, height))

	for _, filename := range msi.Files {
		xh, err := XRIT.ParseFile(filename)
		if err != nil {
			return err, nil
		}

		if xh.PrimaryHeader.FileTypeCode != PacketData.IMAGE {
			return fmt.Errorf("the specified file is not an image container"), nil
		}

		offset := xh.PrimaryHeader.HeaderLength

		f, err := os.Open(filename)

		if err != nil {
			return err, nil
		}

		_, err = f.Seek(int64(offset), io.SeekStart)
		if err != nil {
			return err, nil
		}

		data, err := ioutil.ReadAll(f)
		if err != nil {
			return err, nil
		}

		px := int(xh.SegmentIdentificationHeader.StartColumn)
		py := int(xh.SegmentIdentificationHeader.StartLine)

		// COMS-1 JPG Image
		if xh.Compression() == PacketData.JPEG {
			b := bytes.NewReader(data)
			im, err := jpeg.Decode(b)
			if err != nil {
				return err, nil
			}
			draw.Draw(img, im.Bounds().Add(image.Pt(px, py)), im, im.Bounds().Min, draw.Src)
		} else {
			DrawGray8At(data, px, py, img)
		}
	}

	return nil, img
}

func SaveImage(filename string, img image.Image) error {
	f, err := os.Create(filename)
	if err != nil {
		SLog.Error("Error creating file %s: %s\n", filename, err)
		return err
	}

	defer f.Close()

	err = png.Encode(f, img)

	if err != nil {
		return err
	}

	return nil
}

func GetNoMapName(filename, extra string) string {
	// Remove file timestamp
	if extra != "" {
		return filename[:len(filename)-16] + "-" + extra + "-nomap.png"
	}
	return filename[:len(filename)-16] + "-nomap.png"
}
func GetNoProjName(filename, extra string) string {
	// Remove file timestamp
	if extra != "" {
		return filename[:len(filename)-16] + "-" + extra + "-noproj.png"
	}
	return filename[:len(filename)-16] + "-noproj.png"
}

func cutRegionAndDump(region string, msi *Structs.MultiSegmentImage, mapDrawer *MapDrawer.MapDrawer, mapCutter *MapCutter.MapCutter, visCurve *CurveManipulator, reproject bool, enhance bool, metadata bool) (error, string) {
	s, err := mapCutter.GetSection(region)
	if err != nil {
		return err, ""
	}

	folder := path.Dir(msi.FirstSegmentFilename)

	newFilename := path.Join(folder, fmt.Sprintf("%s-%s.png", msi.Name, region))
	newFilenameEnhanced := path.Join(folder, fmt.Sprintf("%s-%s-enhanced.png", msi.Name, region))

	if Tools.Exists(newFilename) {
		SLog.Info("File %s already exists, skipping...", newFilename)
		return nil, newFilename
	}

	err, img := MultiSegmentAssemble(msi)
	if err != nil {
		return err, ""
	}

	gc, err := Geo.MakeGeoConverterFromXRIT(msi.FirstSegmentHeader)
	if err != nil {
		return err, ""
	}

	if strings.Contains(newFilename, "C02_") { // Only on visible channels
		err = visCurve.ApplyCurve(img)
		if err != nil {
			SLog.Error("Error applying curve to visible image: %s", err)
		}
		enhance = false
	}

	imgRGBA := image.NewRGBA(img.Bounds())
	draw.Draw(imgRGBA, img.Bounds(), img, img.Bounds().Min, draw.Src)

	satLut := msi.FirstSegmentHeader.TemperatureLUT

	enh := MakeImageEnhancer(ImageData.DefaultMinimumTemperature, ImageData.DefaultMaximumTemperature, satLut, ImageData.TemperatureScaleLUT, false)

	if enhance {
		imgRGBA, err = enh.EnhanceWithLUT(imgRGBA)
		if err != nil {
			SLog.Error("Error enhancing image %s: %s", newFilename, err)
		} else {
			img = imgRGBA
		}
	}

	if reproject {
		SLog.Debug("Reprojecting Image to Linear")

		proj := Projector.MakeProjector(gc)
		imgRGBA = proj.ReprojectLinearMultiThread(imgRGBA)
		img = imgRGBA
		gc = Projector.MakeLinearConverter(imgRGBA.Bounds().Dx(), imgRGBA.Bounds().Dy(), gc)
	}

	if mapDrawer != nil {
		SLog.Debug("Map Drawer enabled. Drawing maps...")
		mapDrawer.DrawMap(imgRGBA, gc)
		img = imgRGBA
	}

	img, err = mapCutter.CutMap(region, img, gc)
	imgRGBA = img.(*image.RGBA)

	if err != nil {
		return err, ""
	}

	if metadata {
		if enhance {
			imgRGBA, err = enh.DrawMeta(s.Name, imgRGBA, msi.FirstSegmentHeader)
			if err != nil {
				SLog.Error("Error drawing metadata on %s: %s", newFilenameEnhanced, err)
			}
		} else {
			imgRGBA, err = enh.DrawMetaWithoutScale(s.Name, imgRGBA, msi.FirstSegmentHeader)
			if err != nil {
				SLog.Error("Error drawing metadata on %s: %s", newFilenameEnhanced, err)
			}
		}
		if imgRGBA != nil {
			img = imgRGBA
		}
	}

	if enhance {
		err = SaveImage(newFilenameEnhanced, img)
	} else {
		err = SaveImage(newFilename, img)
	}

	if err != nil {
		return err, ""
	}

	return nil, newFilename
}

func DumpCutRegion(msi *Structs.MultiSegmentImage, mapDrawer *MapDrawer.MapDrawer, mapCutter *MapCutter.MapCutter, visCurve *CurveManipulator, reproject bool, enhance bool, metadata bool, cutRegions []string) error {
	if mapCutter == nil || len(cutRegions) == 0 {
		return fmt.Errorf("no regions to cut")
	}

	for _, region := range cutRegions {
		err, _ := cutRegionAndDump(region, msi, mapDrawer, mapCutter, visCurve, reproject, enhance, metadata)
		if err != nil {
			SLog.Error("Error processing region %s: %s", region, err)
		}
	}

	return nil
}

func DumpMultiSegment(msi *Structs.MultiSegmentImage, mapDrawer *MapDrawer.MapDrawer, visCurve *CurveManipulator, reproject bool, enhance bool, metadata bool) (error, string) {
	folder := path.Dir(msi.FirstSegmentFilename)

	newFilename := path.Join(folder, msi.Name+".png")
	newFilenameNoMap := path.Join(folder, GetNoMapName(msi.Name, ""))
	newFilenameNoProj := path.Join(folder, GetNoProjName(msi.Name, ""))
	newFilenameEnhanced := path.Join(folder, msi.Name+"-enhanced.png")

	if Tools.Exists(newFilename) {
		SLog.Info("File %s already exists, skipping...", newFilename)
		return nil, newFilename
	}

	err, img := MultiSegmentAssemble(msi)
	if err != nil {
		return err, ""
	}

	gc, err := Geo.MakeGeoConverterFromXRIT(msi.FirstSegmentHeader)
	if err != nil {
		return err, ""
	}

	metaName := path.Join(folder, msi.Name+".json")
	err = ioutil.WriteFile(metaName, []byte(msi.FirstSegmentHeader.ToJSON()), os.ModePerm)
	if err != nil {
		SLog.Error("Cannot write Meta file %s: %s", metaName, err)
	}

	if strings.Contains(newFilename, "C02_") { // Only on visible channels
		err = visCurve.ApplyCurve(img)
		if err != nil {
			SLog.Error("Error applying curve to visible image: %s", err)
		}
		enhance = false
	}

	if mapDrawer != nil && saveNoMap && !Tools.Exists(newFilenameNoMap) {
		SLog.Debug("Saving No Map Image: %s", newFilenameNoMap)
		err := SaveImage(newFilenameNoMap, img)
		if err != nil {
			SLog.Error("Error saving %s: %s", newFilenameNoMap, err)
		}
	}

	imgRGBA := image.NewRGBA(img.Bounds())
	draw.Draw(imgRGBA, img.Bounds(), img, img.Bounds().Min, draw.Src)

	satLut := msi.FirstSegmentHeader.TemperatureLUT

	enh := MakeImageEnhancer(ImageData.DefaultMinimumTemperature, ImageData.DefaultMaximumTemperature, satLut, ImageData.TemperatureScaleLUT, false)

	if enhance {
		imgRGBA, err = enh.EnhanceWithLUT(imgRGBA)
		if err != nil {
			SLog.Error("Error enhancing image %s: %s", newFilenameNoMap, err)
		} else {
			img = imgRGBA
		}
	}

	if reproject {
		if saveNoProj && !Tools.Exists(newFilenameNoProj) {
			SLog.Debug("Saving No Projection Image: %s", newFilenameNoProj)
			err := SaveImage(newFilenameNoProj, img)
			if err != nil {
				SLog.Error("Error saving %s: %s", newFilenameNoProj, err)
			}
		}
		SLog.Debug("Reprojecting Image to Linear")

		proj := Projector.MakeProjector(gc)
		imgRGBA = proj.ReprojectLinearMultiThread(imgRGBA)
		img = imgRGBA
		gc = Projector.MakeLinearConverter(imgRGBA.Bounds().Dx(), imgRGBA.Bounds().Dy(), gc)
	}

	if mapDrawer != nil {
		SLog.Debug("Map Drawer enabled. Drawing maps...")
		mapDrawer.DrawMap(imgRGBA, gc)
		img = imgRGBA
	}

	if metadata {
		if enhance {
			imgRGBA, err = enh.DrawMeta("", imgRGBA, msi.FirstSegmentHeader)
			if err != nil {
				SLog.Error("Error drawing metadata on %s: %s", newFilenameEnhanced, err)
			}
		} else {
			imgRGBA, err = enh.DrawMetaWithoutScale("", imgRGBA, msi.FirstSegmentHeader)
			if err != nil {
				SLog.Error("Error drawing metadata on %s: %s", newFilenameEnhanced, err)
			}
		}
		if imgRGBA != nil {
			img = imgRGBA
		}
	}

	if enhance {
		err = SaveImage(newFilenameEnhanced, img)
	} else {
		err = SaveImage(newFilename, img)
	}

	if err != nil {
		return err, ""
	}

	return nil, newFilename
}
