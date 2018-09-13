package Frontend

import (
	"github.com/OpenSatelliteProject/SatHelperApp/Logger"
	"unsafe"
)

const SampleTypeFloatIQ = 0
const SampleTypeS16IQ = 1
const SampleTypeS8IQ = 2

type SampleCallbackData struct {
	ComplexArray []complex64
	Int16Array   []int16
	Int8Array    []int8
	SampleType   int
	NumSamples   int
}

type SamplesCallback func(data SampleCallbackData)

type GoCallback struct {
	callback SamplesCallback
}

func NewGoCallback() GoCallback {
	return GoCallback{}
}

func (p *GoCallback) Info(str string) {
	SLog.Info("%s", str)
}

func (p *GoCallback) Error(str string) {
	SLog.Error("%s", str)
}

func (p *GoCallback) Warn(str string) {
	SLog.Warn("%s", str)
}

func (p *GoCallback) Debug(str string) {
	SLog.Debug("%s", str)
}

func (p *GoCallback) CbFloatIQ(data uintptr, length int) {
	const arrayLen = 1 << 20
	arr := (*[arrayLen]complex64)(unsafe.Pointer(data))[:length:length]
	if p.callback != nil {
		p.callback(SampleCallbackData{
			ComplexArray: arr,
			NumSamples:   length,
			SampleType:   SampleTypeFloatIQ,
		})
	}
}

func (p *GoCallback) CbS16IQ(data uintptr, length int) {
	// Length times two, because each sample contains an I and a Q in S16
	const arrayLen = 1 << 20
	var pairLength = length * 2
	arr := (*[arrayLen]int16)(unsafe.Pointer(data))[:pairLength:pairLength]
	if p.callback != nil {
		p.callback(SampleCallbackData{
			Int16Array: arr,
			NumSamples: length,
			SampleType: SampleTypeS16IQ,
		})
	}
}

func (p *GoCallback) CbS8IQ(data uintptr, length int) {
	// Length times two, because each sample contains an I and a Q in S8
	const arrayLen = 1 << 20
	var pairLength = length * 2
	arr := (*[arrayLen]int8)(unsafe.Pointer(data))[:pairLength:pairLength]
	if p.callback != nil {
		p.callback(SampleCallbackData{
			Int8Array:  arr,
			NumSamples: length,
			SampleType: SampleTypeS8IQ,
		})
	}
}

type BaseFrontend interface {
	SetSampleRate(sampleRate uint32) uint32
	SetCenterFrequency(centerFrequency uint32) uint32
	GetAvailableSampleRates() []uint32
	Start()
	Stop()
	SetAntenna(value string)
	SetAGC(agc bool)
	SetGain1(value uint8)
	SetGain2(value uint8)
	SetGain3(value uint8)
	SetBiasT(value bool)
	GetCenterFrequency() uint32
	GetName() string
	GetShortName() string
	GetSampleRate() uint32
	SetSamplesAvailableCallback(cb SamplesCallback)
	Init() bool
	Destroy()
}
