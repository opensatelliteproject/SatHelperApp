package demuxReplay

import (
	"github.com/OpenSatelliteProject/SatHelperApp/ccsds"
	"os"
)

func main() {
	debugFrames := "/media/ELTN/tmp/demuxdump-1546741011.bin"

	f, err := os.Open(debugFrames)

	if err != nil {
		panic(err)
	}
	defer f.Close()

	finfo, err := f.Stat()

	if err != nil {
		panic(err)
	}

	size := finfo.Size()

	dm := ccsds.MakeDemuxer()

	bytesRead := int64(0)
	buffer := make([]byte, 892)

	for bytesRead < size {
		n, err := f.Read(buffer)
		if err != nil {
			panic(err)
		}

		if n == 892 {
			dm.WriteBytes(buffer)
		} else {
			panic("WAIT")
		}
		bytesRead += int64(n)
	}
}
