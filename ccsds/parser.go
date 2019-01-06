package ccsds

import (
	"sync"
)

type Demuxer struct {
	sync.Mutex
	frameSize   int
	frameBuffer []byte

	transports map[int]*TransportParser
	lastFrame  map[int]int
	framesLost map[int]uint64

	cbNewVCID     func(int)
	cbOnFrameLost func(channelId, currentFrame, lastFrame int)

	fileAssembler *FileAssembler
}

func MakeDemuxer() *Demuxer {
	return &Demuxer{
		frameSize:     892,
		lastFrame:     make(map[int]int),
		framesLost:    make(map[int]uint64),
		frameBuffer:   make([]byte, 0),
		transports:    make(map[int]*TransportParser),
		fileAssembler: MakeFileAssembler(),
	}
}

func (dm *Demuxer) SetTemporaryFolder(folder string) {
	dm.fileAssembler.SetTemporaryFolder(folder)
}

func (dm *Demuxer) SetOutputFolder(folder string) {
	dm.fileAssembler.SetOutputFolder(folder)
}

func (dm *Demuxer) SetOnFrameLost(cb func(channelId, currentFrame, lastFrame int)) {
	dm.Lock()
	dm.cbOnFrameLost = cb
	dm.Unlock()
}

func (dm *Demuxer) SetOnNewVCID(cb func(channelId int)) {
	dm.Lock()
	dm.cbNewVCID = cb
	dm.Unlock()
}

func (dm *Demuxer) WriteBytes(data []byte) {
	dm.frameBuffer = append(dm.frameBuffer, data...)
	dm.parse()
}

func (dm *Demuxer) onMSDU(msdu *MSDU) {
	dm.fileAssembler.PutMSDU(msdu)
}

func (dm *Demuxer) incFrameLost(channelId, count int) {
	if _, ok := dm.framesLost[channelId]; ok {
		dm.framesLost[channelId] += uint64(count)
	} else {
		dm.framesLost[channelId] = uint64(count)
	}
}

func (dm *Demuxer) checkLostFrameAndSave(channelId, currentFrame int) int {
	framesLost := 0
	if v, ok := dm.lastFrame[channelId]; ok {
		d := int(int64(currentFrame) - int64(v) - 1)
		if v > currentFrame {
			// TODO: Frame Backwards
		} else if d > 0 {
			framesLost = d
			if dm.cbOnFrameLost != nil {
				dm.cbOnFrameLost(channelId, currentFrame, v)
			}
		}
	}

	dm.lastFrame[channelId] = currentFrame

	return framesLost
}

func (dm *Demuxer) parse() {
	for len(dm.frameBuffer) >= dm.frameSize {
		frame := dm.frameBuffer[:dm.frameSize]
		dm.frameBuffer = dm.frameBuffer[dm.frameSize:]

		cd := MakeVCDU(frame)

		lostFrames := dm.checkLostFrameAndSave(cd.VCID(), cd.Counter())

		if lostFrames > 0 {
			dm.incFrameLost(cd.VCID(), lostFrames)
		}

		if cd.VCID() != 63 { // Skip Fill Channel
			if dm.transports[cd.VCID()] == nil {
				dm.transports[cd.VCID()] = MakeTransportParser(cd.VCID(), dm.onMSDU)
				if dm.cbNewVCID != nil {
					dm.cbNewVCID(cd.VCID())
				}
			}

			dm.transports[cd.VCID()].WriteChannelData(cd)
		}
	}
}
