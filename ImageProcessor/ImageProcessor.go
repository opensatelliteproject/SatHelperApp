package ImageProcessor

import (
	"github.com/opensatelliteproject/SatHelperApp/ImageProcessor/ImageTools"
	"github.com/opensatelliteproject/SatHelperApp/ImageProcessor/MapDrawer"
	"github.com/opensatelliteproject/SatHelperApp/ImageProcessor/Structs"
	"github.com/opensatelliteproject/SatHelperApp/Logger"
	"github.com/opensatelliteproject/SatHelperApp/XRIT"
	"github.com/opensatelliteproject/SatHelperApp/XRIT/NOAAProductID"
	"github.com/opensatelliteproject/SatHelperApp/XRIT/PacketData"
	"sync"
)

var purgeFiles = false

type ImageProcessor struct {
	sync.Mutex
	MultiSegmentCache map[string]*Structs.MultiSegmentImage
	mapDrawer         *MapDrawer.MapDrawer
	reproject         bool
	drawmap           bool
	falsecolor        bool
}

func MakeImageProcessor() *ImageProcessor {
	return &ImageProcessor{
		MultiSegmentCache: make(map[string]*Structs.MultiSegmentImage),
		mapDrawer:         ImageTools.GetDefaultMapDrawer(),
		reproject:         false,
		drawmap:           false,
		falsecolor:        false,
	}
}

func (ip *ImageProcessor) SetFalseColor(fsclr bool) {
	ip.falsecolor = fsclr
	if fsclr {
		SLog.Warn("False color is enabled, so it will also save plain images with no map")
		ImageTools.SetSaveNoMap(true) // Needed for FSCLR
	}

}

func (ip *ImageProcessor) SetDrawMap(drawMap bool) {
	ip.drawmap = drawMap
}

func (ip *ImageProcessor) SetReproject(reproject bool) {
	ip.reproject = reproject
}

func (ip *ImageProcessor) GetFalseColor() bool {
	return ip.falsecolor
}

func (ip *ImageProcessor) GetDrawMap() bool {
	return ip.drawmap
}

func (ip *ImageProcessor) GetReproject() bool {
	return ip.reproject
}

func (ip *ImageProcessor) GetMapDrawer() *MapDrawer.MapDrawer {
	if ip.drawmap {
		return ip.mapDrawer
	}

	return nil
}

func (ip *ImageProcessor) ProcessImage(filename string) {
	ip.Lock()
	defer ip.Unlock()

	xh, err := XRIT.ParseFile(filename)
	if err != nil {
		SLog.Error("Error parsing file %s: %s", filename, err)
		return
	}

	if xh.PrimaryHeader.FileTypeCode != PacketData.IMAGE {
		return
	}

	switch xh.NOAASpecificHeader.ProductID {
	case NOAAProductID.GOES16_ABI, NOAAProductID.GOES17_ABI:
		ProcessGOESABI(ip, filename, xh)
	}

	ip.checkExpired()
}

func (ip *ImageProcessor) checkExpired() {
	for k, v := range ip.MultiSegmentCache {
		if v.Expired() {
			SLog.Warn("Image %s timed out waiting segments. Removing from cache.", k)
			delete(ip.MultiSegmentCache, k)
			if purgeFiles {
				v.Purge()
			}
		}
	}
}

func SetPurgeFiles(purge bool) {
	purgeFiles = purge
	SLog.Info("Set Purge Files changed to %v", purge)
}
