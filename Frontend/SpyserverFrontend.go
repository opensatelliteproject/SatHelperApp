package Frontend

import (
	"fmt"
	"github.com/opensatelliteproject/SatHelperApp"
	"github.com/opensatelliteproject/SatHelperApp/Logger"
	"github.com/racerxdl/spy2go/spyserver"
	"github.com/racerxdl/spy2go/spytypes"
)

// region Struct Definition
type SpyserverFrontend struct {
	ss *spyserver.Spyserver
	cb SamplesCallback
}

// endregion
// region Constructor
func NewSpyserverFrontend(hostname string, port int) *SpyserverFrontend {
	spyserver.SoftwareID = fmt.Sprintf("SatHelperApp %s.%s", SatHelperApp.GetVersion(), SatHelperApp.GetRevision())
	ss := spyserver.MakeSpyserver(hostname, port)
	afrnt := SpyserverFrontend{
		ss: ss,
	}
	ss.SetCallback(&afrnt)

	return &afrnt
}

func (f *SpyserverFrontend) OnData(dType int, data interface{}) {
	cbData := SampleCallbackData{}

	if dType == spytypes.SamplesComplex64 {
		cbData.SampleType = SampleTypeFloatIQ
		cbData.ComplexArray = data.([]complex64)
		cbData.NumSamples = len(cbData.ComplexArray)
	} else if dType == spytypes.SamplesComplex32 {
		samples := data.([]spytypes.ComplexInt16)
		cbData.SampleType = SampleTypeS16IQ
		cbData.Int16Array = make([]int16, len(samples)*2)
		cbData.NumSamples = len(samples)
		for i := 0; i < len(samples); i++ {
			cbData.Int16Array[i*2] = samples[i].Real
			cbData.Int16Array[i*2+1] = samples[i].Imag
		}
	} else if dType == spytypes.SamplesComplexUInt8 {
		cbData.SampleType = SampleTypeS8IQ
		samples := data.([]spytypes.ComplexUInt8)
		cbData.Int8Array = make([]int8, len(samples)*2)
		for i := 0; i < len(samples); i++ {
			cbData.Int8Array[i*2] = int8(samples[i].Real)
			cbData.Int8Array[i*2+1] = int8(samples[i].Imag)
		}
		cbData.NumSamples = len(samples)
	} else if dType == spytypes.DeviceSync {
		SLog.Info("Got device sync!")
		return
	}

	if f.cb != nil {
		f.cb(cbData)
	}
}

// endregion
// region Getters
func (f *SpyserverFrontend) GetName() string {
	return f.ss.GetName()
}
func (f *SpyserverFrontend) GetShortName() string {
	return f.ss.GetName()
}
func (f *SpyserverFrontend) GetAvailableSampleRates() []uint32 {
	return f.ss.GetAvailableSampleRates()
}
func (f *SpyserverFrontend) GetCenterFrequency() uint32 {
	return f.ss.GetCenterFrequency()
}
func (f *SpyserverFrontend) GetSampleRate() uint32 {
	return f.ss.GetSampleRate()
}

// endregion
// region Setters
func (f *SpyserverFrontend) SetSamplesAvailableCallback(cb SamplesCallback) {
	f.cb = cb
}
func (f *SpyserverFrontend) SetSampleRate(sampleRate uint32) uint32 {
	return f.ss.SetSampleRate(sampleRate)
}
func (f *SpyserverFrontend) SetCenterFrequency(centerFrequency uint32) uint32 {
	return f.ss.SetCenterFrequency(centerFrequency)
}

// endregion
// region Commands
func (f *SpyserverFrontend) Start() {
	f.ss.Start()
}
func (f *SpyserverFrontend) Stop() {
	f.ss.Stop()
}
func (f *SpyserverFrontend) SetAGC(agc bool) {
	SLog.Warn("AGC not supported by SpyServer Frontend")
}
func (f *SpyserverFrontend) SetGain1(gain int) {
	f.ss.SetGain(uint32(gain))
}
func (f *SpyserverFrontend) SetGain2(gain int) {}
func (f *SpyserverFrontend) SetGain3(gain int) {}
func (f *SpyserverFrontend) SetBiasT(biast bool) {
	SLog.Warn("BiasT not supported by SpyServer Frontend")
}
func (f *SpyserverFrontend) Init() bool {
	f.ss.Connect()
	return f.ss.IsConnected
}

func (f *SpyserverFrontend) Destroy() {
	f.ss.Disconnect()
}

func (f *SpyserverFrontend) SetAntenna(string) {}

// endregion
