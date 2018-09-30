package main

import (
	"github.com/OpenSatelliteProject/SatHelperApp/Demuxer"
	"github.com/OpenSatelliteProject/SatHelperApp/Frontend"
	"github.com/OpenSatelliteProject/SatHelperApp/Models"
	"github.com/OpenSatelliteProject/libsathelper"
	"github.com/racerxdl/go.fifo"
	"sync"
)

// region Global Globals
var running = false

// endregion
// region DSP Globals
var samplesFifo *fifo.Queue
var buffer0 []complex64
var buffer1 []complex64

var decimator SatHelper.FirFilter
var agc SatHelper.AGC
var rrcFilter SatHelper.FirFilter
var costasLoop SatHelper.CostasLoop
var clockRecovery SatHelper.ClockRecovery
var device Frontend.BaseFrontend

// endregion
// region Decoder Globals
var symbolsFifo *fifo.Queue

var viterbiData []byte
var decodedData []byte
var lastFrameEnd []byte

var codedData []byte
var rsCorrectedData []byte
var rsWorkBuffer []byte

var syncWord []byte

var viterbi SatHelper.Viterbi27
var reedSolomon SatHelper.ReedSolomon
var correlator SatHelper.Correlator
var packetFixer SatHelper.PacketFixer

var statistics Models.Statistics
var statisticsMutex = &sync.Mutex{}

var demuxer Demuxer.BaseDemuxer
var statisticsServer *TCPServer

var demodFifoUsage uint8
var decodFifoUsage uint8

// endregion

func GetStats() Models.Statistics {
	statisticsMutex.Lock()
	stat := statistics
	statisticsMutex.Unlock()
	return stat
}

func SetStats(stat Models.Statistics) {
	statisticsMutex.Lock()
	statistics = stat
	statisticsMutex.Unlock()
}
