package Frontend

import (
	"github.com/OpenSatelliteProject/SatHelperApp/Frontend/SpyserverDevice"
)

// region Struct Definition
type SpyserverFrontend struct {
	device  SpyserverDevice.SpyserverDevice
	goCb    GoCallback
	goDirCb SpyserverDevice.SpyserverDeviceCallback
}

func MakeSpyserverGoCallbackDirector(callback *GoCallback) SpyserverDevice.SpyserverDeviceCallback {
	return SpyserverDevice.NewDirectorSpyserverDeviceCallback(callback)
}

// endregion
// region Constructor
func NewSpyserverFrontend(hostname string, port int) *SpyserverFrontend {
	goCb := NewGoCallback()
	dirCb := MakeSpyserverGoCallbackDirector(&goCb)
	afrnt := SpyserverFrontend{
		device: SpyserverDevice.NewSpyserverDevice(dirCb, hostname, port),
		goCb:   goCb,
	}

	return &afrnt
}

// endregion
// region Getters
func (f *SpyserverFrontend) GetName() string {
	return f.device.GetName()
}
func (f *SpyserverFrontend) GetShortName() string {
	return f.device.GetName()
}
func (f *SpyserverFrontend) GetAvailableSampleRates() []uint32 {
	var sampleRates = f.device.GetAvailableSampleRates()
	var sr = make([]uint32, sampleRates.Size())
	for i := 0; i < int(sampleRates.Size()); i++ {
		sr[i] = uint32(sampleRates.Get(i))
	}

	return sr
}
func (f *SpyserverFrontend) GetCenterFrequency() uint32 {
	return uint32(f.device.GetCenterFrequency())
}
func (f *SpyserverFrontend) GetSampleRate() uint32 {
	return uint32(f.device.GetSampleRate())
}

// endregion
// region Setters
func (f *SpyserverFrontend) SetSamplesAvailableCallback(cb SamplesCallback) {
	f.goCb.callback = cb
	f.goDirCb = MakeSpyserverGoCallbackDirector(&f.goCb)
	f.device.SetSamplesAvailableCallback(f.goDirCb)
}
func (f *SpyserverFrontend) SetSampleRate(sampleRate uint32) uint32 {
	return uint32(f.device.SetSampleRate(uint(sampleRate)))
}
func (f *SpyserverFrontend) SetCenterFrequency(centerFrequency uint32) uint32 {
	return uint32(f.device.SetCenterFrequency(uint(centerFrequency)))
}

// endregion
// region Commands
func (f *SpyserverFrontend) Start() {
	f.device.Start()
}
func (f *SpyserverFrontend) Stop() {
	f.device.Stop()
}
func (f *SpyserverFrontend) SetAGC(agc bool) {
	f.device.SetAGC(agc)
}
func (f *SpyserverFrontend) SetGain1(gain uint8) {
	f.device.SetLNAGain(gain)
}
func (f *SpyserverFrontend) SetGain2(gain uint8) {
	f.device.SetVGAGain(gain)
}
func (f *SpyserverFrontend) SetGain3(gain uint8) {
	f.device.SetMixerGain(gain)
}
func (f *SpyserverFrontend) SetBiasT(biast bool) {
	val := uint8(0)
	if biast {
		val = 1
	}
	f.device.SetBiasT(val)
}
func (f *SpyserverFrontend) Init() bool {
	return f.device.Init()
}

func (f *SpyserverFrontend) Destroy() {
	f.device.Destroy()
}

func (f *SpyserverFrontend) SetAntenna(string) {}

// endregion
