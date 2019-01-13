package main

import (
	"github.com/OpenSatelliteProject/SatHelperApp/ccsds"
	"log"
	"os"
)

func main() {
	//debugFrames := "/media/ELTN/tmp/demuxdump-1546741011.bin"
	debugFrames := "/media/ELTN/tmp/demuxdump-1490627438.bin"
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
	dm.SetOnFrameLost(func(channelId, currentFrame, lastFrame int) {
		delta := currentFrame - lastFrame
		log.Printf("Lost %d frames in channel %d\n", delta, channelId)
	})

	// region Skip DCS
	dm.AddSkipVCID(30)
	dm.AddSkipVCID(31)
	dm.AddSkipVCID(32)
	// endregion

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
		//time.Sleep(time.Millisecond * 10)
	}
}
