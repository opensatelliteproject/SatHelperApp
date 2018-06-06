package Frontend

import (
	"fmt"
	"log"
	"os"
	"time"
	"bufio"
	"encoding/binary"
	. "github.com/logrusorgru/aurora"
)

const CFileFrontendBufferSize = 65535

// region Struct Definition
type CFileFrontend struct {
	running bool
	filename string
	cb SamplesCallback
	sampleRate uint32
	centerFrequency uint32
	sampleBuffer []complex64
	fileHandler *os.File
	t0 time.Time
}
// endregion
// region Constructor
func NewCFileFrontend(filename string) *CFileFrontend {
	return &CFileFrontend{
		filename: filename,
		running: false,
		sampleRate: 0,
		centerFrequency: 0,
		cb: nil,
	}
}
// endregion
// region Getters
func (f *CFileFrontend) GetName() string {
	return fmt.Sprintf("CFileFrontend (%s)", f.filename)
}
func (f *CFileFrontend) GetShortName() string {
	return "CFileFrontend"
}
func (f *CFileFrontend) GetAvailableSampleRates() []uint32 {
	return make([]uint32, 0)
}
func (f *CFileFrontend) GetCenterFrequency() uint32 {
	return f.centerFrequency
}
func (f *CFileFrontend)  GetSampleRate() uint32 {
	return f.sampleRate
}
// endregion
// region Setters
func (f *CFileFrontend) SetSamplesAvailableCallback(cb SamplesCallback) {
	f.cb = cb
}
func (f *CFileFrontend) SetSampleRate(sampleRate uint32) uint32 {
	f.sampleRate = sampleRate
	return sampleRate
}
func (f *CFileFrontend) SetCenterFrequency(centerFrequency uint32) uint32 {
	f.centerFrequency = centerFrequency
	return centerFrequency
}
// endregion
// region Commands
func (f *CFileFrontend) Init() bool {
	return true
}
func (f *CFileFrontend) Destroy() {}
func (f *CFileFrontend) Start() {
	if f.running {
		log.Println(Red("CFileFrontend is already running."))
		return
	}
	f.running = true
	go func(frontend *CFileFrontend) {
		log.Println(Green("CFileFrontend Routine started"))
		f, err := os.Open(f.filename)

		var period = CFileFrontendBufferSize / float32(frontend.sampleRate)

		if err != nil {
			log.Printf(Red("Error opening file %s: %s").String(), Bold(frontend.filename), Bold(err))
			frontend.running = false
			return
		}
		defer frontend.fileHandler.Close()

		frontend.fileHandler = f
		frontend.t0 = time.Now()
		frontend.sampleBuffer = make([]complex64, CFileFrontendBufferSize)

		var reader = bufio.NewReader(f)

		for frontend.running {
			if float32(time.Now().Sub(frontend.t0).Seconds()) >= period {
				err := binary.Read(reader, binary.LittleEndian, frontend.sampleBuffer)
				if err != nil {
					log.Println(Red("Error reading input CFile: %s").String(), Bold(err))
					frontend.running = false
					break
				}
				if frontend.cb != nil {
					var cbData = SampleCallbackData{
						SampleType:   FrontendSampletypeFloatiq,
						NumSamples:   len(frontend.sampleBuffer),
						ComplexArray: frontend.sampleBuffer,
					}
					frontend.cb(cbData)
				}
				frontend.t0 = time.Now()
			}
			time.Sleep(time.Duration((period / 100) * float32(time.Second)))
		}
		log.Println(Red("CFileFrontend Routine ended"))
	}(f)
}

func (f *CFileFrontend) Stop() {
	if ! f.running {
		log.Println(Red("CFileFrontend is not running"))
		return
	}
	f.running = false
}

func (f *CFileFrontend) SetAGC(bool) {}
func (f *CFileFrontend) SetLNAGain(uint8) {}
func (f *CFileFrontend) SetVGAGain(uint8) {}
func (f *CFileFrontend) SetMixerGain(uint8) {}
func (f *CFileFrontend) SetBiasT(bool) {}
// endregion
