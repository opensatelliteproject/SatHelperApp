package main

import (
	"github.com/OpenSatelliteProject/libsathelper"
	"log"
	"github.com/OpenSatelliteProject/SatHelperApp/Frontend"
	ui "github.com/airking05/termui"
	"github.com/logrusorgru/aurora"
	"github.com/OpenSatelliteProject/SatHelperApp/Display"
	"strings"
	"github.com/OpenSatelliteProject/SatHelperApp/Demuxer"
	"flag"
	"os"
	"runtime/pprof"
)

var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")

func main() {

	flag.Parse()
	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Println(err)
			return
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	log.Printf("%s %s (%s) - %s %s\n",
		aurora.Green(aurora.Bold("SatHelperApp")),
		aurora.Bold(GetVersion()),
		aurora.Bold(GetRevision()),
		aurora.Bold(GetCompilationDate()),
		aurora.Bold(GetCompilationTime()),
	)
	log.Printf("%s %s (%s) - %s %s\n",
		aurora.Green(aurora.Bold("libSatHelper")),
		aurora.Bold(SatHelper.InfoGetVersion()),
		aurora.Bold(SatHelper.InfoGetGitSHA1()),
		aurora.Bold(SatHelper.InfoGetCompilationDate()),
		aurora.Bold(SatHelper.InfoGetCompilationTime()),
	)

	LoadConfig()

	switch strings.ToLower(CurrentConfig.Base.Mode) {
		case "lrit":
			log.Println(aurora.Cyan("Selected LRIT mode. Ignoring parameters from config file."))
			SetLRITMode()
		break
		case "hrit":
			log.Println(aurora.Cyan("Selected HRIT mode. Ignoring parameters from config file."))
			SetHRITMode()
		break
		default:
			log.Println(aurora.Gray("No valid mode selected. Using config file parameters."))
	}

	switch strings.ToLower(CurrentConfig.Base.DeviceType) {
		case "cfile":
			log.Printf(aurora.Cyan("CFile Frontend selected. File Name: %s").String(), aurora.Bold(aurora.Green(CurrentConfig.CFileSource.Filename)))
			device = Frontend.NewCFileFrontend(CurrentConfig.CFileSource.Filename)
			break
		case "lime":
			log.Print(aurora.Cyan("LimeSDR Frontend selected."))
			device = Frontend.NewLimeFrontend()

			device.SetGain1(CurrentConfig.LimeSource.LNAGain)
			device.SetGain2(CurrentConfig.LimeSource.TIAGain)
			device.SetGain3(CurrentConfig.LimeSource.PGAGain)
			device.SetAntenna(CurrentConfig.LimeSource.Antenna)

			log.Printf(aurora.Cyan("LNA Gain: %d").String(), aurora.Bold(aurora.Green(CurrentConfig.LimeSource.LNAGain)))
			log.Printf(aurora.Cyan("TIA Gain: %d").String(), aurora.Bold(aurora.Green(CurrentConfig.LimeSource.TIAGain)))
			log.Printf(aurora.Cyan("PGA Gain: %d").String(), aurora.Bold(aurora.Green(CurrentConfig.LimeSource.LNAGain)))
			log.Printf(aurora.Cyan("Antenna: %s").String(), aurora.Bold(aurora.Green(CurrentConfig.LimeSource.Antenna)))
			break
		case "airspy":
			log.Print(aurora.Cyan("Airspy Frontend selected."))
			Frontend.AirspyInitialize()
			defer Frontend.AirspyDeinitialize()
			device = Frontend.NewAirspyFrontend()

			if ! device.Init() {
				log.Println("Error initializing device")
				return
			}
			defer device.Destroy()

			device.SetGain1(CurrentConfig.AirspySource.LNAGain)
			device.SetGain2(CurrentConfig.AirspySource.VGAGain)
			device.SetGain3(CurrentConfig.AirspySource.MixerGain)
			device.SetBiasT(CurrentConfig.AirspySource.BiasTEnabled)
			break
		case "spyserver":
			log.Printf(aurora.Cyan("Spyserver Frontend Selected. Target: %s:%d").String(), aurora.Bold(CurrentConfig.SpyserverSource.Hostname), aurora.Bold(CurrentConfig.SpyserverSource.Port))
			device = Frontend.NewSpyserverFrontend(CurrentConfig.SpyserverSource.Hostname, CurrentConfig.SpyserverSource.Port)
			if ! device.Init() {
				log.Println("Error initializing device")
				return
			}
			defer device.Destroy()
			device.SetLNAGain(CurrentConfig.SpyserverSource.Gain)
			break
		default:
			log.Println(aurora.Red("Device %s is not currently supported.").String(), aurora.Bold(CurrentConfig.Base.DeviceType))
			return
		break
	}

	switch strings.ToLower(CurrentConfig.Base.DemuxerType) {
	case "tcpserver":
		log.Printf(aurora.Cyan("TCP Server Demuxer selected. Will listen %s:%d").String(), aurora.Bold(CurrentConfig.TCPServerDemuxer.Host), aurora.Bold(CurrentConfig.TCPServerDemuxer.Port))
		demuxer = Demuxer.NewTCPDemuxer(CurrentConfig.TCPServerDemuxer.Host, CurrentConfig.TCPServerDemuxer.Port)
		break
	default:
		log.Printf(aurora.Red("Unknown Demuxer Type %s.").String(), CurrentConfig.Base.DemuxerType)
		return
	}

	if device.SetSampleRate(CurrentConfig.Source.SampleRate) != CurrentConfig.Source.SampleRate {
		log.Println("Cannot set sample rate.")
	}

	if device.SetCenterFrequency(CurrentConfig.Source.Frequency) != CurrentConfig.Source.Frequency {
		log.Printf("Cannot set exact frequency. Current Value: %d\n", device.GetCenterFrequency())
	}

	device.SetSamplesAvailableCallback(newSamplesCallback)

	initDSP()
	initDecoder()

	err := ui.Init()
	if err != nil {
		panic(err)
	}
	defer ui.Close()

	Display.InitDisplay()
	demuxer.Init()

	log.Println(aurora.Cyan("Starting Source"))
	device.Start()
	defer device.Stop()

	log.Println(aurora.Cyan("Starting Main loop"))

	demuxer.Start()
	defer demuxer.Stop()

	go symbolProcessLoop()
	go decoderLoop()

	running = true

	stopFunc := func(ui.Event) {
		log.Println(aurora.Bold(aurora.Red("Got close handler.")))
		running = false
		ui.StopLoop()
	}

	ui.Handle("/sys/kbd/q", stopFunc)
	ui.Handle("/sys/kbd/C-c", stopFunc)
	ui.Handle("/timer/100ms", func (e ui.Event) {
		stat := GetStats()
		Display.UpdateSignalQuality(stat.SignalQuality)
		Display.UpdateLockedState(stat.FrameLock == 1)
		Display.UpdateChannelData(stat.ReceivedPacketsPerChannel)
		Display.UpdateReedSolomon(stat.RsErrors)
		Display.UpdateSyncWord(stat.SyncWord)
		Display.UpdateSCVCID(stat.SCID, stat.VCID)
		Display.UpdateDecoderFifoUsage(stat.DecoderFifoUsage)
		Display.UpdateDemodulatorFifoUsage(stat.DemodulatorFifoUsage)
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

	ui.Loop()
}
