package Display

import (
	"fmt"
	ui "github.com/airking05/termui"
	"github.com/nsf/termbox-go"
	"github.com/opensatelliteproject/SatHelperApp"
	"github.com/opensatelliteproject/SatHelperApp/Logger"
	"log"
	"regexp"
	"strconv"
	"time"
)

const MaxConsoleLines = 10

type Objects struct {
	head                 *ui.Par
	signalLocked         *ui.Par
	signalQuality        *ui.Gauge
	channelData          *ui.BarChart
	rsErrors             *ui.BarChart
	syncWord             *ui.Par
	scid                 *ui.Par
	vcid                 *ui.Par
	decoderFifoUsage     *ui.Gauge
	demodulatorFifoUsage *ui.Gauge
	viterbiErrors        *ui.Par
	syncCorrelation      *ui.Par
	phaseCorrection      *ui.Par
	mode                 *ui.Par
	centerFrequency      *ui.Par
	demuxer              *ui.Par
	device               *ui.Par
	console              *ui.List
}

var state = struct {
	signalQuality        int
	signalLocked         bool
	channelPackets       [256]int64
	rsErrors             [4]int32
	displayObjects       Objects
	syncWord             [4]byte
	scid                 uint8
	vcid                 uint8
	decoderFifoUsage     uint8
	demodulatorFifoUsage uint8
	viterbiErrors        uint
	frameSize            uint
	phaseCorrection      uint16
	syncCorrelation      uint8
	centerFreq           uint32
	mode                 string
	demuxer              string
	device               string
	consoleLines         []string
	cw                   *ConsoleWritter
}{
	signalQuality:        0,
	rsErrors:             [4]int32{0, 0, 0, 0},
	signalLocked:         false,
	syncWord:             [4]byte{0, 0, 0, 0},
	scid:                 0,
	vcid:                 0,
	decoderFifoUsage:     0,
	demodulatorFifoUsage: 0,
	phaseCorrection:      0,
	syncCorrelation:      0,
	centerFreq:           0,
	mode:                 "Not Selected",
	device:               "None",
	demuxer:              "None",
	consoleLines:         []string{},
}

var colorBar []uint32

func AlignCenter(p *ui.Par) {
	totalWidth := p.Width
	originalText := p.Text
	padTotal := totalWidth - len(originalText)
	p.PaddingLeft = padTotal/2 - 1
}

