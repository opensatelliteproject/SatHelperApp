package Display

import (
	ui "github.com/airking05/termui"
	"time"
	"github.com/nsf/termbox-go"
	"strconv"
)

type DisplayObjects struct {
	head *ui.Par
	signalLocked *ui.Par
	signalQuality *ui.Gauge
	channelData *ui.BarChart
	rsErrors *ui.BarChart
}

var state = struct {
	signalQuality int
	signalLocked bool
	channelPackets [256]int64
	rsErrors [4] int
	displayObjects DisplayObjects
}{
	signalQuality: 0,
	rsErrors: [4]int {0,0,0,0},
	signalLocked: false,
}

var colorBar []uint32


func AlignCenter(p *ui.Par) {
	totalWidth := p.Width
	originalText := p.Text
	padTotal := totalWidth - len(originalText)
	p.PaddingLeft = padTotal / 2 - 1
}

func InitDisplay() {
	InitLut()
	termbox.SetOutputMode(termbox.Output256)
	// region Color Bar
	colorBar = make([]uint32, 101)
	for i:=0;i<100;i++ {
		colorBar[i] = uint32(i)
	}
	// endregion
	// region HEAD
	head := ui.NewPar("SatHelper Application")
	head.TextFgColor = ui.ColorWhite
	head.BorderFg = ui.ColorCyan
	head.Height = 3
	head.Text = "SatHelper Application"
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
	rsErrors.Data = []int{0,0,0,0}
	rsErrors.Width = 50
	rsErrors.Height = 10
	rsErrors.BarWidth = 8
	rsErrors.DataLabels = []string{"0", "1", "2", "3"}
	rsErrors.TextColor = ui.ColorGreen | ui.AttrBold
	// endregion


	state.displayObjects.head = head
	state.displayObjects.signalLocked = signalLocked
	state.displayObjects.signalQuality = signalQuality
	state.displayObjects.channelData = channelData

	ui.Body.AddRows(
		ui.NewRow(
			ui.NewCol(12, 0, head),
		),
		ui.NewRow(
			ui.NewCol(3, 0, signalLocked),
			ui.NewCol(4,0,signalQuality),
			ui.NewCol(7, 0, rsErrors),
		),
		ui.NewRow(
			ui.NewCol(12, 0, channelData),
		),
	)
	e := ui.NewTimerCh(10 * time.Millisecond)
	ui.Merge("timer10ms", e)
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
	state.displayObjects.signalQuality.BarColor = ui.Attribute( GetXTermColorHVal(signalQualityColor) + 1  )
	// endregion
	// region Channel Packets
	data := make([]int, 0)
	label := make([]string, 0)

	for i:=0; i< 256; i++ {
		if state.channelPackets[i] != 0 {
			label = append(label, strconv.FormatInt(int64(i), 10))
			data = append(data, int(state.channelPackets[i])) // Overflow Alert!! TODO
		}
	}
	state.displayObjects.channelData.DataLabels = label
	state.displayObjects.channelData.Data = data
	// channelData
	// endregion
	ui.Body.Align()
	AlignCenter(state.displayObjects.signalLocked)
	AlignCenter(state.displayObjects.head)
}

func Render() {
	updateComponents()
	ui.Clear()
	ui.Render(ui.Body)
}

func UpdateLockedState(lck bool) {
	state.signalLocked = lck
}

func UpdateSignalQuality(q uint8) {
	state.signalQuality = int(q)
}

func UpdateChannelData(d [256]int64) {
	state.channelPackets = d
}