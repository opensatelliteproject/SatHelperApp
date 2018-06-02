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

var samplesFifo *fifo.Queue
var buffer0 []complex64
var buffer1 []complex64
var running = false

var decimator SatHelper.FirFilter
var agc SatHelper.AGC
var rrcFilter SatHelper.FirFilter
var costasLoop SatHelper.CostasLoop
var clockRecovery SatHelper.ClockRecovery

var conn net.Conn

func checkAndResizeBuffers(length int) {
	if len(buffer0) < length {
		buffer0 = make([]complex64, length)
	}
	if len(buffer1) < length {
		buffer1 = make([]complex64, length)
	}
}

func swapBuffers(a **complex64, b **complex64) {
	c := *b
	*b = *a
	*a = c
}

func processSamples() {
	if samplesFifo.Len() <= 64 * 1024{
		return
	}

	length := samplesFifo.Len()
	checkAndResizeBuffers(length)

	for i := 0; i < length; i++ {
		buffer0[i] = samplesFifo.Next().(complex64)
	}

	ba := &buffer0[0]
	bb := &buffer1[0]

	if CurrentConfig.Base.Decimation > 1 {
		length /= int(CurrentConfig.Base.Decimation)
		decimator.Work(ba, bb, length)
		swapBuffers(&ba, &bb)
	}

	agc.Work(ba, bb, length)
	swapBuffers(&ba, &bb)

	rrcFilter.Work(ba, bb, length)
	swapBuffers(&ba, &bb)

	costasLoop.Work(ba, bb, length)
	swapBuffers(&ba, &bb)

	symbols := clockRecovery.Work(ba, bb, length)
	swapBuffers(&ba, &bb)

	sendbuffer := make([]byte, symbols)

	var ob *[]complex64

	if ba == &buffer0[0] {
		ob = &buffer0
	} else {
		ob = &buffer1
	}

	for i := 0; i < symbols; i++ {
		z := (*ob)[i]
		v := imag(z) * 127
		if v > 127 {
			v = 127
		} else if v < -128 {
			v = -128
		}

		sendbuffer[i] = byte(v)
	}

	_, err := conn.Write(sendbuffer)
	if err != nil {
		log.Printf("Error writting data: %s", err)
	}
	// log.Printf("Got %d symbols!\n", symbols)

}

func main() {
	var device Frontend.BaseFrontend
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
			log.Println("Selected LRIT mode. Ignoring parameters from config file.")
			SetLRITMode()
		break
		case "hrit":
			log.Println("Selected HRIT mode. Ignoring parameters from config file.")
			SetHRITMode()
		break
		default:
			log.Println("No valid mode selected. Using config file parameters.")
	}

	switch CurrentConfig.Base.DeviceType {
		case "cfile":
			log.Printf("CFile Frontend selected. File Name: %s", Green(CurrentConfig.CFileSource.Filename))
			device = Frontend.NewCFileFrontend(CurrentConfig.CFileSource.Filename)
			device.SetSampleRate(CurrentConfig.Source.SampleRate)
			device.SetCenterFrequency(CurrentConfig.Source.Frequency)
			log.Printf("%d ---", device.GetSampleRate())
			break
		default:
			log.Fatalf("Device %s is not currently supported.", CurrentConfig.Base.DeviceType)
		break
	}

	device.SetSamplesAvailableCallback(func(d Frontend.SampleCallbackData) {
		switch d.SampleType {
		case Frontend.FRONTEND_SAMPLETYPE_FLOATIQ: AddToFifoC64(samplesFifo, d.ComplexArray, d.NumSamples); break
		case Frontend.FRONTEND_SAMPLETYPE_S16IQ: AddToFifoS16(samplesFifo, d.Int16Array, d.NumSamples); break
		case Frontend.FRONTEND_SAMPLETYPE_S8IQ: AddToFifoS8(samplesFifo, d.Int8Array, d.NumSamples); break
		}
	})

	circuitSampleRate := float32(device.GetSampleRate()) / float32(CurrentConfig.Base.Decimation)
	sps := circuitSampleRate / float32(CurrentConfig.Base.SymbolRate)

	log.Printf("Samples per Symbol: %f\n", sps)
	log.Printf("Circuit Sample Rate: %f\n", circuitSampleRate)
	log.Printf("Low Pass Decimator Cut Frequency: %f\n", circuitSampleRate / 2)

	rrcTaps := SatHelper.FiltersRRC(1, float64(circuitSampleRate), float64(CurrentConfig.Base.SymbolRate), float64(CurrentConfig.Base.RRCAlpha), RrcTaps)
	decimatorTaps := SatHelper.FiltersLowPass(1, float64(device.GetSampleRate()), float64(circuitSampleRate / 2), 100e3, SatHelper.FFTWindowsHAMMING, 6.76)

	decimator = SatHelper.NewFirFilter(uint(CurrentConfig.Base.Decimation), decimatorTaps)
	agc = SatHelper.NewAGC(AgcRate, AgcReference, AgcGain, AgcMaxGain)
	costasLoop = SatHelper.NewCostasLoop(PllAlpha, LoopOrder)
	clockRecovery = SatHelper.NewClockRecovery(sps, ClockGainOmega, ClockMu, ClockAlpha, ClockOmegaLimit)
	rrcFilter = SatHelper.NewFirFilter(1, rrcTaps)


	log.Printf("Center Frequency: %d MHz", device.GetCenterFrequency())
	log.Printf("Automatic Gain Control: %t\n", CurrentConfig.Base.AGCEnabled)

	if CurrentConfig.Base.AGCEnabled {
		device.SetAGC(true)
	} else {
		device.SetAGC(false)
		// TODO: Gains
	}

	log.Println("Connecting to localhost:5000")
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
			log.Println("Got Ctrl+C! Closing!!!!")
			running = false
		}
	}()

	log.Println("Starting Source")
	device.Start()

	log.Println("Starting Main loop")

	// go symbolLoopFunc()

	for running {
		processSamples()
		time.Sleep(time.Microsecond)
	}

	log.Println("Stopping Source")
	device.Stop()
}
