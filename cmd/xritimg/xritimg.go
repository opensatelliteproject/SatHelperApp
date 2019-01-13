package main

import (
	"fmt"
	"github.com/OpenSatelliteProject/SatHelperApp"
	"github.com/OpenSatelliteProject/SatHelperApp/XRIT"
	"github.com/OpenSatelliteProject/SatHelperApp/XRIT/PacketData"
	"gopkg.in/alecthomas/kingpin.v2"
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

func dumpDirectly(newFileName string, data []byte) {
	err := ioutil.WriteFile(newFileName, data, os.ModePerm)
	if err != nil {
		fmt.Printf("Error saving file: %s\n", err)
		os.Exit(1)
	}
}

func dumpRaw(newFileName string, data []byte, xh *XRIT.Header) {
	imgStruct := xh.ImageStructureHeader

	f, err := os.Create(newFileName)
	if err != nil {
		fmt.Printf("Error creating file %s: %s\n", newFileName, err)
		os.Exit(1)
	}

	defer f.Close()

	if imgStruct.BitsPerPixel == 8 {
		totalPixels := int(imgStruct.Columns) * int(imgStruct.Lines)
		if totalPixels > len(data) {
			missingBytes := totalPixels - len(data)
			fmt.Printf("Missing %d bytes in image data.\n", missingBytes)
			data = append(data, make([]byte, missingBytes)...)
		}
		img := &image.Gray{Pix: data, Stride: int(imgStruct.Columns), Rect: image.Rect(0, 0, int(imgStruct.Columns), int(imgStruct.Lines))}

		err = jpeg.Encode(f, img, &jpeg.Options{
			Quality: 100,
		})

		if err != nil {
			fmt.Printf("Error encoding file %s: %s\n", newFileName, err)
			os.Exit(1)
		}

		return
	}

	if imgStruct.BitsPerPixel == 1 {
		totalPixels := int(imgStruct.Columns) * int(imgStruct.Lines)
		if totalPixels > len(data)/8 {
			missingBytes := (totalPixels - len(data)) / 8
			fmt.Printf("Missing %d bytes in image data.\n", missingBytes)
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
			fmt.Printf("Error encoding file %s: %s\n", newFileName, err)
			os.Exit(1)
		}
		return
	}

	fmt.Printf("Image Bit Depth not supported: %d\n", xh.ImageStructureHeader.BitsPerPixel)
	os.Exit(1)
}

func dumpImage(filename string) {
	xh, err := XRIT.ParseFile(filename)
	if err != nil {
		fmt.Printf("Error parsing file %s: %s\n", filename, err)
		os.Exit(1)
	}

	if xh.PrimaryHeader.FileTypeCode != PacketData.IMAGE {
		fmt.Printf("The specified file is not an image container.")
		os.Exit(1)
	}

	offset := xh.PrimaryHeader.HeaderLength

	f, err := os.Open(filename)

	if err != nil {
		fmt.Printf("Error parsing file %s: %s\n", filename, err)
		os.Exit(1)
	}

	_, err = f.Seek(int64(offset), io.SeekStart)
	if err != nil {
		fmt.Printf("Error parsing file %s: %s\n", filename, err)
		os.Exit(1)
	}

	data, err := ioutil.ReadAll(f)
	if err != nil {
		fmt.Printf("Error parsing file %s: %s\n", filename, err)
		os.Exit(1)
	}

	newFileName := strings.Replace(filename, ".lrit", compressionToExtension[xh.Compression()], 1)
	fmt.Printf("%s Image, dumping to %s\n", PacketData.GetCompressionTypeString(xh.Compression()), newFileName)

	if xh.Compression() == PacketData.NO_COMPRESSION || xh.Compression() == PacketData.LRIT_RICE {
		dumpRaw(newFileName, data, xh)
	} else {
		dumpDirectly(newFileName, data)
	}
}

func main() {
	kingpin.Version(SatHelperApp.GetVersion())

	files := kingpin.Arg("filename", "File name to dump image").Required().ExistingFiles()

	kingpin.Parse()

	for _, v := range *files {
		dumpImage(v)
	}
}
