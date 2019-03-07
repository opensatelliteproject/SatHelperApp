package ImageData

import (
	"github.com/opensatelliteproject/SatHelperApp/ImageProcessor/MapDrawer"
	"github.com/opensatelliteproject/SatHelperApp/Logger"
	"io/ioutil"
	"os"
	"path"
	"strings"
)

const shpFileName = "ne_50m_admin_0_countries.shp"
const dbfFileName = "ne_50m_admin_0_countries.dbf"

var mapDrawer *MapDrawer.MapDrawer

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
