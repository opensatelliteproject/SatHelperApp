package main

import (
	"github.com/opensatelliteproject/SatHelperApp"
	"github.com/opensatelliteproject/SatHelperApp/ImageProcessor"
	"github.com/opensatelliteproject/SatHelperApp/Logger"
	"github.com/opensatelliteproject/SatHelperApp/Tools"
	"github.com/opensatelliteproject/SatHelperApp/XRIT"
	"gopkg.in/alecthomas/kingpin.v2"
	"io/ioutil"
	"path"
	"strings"
)

func ProcessFile(filename string, ip *ImageProcessor.ImageProcessor) {

	SLog.Debug("Processing %s", filename)
	xh, err := XRIT.ParseFile(filename)

	if err != nil {
		SLog.Error("Error processing file %s: %s", filename, err)
		return
	}

	ImageProcessor.ProcessGOESABI(ip, filename, xh)
}

func main() {
	kingpin.Version(SatHelperApp.GetVersion())

	reproject := kingpin.Flag("linear", "Reproject to linear").Bool()
	drawMap := kingpin.Flag("drawMap", "Draw Map Overlay").Bool()
	falseColor := kingpin.Flag("falsecolor", "Generate False Color Image").Bool()
	purge := kingpin.Flag("purge", "Purge LRIT files after generating").Bool()
	files := kingpin.Arg("filenames", "File names to dump image").Required().ExistingFilesOrDirs()

	kingpin.Parse()

	ip := ImageProcessor.MakeImageProcessor()
	ip.SetDrawMap(*drawMap)
	ip.SetReproject(*reproject)
	ip.SetFalseColor(*falseColor)

	ImageProcessor.SetPurgeFiles(*purge)

	for _, v := range *files {
		if Tools.IsDir(v) {
			ffiles, err := ioutil.ReadDir(v)
			if err != nil {
				SLog.Error("Cannot read folder %s: %s", v, err)
				continue
			}

			for _, v2 := range ffiles {
				if !v2.IsDir() && strings.Contains(v2.Name(), ".lrit") {
					ProcessFile(path.Join(v, v2.Name()), ip)
				} else {
					SLog.Debug("Skipping file %s, does not end with .lrit", v2.Name())
				}
			}
			continue
		}
		if strings.Contains(v, ".lrit") {
			ProcessFile(v, ip)
		} else {
			SLog.Debug("Skipping file %s, does not end with .lrit", v)
		}
	}
}
