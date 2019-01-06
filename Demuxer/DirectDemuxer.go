package Demuxer

import (
	"github.com/OpenSatelliteProject/SatHelperApp/Logger"
	"github.com/OpenSatelliteProject/SatHelperApp/ccsds"
)

type DirectDemuxer struct {
	demux *ccsds.Demuxer
}

func MakeDirectDemuxer() *DirectDemuxer {
	d := &DirectDemuxer{}

	d.demux = ccsds.MakeDemuxer()
	d.demux.SetOnFrameLost(func(channelId, currentFrame, lastFrame int) {
		SLog.Info("Lost Frames for channel %d: %d", channelId, currentFrame-lastFrame-1)
	})

	d.demux.SetOnNewVCID(func(channelId int) {
		SLog.Info("New Channel: %d", channelId)
	})
	return d
}

func (f *DirectDemuxer) Init()  {}
func (f *DirectDemuxer) Start() {}
func (f *DirectDemuxer) Stop()  {}
func (f *DirectDemuxer) SendData(data []byte) {
	f.demux.WriteBytes(data)
}
func (f *DirectDemuxer) GetName() string {
	return "Direct"
}
