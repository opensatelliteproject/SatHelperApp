package Frontend

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

type BaseFrontend interface {
	SetSampleRate(sampleRate uint32) uint32
	SetCenterFrequency(centerFrequency uint32) uint32
	GetAvailableSampleRates() []uint32
	Start()
	Stop()
	SetAGC(agc bool)
	SetLNAGain(value uint8)
	SetVGAGain(value uint8)
	SetMixerGain(value uint8)
	SetBiasT(value bool)
	GetCenterFrequency() uint32
	GetName() string
	GetShortName() string
	GetSampleRate() uint32
	SetSamplesAvailableCallback(cb SamplesCallback)
}