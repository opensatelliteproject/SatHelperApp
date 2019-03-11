package DSP

import (
	. "github.com/logrusorgru/aurora"
	"github.com/opensatelliteproject/SatHelperApp/Frontend"
	"github.com/opensatelliteproject/SatHelperApp/Logger"
	"github.com/opensatelliteproject/libsathelper"
	"github.com/racerxdl/go.fifo"
	"github.com/racerxdl/segdsp/dsp"
	"sync"
	"time"
)

var dspLock = sync.Mutex{}

func InitAll() {
	InitDSP()
	InitDecoder()
}

const defaultTransitionWidth = 200e3

func InitDSP() {
	Device.SetSamplesAvailableCallback(newSamplesCallback)
	lastConstellationSend = time.Now()
	samplesFifo = fifo.NewQueue()
	constellationFifo = fifo.NewQueue()
	constellationBuffer = make([]byte, 1024)
	circuitSampleRate := float32(Device.GetSampleRate()) / float32(CurrentConfig.Base.Decimation)
	sps := circuitSampleRate / float32(CurrentConfig.Base.SymbolRate)

	SLog.Debug("Samples per Symbol: %f", Bold(Green(sps)))
	SLog.Debug("Circuit Sample Rate: %f", Bold(Green(circuitSampleRate)))
	SLog.Debug("Low Pass Decimator Cut Frequency: %f", Bold(Green(circuitSampleRate/2)))

	//rrcTaps := SatHelper.FiltersRRC(1, float64(circuitSampleRate), float64(CurrentConfig.Base.SymbolRate), float64(CurrentConfig.Base.RRCAlpha), RrcTaps)
	//decimatorTaps := SatHelper.FiltersLowPass(1, float64(Device.GetSampleRate()), float64(circuitSampleRate/2) - defaultTransitionWidth / 2, defaultTransitionWidth, SatHelper.FFTWindowsHAMMING, 6.76)

	//decimator = SatHelper.NewFirFilter(uint(CurrentConfig.Base.Decimation), decimatorTaps)
	agc = SatHelper.NewAGC(AgcRate, AgcReference, AgcGain, AgcMaxGain)
	//costasLoop = SatHelper.NewCostasLoop(PllAlpha, LoopOrder)
	clockRecovery = SatHelper.NewClockRecovery(sps, ClockGainOmega, ClockMu, ClockAlpha, ClockOmegaLimit)
	//rrcFilter = SatHelper.NewFirFilter(1, rrcTaps)

	SLog.Debug("Center Frequency: %d MHz", Bold(Green(Device.GetCenterFrequency())))
	SLog.Debug("Automatic Gain Control: %t", Bold(Green(CurrentConfig.Base.AGCEnabled)))

	// region SegDSP
	agcNew = dsp.MakeSimpleAGC(AgcRate, AgcReference, AgcGain, AgcMaxGain)
	rrcFilterNew = dsp.MakeFirFilter(dsp.MakeRRC(1, float64(circuitSampleRate), float64(CurrentConfig.Base.SymbolRate), float64(CurrentConfig.Base.RRCAlpha), RrcTaps))

	newDecimatorTaps := dsp.MakeLowPass(1, float64(Device.GetSampleRate()), float64(circuitSampleRate/2)-defaultTransitionWidth/2, defaultTransitionWidth)
	decimatorNew = dsp.MakeDecimationFirFilter(int(CurrentConfig.Base.Decimation), newDecimatorTaps)

	costasLoopNew = dsp.MakeCostasLoop2(PllAlpha)
	// endregion
}

func newSamplesCallback(d Frontend.SampleCallbackData) {
	switch d.SampleType {
	case Frontend.SampleTypeFloatIQ:
		AddToFifoC64(samplesFifo, d.ComplexArray, d.NumSamples)
	case Frontend.SampleTypeS16IQ:
		AddToFifoS16toC64(samplesFifo, d.Int16Array, d.NumSamples)
	case Frontend.SampleTypeS8IQ:
		AddToFifoS8toC64(samplesFifo, d.Int8Array, d.NumSamples)
	}
}

