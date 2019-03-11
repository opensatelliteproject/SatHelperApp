package ImageTools

import (
	"fmt"
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

		DrawGray8At(data, px, py, img)
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

func GetNoMapName(filename string) string {
	// Remove file timestamp
	return filename[:len(filename)-16] + "-nomap.png"
}
func GetNoProjName(filename string) string {
	// Remove file timestamp
	return filename[:len(filename)-16] + "-noproj.png"
}

func DumpMultiSegment(msi *Structs.MultiSegmentImage, mapDrawer *MapDrawer.MapDrawer, visCurve *CurveManipulator, reproject bool) (error, string) {
	folder := path.Dir(msi.FirstSegmentFilename)

	newFilename := path.Join(folder, msi.Name+".png")
	newFilenameNoMap := path.Join(folder, GetNoMapName(msi.Name))
	newFilenameNoProj := path.Join(folder, GetNoProjName(msi.Name))

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
	}

	if mapDrawer != nil {
		if saveNoMap && !Tools.Exists(newFilenameNoMap) {
			SLog.Debug("Saving No Map Image: %s", newFilenameNoMap)
			err := SaveImage(newFilenameNoMap, img)
			if err != nil {
				SLog.Error("Error saving %s: %s", newFilenameNoMap, err)
			}
		}
		SLog.Debug("Map Drawer enabled. Drawing maps...")
		newImg := image.NewRGBA(img.Bounds())
		draw.Draw(newImg, img.Bounds(), img, img.Bounds().Min, draw.Src)
		mapDrawer.DrawMap(newImg, gc)
		img = newImg
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
		img2 := proj.ReprojectLinearMultiThread(img)
		img = img2
	}

	err = SaveImage(newFilename, img)

	if err != nil {
		return err, ""
	}

	return nil, newFilename
}
