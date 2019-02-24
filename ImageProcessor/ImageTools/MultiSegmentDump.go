package ImageTools

import (
	"fmt"
	"github.com/OpenSatelliteProject/SatHelperApp/ImageProcessor/Structs"
	"github.com/OpenSatelliteProject/SatHelperApp/Logger"
	"github.com/OpenSatelliteProject/SatHelperApp/XRIT"
	"github.com/OpenSatelliteProject/SatHelperApp/XRIT/PacketData"
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

func DumpMultiSegment(msi *Structs.MultiSegmentImage) (error, string) {
	width := int(msi.FirstSegmentHeader.SegmentIdentificationHeader.MaxColumns)
	height := int(msi.FirstSegmentHeader.SegmentIdentificationHeader.MaxRows)

	folder := path.Dir(msi.FirstSegmentFilename)

	newFilename := path.Join(folder, msi.Name+".png")

	f, err := os.Create(newFilename)
	if err != nil {
		SLog.Error("Error creating file %s: %s\n", newFilename, err)
		return err, ""
	}

	defer f.Close()

	img := image.NewGray(image.Rect(0, 0, width, height))

	for _, filename := range msi.Files {
		xh, err := XRIT.ParseFile(filename)
		if err != nil {
			return err, ""
		}

		if xh.PrimaryHeader.FileTypeCode != PacketData.IMAGE {
			return fmt.Errorf("the specified file is not an image container"), ""
		}

		offset := xh.PrimaryHeader.HeaderLength

		f, err := os.Open(filename)

		if err != nil {
			return err, ""
		}

		_, err = f.Seek(int64(offset), io.SeekStart)
		if err != nil {
			return err, ""
		}

		data, err := ioutil.ReadAll(f)
		if err != nil {
			return err, ""
		}

		px := int(xh.SegmentIdentificationHeader.StartColumn)
		py := int(xh.SegmentIdentificationHeader.StartLine)

		DrawGray8At(data, px, py, img)
	}

	err = png.Encode(f, img)

	if err != nil {
		return err, ""
	}

	return nil, newFilename
}
