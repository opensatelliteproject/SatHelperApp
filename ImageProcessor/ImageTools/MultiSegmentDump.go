package ImageTools

import (
	"fmt"
	"github.com/opensatelliteproject/SatHelperApp/ImageProcessor/Projector"
	"github.com/opensatelliteproject/SatHelperApp/ImageProcessor/Structs"
	"github.com/opensatelliteproject/SatHelperApp/Logger"
	"github.com/opensatelliteproject/SatHelperApp/XRIT"
	"github.com/opensatelliteproject/SatHelperApp/XRIT/Geo"
	"github.com/opensatelliteproject/SatHelperApp/XRIT/PacketData"
	"image"
	"image/png"
	"io"
	"io/ioutil"
	"os"
	"path"
)

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

func DumpMultiSegment(msi *Structs.MultiSegmentImage, reproject bool) (error, string) {
	folder := path.Dir(msi.FirstSegmentFilename)

	newFilename := path.Join(folder, msi.Name+".png")

	f, err := os.Create(newFilename)
	if err != nil {
		SLog.Error("Error creating file %s: %s\n", newFilename, err)
		return err, ""
	}

	defer f.Close()

	err, img := MultiSegmentAssemble(msi)

	if err != nil {
		return err, ""
	}

	if reproject {
		SLog.Info("Reprojecting Image to Linear")
		gc, err := Geo.MakeGeoConverterFromXRIT(msi.FirstSegmentHeader)
		if err != nil {
			return err, ""
		}

		proj := Projector.MakeProjector(gc)
		img2 := proj.ReprojectLinearMultiThread(img)
		img = img2
	}

	err = png.Encode(f, img)

	if err != nil {
		return err, ""
	}

	return nil, newFilename
}
