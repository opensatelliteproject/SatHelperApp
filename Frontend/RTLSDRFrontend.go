package Frontend

import (
	"fmt"
	"github.com/opensatelliteproject/SatHelperApp/Frontend/RTLSDRDevice"
)

type RTLSDRTuner uint8

const (
	RTLSDR_TUNER_UNKNOWN RTLSDRTuner = 0
	RTLSDR_TUNER_E4000   RTLSDRTuner = 1
	RTLSDR_TUNER_FC0012  RTLSDRTuner = 2
	RTLSDR_TUNER_FC0013  RTLSDRTuner = 3
	RTLSDR_TUNER_FC2580  RTLSDRTuner = 4
	RTLSDR_TUNER_R820T2  RTLSDRTuner = 5
	RTLSDR_TUNER_R828D   RTLSDRTuner = 6
)

var rtlsdrTunerName = map[RTLSDRTuner]string{
	RTLSDR_TUNER_UNKNOWN: "Unknown",
	RTLSDR_TUNER_E4000:   "E4000",
	RTLSDR_TUNER_FC0012:  "FC0012",
	RTLSDR_TUNER_FC0013:  "FC0013",
	RTLSDR_TUNER_FC2580:  "FC2580",
	RTLSDR_TUNER_R820T2:  "R820T/2",
	RTLSDR_TUNER_R828D:   "R828D",
}

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
	return fmt.Sprintf("RTLSDR with %s Tuner", f.GetTunerName())
}

func (f *RTLSDRFrontend) GetShortName() string {
	return fmt.Sprintf("RTLSDR %s", f.GetTunerName())
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

func (f *RTLSDRFrontend) GetTuner() RTLSDRTuner {
	return RTLSDRTuner(f.device.GetTuner())
}

func (f *RTLSDRFrontend) GetTunerName() string {
	return rtlsdrTunerName[f.GetTuner()]
}

func (f *RTLSDRFrontend) SetOffsetTunning(enable bool) {
	b := byte(0)
	if enable {
		b = 1
	}
	f.device.SetOffsetTunning(b)
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

func (f *RTLSDRFrontend) SetGain1(gain int) {
	f.device.SetLNAGain(gain)
}

func (f *RTLSDRFrontend) SetGain2(gain int) {
	f.device.SetVGAGain(gain)
}

func (f *RTLSDRFrontend) SetGain3(gain int) {
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
