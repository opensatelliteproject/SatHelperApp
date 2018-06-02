package main

import (
	"github.com/foize/go.fifo"
	"github.com/OpenSatelliteProject/libsathelper"
	"net"
	"github.com/OpenSatelliteProject/SatHelperApp/Frontend"
)

var samplesFifo *fifo.Queue
var symbolsFifo *fifo.Queue
var buffer0 []complex64
var buffer1 []complex64
var running = false

var decimator SatHelper.FirFilter
var agc SatHelper.AGC
var rrcFilter SatHelper.FirFilter
var costasLoop SatHelper.CostasLoop
var clockRecovery SatHelper.ClockRecovery
var device Frontend.BaseFrontend

var conn net.Conn


func checkAndResizeBuffers(length int) {
	if len(buffer0) < length {
		buffer0 = make([]complex64, length)
	}
	if len(buffer1) < length {
		buffer1 = make([]complex64, length)
	}
}
