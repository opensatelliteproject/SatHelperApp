package Frontend

import (
	"github.com/OpenSatelliteProject/SatHelperApp/Frontend/SpyserverDevice"
	"unsafe"
)


// region Struct Definition
type SpyserverFrontend struct {
	device SpyserverDevice.SpyserverDevice
	goCb SpyserverGoCallback
	goDirCb SpyserverDevice.SpyserverDeviceCallback
}

type SpyserverGoCallback struct {
	callback SamplesCallback
}

func (p *SpyserverGoCallback) CbFloatIQ(data uintptr, length int) {
	const arrayLen = 1 << 30
	arr := (*[arrayLen]complex64)(unsafe.Pointer(data))[:length:length]
	if p.callback != nil {
		p.callback(SampleCallbackData{
			ComplexArray: arr,
			NumSamples: length,
			SampleType: FrontendSampletypeFloatiq,
		})
	}
}

func (p *SpyserverGoCallback) CbS16IQ(data uintptr, length int) {
	// Length times two, because each sample contains an I and a Q in S16
	const arrayLen = 1 << 30
	var pairLength = length * 2
	arr := (*[arrayLen]int16)(unsafe.Pointer(data))[:pairLength:pairLength]
	if p.callback != nil {
		p.callback(SampleCallbackData{
			Int16Array: arr,
			NumSamples: length,
			SampleType: FrontendSampletypeS16iq,
		})
	}
}

func (p *SpyserverGoCallback) CbS8IQ(data uintptr, length int) {
	// Length times two, because each sample contains an I and a Q in S8
	const arrayLen = 1 << 30
	var pairLength = length * 2
	arr := (*[arrayLen]int8)(unsafe.Pointer(data))[:pairLength:pairLength]
	if p.callback != nil {
		p.callback(SampleCallbackData{
			Int8Array: arr,
			NumSamples: length,
			SampleType: FrontendSampletypeS8iq,
		})
	}
}

func MakeSpyserverGoCallbackDirector(callback *SpyserverGoCallback) SpyserverDevice.SpyserverDeviceCallback {
	return SpyserverDevice.NewDirectorSpyserverDeviceCallback(callback)
}

// endregion
// region Constructor
func NewSpyserverFrontend(hostname string, port int) *SpyserverFrontend {
	goCb := SpyserverGoCallback{}
	afrnt := SpyserverFrontend{
		device: SpyserverDevice.NewSpyserverDevice(hostname, port),
		goCb: goCb,
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
func (f *SpyserverFrontend)  GetSampleRate() uint32 {
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
func (f *SpyserverFrontend) SetLNAGain(gain uint8) {
	f.device.SetLNAGain(gain)
}
func (f *SpyserverFrontend) SetVGAGain(gain uint8) {
	f.device.SetVGAGain(gain)
}
func (f *SpyserverFrontend) SetMixerGain(gain uint8) {
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
// endregion
