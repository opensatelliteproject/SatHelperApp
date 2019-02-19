package Frontend

import (
	"github.com/OpenSatelliteProject/SatHelperApp/Frontend/RTLSDRDevice"
)

// region Struct Definition
type RTLSDRFrontend struct {
	device  RTLSDRDevice.RtlFrontend
	goCb    GoCallback
	goDirCb RTLSDRDevice.RTLSDRDeviceCallback
}

func MakeRTLSDRGoCallbackDirector(callback *GoCallback) RTLSDRDevice.RTLSDRDeviceCallback {
	return RTLSDRDevice.NewDirectorRTLSDRDeviceCallback(callback)
}

// endregion
// region Constructor
func NewRTLSDRFrontend() *RTLSDRFrontend {
	goCb := NewGoCallback()
	dirCb := MakeRTLSDRGoCallbackDirector(&goCb)
	afrnt := RTLSDRFrontend{
		device: RTLSDRDevice.NewRtlFrontend(dirCb),
		goCb:   goCb,
	}

	return &afrnt
}

// endregion
// region Getters
func (f *RTLSDRFrontend) GetName() string {
	return f.device.GetName()
}

func (f *RTLSDRFrontend) GetShortName() string {
	return f.device.GetName()
}

func (f *RTLSDRFrontend) GetAvailableSampleRates() []uint32 {
	var sampleRates = f.device.GetAvailableSampleRates()
	var sr = make([]uint32, sampleRates.Size())
	for i := 0; i < int(sampleRates.Size()); i++ {
		sr[i] = uint32(sampleRates.Get(i))
	}

	return sr
}

func (f *RTLSDRFrontend) GetCenterFrequency() uint32 {
	return uint32(f.device.GetCenterFrequency())
}

func (f *RTLSDRFrontend) GetSampleRate() uint32 {
	return uint32(f.device.GetSampleRate())
}

// endregion
// region Setters
func (f *RTLSDRFrontend) SetSamplesAvailableCallback(cb SamplesCallback) {
	f.goCb.callback = cb
	f.goDirCb = MakeRTLSDRGoCallbackDirector(&f.goCb)
	f.device.SetSamplesAvailableCallback(f.goDirCb)
}

func (f *RTLSDRFrontend) SetSampleRate(sampleRate uint32) uint32 {
	return uint32(f.device.SetSampleRate(uint(sampleRate)))
}

func (f *RTLSDRFrontend) SetCenterFrequency(centerFrequency uint32) uint32 {
	return uint32(f.device.SetCenterFrequency(uint(centerFrequency)))
}

// endregion
// region Commands
func (f *RTLSDRFrontend) Start() {
	f.device.Start()
}

func (f *RTLSDRFrontend) Stop() {
	f.device.Stop()
}

func (f *RTLSDRFrontend) SetAGC(agc bool) {
	f.device.SetAGC(agc)
}

func (f *RTLSDRFrontend) SetGain1(gain uint8) {
	f.device.SetLNAGain(gain)
}

func (f *RTLSDRFrontend) SetGain2(gain uint8) {
	f.device.SetVGAGain(gain)
}

func (f *RTLSDRFrontend) SetGain3(gain uint8) {
	f.device.SetMixerGain(gain)
}

func (f *RTLSDRFrontend) SetBiasT(biast bool) {
	val := uint8(0)
	if biast {
		val = 1
	}
	f.device.SetBiasT(val)
}

func (f *RTLSDRFrontend) Init() bool {
	return f.device.Init()
}

func (f *RTLSDRFrontend) Destroy() {}

func (f *RTLSDRFrontend) SetAntenna(string) {}

// endregion
