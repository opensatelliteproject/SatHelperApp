package Frontend

const FRONTEND_SAMPLETYPE_FLOATIQ = 0
const FRONTEND_SAMPLETYPE_S16IQ = 1
const FRONTEND_SAMPLETYPE_S8IQ = 2

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
	SetBiasT(value uint8)
	GetCenterFrequency() uint32
	GetName() string
	GetShortName() string
	GetSampleRate() uint32
	SetSamplesAvailableCallback(cb SamplesCallback)
}