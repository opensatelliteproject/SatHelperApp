package main

import (
	"github.com/opensatelliteproject/SatHelperApp"
	"github.com/opensatelliteproject/SatHelperApp/ImageProcessor/ImageTools"
	"github.com/opensatelliteproject/SatHelperApp/ImageProcessor/Structs"
	"github.com/opensatelliteproject/SatHelperApp/Logger"
	"github.com/opensatelliteproject/SatHelperApp/XRIT"
	"gopkg.in/alecthomas/kingpin.v2"
	"os"
	"path"
	"path/filepath"
	"strings"
)

func main() {
	kingpin.Version(SatHelperApp.GetVersion())

	files := kingpin.Arg("filenames", "File names to dump image").Required().ExistingFiles()

	kingpin.Parse()

	f := (*files)[0]

	xh, err := XRIT.ParseFile(f)

	if err != nil {
		SLog.Error("Error processing file %s: %s", f, err)
		os.Exit(1)
	}

	basename := path.Base(f)
	name := strings.TrimSuffix(basename, filepath.Ext(basename))

	msi := Structs.MakeMultiSegmentImage(name, int(xh.NOAASpecificHeader.ProductSubID), int(xh.SegmentIdentificationHeader.ImageID))
	msi.PutSegment(f, xh)

	for _, v := range *files {
		if v != f {
			basename := path.Base(v)
			_name := strings.TrimSuffix(basename, filepath.Ext(basename))

			if _name != name {
				SLog.Warn("Skipping %s since its not the same group of %s", _name, name)
				continue
			}

			xh, err = XRIT.ParseFile(v)

			if err != nil {
				SLog.Error("Error processing file %s: %s", f, err)
				os.Exit(1)
			}

			msi.PutSegment(v, xh)
		}
	}

	if msi.Done() {
		SLog.Info("Got all segments, generating image.")
		err, outname := ImageTools.DumpMultiSegment(msi)
		if err != nil {
			SLog.Error("Error dumping image: %s", err)
			os.Exit(1)
		}
		SLog.Info("Output Image: %s", outname)
	} else {
		SLog.Error("Not all segments arrived. Expected %d got %d", msi.MaxSegments, len(msi.Files))
	}
}
