package ImageTools

import (
	"fmt"
	"github.com/opensatelliteproject/SatHelperApp/ImageProcessor/ImageData"
	"github.com/opensatelliteproject/SatHelperApp/ImageProcessor/MapCutter"
	"github.com/opensatelliteproject/SatHelperApp/ImageProcessor/MapDrawer"
	"github.com/opensatelliteproject/SatHelperApp/Logger"
	"github.com/opensatelliteproject/SatHelperApp/XRIT"
	"github.com/opensatelliteproject/SatHelperApp/XRIT/PacketData"
	"image"
	"image/draw"
	"image/jpeg"
	"io"
	"io/ioutil"
	"os"
	"path"
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
		SLog.Error("Error saving file: %s\n", err)
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

func LoadImage(filename string) (image.Image, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}

	defer f.Close()

	src, _, err := image.Decode(f)
	if err != nil {
		return nil, err
	}

	return src, nil
}

func LoadImageGrayScale(filename string) (*image.Gray, error) {
	img, err := LoadImage(filename)
	if err != nil {
		return nil, err
	}

	switch i := img.(type) {
	case *image.Gray:
		return i, nil
	default:
		return Image2Gray(img), nil
	}
}

func Image2Gray(img image.Image) *image.Gray {
	gray := image.NewGray(img.Bounds())
	draw.Draw(gray, img.Bounds(), img, img.Bounds().Min, draw.Src)
	return gray
}

const shpFileName = "ne_50m_admin_0_countries.shp"
const dbfFileName = "ne_50m_admin_0_countries.dbf"
const statesShpFileName = "ne_50m_admin_1_states_provinces.shp"
const statesDbfFileName = "ne_50m_admin_1_states_provinces.dbf"
const falseColorLutName = "wx-star.com_GOES-R_ABI_False-Color-LUT.png"

var mapDrawer *MapDrawer.MapDrawer
var fsclrLut *Lut2D
var visCurve *CurveManipulator
var mapCutter *MapCutter.MapCutter

var tempLut map[string]*Lut1D

// ExtractShapeFiles extracts the shapefiles to temp folder and return path for shp file
func ExtractShapeFiles() (string, error) {

	folder, err := ioutil.TempDir("", "satHelperShapes")
	if err != nil {
		return "", err
	}
	SLog.Debug("Extracting ShapeFiles to %s", folder)

	shpFileData, err := ImageData.Asset(shpFileName)
	if err != nil {
		return "", err
	}
	dbfFileData, err := ImageData.Asset(dbfFileName)
	if err != nil {
		return "", err
	}

	err = ioutil.WriteFile(path.Join(folder, shpFileName), shpFileData, os.ModePerm)
	if err != nil {
		return "", err
	}
	err = ioutil.WriteFile(path.Join(folder, dbfFileName), dbfFileData, os.ModePerm)
	if err != nil {
		return "", err
	}

	return path.Join(folder, shpFileName), nil
}

// ExtractStateShapeFiles extracts the shapefiles for states to temp folder and return path for shp file
func ExtractStateShapeFiles() (string, error) {

	folder, err := ioutil.TempDir("", "satHelperShapes")
	if err != nil {
		return "", err
	}
	SLog.Debug("Extracting ShapeFiles to %s", folder)

	shpFileData, err := ImageData.Asset(statesShpFileName)
	if err != nil {
		return "", err
	}
	dbfFileData, err := ImageData.Asset(statesDbfFileName)
	if err != nil {
		return "", err
	}

	err = ioutil.WriteFile(path.Join(folder, statesShpFileName), shpFileData, os.ModePerm)
	if err != nil {
		return "", err
	}
	err = ioutil.WriteFile(path.Join(folder, statesDbfFileName), dbfFileData, os.ModePerm)
	if err != nil {
		return "", err
	}

	return path.Join(folder, statesShpFileName), nil
}

func CleanShapeFiles(shpFile string) {
	SLog.Debug("Cleaning ShapeFiles from %s", shpFile)
	_ = os.Remove(shpFile)
	_ = os.Remove(strings.Replace(shpFile, ".shp", ".dbf", -1))
}

func GetDefaultMapDrawer() *MapDrawer.MapDrawer {
	if mapDrawer == nil {
		shpFile, err := ExtractShapeFiles()
		if err != nil {
			SLog.Error("Error extracting Shape Files: %s", err)
			SLog.Error("Map Drawer will be disabled")
		}

		if shpFile != "" {
			err, mapDrawer = MapDrawer.MakeMapDrawer(shpFile)
			if err != nil {
				SLog.Error("Error creating Map Drawer: %s", err)
			} else {
				CleanShapeFiles(shpFile)
			}
		}
	}

	return mapDrawer
}

func GetDefaultMapCutter() *MapCutter.MapCutter {
	if mapCutter == nil {
		shpFiles := make([]string, 0)

		shpFile, err := ExtractShapeFiles()
		if err != nil {
			SLog.Error("Error extracting Shape Files: %s", err)
			SLog.Error("Map Cutter will be disabled")
		} else {
			shpFiles = append(shpFiles, shpFile)
		}

		stateShpFile, err := ExtractStateShapeFiles()
		if err != nil {
			SLog.Error("Error extracting State Shape Files: %s", err)
			SLog.Error("Map Cutter wont have states")
		} else {
			shpFiles = append(shpFiles, stateShpFile)
		}

		if len(shpFiles) > 0 {
			mapCutter, err = MapCutter.MakeMapCutterFromFiles(shpFiles)
			if err != nil {
				SLog.Error("Error creating Map Cutter: %s", err)
			}
		}

		for _, v := range shpFiles {
			CleanShapeFiles(v)
		}
	}

	return mapCutter
}

func GetVisibleCurveManipulator() *CurveManipulator {
	if visCurve == nil {
		visCurve = MakeDefaultCurveManipulator()
	}

	return visCurve
}

func GetFalseColorLUT() *Lut2D {
	if fsclrLut == nil {
		lutData, err := ImageData.Asset(falseColorLutName)
		if err != nil {
			SLog.Error("Cannot load False Color LUT data: %s", err)
			return nil
		}

		lut2d, err := MakeLut2DFromMemory(lutData)

		if err != nil {
			SLog.Error("Error creating False Color LUT: %s", err)
			return nil
		}

		fsclrLut = lut2d
	}

	return fsclrLut
}

const minV = 173 // Kelvin
const scaleFact = 256.0 / (340 - minV)

func GetTemperatureLUT(xh *XRIT.Header) *Lut1D {
	if xh.ImageDataFunctionHash == "" {
		return nil
	}

	if tempLut == nil {
		tempLut = map[string]*Lut1D{}
	}

	if tempLut[xh.ImageDataFunctionHash] == nil {
		colorLut := ImageData.ScaleLutToColor(minV, scaleFact, xh.GetTemperatureLUT(), ImageData.TemperatureScaleLUT)
		lut1d, err := MakeLut1DFromColors(colorLut)
		if err != nil {
			SLog.Error("Error creating Temperature LUT: %s", err)
			return nil
		}
		tempLut[xh.ImageDataFunctionHash] = lut1d
	}

	return tempLut[xh.ImageDataFunctionHash]
}
