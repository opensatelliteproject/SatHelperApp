package main

import (
	"fmt"
	"github.com/opensatelliteproject/SatHelperApp"
	"github.com/opensatelliteproject/SatHelperApp/ImageProcessor"
	"github.com/opensatelliteproject/SatHelperApp/ImageProcessor/ImageTools"
	"github.com/opensatelliteproject/SatHelperApp/Logger"
	"github.com/opensatelliteproject/SatHelperApp/Tools"
	"github.com/opensatelliteproject/SatHelperApp/XRIT"
	"github.com/richardwilkes/toolbox/xio/fs"
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

	b := path.Base(filename)

	if b != xh.Filename() {
		p := path.Join(path.Dir(filename), xh.Filename())
		SLog.Info("Moving file %s to %s", filename, p)
		err = fs.MoveFile(filename, p)
		if err != nil {
			SLog.Error("Error moving %s to %s: %s", filename, p, err)
		}
		filename = p
	}

	ImageProcessor.ProcessGOESABI(ip, filename, xh)
}

// ReplaceAll returns a copy of the string s with all
// non-overlapping instances of old replaced by new.
// If old is empty, it matches at the beginning of the string
// and after each UTF-8 sequence, yielding up to k+1 replacements
// for a k-rune string.
func ReplaceAll(s, old, new string) string {
	return strings.Replace(s, old, new, -1)
}

func main() {
	kingpin.Version(SatHelperApp.GetVersion())

	reproject := kingpin.Flag("linear", "Reproject to linear").Bool()
	drawMap := kingpin.Flag("drawMap", "Draw Map Overlay").Bool()
	falseColor := kingpin.Flag("falsecolor", "Generate False Color Image").Bool()
	metadata := kingpin.Flag("metadata", "Generate Overlays with Metadata").Bool()
	enhance := kingpin.Flag("enhance", "Output Enhanced Infrared Images").Bool()
	purge := kingpin.Flag("purge", "Purge LRIT files after generating").Bool()
	regions := kingpin.Flag("region", "Regions to cut by name (use --list-regions to see all available)").Strings()
	listRegions := kingpin.Flag("list-regions", "List all available regions to cut image").Bool()
	searchRegion := kingpin.Flag("search-region", "Search for a region").String()
	marginPixels := kingpin.Flag("margin-pixels", "Margin Pixels for MapCutter").Default("5").Int()
	files := kingpin.Arg("filenames", "File names to dump image").Required().ExistingFilesOrDirs()

	kingpin.Parse()

	mapCutter := ImageTools.GetDefaultMapCutter()

	if *listRegions {
		sections := mapCutter.ListSections()
		for k, v := range sections {
			fmt.Printf("Section Name: \"%s\"\n\t%s\n", k, ReplaceAll(v.String(), "\n", "\n\t"))
		}
		return
	}

	if searchRegion != nil && *searchRegion != "" {
		sections := mapCutter.SearchSection(*searchRegion)
		for _, k := range sections {
			v, _ := mapCutter.GetSection(k)
			fmt.Printf("Section Name: \"%s\"\n\t%s\n", k, ReplaceAll(v.String(), "\n", "\n\t"))
		}
		return
	}

	ip := ImageProcessor.MakeImageProcessor()
	ip.SetDrawMap(*drawMap)
	ip.SetReproject(*reproject)
	ip.SetFalseColor(*falseColor)
	ip.SetMetadata(*metadata)
	ip.SetEnhance(*enhance)
	ip.SetCutRegions(*regions)
	ip.GetMapCutter().SetMarginPixels(*marginPixels)

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
