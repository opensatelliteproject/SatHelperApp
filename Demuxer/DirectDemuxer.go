package Demuxer

import (
	"github.com/OpenSatelliteProject/SatHelperApp/Logger"
	"github.com/OpenSatelliteProject/SatHelperApp/ccsds"
	"github.com/logrusorgru/aurora"
)

type DirectDemuxer struct {
	demux *ccsds.Demuxer
}

func MakeDirectDemuxer(outFolder, tmpFolder string) *DirectDemuxer {
	d := &DirectDemuxer{}

	d.demux = ccsds.MakeDemuxer()
	d.demux.SetOutputFolder(outFolder)
	d.demux.SetTemporaryFolder(tmpFolder)

	SLog.Info("Starting direct Demuxer with: ")
	SLog.Info(" Output Folder: %s", aurora.Bold(outFolder).Green())
	SLog.Info(" Temporary Folder: %s", aurora.Bold(tmpFolder).Green())

	d.demux.SetOnFrameLost(func(channelId, currentFrame, lastFrame int) {
		SLog.Info("Lost Frames for channel %d: %d", channelId, currentFrame-lastFrame-1)
	})

	d.demux.SetOnNewVCID(func(channelId int) {
		SLog.Debug("New Channel: %d", channelId)
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
