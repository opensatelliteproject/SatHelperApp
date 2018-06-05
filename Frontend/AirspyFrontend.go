package Frontend

import (
	"github.com/OpenSatelliteProject/SatHelperApp/Frontend/AirspyDevice"
	"unsafe"
)


// region Struct Definition
type AirspyFrontend struct {
	device AirspyDevice.AirspyDevice
	goCb GoCallback
	goDirCb AirspyDevice.GoDeviceCallback
}

type GoCallback struct {
	callback SamplesCallback
}

func (p *GoCallback) CbFloatIQ(data uintptr, length int) {
	const arrayLen = 1 << 30
	arr := (*[arrayLen]complex64)(unsafe.Pointer(data))[:length:length]
	if p.callback != nil {
		p.callback(SampleCallbackData{
			ComplexArray: arr,
			NumSamples:   length,
			SampleType:   FrontendSampletypeFloatiq,
		})
	}
}

func MakeGoCallbackDirector(callback *GoCallback) AirspyDevice.GoDeviceCallback {
	return AirspyDevice.NewDirectorGoDeviceCallback(callback)
}

// endregion
// region Constructor
func NewAirspyFrontend() *AirspyFrontend {
	goCb := GoCallback{}
	afrnt := AirspyFrontend{
		device: AirspyDevice.NewAirspyDevice(),
		goCb: goCb,
	}

	return &afrnt
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
func (f *AirspyFrontend)  GetSampleRate() uint32 {
	return uint32(f.device.GetSampleRate())
}
// endregion
// region Setters
func (f *AirspyFrontend) SetSamplesAvailableCallback(cb SamplesCallback) {
	f.goCb.callback = cb
	f.goDirCb = MakeGoCallbackDirector(&f.goCb)
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
func (f *AirspyFrontend) SetLNAGain(gain uint8) {
	f.device.SetLNAGain(gain)
}
func (f *AirspyFrontend) SetVGAGain(gain uint8) {
	f.device.SetVGAGain(gain)
}
func (f *AirspyFrontend) SetMixerGain(gain uint8) {
	f.device.SetMixerGain(gain)
}
func (f *AirspyFrontend) SetBiasT(biast bool) {
	val := uint8(0)
	if biast {
		val = 1
	}
	f.device.SetBiasT(val)
}
// endregion
