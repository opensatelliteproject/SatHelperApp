package main

import (
	"github.com/opensatelliteproject/SatHelperApp"
	"github.com/opensatelliteproject/SatHelperApp/ImageProcessor"
	"github.com/opensatelliteproject/SatHelperApp/Logger"
	"github.com/opensatelliteproject/SatHelperApp/XRIT"
	"gopkg.in/alecthomas/kingpin.v2"
)

func main() {
	kingpin.Version(SatHelperApp.GetVersion())

	reproject := kingpin.Flag("linear", "Reproject to linear").Bool()
	drawMap := kingpin.Flag("drawMap", "Draw Map Overlay").Bool()
	falseColor := kingpin.Flag("falsecolor", "Generate False Color Image").Bool()
	files := kingpin.Arg("filenames", "File names to dump image").Required().ExistingFiles()

	kingpin.Parse()

	ip := ImageProcessor.MakeImageProcessor()
	ip.SetDrawMap(*drawMap)
	ip.SetReproject(*reproject)
	ip.SetFalseColor(*falseColor)

	for _, v := range *files {
		SLog.Debug("Processing %s", v)
		xh, err := XRIT.ParseFile(v)

		if err != nil {
			SLog.Error("Error processing file %s: %s", v, err)
			continue
		}

		ImageProcessor.ProcessGOESABI(ip, v, xh)
	}
}