func InitDisplay() {
	InitLut()
	termbox.SetOutputMode(termbox.Output256)
	// region Color Bar
	colorBar = make([]uint32, 256)
	for i := 0; i < 100; i++ {
		colorBar[i] = uint32(i)
	}
	for i := 100; i < 256; i++ { // Invalid values are no signal.
		colorBar[i] = 0
	}
	// endregion
	// region HEAD
	headStr := fmt.Sprintf("SatHelperApp - %s.%s", SatHelperApp.GetVersion(), SatHelperApp.GetRevision())
	head := ui.NewPar(headStr)
	head.TextFgColor = ui.ColorWhite
	head.BorderFg = ui.ColorCyan
	head.Height = 3
	head.Text = headStr
	// endregion
	// region Locked Bar
	signalLocked := ui.NewPar("NOT LOCKED")
	signalLocked.Bg = ui.ColorRed
	signalLocked.TextBgColor = ui.ColorRed
	signalLocked.BorderFg = ui.ColorRed
	signalLocked.TextFgColor = ui.ColorWhite
	signalLocked.Text = "NOT LOCKED"
	signalLocked.Height = 3
	signalLocked.Align()
	// endregion
	// region Signal Quality
	signalQuality := ui.NewGauge()
	signalQuality.Percent = 0
	signalQuality.Height = 3
	signalQuality.BorderLabel = "Signal Quality"
	signalQuality.BarColor = ui.ColorRed
	signalQuality.BorderFg = ui.ColorWhite
	signalQuality.PercentColor = ui.ColorBlue
	signalQuality.PercentColorHighlighted = ui.ColorBlue
	signalQuality.BorderLabelFg = ui.ColorCyan
	// endregion
	// region Channel Data
	channelData := ui.NewBarChart()
	channelData.BorderLabel = "Channel Packets"
	channelData.Data = []int{}
	channelData.Width = 50
	channelData.Height = 10
	channelData.BarWidth = 8
	channelData.DataLabels = []string{}
	channelData.TextColor = ui.ColorGreen | ui.AttrBold
	//channelData.BarColor = ui.ColorRed
	//channelData.NumColor = ui.ColorYellow
	// endregion
	// region RS Errors
	rsErrors := ui.NewBarChart()
	rsErrors.BorderLabel = "ReedSolomon Errors"
	rsErrors.Data = []int{0, 0, 0, 0}
	rsErrors.Width = 50
	rsErrors.Height = 10
	rsErrors.BarWidth = 8
	rsErrors.DataLabels = []string{"0", "1", "2", "3"}
	rsErrors.TextColor = ui.ColorGreen | ui.AttrBold
	// endregion
	// region Sync Word
	syncWord := ui.NewPar("Sync Word")
	syncWord.BorderLabel = "Sync Word"
	syncWord.TextFgColor = ui.ColorWhite
	syncWord.Text = "00000000"
	syncWord.Height = 3
	syncWord.Align()
	// endregion
	// region VCID / SCID
	vcid := ui.NewPar("VCID")
	vcid.BorderLabel = "VCID"
	vcid.TextFgColor = ui.ColorWhite
	vcid.Text = "0"
	vcid.Height = 3
	vcid.Align()
	scid := ui.NewPar("SCID")
	scid.BorderLabel = "SCID"
	scid.TextFgColor = ui.ColorWhite
	scid.Text = "0"
	scid.Height = 3
	scid.Align()
	// endregion
	// region Decoder Fifo Usage
	decoderFifoUsage := ui.NewGauge()
	decoderFifoUsage.Percent = 0
	decoderFifoUsage.Height = 3
	decoderFifoUsage.BorderLabel = "Decoder FIFO"
	decoderFifoUsage.BarColor = ui.ColorRed
	decoderFifoUsage.BorderFg = ui.ColorWhite
	decoderFifoUsage.PercentColor = ui.ColorBlue
	decoderFifoUsage.PercentColorHighlighted = ui.ColorBlue
	decoderFifoUsage.BorderLabelFg = ui.ColorCyan
	// endregion
	// region Demodulator Fifo Usage
	demodulatorFifoUsage := ui.NewGauge()
	demodulatorFifoUsage.Percent = 0
	demodulatorFifoUsage.Height = 3
	demodulatorFifoUsage.BorderLabel = "Demodulator FIFO"
	demodulatorFifoUsage.BarColor = ui.ColorRed
	demodulatorFifoUsage.BorderFg = ui.ColorWhite
	demodulatorFifoUsage.PercentColor = ui.ColorBlue
	demodulatorFifoUsage.PercentColorHighlighted = ui.ColorBlue
	demodulatorFifoUsage.BorderLabelFg = ui.ColorCyan
	// endregion
	// region Viterbi Errors
	viterbiErrors := ui.NewPar("Viterbi Err")
	viterbiErrors.BorderLabel = "Viterbi Err"
	viterbiErrors.TextFgColor = ui.ColorWhite
	viterbiErrors.Text = "   0 /    0 bits"
	viterbiErrors.Height = 3
	viterbiErrors.Align()
	// endregion
	// region Phase Correction
	phaseCorrection := ui.NewPar("Phase Corr")
	phaseCorrection.BorderLabel = "Phase Corr"
	phaseCorrection.TextFgColor = ui.ColorWhite
	phaseCorrection.Text = "  0 deg"
	phaseCorrection.Height = 3
	phaseCorrection.Align()
	// endregion
	// region Sync Correlation
	syncCorrelation := ui.NewPar("Sync Corr")
	syncCorrelation.BorderLabel = "Sync Corr"
	syncCorrelation.TextFgColor = ui.ColorWhite
	syncCorrelation.Text = " 0"
	syncCorrelation.Height = 3
	syncCorrelation.Align()
	// endregion
	// region Center Frequency
	centerFrequency := ui.NewPar("Center Frequency")
	centerFrequency.BorderLabel = "Center Freq."
	centerFrequency.TextFgColor = ui.ColorWhite
	centerFrequency.Text = "0 MHz"
	centerFrequency.Height = 3
	centerFrequency.Align()
	// endregion
	// region Mode
	mode := ui.NewPar("Mode")
	mode.BorderLabel = "Mode"
	mode.TextFgColor = ui.ColorWhite
	mode.Text = "None"
	mode.Height = 3
	mode.Align()
	// endregion
	// region Device
	device := ui.NewPar("Device")
	device.BorderLabel = "Device"
	device.TextFgColor = ui.ColorWhite
	device.Text = "None"
	device.Height = 3
	device.Align()
	// endregion
	// region Demuxer
	demuxer := ui.NewPar("Demuxer")
	demuxer.BorderLabel = "Demuxer"
	demuxer.TextFgColor = ui.ColorWhite
	demuxer.Text = "None"
	demuxer.Height = 3
	demuxer.Align()
	// endregion
	// region Console
	console := ui.NewList()
	console.Overflow = "wrap"
	console.Items = []string{}
	console.BorderLabel = "Console"
	console.Height = MaxConsoleLines
	// endregion
	// region Save Objects
	state.displayObjects.head = head
	state.displayObjects.signalLocked = signalLocked
	state.displayObjects.signalQuality = signalQuality
	state.displayObjects.channelData = channelData
	state.displayObjects.rsErrors = rsErrors
	state.displayObjects.syncWord = syncWord
	state.displayObjects.vcid = vcid
	state.displayObjects.scid = scid
	state.displayObjects.decoderFifoUsage = decoderFifoUsage
	state.displayObjects.demodulatorFifoUsage = demodulatorFifoUsage
	state.displayObjects.viterbiErrors = viterbiErrors
	state.displayObjects.phaseCorrection = phaseCorrection
	state.displayObjects.syncCorrelation = syncCorrelation
	state.displayObjects.centerFrequency = centerFrequency
	state.displayObjects.mode = mode
	state.displayObjects.device = device
	state.displayObjects.demuxer = demuxer
	state.displayObjects.console = console
	// endregion
	// region Configure Body
	ui.Body.AddRows(
		ui.NewRow(
			ui.NewCol(12, 0, head),
		),
		ui.NewRow(
			ui.NewCol(3, 0, signalLocked),
			ui.NewCol(3, 0, signalQuality),
			ui.NewCol(2, 0, syncWord),
			ui.NewCol(2, 0, vcid),
			ui.NewCol(2, 0, scid),
		),
		ui.NewRow(
			ui.NewCol(3, 0, decoderFifoUsage),
			ui.NewCol(3, 0, demodulatorFifoUsage),
			ui.NewCol(2, 0, viterbiErrors),
			ui.NewCol(2, 0, syncCorrelation),
			ui.NewCol(2, 0, phaseCorrection),
		),
		ui.NewRow(
			ui.NewCol(3, 0, centerFrequency),
			ui.NewCol(2, 0, mode),
			ui.NewCol(3, 0, demuxer),
			ui.NewCol(4, 0, device),
		),
		ui.NewRow(
			ui.NewCol(6, 0, rsErrors),
			ui.NewCol(6, 0, channelData),
		),
		ui.NewRow(
			ui.NewCol(12, 0, console),
		),
	)
	// endregion
	// region Create Timers
	e := ui.NewTimerCh(100 * time.Millisecond)
	ui.Merge("refresh", e)
	// endregion
	// region Callbacks
	state.cw = NewConsoleWritter(func(data string) (int, error) {
		AddConsoleLine(data[:len(data)-1])
		return len(data) - 1, nil
	})
	log.SetOutput(state.cw)
	SLog.SetTermUiDisplay(true)
	// endregion
}

