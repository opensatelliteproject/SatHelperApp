package main

import (
	"log"
	"github.com/BurntSushi/toml"
	"bytes"
	"io/ioutil"
)

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
const ClockGainOmega float32 = (ClockAlpha * ClockAlpha) / 4.0
const AgcRate float32 = 0.01
const AgcReference float32 = 0.5
const AgcGain float32 = 1.0
const AgcMaxGain float32 = 4000

const AirspyMiniDefaultSamplerate = 3000000
const AirspyR2DefaultSamplerate = 2500000
const DefaultSampleRate = AirspyMiniDefaultSamplerate
const DefaultDecimation = 1
const DefaultDeviceNumber = 0

const DefaultDecoderAddress = "127.0.0.1"
const DefaultDecoderPort = 5000

const DefaultLnaGain = 5
const DefaultVgaGain = 5
const DefaultMixGain = 5

const DefaultBiast = false

// FIFO Size in Samples
// 1024 * 1024 samples is about 4Mb of ram.
// This should be more than enough
const FifoSize = 1024 * 1024

type BaseConfig struct {
	SymbolRate uint32
	RRCAlpha float32
	Mode string
	Decimation uint8
	AGCEnabled bool
	DeviceType string
	SendConstellation bool
	PLLAlpha float32
}

type CFileSourceConfig struct {
	Filename string
}

type AirspySourceConfig struct {
	MixerGain uint8
	LNAGain uint8
	VGAGain uint8
	BiasTEnabled bool
}

type SpyServerConfig struct {
	SpyServerHost string
	SpyServerPort int
	BiasTEnabled bool
}

type SourceConfig struct {
	SampleRate uint32
	Frequency uint32
}

type AppConfig struct {
	Title string
	Base BaseConfig
	Source SourceConfig
	AirspySource AirspySourceConfig
	SpyServerSource SpyServerConfig
	CFileSource CFileSourceConfig
}

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

func LoadDefaults() {
	CurrentConfig.Title = "SatHelperApp"

	SetHRITMode()

	// Other options
	CurrentConfig.Source.SampleRate = DefaultSampleRate
	CurrentConfig.Base.Decimation = DefaultDecimation
	CurrentConfig.Base.AGCEnabled = true
	CurrentConfig.AirspySource.LNAGain = DefaultLnaGain
	CurrentConfig.AirspySource.MixerGain = DefaultMixGain
	CurrentConfig.AirspySource.VGAGain = DefaultVgaGain
	CurrentConfig.Base.DeviceType = "airspy"
	CurrentConfig.Base.SendConstellation = true
	CurrentConfig.AirspySource.BiasTEnabled = DefaultBiast
	CurrentConfig.SpyServerSource.BiasTEnabled = DefaultBiast

	SaveConfig()
}

func SaveConfig() {
	var firstBuffer bytes.Buffer
	e := toml.NewEncoder(&firstBuffer)
	err := e.Encode(CurrentConfig)
	if err != nil {
		log.Fatalf("Cannot save config: %s", err)
	}
	err = ioutil.WriteFile("SatHelperApp.cfg", firstBuffer.Bytes(), 0644)
	if err != nil {
		log.Fatalf("Cannot save config: %s", err)
	}
}

func LoadConfig() {
	if _, err := toml.DecodeFile("SatHelperApp.cfg", &CurrentConfig); err != nil {
		log.Println("Cannot load file SatHelperApp.cfg. Loading default values.")
		LoadDefaults()
	}
}