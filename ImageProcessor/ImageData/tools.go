package ImageData

import (
	"github.com/opensatelliteproject/SatHelperApp/ImageProcessor/ImageTools"
	"github.com/opensatelliteproject/SatHelperApp/ImageProcessor/MapDrawer"
	"github.com/opensatelliteproject/SatHelperApp/Logger"
	"github.com/opensatelliteproject/SatHelperApp/XRIT"
	"io/ioutil"
	"os"
	"path"
	"strings"
)

const shpFileName = "ne_50m_admin_0_countries.shp"
const dbfFileName = "ne_50m_admin_0_countries.dbf"
const falseColorLutName = "wx-star.com_GOES-R_ABI_False-Color-LUT.png"

var mapDrawer *MapDrawer.MapDrawer
var fsclrLut *ImageTools.Lut2D
var visCurve *ImageTools.CurveManipulator

var tempLut map[string]*ImageTools.Lut1D

// ExtractShapeFiles extracts the shapefiles to temp folder and return path for shp file
func ExtractShapeFiles() (string, error) {

	folder, err := ioutil.TempDir("", "satHelperShapes")
	if err != nil {
		return "", err
	}
	SLog.Debug("Extracting ShapeFiles to %s", folder)

	shpFileData, err := Asset(shpFileName)
	if err != nil {
		return "", err
	}
	dbfFileData, err := Asset(dbfFileName)
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

func GetVisibleCurveManipulator() *ImageTools.CurveManipulator {
	if visCurve == nil {
		visCurve = ImageTools.MakeDefaultCurveManipulator()
	}

	return visCurve
}

func GetFalseColorLUT() *ImageTools.Lut2D {
	if fsclrLut == nil {
		lutData, err := Asset(falseColorLutName)
		if err != nil {
			SLog.Error("Cannot load False Color LUT data: %s", err)
			return nil
		}

		lut2d, err := ImageTools.MakeLut2DFromMemory(lutData)

		if err != nil {
			SLog.Error("Error creating False Color LUT: %s", err)
			return nil
		}

		fsclrLut = lut2d
	}

	return fsclrLut
}

func GetTemperatureLUT(xh *XRIT.Header) *ImageTools.Lut1D {
	if xh.ImageDataFunctionHash == "" {
		return nil
	}

	if tempLut == nil {
		tempLut = map[string]*ImageTools.Lut1D{}
	}

	if tempLut[xh.ImageDataFunctionHash] == nil {
		colorLut := ScaleLutToColor(minV, scaleFact, xh.GetTemperatureLUT(), TemperatureScaleLUT)
		lut1d, err := ImageTools.MakeLut1DFromColors(colorLut)
		if err != nil {
			SLog.Error("Error creating Temperature LUT: %s", err)
			return nil
		}
		tempLut[xh.ImageDataFunctionHash] = lut1d
	}

	return tempLut[xh.ImageDataFunctionHash]
}
