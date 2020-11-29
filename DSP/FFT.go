package DSP

import (
	"bytes"
	"github.com/llgcode/draw2d"
	"github.com/llgcode/draw2d/draw2dimg"
	"github.com/prometheus/common/log"
	"github.com/racerxdl/segdsp/dsp"
	"github.com/racerxdl/segdsp/dsp/fft"
	"github.com/racerxdl/segdsp/tools"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"math"
	"sync"
	"time"
)

var lastFFT = time.Now()
var fftWindow = dsp.BlackmanHarris(fftWidth, 61)
var fftSamples = make([]complex64, fftWidth)
var lastRealFFT = make([]float32, fftWidth)
var img = image.NewRGBA(image.Rect(0, 0, imgWidth, imgHeight))
var gc = draw2dimg.NewGraphicContext(img)
var imgMtx = sync.Mutex{}

func GetFFTImage() []byte {
	imgMtx.Lock()
	defer imgMtx.Unlock()
	b := bytes.NewBuffer(nil)
	err := jpeg.Encode(b, img, nil)
	if err != nil {
		log.Errorf("Error generating FFT Image: %s", err)
		return nil
	}
	return b.Bytes()
}

func GenerateImage() {
	imgMtx.Lock()
	defer imgMtx.Unlock()
	widthScale := float32(len(lastRealFFT)) / float32(imgWidth)
	gc.SetFontData(draw2d.FontData{
		Name:   "FreeMono",
		Family: draw2d.FontFamilyMono,
		Style:  draw2d.FontStyleNormal,
	})
	gc.SetLineWidth(2)
	gc.SetStrokeColor(color.RGBA{R: 255, A: 255})
	gc.SetFillColor(color.RGBA{R: 0, A: 255})
	gc.Clear()
	gc.SetFontSize(32)
	gc.FillStringAt("FFT", 10, 10)
	draw.Draw(img, img.Bounds(), &image.Uniform{C: color.NRGBA{A: 255}}, image.Point{}, draw.Src)

	var lastX = float32(0)
	var lastY = float32(0)

	for i := 0; i < len(lastRealFFT); i++ {
		var iPos = (i + len(lastRealFFT)/2) % len(lastRealFFT)
		var s = lastRealFFT[iPos]
		var v = ((fftOffset) - s) * float32(fftScale)
		var x = float32(i) / widthScale
		if i != 0 {
			DrawLine(lastX, lastY, x, v, color.NRGBA{R: 0, G: 127, B: 127, A: 255}, img)
		}

		lastX = x
		lastY = v
	}
}
func combine(c1, c2 color.Color) color.Color {
	r, g, b, a := c1.RGBA()
	r2, g2, b2, a2 := c2.RGBA()

	return color.RGBA{
		R: uint8((r + r2) >> 9), // div by 2 followed by ">> 8"  is ">> 9"
		G: uint8((g + g2) >> 9),
		B: uint8((b + b2) >> 9),
		A: uint8((a + a2) >> 9),
	}
}

func DrawLine(x0, y0, x1, y1 float32, color color.Color, img *image.RGBA) {
	// DDA
	_, _, _, a := color.RGBA()
	needsCombine := a != 255 && a != 0
	var dx = x1 - x0
	var dy = y1 - y0
	var steps float32
	if tools.Abs(dx) > tools.Abs(dy) {
		steps = tools.Abs(dx)
	} else {
		steps = tools.Abs(dy)
	}

	var xinc = dx / steps
	var yinc = dy / steps

	var x = x0
	var y = y0
	for i := 0; i < int(steps); i++ {
		if needsCombine {
			var p = img.At(int(x), int(y))
			img.Set(int(x), int(y), combine(p, color))
		} else {
			img.Set(int(x), int(y), color)
		}
		x = x + xinc
		y = y + yinc
	}
}

func computeFFT(samples []complex64) {
	if time.Since(lastFFT) > fftInterval {
		lastFFT = time.Now()
		copy(fftSamples, samples)
		for j := 0; j < fftWidth; j++ {
			var s = fftSamples[j]
			var r = real(s) * float32(fftWindow[j])
			var i = imag(s) * float32(fftWindow[j])
			fftSamples[j] = complex(r, i)
		}
		fftResult := fft.FFT(fftSamples)
		fftReal := make([]float32, len(fftResult))
		for i := 0; i < len(fftResult); i++ {
			// Convert FFT to Power in dB
			var v = tools.ComplexAbsSquared(fftResult[i]) * (1.0 / float32(Device.GetSampleRate()))
			fftReal[i] = float32(10 * math.Log10(float64(v)))
		}

		// Filter FFT
		for i := 0; i < len(fftReal); i++ {
			fftReal[i] = (lastRealFFT[i]*(fftFilterAlpha-1) + fftReal[i]) / fftFilterAlpha
		}
		copy(lastRealFFT, fftReal)
		GenerateImage()
	}
}
