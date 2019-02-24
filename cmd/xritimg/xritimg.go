package main

import (
	"github.com/OpenSatelliteProject/SatHelperApp"
	"github.com/OpenSatelliteProject/SatHelperApp/ImageProcessor/ImageTools"
	"github.com/OpenSatelliteProject/SatHelperApp/Logger"
	"gopkg.in/alecthomas/kingpin.v2"
)

func main() {
	kingpin.Version(SatHelperApp.GetVersion())

	files := kingpin.Arg("filename", "File name to dump image").Required().ExistingFiles()

	kingpin.Parse()

	for _, v := range *files {
		err := ImageTools.DumpImage(v)
		if err != nil {
			SLog.Error("Error processing %s: %s", v, err)
		}
	}
}
