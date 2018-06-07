package Frontend

import (
	"github.com/OpenSatelliteProject/SatHelperApp/Frontend/LimeDevice"
)

const LimeFrontendBufferSize = 65535

// region Struct Definition
type LimeFrontend struct {
	running bool
	device LimeDevice.LimeDevice
	goCb GoCallback
	goDirCb LimeDevice.LimeDeviceCallback
}

func LimeMakeGoCallbackDirector(callback *GoCallback) LimeDevice.LimeDeviceCallback {
	return LimeDevice.NewDirectorLimeDeviceCallback(callback)
}

// endregion
// region Constructor
func NewLimeFrontend() *LimeFrontend {
	goCb := GoCallback{}
	afrnt := LimeFrontend{
		running: false,
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
	f.running = true
	go func(frontend *LimeFrontend) {
		for frontend.running {
			f.device.GetSamples(LimeFrontendBufferSize)
		}
	}(f)
}
func (f *LimeFrontend) Stop() {
	f.running = false
	f.device.Stop()
}
func (f *LimeFrontend) SetGain1(gain uint8) {
	f.device.SetLNAGain(gain)
}
func (f *LimeFrontend) SetGain2(gain uint8) {
	f.device.SetTIAGain(gain)
}
func (f *LimeFrontend) SetGain3(gain uint8) {
	f.device.SetPGAGain(gain)
}
func (f *LimeFrontend) SetAntenna(value string) {
	f.device.SetAntenna(value)
}
func (f *LimeFrontend) Init() bool {
	return f.device.Init()
}
func (f *LimeFrontend) Destroy() {
	f.device.Destroy()
}

func (f *LimeFrontend) GetAvailableSampleRates() []uint32 { return nil }
func (f *LimeFrontend) SetAGC(agc bool) {}
func (f *LimeFrontend) SetBiasT(biast bool) {}
// endregion
