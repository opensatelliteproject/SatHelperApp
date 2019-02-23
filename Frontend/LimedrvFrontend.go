package Frontend

import (
	"github.com/OpenSatelliteProject/SatHelperApp/Logger"
	"github.com/myriadrf/limedrv"
	"math"
)

const limeDrvSampleRateDelta = 100e3
const limeDrvLPFBW = 3e6
const limeDrvDefaultOversample = 4

// region Struct Definition
type LimeDrvFrontend struct {
	running         bool
	device          *limedrv.LMSDevice
	activeChannel   int
	currentAntenna  string
	oversample      int
	sampleRate      float64
	gain            uint
	centerFrequency float64
	cb              SamplesCallback
}

// endregion
// region Constructor
func NewLimeDrvFrontend() *LimeDrvFrontend {
	afrnt := LimeDrvFrontend{
		running:         false,
		activeChannel:   -1,
		device:          nil,
		oversample:      1,
		gain:            0,
		centerFrequency: 0,
		cb:              nil,
		sampleRate:      2.5e6,
	}

	return &afrnt
}

// endregion
// region Getters
func (f *LimeDrvFrontend) GetName() string {
	if f.device != nil {
		return f.device.DeviceInfo.DeviceName
	}

	return "LimeSDR (Not Connected)"
}
func (f *LimeDrvFrontend) GetShortName() string {
	if f.device != nil {
		return f.device.DeviceInfo.DeviceName
	}

	return "LimeSDR (Not Connected)"
}
func (f *LimeDrvFrontend) GetCenterFrequency() uint32 {
	if f.device != nil {
		return uint32(f.device.GetCenterFrequency(f.activeChannel, true))
	}

	SLog.Error("LimeDRV Frontend not initialized!")

	return 0
}
func (f *LimeDrvFrontend) GetSampleRate() uint32 {
	if f.device != nil {
		host, _ := f.device.GetSampleRate()
		return uint32(math.Round(host))
	}
	SLog.Error("LimeDRV Frontend not initialized!")

	return 0
}

// endregion
// region Setters
func (f *LimeDrvFrontend) SetSamplesAvailableCallback(cb SamplesCallback) {
	f.cb = cb
}
func (f *LimeDrvFrontend) SetSampleRate(sampleRate uint32) uint32 {
	f.sampleRate = float64(sampleRate)
	if f.device != nil {
		SLog.Debug("Setting Sample Rate to %f with %d oversampling.", f.sampleRate, f.oversample)
		f.device.SetSampleRate(f.sampleRate, f.oversample)
		return f.GetSampleRate()
	}

	SLog.Error("LimeDRV Frontend not initialized!")
	return 0
}

func (f *LimeDrvFrontend) SetCenterFrequency(centerFrequency uint32) uint32 {
	f.centerFrequency = float64(centerFrequency)
	if f.device != nil {
		f.device.SetCenterFrequency(f.activeChannel, true, float64(centerFrequency))
		return f.GetCenterFrequency()
	}

	return uint32(f.centerFrequency)
}

// endregion
// region Commands
func (f *LimeDrvFrontend) Start() {
	if f.device != nil {
		f.device.SetSampleRate(f.sampleRate, f.oversample)
		SLog.Debug("Enabling channel %d (%s) - %d %f", f.activeChannel, f.currentAntenna, f.gain, f.centerFrequency)
		f.device.RXChannels[f.activeChannel].
			Enable().
			SetAntennaByName(f.currentAntenna).
			SetGainDB(f.gain).
			SetLPF(limeDrvLPFBW).
			EnableLPF().
			SetCenterFrequency(f.centerFrequency)
		f.device.Start()
		f.running = true
		return
	}

	SLog.Error("LimeDRV Frontend not initialized!")
}
func (f *LimeDrvFrontend) Stop() {
	if f.device != nil {
		f.running = false
		f.device.Stop()
		SLog.Debug("Disabling channel %d", f.activeChannel)
		f.device.DisableChannel(f.activeChannel, true)
		return
	}

	SLog.Error("LimeDRV Frontend not initialized!")
}
func (f *LimeDrvFrontend) SetGain1(gain uint8) {
	f.gain = uint(gain)
	if f.device != nil {
		f.device.SetGainDB(f.activeChannel, true, f.gain)
		return
	}
	SLog.Error("LimeDRV Frontend not initialized!")
}
func (f *LimeDrvFrontend) SetGain2(gain uint8) {}
func (f *LimeDrvFrontend) SetGain3(gain uint8) {}
func (f *LimeDrvFrontend) SetChannel(value int) {
	if f.device != nil {
		if f.running {
			f.device.Stop()
			f.device.DisableChannel(f.activeChannel, true)
		}

		f.activeChannel = value
		f.device.EnableChannel(f.activeChannel, true)

		if f.running {
			f.device.Start()
		}
		return
	}

	SLog.Error("LimeDRV Frontend not initialized!")
}
func (f *LimeDrvFrontend) SetAntenna(value string) {
	f.currentAntenna = value
	if f.device != nil {
		f.device.SetAntennaByName(value, f.activeChannel, true)
		return
	}
	SLog.Error("LimeDRV Frontend not initialized!")
}

func (f *LimeDrvFrontend) samplesCallback(samples []complex64, _ int, _ uint64) {
	//SLog.Info("Received %d samples!", len(samples))
	if f.cb != nil {
		f.cb(SampleCallbackData{
			ComplexArray: samples,
			SampleType:   SampleTypeFloatIQ,
			NumSamples:   len(samples),
		})
	}
}

func (f *LimeDrvFrontend) Init() (ret bool) {
	defer func() {
		if r := recover(); r != nil {
			ret = false
			SLog.Error("Error initializing device: %s", r)
		}
	}()

	devices := limedrv.GetDevices()
	if len(devices) == 0 {
		SLog.Error("No LimeSDR devices available")
		return false
	}

	f.device = limedrv.Open(devices[0])

	SLog.Debug("Got device %s", f.device.DeviceInfo.DeviceName)
	f.SetChannel(limedrv.ChannelA)
	f.SetAntenna(limedrv.LNAH)
	f.SetOversample(limeDrvDefaultOversample)

	f.device.SetCallback(f.samplesCallback)

	return true
}

func (f *LimeDrvFrontend) Destroy() {
	if f.device != nil {
		limedrv.Close(f.device)
	}
}

func (f *LimeDrvFrontend) GetAvailableSampleRates() []uint32 {
	if f.device != nil {
		max := f.device.MaximumSampleRate
		min := f.device.MinimumSampleRate

		rates := make([]uint32, int((max-min)/limeDrvSampleRateDelta))
		for i := 0; i < len(rates); i++ {
			rates[i] = uint32(min + float64(i)*limeDrvSampleRateDelta)
		}
		return rates
	}
	return nil
}
func (f *LimeDrvFrontend) SetAGC(agc bool) {
	SLog.Warn("AGC not supported by LimeDrv Frontend")
}
func (f *LimeDrvFrontend) SetBiasT(biast bool) {
	SLog.Warn("BiasT not supported by LimeDrv Frontend")
}

func (f *LimeDrvFrontend) SetOversample(oversample int) {
	f.oversample = oversample

	if f.device != nil {
		if f.running {
			f.device.Stop()
		}

		f.SetSampleRate(f.GetSampleRate())

		if f.running {
			f.device.Start()
		}
	}
}

func (f *LimeDrvFrontend) GetTemperature() float64 {
	if f.device != nil {
		return f.device.GetTemperature()
	}
	SLog.Error("LimeDRV Frontend not initialized!")

	return 0
}

// endregion
