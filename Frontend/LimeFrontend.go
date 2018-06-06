package Frontend

import (
	"github.com/OpenSatelliteProject/SatHelperApp/Frontend/LimeDevice"
)

// region Struct Definition
type LimeFrontend struct {
	device LimeDevice.LimeDevice
	goCb GoCallback
	goDirCb LimeDevice.LimeCallback
}

func LimeMakeGoCallbackDirector(callback *GoCallback) LimeDevice.LimeCallback {
	return LimeDevice.NewDirectorLimeCallback(callback)
}

// endregion
// region Constructor
func NewLimeFrontend() *LimeFrontend {
	goCb := GoCallback{}
	afrnt := LimeFrontend{
		device: LimeDevice.NewLimeDevice(),
		goCb: goCb,
	}

	return &afrnt
}
// endregion
// region Getters
func (f *LimeFrontend) GetName() string {
	return f.device.GetName()
}
func (f *LimeFrontend) GetShortName() string {
	return f.device.GetName()
}
func (f *LimeFrontend) GetCenterFrequency() uint32 {
	return uint32(f.device.GetCenterFrequency())
}
func (f *LimeFrontend)  GetSampleRate() uint32 {
	return uint32(f.device.GetSampleRate())
}
func (f *LimeFrontend) GetAvailableSampleRates() []uint32 {
	var sampleRates = f.device.GetAvailableSampleRates()
	var sr = make([]uint32, sampleRates.Size())
	for i := 0; i < int(sampleRates.Size()); i++ {
		sr[i] = uint32(sampleRates.Get(i))
	}

	return sr
}
// endregion
// region Setters
func (f *LimeFrontend) SetSamplesAvailableCallback(cb SamplesCallback) {
	f.goCb.callback = cb
	f.goDirCb = LimeMakeGoCallbackDirector(&f.goCb)
	f.device.SetSamplesAvailableCallback(f.goDirCb)
}
func (f *LimeFrontend) SetSampleRate(sampleRate uint32) uint32 {
	return uint32(f.device.SetSampleRate(uint(sampleRate)))
}
func (f *LimeFrontend) SetCenterFrequency(centerFrequency uint32) uint32 {
	return uint32(f.device.SetCenterFrequency(uint(centerFrequency)))
}
// endregion
// region Commands
func (f *LimeFrontend) Start() {
	f.device.Start()
}
func (f *LimeFrontend) Stop() {
	f.device.Stop()
}
func (f *LimeFrontend) SetAGC(agc bool) {
	f.device.SetAGC(agc)
}
func (f *LimeFrontend) SetLNAGain(gain uint8) {
	f.device.SetLNAGain(gain)
}
func (f *LimeFrontend) SetVGAGain(gain uint8) {}
func (f *LimeFrontend) SetMixerGain(gain uint8) {}
func (f *LimeFrontend) SetBiasT(biast bool) {}
// endregion
