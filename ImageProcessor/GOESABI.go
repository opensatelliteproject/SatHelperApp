package ImageProcessor

import (
	"github.com/opensatelliteproject/SatHelperApp/ImageProcessor/ImageTools"
	"github.com/opensatelliteproject/SatHelperApp/ImageProcessor/Projector"
	"github.com/opensatelliteproject/SatHelperApp/ImageProcessor/Structs"
	"github.com/opensatelliteproject/SatHelperApp/Logger"
	"github.com/opensatelliteproject/SatHelperApp/Tools"
	"github.com/opensatelliteproject/SatHelperApp/XRIT"
	"github.com/opensatelliteproject/SatHelperApp/XRIT/Geo"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"
)

var NOAANameRegex = regexp.MustCompile(`OR_ABI-(.*)-(.*)_(.*)_(s.*).*`)

const visChan = "C02"
const irChan = "C14"

func ProcessGOESABI(ip *ImageProcessor, filename string, xh *XRIT.Header) {
	if xh.SubProduct().ID == 0 && !(xh.SegmentIdentificationHeader != nil && xh.SegmentIdentificationHeader.COMS1) { // Mesoscales and unknown data
		PlainLRITImage(ip, filename, xh)
		return
	}
	curveManipulator := ImageTools.GetVisibleCurveManipulator()

	basename := path.Base(filename)
	name := strings.TrimSuffix(basename, filepath.Ext(basename))

	if xh.SegmentIdentificationHeader != nil && xh.SegmentIdentificationHeader.COMS1 {
		// Remove last 3 bytes from the name
		name = name[:len(name)-3]
	}

	if ip.MultiSegmentCache[name] == nil {
		ip.MultiSegmentCache[name] = Structs.MakeMultiSegmentImage(name, xh.SubProduct().ID, int(xh.SegmentIdentificationHeader.ImageID))
	}

	ms := ip.MultiSegmentCache[name]
	ms.PutSegment(filename, xh)

	if ms.Done() {
		SLog.Info("Got all segments for %s", name)
		err, outname := ImageTools.DumpMultiSegment(ms, ip.GetMapDrawer(), curveManipulator, ip.reproject, ip.enhance, ip.metadata)
		if err != nil {
			SLog.Error("Error dumping Multi Segment Image %s: %s", name, err)
		}

		SLog.Info("New image %s", outname)

		if len(ip.cutRegions) > 0 {
			err = ImageTools.DumpCutRegion(ms, ip.GetMapDrawer(), ip.GetMapCutter(), curveManipulator, ip.reproject, ip.enhance, ip.metadata, ip.cutRegions)
			if err != nil {
				SLog.Error("Error processing region cuts: %s", err)
			}
		}

		delete(ip.MultiSegmentCache, name)

		if purgeFiles {
			if !ip.GetFalseColor() || !ms.FirstSegmentHeader.IsFalseColorPiece() {
				ms.Purge()
			}
		}

		if ms.FirstSegmentHeader.IsFalseColorPiece() {
			folder := path.Dir(ms.FirstSegmentFilename)
			nomapFile := path.Join(folder, ImageTools.GetNoMapName(ms.Name, ""))
			ProcessFalseColor(ip, ms.FirstSegmentHeader, nomapFile)
		}
	}
}

func ProcessFalseColor(ip *ImageProcessor, xh *XRIT.Header, filename string) {
	if !NOAANameRegex.MatchString(path.Base(filename)) {
		SLog.Debug("Filename %s does not match noaa name. Not continuing...", filename)
		return
	}

	// 0 => full string, 1 => Level-Product, 2 => Mode/Channel, 3 => Satellite Name, 4 => Group, 5 => File Stamp
	groups := NOAANameRegex.FindStringSubmatch(path.Base(filename))

	mdch := groups[2]
	md := mdch[:2]
	name := groups[4]

	vismdch := md + visChan
	irmdch := md + irChan

	visFilename := strings.Replace(filename, mdch, vismdch, -1)
	irFilename := strings.Replace(filename, mdch, irmdch, -1)
	fsclrFileName := strings.Replace(filename, mdch, md+"C99", -1)
	fsclrFileName = strings.Replace(fsclrFileName, "-nomap", "", -1)

	if !Tools.Exists(visFilename) || !Tools.Exists(irFilename) {
		// Not Ready
		return
	}

	if Tools.Exists(fsclrFileName) {
		SLog.Debug("Skipping generating false color. File exists...")
		return
	}

	SLog.Info("Generating false color for %s", name)

	vis, err := ImageTools.LoadImageGrayScale(visFilename)
	if err != nil {
		SLog.Error("Error loading visible image at %s: %s", visFilename, err)
		return
	}

	ir, err := ImageTools.LoadImageGrayScale(irFilename)
	if err != nil {
		SLog.Error("Error loading infrared image at %s: %s", irFilename, err)
		return
	}

	falseLut := ImageTools.GetFalseColorLUT()

	enhancer := ImageTools.MakeImageEnhancerEmpty(false)

	fsclr, err := falseLut.Apply(vis, ir)
	if err != nil {
		SLog.Error("Error applying false color LUT: %s", err)
		return
	}

	gc, err := Geo.MakeGeoConverterFromXRIT(xh)

	if err == nil {
		if ip.GetReproject() {
			SLog.Debug("Reprojection Enabled, reprojecting FalseColor")
			proj := Projector.MakeProjector(gc)
			fsclr = proj.ReprojectLinearMultiThread(fsclr)
			gc = Projector.MakeLinearConverter(fsclr.Bounds().Dx(), fsclr.Bounds().Dy(), gc)
		}
		mapDrawer := ip.GetMapDrawer()

		if mapDrawer != nil {
			SLog.Debug("Map Drawer Enabled, drawing at FalseColor")
			mapDrawer.DrawMap(fsclr, gc)
		}
	} else {
		SLog.Error("Cannot crate GeoConverter: %s", err)
	}
	xh.NOAASpecificHeader.ProductSubID = 99 // False Color

	if ip.metadata {
		im2, err := enhancer.DrawMetaWithoutScale("", fsclr, xh)
		if err != nil {
			SLog.Error("Error drawing metadata: %s", err)
		}

		if im2 != nil {
			fsclr = im2
		}
	}

	metaName := strings.Replace(fsclrFileName, ".png", ".json", -1)
	err = ioutil.WriteFile(metaName, []byte(xh.ToJSON()), os.ModePerm)
	if err != nil {
		SLog.Error("Cannot write Meta file %s: %s", metaName, err)
	}

	err = ImageTools.SaveImage(fsclrFileName, fsclr)
	if err != nil {
		SLog.Error("Error saving false color image to %s: %s", fsclrFileName, err)
		return
	}

	SLog.Debug("Removing %s", visFilename)
	err = os.Remove(visFilename)
	if err != nil {
		SLog.Error("Error erasing %s: %s", visFilename, err)
	}

	SLog.Debug("Removing %s", irFilename)
	err = os.Remove(irFilename)
	if err != nil {
		SLog.Error("Error erasing %s: %s", irFilename, err)
	}

	SLog.Info("New image %s", fsclrFileName)
}
