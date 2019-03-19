package main

import (
	"github.com/opensatelliteproject/SatHelperApp/DSP"
	"github.com/opensatelliteproject/SatHelperApp/Display"
	"github.com/opensatelliteproject/SatHelperApp/RPC/sathelperapp"
)

type baseSource struct{}

var rpcSource = &baseSource{}

func (s *baseSource) GetStatistics() (*sathelperapp.StatData, error) {
	stats := DSP.GetStats()
	return &sathelperapp.StatData{
		SignalQuality: uint32(stats.SignalQuality),

		SignalLocked:         stats.FrameLock == 1,
		ChannelPackets:       stats.ReceivedPacketsPerChannel[:],
		RsErrors:             stats.RsErrors[:],
		SyncWord:             stats.SyncWord[:],
		Scid:                 int32(stats.SCID),
		Vcid:                 int32(stats.VCID),
		DecoderFifoUsage:     int32(stats.DecoderFifoUsage),
		DemodulatorFifoUsage: int32(stats.DemodulatorFifoUsage),
		ViterbiErrors:        int32(stats.VitErrors),
		FrameSize:            int32(stats.FrameBits),
		PhaseCorrection:      int32(stats.PhaseCorrection),
		SyncCorrelation:      int32(stats.SyncCorrelation),
		CenterFrequency:      DSP.Device.GetCenterFrequency() + uint32(DSP.GetCostasFrequency()),
		Demuxer:              DSP.SDemuxer.GetName(),
		Mode:                 DSP.CurrentConfig.Base.Mode,
		Device:               DSP.Device.GetName(),
	}, nil
}

func (s *baseSource) GetConsoleLines() (*sathelperapp.ConsoleData, error) {
	return &sathelperapp.ConsoleData{
		ConsoleLines: Display.GetConsoleLines(),
	}, nil
}
