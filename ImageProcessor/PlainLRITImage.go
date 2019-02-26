package ImageProcessor

import (
	"github.com/opensatelliteproject/SatHelperApp/ImageProcessor/ImageTools"
	"github.com/opensatelliteproject/SatHelperApp/Logger"
	"github.com/opensatelliteproject/SatHelperApp/XRIT"
	"os"
)

func PlainLRITImage(_ *ImageProcessor, filename string, _ *XRIT.Header) {
	// Plain images we just need to dump to jpeg.
	err := ImageTools.DumpImage(filename)
	if err != nil {
		SLog.Error("Error processing %s: %s", filename, err)
	}

	if purgeFiles {
		err = os.Remove(filename)
		if err != nil {
			SLog.Error("Error erasing %s: %s", filename, err)
		}
	}
}
