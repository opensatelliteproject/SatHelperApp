package main

import (
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/mewkiz/pkg/osutil"
	"github.com/mitchellh/go-homedir"
	"github.com/opensatelliteproject/SatHelperApp/DSP"
	"github.com/opensatelliteproject/SatHelperApp/Models"
	"github.com/prometheus/common/log"
	"github.com/quan-to/slog"
	"io"
	"os"
)

var finalConfigFilePath string
var configLoaded bool
var satlog = slog.Scope("SatUI")

func isConfigLoaded() bool {
	log.Info("isConfigLoaded()")
	return configLoaded
}

func getConfig() Models.AppConfig {
	log.Info("getConfig()")
	return DSP.CurrentConfig
}

func setConfig(config Models.AppConfig) {
	log.Info("setConfig()")
	DSP.CurrentConfig = config
}

func saveConfig() error {
	log.Info("saveConfig()")
	f, err := os.OpenFile(finalConfigFilePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	_ = f.Truncate(0)
	_, _ = f.Seek(0, io.SeekStart)

	enc := toml.NewEncoder(f)
	enc.Indent = "  "
	return enc.Encode(&DSP.CurrentConfig)
}

func loadConfig() {
	log.Info("loadConfig()")
	home, _ := homedir.Dir()

	err := os.MkdirAll(fmt.Sprintf("%s/SatHelperApp", home), os.ModePerm)
	if err != nil {
		panic(err)
	}

	finalConfigFilePath = fmt.Sprintf("%s/SatHelperApp/%s", home, "SatHelperApp.cfg")

	if !osutil.Exists(finalConfigFilePath) {
		satlog.Warn("Config file %s does not exists. Creating one with defaults.", finalConfigFilePath)
		configLoaded = false
	} else {
		satlog.Info("Loading config file from %s", finalConfigFilePath)
		_, err = toml.DecodeFile(finalConfigFilePath, &DSP.CurrentConfig)

		if err != nil {
			satlog.Warn("Cannot load file SatHelperApp.cfg. Loading default values.")
			configLoaded = false
		} else {
			configLoaded = true
		}
	}

	if !configLoaded {
		loadDefaults()
	}
}

func loadDefaults() {
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

	// RTLSDR Source Defaults
	DSP.CurrentConfig.RtlsdrSource.LNAGain = DSP.DefaultLnaGain
	DSP.CurrentConfig.RtlsdrSource.MixerGain = DSP.DefaultMixGain
	DSP.CurrentConfig.RtlsdrSource.VGAGain = DSP.DefaultVgaGain
	DSP.CurrentConfig.RtlsdrSource.BiasTEnabled = DSP.DefaultBiast
	DSP.CurrentConfig.RtlsdrSource.OffsetTunning = DSP.DefaultOffsetTunning

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
	DSP.CurrentConfig.DirectDemuxer.PurgeFilesAfterProcess = false
	DSP.CurrentConfig.DirectDemuxer.SkipVCID = make([]int, 0)
	DSP.CurrentConfig.DirectDemuxer.DrawMap = false
	DSP.CurrentConfig.DirectDemuxer.ReprojectImages = false
	DSP.CurrentConfig.DirectDemuxer.FalseColor = false
	DSP.CurrentConfig.DirectDemuxer.Enhanced = false
	DSP.CurrentConfig.DirectDemuxer.MetaFrame = true

	// RPC
	DSP.CurrentConfig.RPC.Enable = true
	DSP.CurrentConfig.RPC.ListenAddr = ""
	DSP.CurrentConfig.RPC.ListenPort = DSP.DefaultRPCPort

	// Prometheus
	DSP.CurrentConfig.Prometheus.Enable = true
	DSP.CurrentConfig.Prometheus.ListenAddr = ""
	DSP.CurrentConfig.Prometheus.ListenPort = DSP.DefaultPrometheusPort
}