func updateComponents() {
	// region Signal Locked
	if state.signalLocked {
		state.displayObjects.signalLocked.PaddingLeft = 0
		state.displayObjects.signalLocked.TextBgColor = ui.ColorGreen
		state.displayObjects.signalLocked.BorderFg = ui.ColorGreen
		state.displayObjects.signalLocked.TextFgColor = ui.ColorWhite
		state.displayObjects.signalLocked.Bg = ui.ColorGreen
		state.displayObjects.signalLocked.Text = "LOCKED"
	} else {
		state.displayObjects.signalLocked.PaddingLeft = 0
		state.displayObjects.signalLocked.TextBgColor = ui.ColorRed
		state.displayObjects.signalLocked.BorderFg = ui.ColorRed
		state.displayObjects.signalLocked.TextFgColor = ui.ColorWhite
		state.displayObjects.signalLocked.Bg = ui.ColorRed
		state.displayObjects.signalLocked.Text = "NOT LOCKED"
	}
	state.displayObjects.signalLocked.Align()
	state.displayObjects.signalQuality.Percent = state.signalQuality
	// endregion
	// region Signal Quality
	signalQualityColor := colorBar[state.signalQuality]
	state.displayObjects.signalQuality.BarColor = ui.Attribute(GetXTermColorHVal(signalQualityColor) + 1)
	// endregion
	// region Channel Packets
	data := make([]int, 0)
	label := make([]string, 0)

	for i := 0; i < 256; i++ {
		if state.channelPackets[i] != 0 {
			label = append(label, strconv.FormatInt(int64(i), 10))
			data = append(data, int(state.channelPackets[i])) // Overflow Alert!! TODO
		}
	}
	state.displayObjects.channelData.DataLabels = label
	state.displayObjects.channelData.Data = data
	// endregion
	// region RS Errors
	for i := 0; i < 4; i++ {
		state.displayObjects.rsErrors.Data[i] = int(state.rsErrors[i])
	}
	// endregion
	// region Sync Word
	state.displayObjects.syncWord.Text = ""
	for i := 0; i < 4; i++ {
		state.displayObjects.syncWord.Text += fmt.Sprintf("%02X", state.syncWord[i])
	}
	state.displayObjects.syncWord.Align()
	// endregion
	// region SCID / VCID
	state.displayObjects.scid.Text = strconv.FormatUint(uint64(state.scid), 10)
	state.displayObjects.vcid.Text = strconv.FormatUint(uint64(state.vcid), 10)
	// endregion
	// region FIFO
	state.displayObjects.decoderFifoUsage.Percent = int(state.decoderFifoUsage)
	state.displayObjects.demodulatorFifoUsage.Percent = int(state.demodulatorFifoUsage)
	// endregion
	// region Viterbi Errors
	state.displayObjects.viterbiErrors.Text = fmt.Sprintf("%4d / %4d bits", state.viterbiErrors, state.frameSize)
	// endregion
	// region Sync Correlation
	state.displayObjects.syncCorrelation.Text = fmt.Sprintf("%2d", state.syncCorrelation)
	// endregion
	// region Phase Correction
	state.displayObjects.phaseCorrection.Text = fmt.Sprintf("%3d deg", state.phaseCorrection)
	// endregion
	// region Center Frequency
	state.displayObjects.centerFrequency.Text = fmt.Sprintf("%d Hz", state.centerFreq)
	// endregion
	// region Mode
	state.displayObjects.mode.Text = state.mode
	// endregion
	// region Device
	state.displayObjects.device.Text = state.device
	// endregion
	// region Device
	state.displayObjects.demuxer.Text = state.demuxer
	// endregion
	// region Console
	state.displayObjects.console.Items = state.consoleLines
	// endregion
	// region UI Alignments
	ui.Body.Align()
	AlignCenter(state.displayObjects.scid)
	AlignCenter(state.displayObjects.vcid)
	AlignCenter(state.displayObjects.signalLocked)
	AlignCenter(state.displayObjects.head)
	AlignCenter(state.displayObjects.syncWord)
	AlignCenter(state.displayObjects.viterbiErrors)
	AlignCenter(state.displayObjects.phaseCorrection)
	AlignCenter(state.displayObjects.syncCorrelation)
	AlignCenter(state.displayObjects.mode)
	AlignCenter(state.displayObjects.centerFrequency)
	AlignCenter(state.displayObjects.demuxer)
	AlignCenter(state.displayObjects.device)
	// endregion
}

