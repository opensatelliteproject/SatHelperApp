package ImageProcessor

import (
	"github.com/OpenSatelliteProject/SatHelperApp/ImageProcessor/ImageTools"
	"github.com/OpenSatelliteProject/SatHelperApp/Logger"
	"github.com/OpenSatelliteProject/SatHelperApp/XRIT"
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
