package main

import (
	"flag"
	"fmt"
	"github.com/OpenSatelliteProject/SatHelperApp/Demuxer"
	"github.com/OpenSatelliteProject/SatHelperApp/Display"
	"github.com/OpenSatelliteProject/SatHelperApp/Frontend"
	"github.com/OpenSatelliteProject/SatHelperApp/Logger"
	"github.com/OpenSatelliteProject/libsathelper"
	ui "github.com/airking05/termui"
	"github.com/logrusorgru/aurora"
	"log"
	"os"
	"runtime/pprof"
	"strings"
	"time"
)

var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")

func main() {
	SLog.StartLog()
	defer SLog.EndLog()
	err := ui.Init()
	if err != nil {
		panic(err)
	}
	defer ui.Close()

	flag.Parse()
	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Println(err)
			return
		}
		err = pprof.StartCPUProfile(f)
		if err != nil {
			panic(err)
		}
		defer pprof.StopCPUProfile()
	}

	startTime = uint32(time.Now().Unix() & 0xFFFFFFFF)

	SLog.Log("%s %s (%s) - %s %s",
		aurora.Green(aurora.Bold("SatHelperApp")),
		aurora.Bold(GetVersion()),
		aurora.Bold(GetRevision()),
		aurora.Bold(GetCompilationDate()),
		aurora.Bold(GetCompilationTime()),
	)
	SLog.Log("%s %s (%s) - %s %s",
		aurora.Green(aurora.Bold("libSatHelper")),
		aurora.Bold(SatHelper.InfoGetVersion()),
		aurora.Bold(SatHelper.InfoGetGitSHA1()),
		aurora.Bold(SatHelper.InfoGetCompilationDate()),
		aurora.Bold(SatHelper.InfoGetCompilationTime()),
	)

	LoadConfig()

	statisticsServer = Demuxer.NewTCPServer("", CurrentConfig.Base.StatisticsPort)

	statisticsServer.Start()
	defer statisticsServer.Stop()

	constellationServer = Demuxer.NewUDPServer("localhost", 9000)
	constellationServer.Start()
	defer constellationServer.Stop()

	switch strings.ToLower(CurrentConfig.Base.Mode) {
	case "lrit":
		SLog.Info(aurora.Cyan("Selected LRIT mode. Ignoring parameters from config file.").String())
		SetLRITMode()
	case "hrit":
		SLog.Info(aurora.Cyan("Selected HRIT mode. Ignoring parameters from config file.").String())
		SetHRITMode()
	default:
		SLog.Info(aurora.Gray("No valid mode selected. Using config file parameters.").String())
	}

	switch strings.ToLower(CurrentConfig.Base.DeviceType) {
	case "cfile":
		SLog.Info(aurora.Cyan("CFile Frontend selected. File Name: %s").String(), aurora.Bold(aurora.Green(CurrentConfig.CFileSource.Filename)))
		device = Frontend.NewCFileFrontend(CurrentConfig.CFileSource.Filename)
		if CurrentConfig.CFileSource.FastAsPossible {
			SLog.Info(aurora.Cyan("Fast as possible enabled!").String())
			device.(*Frontend.CFileFrontend).EnableFastAsPossible()
		}
	case "lime":
		SLog.Info(aurora.Cyan("LimeSDR Frontend selected.").String())
		device = Frontend.NewLimeFrontend()

		if !device.Init() {
			SLog.Error("Error initializing device")
			return
		}
		defer device.Destroy()

		device.SetGain1(CurrentConfig.LimeSource.LNAGain)
		device.SetAntenna(CurrentConfig.LimeSource.Antenna)

		SLog.Info(aurora.Cyan("	LNA Gain: %d").String(), aurora.Bold(aurora.Green(CurrentConfig.LimeSource.LNAGain)))
		SLog.Info(aurora.Cyan("	Antenna: %s").String(), aurora.Bold(aurora.Green(CurrentConfig.LimeSource.Antenna)))
	case "airspy":
		SLog.Info(aurora.Cyan("Airspy Frontend selected.").String())
		Frontend.AirspyInitialize()
		defer Frontend.AirspyDeinitialize()
		device = Frontend.NewAirspyFrontend()

		if !device.Init() {
			SLog.Error("Error initializing device")
			return
		}
		defer device.Destroy()

		device.SetGain1(CurrentConfig.AirspySource.LNAGain)
		device.SetGain2(CurrentConfig.AirspySource.VGAGain)
		device.SetGain3(CurrentConfig.AirspySource.MixerGain)
		device.SetBiasT(CurrentConfig.AirspySource.BiasTEnabled)
	case "spyserver":
		SLog.Info(aurora.Cyan("Spyserver Frontend Selected. Target: %s:%d").String(), aurora.Bold(CurrentConfig.SpyserverSource.Hostname), aurora.Bold(CurrentConfig.SpyserverSource.Port))
		device = Frontend.NewSpyserverFrontend(CurrentConfig.SpyserverSource.Hostname, CurrentConfig.SpyserverSource.Port)
		if !device.Init() {
			SLog.Error("Error initializing device")
			return
		}
		defer device.Destroy()
		device.SetGain1(CurrentConfig.SpyserverSource.Gain)
	default:
		SLog.Error(aurora.Red("Device %s is not currently supported.").String(), aurora.Bold(CurrentConfig.Base.DeviceType))
		return
	}

	switch strings.ToLower(CurrentConfig.Base.DemuxerType) {
	case "direct":
		SLog.Info(aurora.Cyan("Direct Internal Demuxer selected.").String())
		demuxer = Demuxer.MakeDirectDemuxer()
	case "tcpserver":
		SLog.Info(aurora.Cyan("TCP Server Demuxer selected. Will listen %s:%d").String(), aurora.Bold(CurrentConfig.TCPServerDemuxer.Host), aurora.Bold(CurrentConfig.TCPServerDemuxer.Port))
		demuxer = Demuxer.NewTCPServer(CurrentConfig.TCPServerDemuxer.Host, CurrentConfig.TCPServerDemuxer.Port)
	case "file":
		if CurrentConfig.FileDemuxer.Filename == "" {
			CurrentConfig.FileDemuxer.Filename = fmt.Sprintf("demuxdump-%d.bin", time.Now().Unix())
		}
		SLog.Info(aurora.Cyan("File Demuxer selected. Will write to %s").String(), aurora.Bold(CurrentConfig.FileDemuxer.Filename))
		demuxer = Demuxer.NewFileDemuxer(CurrentConfig.FileDemuxer.Filename)
	default:
		SLog.Error("Unknown Demuxer Type %s.", CurrentConfig.Base.DemuxerType)
		return
	}

	if device.SetSampleRate(CurrentConfig.Source.SampleRate) != CurrentConfig.Source.SampleRate {
		SLog.Warn("Cannot set sample rate.")
	}

	if device.SetCenterFrequency(CurrentConfig.Source.Frequency) != CurrentConfig.Source.Frequency {
		SLog.Warn("Cannot set exact frequency. Current Value: %d", device.GetCenterFrequency())
	}

	device.SetSamplesAvailableCallback(newSamplesCallback)

	initDSP()
	initDecoder()

	if CurrentConfig.Decoder.Display {
		Display.InitDisplay()
	}

	demuxer.Init()
	device.Start()
	defer device.Stop()

	SLog.Info("Starting main loop")

	demuxer.Start()
	defer demuxer.Stop()

	go symbolProcessLoop()
	go decoderLoop()

	running = true

	stopFunc := func(ui.Event) {
		SLog.Warn(aurora.Bold("Got close handler").String())
		running = false
		SLog.SetTermUiDisplay(false)
		ui.StopLoop()
		time.Sleep(1 * time.Second)
	}

	ui.Handle("/sys/kbd/q", stopFunc)
	ui.Handle("/sys/kbd/C-c", stopFunc)

	if CurrentConfig.Decoder.Display {
		ui.Handle("/timer/100ms", func(e ui.Event) {
			stat := GetStats()
			Display.UpdateSignalQuality(stat.SignalQuality)
			Display.UpdateLockedState(stat.FrameLock == 1)
			Display.UpdateChannelData(stat.ReceivedPacketsPerChannel)
			Display.UpdateReedSolomon(stat.RsErrors)
			Display.UpdateSyncWord(stat.SyncWord)
			Display.UpdateSCVCID(stat.SCID, stat.VCID)
			Display.UpdateDecoderFifoUsage(decodFifoUsage)
			Display.UpdateDemodulatorFifoUsage(demodFifoUsage)
			Display.UpdateViterbiErrors(uint(stat.VitErrors), uint(stat.FrameBits))
			Display.UpdatePhaseCorr(stat.PhaseCorrection)
			Display.UpdateSyncCorrelation(stat.SyncCorrelation)
			Display.UpdateMode(strings.ToUpper(CurrentConfig.Base.Mode))
			Display.UpdateCenterFrequency(device.GetCenterFrequency())
			Display.UpdateDevice(device.GetShortName())
			Display.UpdateDemuxer(demuxer.GetName())
			Display.Render()
		})
		CallClear()
	}
	ui.Loop()
}
