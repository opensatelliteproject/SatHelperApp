package Frontend

import (
	"unsafe"
)

const FrontendSampletypeFloatiq = 0
const FrontendSampletypeS16iq = 1
const FrontendSampletypeS8iq = 2

type SampleCallbackData struct {
	ComplexArray []complex64
	Int16Array []int16
	Int8Array []int8
	SampleType int
	NumSamples int
}

type SamplesCallback func(data SampleCallbackData)

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
}