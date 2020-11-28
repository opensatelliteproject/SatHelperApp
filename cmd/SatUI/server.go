package main

import (
	"fmt"
	"github.com/alecthomas/kong"
	"github.com/logrusorgru/aurora"
	"github.com/opensatelliteproject/SatHelperApp"
	"github.com/opensatelliteproject/SatHelperApp/DSP"
	"github.com/opensatelliteproject/SatHelperApp/Demuxer"
	"github.com/opensatelliteproject/SatHelperApp/Display"
	"github.com/opensatelliteproject/SatHelperApp/Frontend"
	"github.com/opensatelliteproject/SatHelperApp/ImageProcessor"
	SLog "github.com/opensatelliteproject/SatHelperApp/Logger"
	SatHelper "github.com/opensatelliteproject/libsathelper"
	"os"
	"os/exec"
	"os/signal"
	"runtime"
	"strings"
	"syscall"
	"time"
)

func startServer(pctx *kong.Context) error {
	SLog.StartLog()
	defer SLog.EndLog()

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

	loadConfig()

	if !isConfigLoaded() {
		SLog.Error("SatUI is not configured. Please run the UI or create the SatHelperApp.cfg manually")
		return fmt.Errorf("SatUI is not configured. Please run the UI or create the SatHelperApp.cfg manually")
	}

	DSP.StatisticsServer = Demuxer.NewTCPServer("", DSP.CurrentConfig.Base.StatisticsPort)

	DSP.StatisticsServer.Start()
	defer DSP.StatisticsServer.Stop()

	DSP.ConstellationServer = Demuxer.NewUDPServer("localhost", 9000)
	DSP.ConstellationServer.Start()
	defer DSP.ConstellationServer.Stop()

	//if DSP.CurrentConfig.RPC.Enable {
	//	addr := fmt.Sprintf("%s:%d", DSP.CurrentConfig.RPC.ListenAddr, DSP.CurrentConfig.RPC.ListenPort)
	//	SLog.Info("Enabling gRPC at %s", addr)
	//	rpc := RPC.MakeRPCServer(rpcSource)
	//	defer rpc.Stop()
	//	err := rpc.Listen(addr)
	//	if err != nil {
	//		SLog.Error("Error starting gRPC: %s", err)
	//	}
	//}

	//if DSP.CurrentConfig.Prometheus.Enable {
	//	addr := fmt.Sprintf("%s:%d", DSP.CurrentConfig.Prometheus.ListenAddr, DSP.CurrentConfig.Prometheus.ListenPort)
	//	SLog.Info("Enabling Prometheus Metrics at %s", addr)
	//	metrics.EnablePrometheus()
	//
	//	srv := &http.Server{}
	//	srv.Handler = metrics.GetHandler()
	//	srv.Addr = addr
	//
	//	lis, err := net.Listen("tcp", addr)
	//	if err != nil {
	//		SLog.Error("Error starting prometheus server: %s", err)
	//		return
	//	}
	//
	//	go func() {
	//		_ = srv.Serve(lis)
	//	}()
	//}

	ImageProcessor.SetPurgeFiles(DSP.CurrentConfig.DirectDemuxer.PurgeFilesAfterProcess)

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
			return fmt.Errorf("error initializing device")
		}
		defer DSP.Device.Destroy()

		DSP.Device.SetGain1(DSP.CurrentConfig.LimeSource.LNAGain)
		DSP.Device.SetAntenna(DSP.CurrentConfig.LimeSource.Antenna)

		SLog.Info(aurora.Cyan("	LNA Gain: %d").String(), aurora.Bold(aurora.Green(DSP.CurrentConfig.LimeSource.LNAGain)))
		SLog.Info(aurora.Cyan("	Antenna: %s").String(), aurora.Bold(aurora.Green(DSP.CurrentConfig.LimeSource.Antenna)))
	case "rtlsdr":
		SLog.Info(aurora.Cyan("RTLSDR Frontend selected.").String())
		DSP.Device = Frontend.NewRTLSDRFrontend()

		if !DSP.Device.Init() {
			SLog.Error("Error initializing device")
			return fmt.Errorf("error initializing device")
		}
		defer DSP.Device.Destroy()

		DSP.Device.SetGain1(DSP.CurrentConfig.RtlsdrSource.LNAGain)
		DSP.Device.SetGain2(DSP.CurrentConfig.RtlsdrSource.VGAGain)
		DSP.Device.SetGain3(DSP.CurrentConfig.RtlsdrSource.MixerGain)
		DSP.Device.SetBiasT(DSP.CurrentConfig.RtlsdrSource.BiasTEnabled)
	case "airspy":
		SLog.Info(aurora.Cyan("Airspy Frontend selected.").String())
		Frontend.AirspyInitialize()
		defer Frontend.AirspyDeinitialize()
		DSP.Device = Frontend.NewAirspyFrontend()

		if !DSP.Device.Init() {
			SLog.Error("Error initializing device")
			return fmt.Errorf("error initializing device")
		}
		defer DSP.Device.Destroy()

		DSP.Device.SetGain1(DSP.CurrentConfig.AirspySource.LNAGain)
		DSP.Device.SetGain2(DSP.CurrentConfig.AirspySource.VGAGain)
		DSP.Device.SetGain3(DSP.CurrentConfig.AirspySource.MixerGain)
		DSP.Device.SetBiasT(DSP.CurrentConfig.AirspySource.BiasTEnabled)
	case "spyserver":
		SLog.Info(aurora.Cyan("Spyserver Frontend Selected. Target: %s:%d").String(), aurora.Bold(DSP.CurrentConfig.SpyserverSource.Hostname), aurora.Bold(DSP.CurrentConfig.SpyserverSource.Port))
		DSP.Device = Frontend.NewSpyserverFrontend(DSP.CurrentConfig.SpyserverSource.Hostname, DSP.CurrentConfig.SpyserverSource.Port)
		if !DSP.Device.Init() {
			SLog.Error("Error initializing device")
			return fmt.Errorf("error initializing device")
		}
		defer DSP.Device.Destroy()
		DSP.Device.SetGain1(DSP.CurrentConfig.SpyserverSource.Gain)
	default:
		SLog.Error(aurora.Red("Device %s is not currently supported.").String(), aurora.Bold(DSP.CurrentConfig.Base.DeviceType))
		return fmt.Errorf("device %s is not currently supported", DSP.CurrentConfig.Base.DeviceType)
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
		return fmt.Errorf("unknown Demuxer Type %s", DSP.CurrentConfig.Base.DemuxerType)
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

	stopFunc := func() {
		SLog.Warn(aurora.Bold("Got close handler").String())
		DSP.SetRunning(false)
		time.Sleep(1 * time.Second)
	}

	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		stopFunc()
	}()

	for DSP.IsRunning() {
		time.Sleep(time.Millisecond * 100)
	}

	return nil
}

var clear = map[string]func(){
	"linux": func() {
		cmd := exec.Command("clear") //Linux example, its tested
		cmd.Stdout = os.Stdout
		err := cmd.Run()
		if err != nil {
			panic(err)
		}
	},
	"windows": func() {
		cmd := exec.Command("cmd", "/c", "cls") //Windows example, its tested
		cmd.Stdout = os.Stdout
		err := cmd.Run()
		if err != nil {
			panic(err)
		}
	},
}

func callClear() {
	value, ok := clear[runtime.GOOS] //runtime.GOOS -> linux, windows, darwin etc.
	if ok {                          //if we defined a clear func for that platform:
		value() //we execute it
	} else { //unsupported platform
		clear["linux"]()
	}
}