func sendConstellation() {
	if ConstellationServer != nil && constellationFifo.UnsafeLen() >= 1024 && time.Since(lastConstellationSend) > (time.Millisecond*10) {
		for i := 0; i < 1024; i++ {
			constellationBuffer[i] = constellationFifo.UnsafeNext().(uint8)
		}
		ConstellationServer.SendData(constellationBuffer)
		lastConstellationSend = time.Now()
	}
}

func GetCostasFrequency() float32 {
	dspLock.Lock()
	defer dspLock.Unlock()
	return costasLoopNew.GetFrequencyHz()
}

func processSamples() {
	dspLock.Lock()
	length := samplesFifo.Len()
	demodFifoUsage = uint8(100 * float32(length) / float32(FifoSize))
	dspLock.Unlock()

	if length <= 64*1024 {
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

	ban := buffer0
	bbn := buffer1

	ba := &ban[0]
	bb := &bbn[0]

	if CurrentConfig.Base.Decimation > 1 {
		//length /= int(CurrentConfig.Base.Decimation)
		length = decimatorNew.WorkBuffer(ban, bbn)
		//decimator.Work(ba, bb, length)
		swapAndTrimSlices(&ban, &bbn, length)
		ba = &ban[0]
		bb = &bbn[0]
	}

	//length = agcNew.WorkBuffer(ban, bbn)
	agc.Work(ba, bb, length)
	swapAndTrimSlices(&ban, &bbn, length)
	//ba = &ban[0]
	//bb = &bbn[0]

	length = rrcFilterNew.WorkBuffer(ban, bbn)
	//rrcFilter.Work(ba, bb, length)
	swapAndTrimSlices(&ban, &bbn, length)
	//ba = &ban[0]
	//bb = &bbn[0]

	//costasLoop.Work(ba, bb, length)
	length = costasLoopNew.WorkBuffer(ban, bbn)
	swapAndTrimSlices(&ban, &bbn, length)
	ba = &ban[0]
	bb = &bbn[0]

	symbols := clockRecovery.Work(ba, bb, length)
	swapAndTrimSlices(&ban, &bbn, length)
	ba = &ban[0]
	//bb = &bbn[0]

	var ob *[]complex64

	if ba == &buffer0[0] {
		ob = &buffer0
	} else {
		ob = &buffer1
	}

	symbolsFifo.UnsafeLock()
	constellationFifo.UnsafeLock()
	defer symbolsFifo.UnsafeUnlock()
	defer constellationFifo.UnsafeUnlock()
	for i := 0; i < symbols; i++ {
		z := (*ob)[i]
		v := real(z) * 127
		if v > 127 {
			v = 127
		} else if v < -128 {
			v = -128
		}

		symbolsFifo.UnsafeAdd(byte(v))

		if CurrentConfig.Base.SendConstellation {
			v2 := real(z) * 127
			if v2 > 127 {
				v2 = 127
			} else if v2 < -128 {
				v2 = -128
			}

			constellationFifo.UnsafeAdd(uint8(v2))
			constellationFifo.UnsafeAdd(uint8(v))
			if constellationFifo.UnsafeLen() > 1024 {
				_ = constellationFifo.UnsafeNext()
				_ = constellationFifo.UnsafeNext()
			}
		}
	}
	if CurrentConfig.Base.SendConstellation {
		sendConstellation()
	}
}

func GetDemodFIFOUsage() uint8 {
	dspLock.Lock()
	defer dspLock.Unlock()
	return demodFifoUsage
}

func GetDecoderFIFOUsage() uint8 {
	dspLock.Lock()
	defer dspLock.Unlock()
	return decodFifoUsage
}

func SetRunning(r bool) {
	dspLock.Lock()
	running = r
	dspLock.Unlock()
}

func IsRunning() bool {
	dspLock.Lock()
	defer dspLock.Unlock()
	return running
}

func symbolProcessLoop() {
	SLog.Info("Symbol Process Routine started")

	for IsRunning() {
		processSamples()
		time.Sleep(time.Microsecond)
	}

	SLog.Error("Symbol Process Routine stopped")
}

func StartDSPLoops() {
	startTime = uint32(time.Now().Unix() & 0xFFFFFFFF)
	go symbolProcessLoop()
	go decoderLoop()
}
