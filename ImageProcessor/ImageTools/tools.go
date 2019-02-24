package ImageTools

import (
	"fmt"
	"github.com/OpenSatelliteProject/SatHelperApp/Logger"
	"github.com/OpenSatelliteProject/SatHelperApp/XRIT"
	"github.com/OpenSatelliteProject/SatHelperApp/XRIT/PacketData"
	"image"
	"image/jpeg"
	"io"
	"io/ioutil"
	"os"
	"strings"
)

var compressionToExtension = map[int]string{
	PacketData.NO_COMPRESSION: ".jpg",
	PacketData.LRIT_RICE:      ".jpg",
	PacketData.JPEG:           ".jpg",
	PacketData.GIF:            ".gif",
}

func DumpRaw(newFileName string, data []byte, xh *XRIT.Header) error {
	imgStruct := xh.ImageStructureHeader

	f, err := os.Create(newFileName)
	if err != nil {
		SLog.Error("Error creating file %s: %s\n", newFileName, err)
		return err
	}

	defer f.Close()

	if imgStruct.BitsPerPixel == 8 {
		totalPixels := int(imgStruct.Columns) * int(imgStruct.Lines)
		if totalPixels > len(data) {
			missingBytes := totalPixels - len(data)
			SLog.Warn("Missing %d bytes in image data.\n", missingBytes)
			data = append(data, make([]byte, missingBytes)...)
		}
		img := &image.Gray{Pix: data, Stride: int(imgStruct.Columns), Rect: image.Rect(0, 0, int(imgStruct.Columns), int(imgStruct.Lines))}

		err = jpeg.Encode(f, img, &jpeg.Options{
			Quality: 100,
		})

		if err != nil {
			return err
		}

		return nil
	}

	if imgStruct.BitsPerPixel == 1 {
		totalPixels := int(imgStruct.Columns) * int(imgStruct.Lines)
		if totalPixels > len(data)/8 {
			missingBytes := (totalPixels - len(data)) / 8
			SLog.Warn("Missing %d bytes in image data.\n", missingBytes)
			data = append(data, make([]byte, missingBytes)...)
		}
		img := image.NewGray(image.Rect(0, 0, int(imgStruct.Columns), int(imgStruct.Lines)))

		for i := 0; i < len(data); i++ {
			d := data[i]
			for b := 0; b < 8; b++ {
				v := d&(1<<uint(b)) > 0
				pos := i*8 + b
				if v {
					img.Pix[pos] = 255
				} else {
					img.Pix[pos] = 0
				}
			}
		}

		err = jpeg.Encode(f, img, &jpeg.Options{
			Quality: 100,
		})

		if err != nil {
			return err
		}

		return nil
	}

	return fmt.Errorf("image Bit Depth not supported: %d\n", xh.ImageStructureHeader.BitsPerPixel)
}

func DumpDirectly(newFileName string, data []byte) error {
	err := ioutil.WriteFile(newFileName, data, os.ModePerm)
	if err != nil {
		fmt.Printf("Error saving file: %s\n", err)
		return err
	}

	return nil
}

func DumpImage(filename string) error {
	xh, err := XRIT.ParseFile(filename)
	if err != nil {
		return err
	}

	if xh.PrimaryHeader.FileTypeCode != PacketData.IMAGE {
		return fmt.Errorf("the specified file is not an image container")
	}

	offset := xh.PrimaryHeader.HeaderLength

	f, err := os.Open(filename)

	if err != nil {
		return err
	}

	_, err = f.Seek(int64(offset), io.SeekStart)
	if err != nil {
		return err
	}

	data, err := ioutil.ReadAll(f)
	if err != nil {
		return err
	}

	newFileName := strings.Replace(filename, ".lrit", compressionToExtension[xh.Compression()], 1)
	SLog.Info("%s Image, dumping to %s\n", PacketData.GetCompressionTypeString(xh.Compression()), newFileName)

	if xh.Compression() == PacketData.NO_COMPRESSION || xh.Compression() == PacketData.LRIT_RICE {
		return DumpRaw(newFileName, data, xh)
	} else {
		return DumpDirectly(newFileName, data)
	}
}
