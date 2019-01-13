package main

import (
	"bytes"
	"flag"
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/OpenSatelliteProject/SatHelperApp/DSP"
	"github.com/OpenSatelliteProject/SatHelperApp/Logger"
	"github.com/mitchellh/go-homedir"
	"io/ioutil"
	"log"
	"os"
)

var configFile = flag.String("config", "", "write cpu profile to file")
var finalConfigFilePath string

func LoadDefaults() {
	DSP.CurrentConfig.Title = "SatHelperApp"

	DSP.SetHRITMode()

	// Other options
	DSP.CurrentConfig.Source.SampleRate = DSP.DefaultSampleRate
	DSP.CurrentConfig.Base.Decimation = DSP.DefaultDecimation
	DSP.CurrentConfig.Base.AGCEnabled = true
	DSP.CurrentConfig.Base.DeviceType = "airspy"
	DSP.CurrentConfig.Base.SendConstellation = true
	DSP.CurrentConfig.Base.StatisticsPort = DSP.DefaultStatisticsPort

	// Airspy Source Defaults
	DSP.CurrentConfig.AirspySource.LNAGain = DSP.DefaultLnaGain
	DSP.CurrentConfig.AirspySource.MixerGain = DSP.DefaultMixGain
	DSP.CurrentConfig.AirspySource.VGAGain = DSP.DefaultVgaGain
	DSP.CurrentConfig.AirspySource.BiasTEnabled = DSP.DefaultBiast

	// CFile Source Defaults
	DSP.CurrentConfig.CFileSource.Filename = ""
	DSP.CurrentConfig.CFileSource.FastAsPossible = false

	// LimeSDR Source Defaults
	DSP.CurrentConfig.LimeSource.LNAGain = 10
	DSP.CurrentConfig.LimeSource.Antenna = "LNAH"

	// Spyserver
	DSP.CurrentConfig.SpyserverSource.Hostname = "127.0.0.1"
	DSP.CurrentConfig.SpyserverSource.Port = 5555
	DSP.CurrentConfig.SpyserverSource.Gain = 20

	// Decoder
	DSP.CurrentConfig.Decoder.Display = true
	DSP.CurrentConfig.Decoder.UseLastFrameData = true

	// Others
	DSP.CurrentConfig.Base.DemuxerType = "tcpserver"

	// TCPDemuxer
	DSP.CurrentConfig.TCPServerDemuxer.Port = DSP.DefaultVchannelPort
	DSP.CurrentConfig.TCPServerDemuxer.Host = ""

	// FileDemuxer
	DSP.CurrentConfig.FileDemuxer.Filename = ""

	// Direct Demuxer
	DSP.CurrentConfig.DirectDemuxer.OutputFolder = "out"
	DSP.CurrentConfig.DirectDemuxer.TemporaryFolder = "tmp"

	SaveConfig()
}

func SaveConfig() {
	var firstBuffer bytes.Buffer
	e := toml.NewEncoder(&firstBuffer)
	err := e.Encode(DSP.CurrentConfig)
	if err != nil {
		log.Printf("Cannot save config: %s", err)
		return
	}
	SLog.Info("Saving config file to %s", finalConfigFilePath)
	err = ioutil.WriteFile(finalConfigFilePath, firstBuffer.Bytes(), 0644)
	if err != nil {
		log.Printf("Cannot save config: %s", err)
		return
	}
}

func LoadConfig() {
	home, _ := homedir.Dir()

	err := os.MkdirAll(fmt.Sprintf("%s/SatHelperApp", home), os.ModePerm)
	if err != nil {
		panic(err)
	}

	finalConfigFilePath = fmt.Sprintf("%s/SatHelperApp/%s", home, "SatHelperApp.cfg")

	if *configFile != "" {
		finalConfigFilePath = *configFile
	}

	SLog.Info("Loading config file from %s", finalConfigFilePath)
	_, err = toml.DecodeFile(finalConfigFilePath, &DSP.CurrentConfig)

	if err != nil {
		SLog.Warn("Cannot load file SatHelperApp.cfg. Loading default values.")
		LoadDefaults()
	}
}
