package ImageProcessor

import (
	"github.com/opensatelliteproject/SatHelperApp/ImageProcessor/ImageTools"
	"github.com/opensatelliteproject/SatHelperApp/ImageProcessor/Structs"
	"github.com/opensatelliteproject/SatHelperApp/Logger"
	"github.com/opensatelliteproject/SatHelperApp/XRIT"
	"path"
	"path/filepath"
	"strings"
)

func ProcessGOESABI(ip *ImageProcessor, filename string, xh *XRIT.Header) {
	if xh.NOAASpecificHeader.ProductSubID == 0 { // Mesoscales and unknown data
		PlainLRITImage(ip, filename, xh)
		return
	}

	basename := path.Base(filename)
	name := strings.TrimSuffix(basename, filepath.Ext(basename))

	if ip.MultiSegmentCache[name] == nil {
		ip.MultiSegmentCache[name] = Structs.MakeMultiSegmentImage(name, int(xh.NOAASpecificHeader.ProductSubID), int(xh.SegmentIdentificationHeader.ImageID))
	}

	ms := ip.MultiSegmentCache[name]
	ms.PutSegment(filename, xh)

	if ms.Done() {
		SLog.Info("Got all segments for %s", name)
		err, outname := ImageTools.DumpMultiSegment(ms, ip.GetMapDrawer(), ip.reproject)
		if err != nil {
			SLog.Error("Error dumping Multi Segment Image %s: %s", name, err)
		}

		SLog.Info("New image %s", outname)

		delete(ip.MultiSegmentCache, name)

		if purgeFiles {
			ms.Purge()
		}
	}
}
