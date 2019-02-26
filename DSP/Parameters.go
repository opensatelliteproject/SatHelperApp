package DSP

import (
	. "github.com/opensatelliteproject/SatHelperApp/Models"
)

// region Demodulator Parameters
// These are the parameters used by the demodulator. Change with care.

// GOES HRIT Settings
const HritCenterFrequency = 1694100000
const HritSymbolRate = 927000
const HritRrcAlpha float32 = 0.3

// GOES LRIT Settings
const LritCenterFrequency = 1691000000
const LritSymbolRate = 293883
const LritRrcAlpha = 0.5

// Loop Settings
const LoopOrder = 2
const RrcTaps = 31
const PllAlpha float32 = 0.001
const ClockAlpha float32 = 0.0037
const ClockMu float32 = 0.5
const ClockOmegaLimit float32 = 0.005
const ClockGainOmega = (ClockAlpha * ClockAlpha) / 4.0
const AgcRate float32 = 0.01
const AgcReference float32 = 0.5
const AgcGain float32 = 1.0
const AgcMaxGain float32 = 4000

const AirspyMiniDefaultSamplerate = 3000000

//const AirspyR2DefaultSamplerate = 2500000
const DefaultSampleRate = AirspyMiniDefaultSamplerate
const DefaultDecimation = 1

//const DefaultDeviceNumber = 0

const DefaultLnaGain = 5
const DefaultVgaGain = 5
const DefaultMixGain = 5

const DefaultBiast = false

// FIFO Size in Samples
// 10 * 1024 * 1024 samples is about 40Mb of ram.
// This should be more than enough
const FifoSize = 10 * 1024 * 1024

// endregion
// region Decoder Parameters

const HritUw0 uint64 = 0xfc4ef4fd0cc2df89
const HritUw2 uint64 = 0x25010b02f33d2076
const LritUw0 uint64 = 0xfca2b63db00d9794
const LritUw2 uint64 = 0x035d49c24ff2686b

const SyncWordSize = 4
const FrameSize = 1024
const FrameBits = FrameSize * 8
const CodedFrameSize = FrameBits * 2
const MinCorrelationBits = 46
const RsBlocks = 4
const RsParitySize = 32
const RsParityBlockSize = RsParitySize * RsBlocks
const LastFrameDataBits = 64
const LastFrameData = LastFrameDataBits / 8

const DefaultFlywheelRecheck = 100
const DefaultVchannelPort = 5001

const DefaultStatisticsPort = 5002

const AverageLastNSamples = 10000

// endregion
// region Current Config Stuff
var CurrentConfig AppConfig

func SetHRITMode() {
	// HRIT Mode
	CurrentConfig.Base.SymbolRate = HritSymbolRate
	CurrentConfig.Base.Mode = "hrit"
	CurrentConfig.Base.RRCAlpha = HritRrcAlpha
	CurrentConfig.Source.Frequency = HritCenterFrequency
}

func SetLRITMode() {
	// LRIT Mode
	CurrentConfig.Base.SymbolRate = LritSymbolRate
	CurrentConfig.Base.Mode = "lrit"
	CurrentConfig.Base.RRCAlpha = LritRrcAlpha
	CurrentConfig.Source.Frequency = LritCenterFrequency
}

// endregion
