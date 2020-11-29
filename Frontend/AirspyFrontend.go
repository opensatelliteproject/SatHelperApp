package Frontend

import (
	"github.com/opensatelliteproject/SatHelperApp/Frontend/AirspyDevice"
)

// region Struct Definition
type AirspyFrontend struct {
	device  AirspyDevice.AirspyDevice
	goCb    GoCallback
	goDirCb AirspyDevice.AirspyDeviceCallback
}

func MakeAirspyGoCallbackDirector(callback *GoCallback) AirspyDevice.AirspyDeviceCallback {
	return AirspyDevice.NewDirectorAirspyDeviceCallback(callback)
}

// endregion
// region Constructor
func NewAirspyFrontend() *AirspyFrontend {
	goCb := NewGoCallback()
	dirCb := MakeAirspyGoCallbackDirector(&goCb)
	afrnt := AirspyFrontend{
		device: AirspyDevice.NewAirspyDevice(dirCb),
		goCb:   goCb,
	}

	return &afrnt
}
func AirspyInitialize() {
	AirspyDevice.AirspyDeviceInitialize()
}
func AirspyDeinitialize() {
	AirspyDevice.AirspyDeviceDeInitialize()
}

// endregion
// region Getters
func (f *AirspyFrontend) GetName() string {
	return f.device.GetName()
}
func (f *AirspyFrontend) GetShortName() string {
	return f.device.GetName()
}
func (f *AirspyFrontend) GetAvailableSampleRates() []uint32 {
	var sampleRates = f.device.GetAvailableSampleRates()
	var sr = make([]uint32, sampleRates.Size())
	for i := 0; i < int(sampleRates.Size()); i++ {
		sr[i] = uint32(sampleRates.Get(i))
	}

	return sr
}
func (f *AirspyFrontend) GetCenterFrequency() uint32 {
	return uint32(f.device.GetCenterFrequency())
}
func (f *AirspyFrontend) GetSampleRate() uint32 {
	return uint32(f.device.GetSampleRate())
}

// endregion
// region Setters
func (f *AirspyFrontend) SetSamplesAvailableCallback(cb SamplesCallback) {
	f.goCb.callback = cb
	f.goDirCb = MakeAirspyGoCallbackDirector(&f.goCb)
	f.device.SetSamplesAvailableCallback(f.goDirCb)
}
func (f *AirspyFrontend) SetSampleRate(sampleRate uint32) uint32 {
	return uint32(f.device.SetSampleRate(uint(sampleRate)))
}
func (f *AirspyFrontend) SetCenterFrequency(centerFrequency uint32) uint32 {
	return uint32(f.device.SetCenterFrequency(uint(centerFrequency)))
}

// endregion
// region Commands
func (f *AirspyFrontend) Start() {
	f.device.Start()
}
func (f *AirspyFrontend) Stop() {
	f.device.Stop()
}
func (f *AirspyFrontend) SetAGC(agc bool) {
	f.device.SetAGC(agc)
}
func (f *AirspyFrontend) SetGain1(gain int) {
	f.device.SetLNAGain(uint8(gain))
}
func (f *AirspyFrontend) SetGain2(gain int) {
	f.device.SetVGAGain(uint8(gain))
}
func (f *AirspyFrontend) SetGain3(gain int) {
	f.device.SetMixerGain(uint8(gain))
}
func (f *AirspyFrontend) SetBiasT(biast bool) {
	val := uint8(0)
	if biast {
		val = 1
	}
	f.device.SetBiasT(val)
}
func (f *AirspyFrontend) Init() bool {
	return f.device.Init()
}
func (f *AirspyFrontend) Destroy() {
	f.device.Destroy()
}

func (f *AirspyFrontend) SetAntenna(string) {}

// endregion
