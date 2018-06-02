package main

import (
	"github.com/OpenSatelliteProject/libsathelper"
	. "github.com/logrusorgru/aurora"
	"log"
	"github.com/OpenSatelliteProject/SatHelperApp/Frontend"
	"time"
	"github.com/foize/go.fifo"
	"os"
	"os/signal"
	"net"
)

func main() {
	log.Printf("%s %s (%s) - %s %s\n",
		Green(Bold("SatHelperApp")),
		Bold(GetVersion()),
		Bold(GetRevision()),
		Bold(GetCompilationDate()),
		Bold(GetCompilationTime()),
	)
	log.Printf("%s %s (%s) - %s %s\n",
		Green(Bold("libSatHelper")),
		Bold(SatHelper.InfoGetVersion()),
		Bold(SatHelper.InfoGetGitSHA1()),
		Bold(SatHelper.InfoGetCompilationDate()),
		Bold(SatHelper.InfoGetCompilationTime()),
	)

	LoadConfig()

	samplesFifo = fifo.NewQueue()

	switch CurrentConfig.Base.Mode {
		case "lrit":
			log.Println(Cyan("Selected LRIT mode. Ignoring parameters from config file."))
			SetLRITMode()
		break
		case "hrit":
			log.Println(Cyan("Selected HRIT mode. Ignoring parameters from config file."))
			SetHRITMode()
		break
		default:
			log.Println(Gray("No valid mode selected. Using config file parameters."))
	}

	switch CurrentConfig.Base.DeviceType {
		case "cfile":
			log.Printf(Cyan("CFile Frontend selected. File Name: %s").String(), Bold(Green(CurrentConfig.CFileSource.Filename)))
			device = Frontend.NewCFileFrontend(CurrentConfig.CFileSource.Filename)
			device.SetSampleRate(CurrentConfig.Source.SampleRate)
			device.SetCenterFrequency(CurrentConfig.Source.Frequency)
			break
		default:
			log.Fatalf(Red("Device %s is not currently supported.").String(), Bold(CurrentConfig.Base.DeviceType))
		break
	}

	device.SetSamplesAvailableCallback(newSamplesCallback)

	initDSP()

	log.Println(Cyan("Connecting to localhost:5000"))
	cn, err := net.Dial("tcp", "127.0.0.1:5000")

	conn = cn

	if err != nil {
		log.Fatal(err)
	}

	running = true

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func(){
		for range c {
			log.Println(Bold(Red("Got Ctrl+C! Closing!!!!")))
			running = false
		}
	}()

	log.Println(Cyan("Starting Source"))
	device.Start()

	log.Println(Cyan("Starting Main loop"))

	// go symbolLoopFunc()

	for running {
		processSamples()
		time.Sleep(time.Microsecond)
	}

	log.Println(Red("Stopping Source"))
	device.Stop()
}
