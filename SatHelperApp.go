package main

import (
	"github.com/OpenSatelliteProject/libsathelper"
	"log"
	"github.com/OpenSatelliteProject/SatHelperApp/Frontend"
	"github.com/foize/go.fifo"
	ui "github.com/airking05/termui"
	"github.com/logrusorgru/aurora"
	"github.com/OpenSatelliteProject/SatHelperApp/Display"
)

func main() {
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

	err := ui.Init()
	if err != nil {
		panic(err)
	}
	defer ui.Close()

	LoadConfig()

	samplesFifo = fifo.NewQueue()

	switch CurrentConfig.Base.Mode {
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

	switch CurrentConfig.Base.DeviceType {
		case "cfile":
			log.Printf(aurora.Cyan("CFile Frontend selected. File Name: %s").String(), aurora.Bold(aurora.Green(CurrentConfig.CFileSource.Filename)))
			device = Frontend.NewCFileFrontend(CurrentConfig.CFileSource.Filename)
			device.SetSampleRate(CurrentConfig.Source.SampleRate)
			device.SetCenterFrequency(CurrentConfig.Source.Frequency)
			break
		default:
			log.Fatalf(aurora.Red("Device %s is not currently supported.").String(), aurora.Bold(CurrentConfig.Base.DeviceType))
		break
	}

	device.SetSamplesAvailableCallback(newSamplesCallback)

	initDSP()
	initDecoder()
	Display.InitDisplay()

	log.Println(aurora.Cyan("Starting Source"))
	device.Start()

	log.Println(aurora.Cyan("Starting Main loop"))

	go symbolProcessLoop()
	go decoderLoop()


	//log.Println(Cyan("Connecting to localhost:5000"))
	//cn, err := net.Dial("tcp", "127.0.0.1:5000")
	//
	//conn = cn
	//
	//if err != nil {
	//	log.Fatal(err)
	//}

	// Display.Render()

	running = true

	stopFunc := func(ui.Event) {
		log.Println(aurora.Bold(aurora.Red("Got close handler.")))
		running = false
		ui.StopLoop()
	}

	ui.Handle("/sys/kbd/q", stopFunc)
	ui.Handle("/sys/kbd/C-c", stopFunc)
	ui.Handle("/timer/10ms", func (e ui.Event) {
		stat := GetStats()
		Display.UpdateSignalQuality(stat.SignalQuality)
		Display.UpdateLockedState(stat.FrameLock == 1)
		Display.UpdateChannelData(stat.ReceivedPacketsPerChannel)
		Display.UpdateReedSolomon(stat.RsErrors)
		Display.UpdateSyncWord(stat.SyncWord)
		Display.UpdateSCVCID(stat.SCID, stat.VCID)
		Display.Render()
	})

	CallClear()

	ui.Loop()

	log.Println(aurora.Red("Stopping Source"))
	device.Stop()
}
