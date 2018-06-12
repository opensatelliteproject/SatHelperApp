package main

import (
	"github.com/OpenSatelliteProject/libsathelper"
	"github.com/OpenSatelliteProject/SatHelperApp/Frontend"
	. "github.com/logrusorgru/aurora"
	"time"
	"github.com/racerxdl/go.fifo"
	"github.com/OpenSatelliteProject/SatHelperApp/Logger"
)

func initDSP() {
	samplesFifo = fifo.NewQueue()
	circuitSampleRate := float32(device.GetSampleRate()) / float32(CurrentConfig.Base.Decimation)
	sps := circuitSampleRate / float32(CurrentConfig.Base.SymbolRate)

	SLog.Debug("Samples per Symbol: %f", Bold(Green(sps)))
	SLog.Debug("Circuit Sample Rate: %f", Bold(Green(circuitSampleRate)))
	SLog.Debug("Low Pass Decimator Cut Frequency: %f", Bold(Green(circuitSampleRate / 2)))

	rrcTaps := SatHelper.FiltersRRC(1, float64(circuitSampleRate), float64(CurrentConfig.Base.SymbolRate), float64(CurrentConfig.Base.RRCAlpha), RrcTaps)
	decimatorTaps := SatHelper.FiltersLowPass(1, float64(device.GetSampleRate()), float64(circuitSampleRate / 2), 100e3, SatHelper.FFTWindowsHAMMING, 6.76)

	decimator = SatHelper.NewFirFilter(uint(CurrentConfig.Base.Decimation), decimatorTaps)
	agc = SatHelper.NewAGC(AgcRate, AgcReference, AgcGain, AgcMaxGain)
	costasLoop = SatHelper.NewCostasLoop(PllAlpha, LoopOrder)
	clockRecovery = SatHelper.NewClockRecovery(sps, ClockGainOmega, ClockMu, ClockAlpha, ClockOmegaLimit)
	rrcFilter = SatHelper.NewFirFilter(1, rrcTaps)


	SLog.Debug("Center Frequency: %d MHz", Bold(Green(device.GetCenterFrequency())))
	SLog.Debug("Automatic Gain Control: %t", Bold(Green(CurrentConfig.Base.AGCEnabled)))
}

func newSamplesCallback(d Frontend.SampleCallbackData) {
	switch d.SampleType {
	case Frontend.SampleTypeFloatIQ: AddToFifoC64(samplesFifo, d.ComplexArray, d.NumSamples); break
	case Frontend.SampleTypeS16IQ: AddToFifoS16toC64(samplesFifo, d.Int16Array, d.NumSamples); break
	case Frontend.SampleTypeS8IQ: AddToFifoS8toC64(samplesFifo, d.Int8Array, d.NumSamples); break
	}
}

func processSamples() {
	length := samplesFifo.Len()
	demodFifoUsage = uint8(100 * float32(length) / float32(FifoSize))
	if length <= 64 * 1024 {
		return
	}

	samplesFifo.UnsafeLock()
	// region Unsafe Locked Section
	checkAndResizeBuffers(length)

	for i := 0; i < length; i++ {
		buffer0[i] = samplesFifo.UnsafeNext().(complex64)
	}
	// endregion
	samplesFifo.UnsafeUnlock()

	ba := &buffer0[0]
	bb := &buffer1[0]

	if CurrentConfig.Base.Decimation > 1 {
		length /= int(CurrentConfig.Base.Decimation)
		decimator.Work(ba, bb, length)
		swapBuffers(&ba, &bb)
	}

	agc.Work(ba, bb, length)
	swapBuffers(&ba, &bb)

	rrcFilter.Work(ba, bb, length)
	swapBuffers(&ba, &bb)

	costasLoop.Work(ba, bb, length)
	swapBuffers(&ba, &bb)

	symbols := clockRecovery.Work(ba, bb, length)
	swapBuffers(&ba, &bb)

	var ob *[]complex64

	if ba == &buffer0[0] {
		ob = &buffer0
	} else {
		ob = &buffer1
	}

	symbolsFifo.UnsafeLock()
	defer symbolsFifo.UnsafeUnlock()
	for i := 0; i < symbols; i++ {
		z := (*ob)[i]
		v := imag(z) * 127
		if v > 127 {
			v = 127
		} else if v < -128 {
			v = -128
		}

		symbolsFifo.UnsafeAdd(byte(v))
	}
}

func symbolProcessLoop() {
	SLog.Info("Symbol Process Routine started")
	for running {
		processSamples()
		time.Sleep(time.Microsecond)
	}
	SLog.Error("Symbol Process Routine stopped")
}