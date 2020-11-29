package main

import (
	"flag"
	"fmt"
	ui "github.com/airking05/termui"
	"github.com/logrusorgru/aurora"
	"github.com/opensatelliteproject/SatHelperApp"
	"github.com/opensatelliteproject/SatHelperApp/DSP"
	"github.com/opensatelliteproject/SatHelperApp/Demuxer"
	"github.com/opensatelliteproject/SatHelperApp/Display"
	"github.com/opensatelliteproject/SatHelperApp/Frontend"
	"github.com/opensatelliteproject/SatHelperApp/ImageProcessor"
	"github.com/opensatelliteproject/SatHelperApp/Logger"
	"github.com/opensatelliteproject/SatHelperApp/RPC"
	"github.com/opensatelliteproject/SatHelperApp/metrics"
	"github.com/opensatelliteproject/libsathelper"
	"log"
	"net"
	"net/http"
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

	SLog.Log("%s %s (%s) - %s %s",
		aurora.Green(aurora.Bold("SatHelperApp")),
		aurora.Bold(SatHelperApp.GetVersion()),
		aurora.Bold(SatHelperApp.GetRevision()),
		aurora.Bold(SatHelperApp.GetCompilationDate()),
		aurora.Bold(SatHelperApp.GetCompilationTime()),
	)
	SLog.Log("%s %s (%s) - %s %s",
		aurora.Green(aurora.Bold("libSatHelper")),
		aurora.Bold(SatHelper.InfoGetVersion()),
		aurora.Bold(SatHelper.InfoGetGitSHA1()),
		aurora.Bold(SatHelper.InfoGetCompilationDate()),
		aurora.Bold(SatHelper.InfoGetCompilationTime()),
	)

	LoadConfig()

	DSP.StatisticsServer = Demuxer.NewTCPServer("", DSP.CurrentConfig.Base.StatisticsPort)

	DSP.StatisticsServer.Start()
	defer DSP.StatisticsServer.Stop()

	DSP.ConstellationServer = Demuxer.NewUDPServer("localhost", 9000)
	DSP.ConstellationServer.Start()
	defer DSP.ConstellationServer.Stop()

	if DSP.CurrentConfig.RPC.Enable {
		addr := fmt.Sprintf("%s:%d", DSP.CurrentConfig.RPC.ListenAddr, DSP.CurrentConfig.RPC.ListenPort)
		SLog.Info("Enabling gRPC at %s", addr)
		rpc := RPC.MakeRPCServer(rpcSource)
		defer rpc.Stop()
		err := rpc.Listen(addr)
		if err != nil {
			SLog.Error("Error starting gRPC: %s", err)
		}
	}

	if DSP.CurrentConfig.Prometheus.Enable {
		addr := fmt.Sprintf("%s:%d", DSP.CurrentConfig.Prometheus.ListenAddr, DSP.CurrentConfig.Prometheus.ListenPort)
		SLog.Info("Enabling Prometheus Metrics at %s", addr)
		metrics.EnablePrometheus()

		srv := &http.Server{}
		srv.Handler = metrics.GetHandler()
		srv.Addr = addr

		lis, err := net.Listen("tcp", addr)
		if err != nil {
			SLog.Error("Error starting prometheus server: %s", err)
			return
		}

		go func() {
			_ = srv.Serve(lis)
		}()
	}

	ImageProcessor.SetPurgeFiles(DSP.CurrentConfig.DirectDemuxer.PurgeFilesAfterProcess)
	DSP.CurrentConfig.Base.GenerateFFTImage = false // No use on SatHelperApp

	switch strings.ToLower(DSP.CurrentConfig.Base.Mode) {
	case "lrit":
		SLog.Info(aurora.Cyan("Selected LRIT mode. Ignoring parameters from config file.").String())
		DSP.SetLRITMode()
	case "hrit":
		SLog.Info(aurora.Cyan("Selected HRIT mode. Ignoring parameters from config file.").String())
		DSP.SetHRITMode()
	default:
		SLog.Info(aurora.Gray("No valid mode selected. Using config file parameters.").String())
	}

	switch strings.ToLower(DSP.CurrentConfig.Base.DeviceType) {
	case "cfile":
		SLog.Info(aurora.Cyan("CFile Frontend selected. File Name: %s").String(), aurora.Bold(aurora.Green(DSP.CurrentConfig.CFileSource.Filename)))
		DSP.Device = Frontend.NewCFileFrontend(DSP.CurrentConfig.CFileSource.Filename)
		if DSP.CurrentConfig.CFileSource.FastAsPossible {
			SLog.Info(aurora.Cyan("Fast as possible enabled!").String())
			DSP.Device.(*Frontend.CFileFrontend).EnableFastAsPossible()
		}
	case "limedrv":
		fallthrough
	case "lime":
		SLog.Info(aurora.Cyan("Lime Frontend selected.").String())
		DSP.Device = Frontend.NewLimeDrvFrontend()

		if !DSP.Device.Init() {
			SLog.Error("Error initializing device")
			return
		}
		defer DSP.Device.Destroy()

		DSP.Device.SetGain1(int(DSP.CurrentConfig.LimeSource.LNAGain))
		DSP.Device.SetAntenna(DSP.CurrentConfig.LimeSource.Antenna)

		SLog.Info(aurora.Cyan("	LNA Gain: %d").String(), aurora.Bold(aurora.Green(DSP.CurrentConfig.LimeSource.LNAGain)))
		SLog.Info(aurora.Cyan("	Antenna: %s").String(), aurora.Bold(aurora.Green(DSP.CurrentConfig.LimeSource.Antenna)))
	case "rtlsdr":
		SLog.Info(aurora.Cyan("RTLSDR Frontend selected.").String())
		DSP.Device = Frontend.NewRTLSDRFrontend()

		if !DSP.Device.Init() {
			SLog.Error("Error initializing device")
			return
		}
		defer DSP.Device.Destroy()

		DSP.Device.SetGain1(int(DSP.CurrentConfig.RtlsdrSource.LNAGain))
		DSP.Device.SetGain2(int(DSP.CurrentConfig.RtlsdrSource.VGAGain))
		DSP.Device.SetGain3(int(DSP.CurrentConfig.RtlsdrSource.MixerGain))
		DSP.Device.SetBiasT(DSP.CurrentConfig.RtlsdrSource.BiasTEnabled)
		DSP.Device.(*Frontend.RTLSDRFrontend).SetOffsetTunning(DSP.CurrentConfig.RtlsdrSource.OffsetTunning)
	case "airspy":
		SLog.Info(aurora.Cyan("Airspy Frontend selected.").String())
		Frontend.AirspyInitialize()
		defer Frontend.AirspyDeinitialize()
		DSP.Device = Frontend.NewAirspyFrontend()

		if !DSP.Device.Init() {
			SLog.Error("Error initializing device")
			return
		}
		defer DSP.Device.Destroy()

		DSP.Device.SetGain1(int(DSP.CurrentConfig.AirspySource.LNAGain))
		DSP.Device.SetGain2(int(DSP.CurrentConfig.AirspySource.VGAGain))
		DSP.Device.SetGain3(int(DSP.CurrentConfig.AirspySource.MixerGain))
		DSP.Device.SetBiasT(DSP.CurrentConfig.AirspySource.BiasTEnabled)
	case "spyserver":
		SLog.Info(aurora.Cyan("Spyserver Frontend Selected. Target: %s:%d").String(), aurora.Bold(DSP.CurrentConfig.SpyserverSource.Hostname), aurora.Bold(DSP.CurrentConfig.SpyserverSource.Port))
		DSP.Device = Frontend.NewSpyserverFrontend(DSP.CurrentConfig.SpyserverSource.Hostname, DSP.CurrentConfig.SpyserverSource.Port)
		if !DSP.Device.Init() {
			SLog.Error("Error initializing device")
			return
		}
		defer DSP.Device.Destroy()
		DSP.Device.SetGain1(int(DSP.CurrentConfig.SpyserverSource.Gain))
	default:
		SLog.Error(aurora.Red("Device %s is not currently supported.").String(), aurora.Bold(DSP.CurrentConfig.Base.DeviceType))
		return
	}

	switch strings.ToLower(DSP.CurrentConfig.Base.DemuxerType) {
	case "direct":
		SLog.Info(aurora.Cyan("Direct Internal Demuxer selected.").String())
		dd := Demuxer.MakeDirectDemuxer(
			DSP.CurrentConfig.DirectDemuxer.OutputFolder,
			DSP.CurrentConfig.DirectDemuxer.TemporaryFolder,
			DSP.CurrentConfig.DirectDemuxer.DrawMap,
			DSP.CurrentConfig.DirectDemuxer.ReprojectImages,
			DSP.CurrentConfig.DirectDemuxer.FalseColor,
			DSP.CurrentConfig.DirectDemuxer.MetaFrame,
			DSP.CurrentConfig.DirectDemuxer.Enhanced)
		for _, v := range DSP.CurrentConfig.DirectDemuxer.SkipVCID {
			dd.AddSkipVCID(v)
		}
		DSP.SDemuxer = dd
	case "tcpserver":
		SLog.Info(aurora.Cyan("TCP Server Demuxer selected. Will listen %s:%d").String(), aurora.Bold(DSP.CurrentConfig.TCPServerDemuxer.Host), aurora.Bold(DSP.CurrentConfig.TCPServerDemuxer.Port))
		DSP.SDemuxer = Demuxer.NewTCPServer(DSP.CurrentConfig.TCPServerDemuxer.Host, DSP.CurrentConfig.TCPServerDemuxer.Port)
	case "file":
		if DSP.CurrentConfig.FileDemuxer.Filename == "" {
			DSP.CurrentConfig.FileDemuxer.Filename = fmt.Sprintf("demuxdump-%d.bin", time.Now().Unix())
		}
		SLog.Info(aurora.Cyan("File Demuxer selected. Will write to %s").String(), aurora.Bold(DSP.CurrentConfig.FileDemuxer.Filename))
		DSP.SDemuxer = Demuxer.NewFileDemuxer(DSP.CurrentConfig.FileDemuxer.Filename)
	default:
		SLog.Error("Unknown Demuxer Type %s.", DSP.CurrentConfig.Base.DemuxerType)
		return
	}

	if DSP.Device.SetSampleRate(DSP.CurrentConfig.Source.SampleRate) != DSP.CurrentConfig.Source.SampleRate {
		SLog.Warn("Cannot set sample rate.")
	}

	if DSP.Device.SetCenterFrequency(DSP.CurrentConfig.Source.Frequency) != DSP.CurrentConfig.Source.Frequency {
		SLog.Warn("Cannot set exact frequency. Current Value: %d", DSP.Device.GetCenterFrequency())
	}

	DSP.InitAll()

	if DSP.CurrentConfig.Decoder.Display {
		Display.InitDisplay()
	}

	DSP.SDemuxer.Init()
	DSP.Device.Start()
	defer DSP.Device.Stop()

	SLog.Info("Starting main loop")

	DSP.SDemuxer.Start()
	defer DSP.SDemuxer.Stop()

	DSP.StartDSPLoops()
	DSP.SetRunning(true)

	stopFunc := func(ui.Event) {
		SLog.Warn(aurora.Bold("Got close handler").String())
		DSP.SetRunning(false)
		SLog.SetTermUiDisplay(false)
		ui.StopLoop()
		time.Sleep(1 * time.Second)
	}

	ui.Handle("/sys/kbd/q", stopFunc)
	ui.Handle("/sys/kbd/C-c", stopFunc)

	if DSP.CurrentConfig.Decoder.Display {
		ui.Handle("/timer/100ms", func(e ui.Event) {
			freq := float32(DSP.Device.GetCenterFrequency()) + DSP.GetCostasFrequency()
			stat := DSP.GetStats()
			Display.UpdateSignalQuality(stat.SignalQuality)
			Display.UpdateLockedState(stat.FrameLock == 1)
			Display.UpdateChannelData(stat.ReceivedPacketsPerChannel)
			Display.UpdateReedSolomon(stat.RsErrors)
			Display.UpdateSyncWord(stat.SyncWord)
			Display.UpdateSCVCID(stat.SCID, stat.VCID)
			Display.UpdateDecoderFifoUsage(DSP.GetDecoderFIFOUsage())
			Display.UpdateDemodulatorFifoUsage(DSP.GetDemodFIFOUsage())
			Display.UpdateViterbiErrors(uint(stat.VitErrors), uint(stat.FrameBits))
			Display.UpdatePhaseCorr(stat.PhaseCorrection)
			Display.UpdateSyncCorrelation(stat.SyncCorrelation)
			Display.UpdateMode(strings.ToUpper(DSP.CurrentConfig.Base.Mode))
			Display.UpdateCenterFrequency(uint32(freq))
			Display.UpdateDevice(DSP.Device.GetShortName())
			Display.UpdateDemuxer(DSP.SDemuxer.GetName())
			Display.Render()
		})
		CallClear()
	}
	ui.Loop()
}