func Render() {
	updateComponents()
	ui.Clear()
	ui.Render(ui.Body)
}

// region Update Functions
func UpdateLockedState(lck bool) {
	state.signalLocked = lck
}

func UpdateSignalQuality(q uint8) {
	state.signalQuality = int(q)
}

func UpdateChannelData(d [256]int64) {
	state.channelPackets = d
}

func UpdateReedSolomon(d [4]int32) {
	state.rsErrors = d
}

func UpdateSyncWord(d [4]byte) {
	state.syncWord = d
}

func UpdateSCVCID(scid byte, vcid byte) {
	state.vcid = vcid
	state.scid = scid
}

func UpdateDecoderFifoUsage(percent uint8) {
	state.decoderFifoUsage = percent
}

func UpdateDemodulatorFifoUsage(percent uint8) {
	state.demodulatorFifoUsage = percent
}

func UpdateViterbiErrors(errors uint, frameSize uint) {
	state.viterbiErrors = errors
	state.frameSize = frameSize
}

func UpdateSyncCorrelation(corr uint8) {
	state.syncCorrelation = corr
}

func UpdatePhaseCorr(corr uint8) {
	switch corr {
	case 0:
		state.phaseCorrection = 0
	case 1:
		state.phaseCorrection = 90
	case 2:
		state.phaseCorrection = 180
	case 3:
		state.phaseCorrection = 270
	}
}

func UpdateCenterFrequency(freq uint32) {
	state.centerFreq = freq
}

func UpdateMode(mode string) {
	state.mode = mode
}

func UpdateDevice(device string) {
	state.device = device
}

func UpdateDemuxer(demuxer string) {
	state.demuxer = demuxer
}

var reg, _ = regexp.Compile("\x1b\\[[0-9;]*m")

func AddConsoleLine(line string) {
	line = reg.ReplaceAllString(line, "")

	state.consoleLines = append(state.consoleLines, line)
	if len(state.consoleLines) > MaxConsoleLines {
		state.consoleLines = state.consoleLines[1:]
	}
}

// endregion
