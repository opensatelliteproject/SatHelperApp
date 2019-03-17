package ImageTools

import (
	"fmt"
	"github.com/llgcode/draw2d"
	"github.com/llgcode/draw2d/draw2dimg"
	"github.com/opensatelliteproject/SatHelperApp"
	"github.com/opensatelliteproject/SatHelperApp/ImageProcessor/ImageData"
	"github.com/opensatelliteproject/SatHelperApp/XRIT"
	"image"
	"image/color"
	"image/draw"
	"math"
)

const (
	lineWidth       = 18
	fontSpaceFactor = 1.5

	baseImageWidth = 2732

	baseScaleHeight       = float64(30)
	baseImgPad            = 10
	baseScaleTextFontSize = 38
	baseTitleFontSize     = 50
	baseFootFontSize      = 34

	minFontSize = 10
)

type ImageEnhancer struct {
	overLay          bool
	colorMap         []color.Color
	originalColorMap []color.Color
	colorLut         *Lut1D
	minTemp          float32
	maxTemp          float32
}

func scaleFont(sf, fontSize float64) float64 {
	fontSize *= sf
	if fontSize < minFontSize {
		return minFontSize
	}

	return fontSize
}

func MakeImageEnhancer(minTemp, maxTemp float32, satLut []float32, enhanceLut []color.Color, overlay bool) *ImageEnhancer {
	colorMap := ImageData.ScaleLutToColor(minTemp, maxTemp, satLut, enhanceLut)

	l, err := MakeLut1DFromColors(colorMap)

	if err != nil {
		panic(err)
	}

	return &ImageEnhancer{
		colorMap:         colorMap,
		originalColorMap: enhanceLut,
		minTemp:          minTemp,
		maxTemp:          maxTemp,
		overLay:          overlay,
		colorLut:         l,
	}
}

func (ie *ImageEnhancer) EnhanceWithLUT(img *image.RGBA) (*image.RGBA, error) {
	// First apply LUT
	img, err := ie.colorLut.ApplyFromRGBA(img)

	if err != nil {
		return nil, err
	}

	return img, nil
}

func (ie *ImageEnhancer) DrawMeta(section string, img *image.RGBA, xh *XRIT.Header) (*image.RGBA, error) {
	// Compute Scale of everything
	scaleFactor := float64(img.Bounds().Dx()) / baseImageWidth

	// Scale Constants used here
	scaleHeight := scaleFactor * baseScaleHeight
	imgPad := scaleFactor * baseImgPad
	scaleTextFontSize := scaleFont(scaleFactor, baseScaleTextFontSize)
	titleFontSize := scaleFont(scaleFactor, baseTitleFontSize)
	footFontSize := scaleFont(scaleFactor, baseFootFontSize)

	scaleTotalHeight := imgPad + scaleHeight + scaleTextFontSize*fontSpaceFactor + imgPad
	headHeight := imgPad + titleFontSize*fontSpaceFactor + imgPad
	footHeight := imgPad + footFontSize*fontSpaceFactor + imgPad
	outLayDiff := scaleTotalHeight + headHeight + footHeight

	// Check Title
	title := xh.ToBaseNameString()
	if section != "" {
		title += " - " + section
	}
	// TODO: Check if title requires more space than image width

	// Now check if its a overlay or outlay
	if !ie.overLay {
		b := img.Bounds()

		oldW := b.Dx()
		oldH := b.Dy()

		newW := oldW + 2*int(imgPad)
		newH := oldH + int(outLayDiff)

		targetRect := image.Rect(int(imgPad), int(headHeight), newW-int(imgPad), int(headHeight)+oldH)

		// Create Image that fits head and scale
		n := image.NewRGBA(image.Rect(0, 0, newW, newH))

		draw.Draw(n, targetRect, img, b.Min, draw.Src)
		img = n
	}

	w := float64(img.Bounds().Dx())
	h := float64(img.Bounds().Dy())

	dc := draw2dimg.NewGraphicContext(img)
	dc.SetFontData(draw2d.FontData{
		Name:   "FreeMono",
		Family: draw2d.FontFamilyMono,
		Style:  draw2d.FontStyleNormal,
	})
	dc.SetFontSize(titleFontSize)
	dc.SetFillColor(color.White)
	dc.SetStrokeColor(color.Black)
	dc.SetLineWidth(lineWidth)

	ie.drawTitle(title, w, scaleFactor, dc)
	ie.drawScale(w, h, scaleFactor, img, dc)
	ie.drawFoot(w, h, scaleFactor, dc, xh)

	return img, nil
}

func (ie *ImageEnhancer) DrawMetaWithoutScale(section string, img *image.RGBA, xh *XRIT.Header) (*image.RGBA, error) {
	// Compute Scale of everything
	scaleFactor := float64(img.Bounds().Dx()) / baseImageWidth

	// Scale Constants used here
	scaleHeight := scaleFactor * baseScaleHeight
	imgPad := scaleFactor * baseImgPad
	scaleTextFontSize := scaleFont(scaleFactor, baseScaleTextFontSize)
	titleFontSize := scaleFont(scaleFactor, baseTitleFontSize)
	footFontSize := scaleFont(scaleFactor, baseFootFontSize)

	scaleTotalHeight := imgPad + scaleHeight + scaleTextFontSize*fontSpaceFactor + imgPad
	headHeight := imgPad + titleFontSize*fontSpaceFactor + imgPad
	footHeight := imgPad + footFontSize*fontSpaceFactor + imgPad
	outLayDiff := scaleTotalHeight + headHeight + footHeight

	// Check Title
	title := xh.ToBaseNameString()
	if section != "" {
		title += " - " + section
	}
	// TODO: Check if title requires more space than image width

	// Now check if its a overlay or outlay
	if !ie.overLay {
		b := img.Bounds()

		oldW := b.Dx()
		oldH := b.Dy()

		newW := oldW + 2*int(imgPad)
		newH := oldH + int(outLayDiff)

		targetRect := image.Rect(int(imgPad), int(headHeight), newW-int(imgPad), int(headHeight)+oldH)

		// Create Image that fits head and scale
		n := image.NewRGBA(image.Rect(0, 0, newW, newH))

		draw.Draw(n, targetRect, img, b.Min, draw.Src)
		img = n
	}

	w := float64(img.Bounds().Dx())
	h := float64(img.Bounds().Dy())

	dc := draw2dimg.NewGraphicContext(img)
	dc.SetFontData(draw2d.FontData{
		Name:   "FreeMono",
		Family: draw2d.FontFamilyMono,
		Style:  draw2d.FontStyleNormal,
	})
	dc.SetFontSize(titleFontSize)
	dc.SetFillColor(color.White)
	dc.SetStrokeColor(color.Black)
	dc.SetLineWidth(lineWidth)

	ie.drawTitle(title, w, scaleFactor, dc)
	ie.drawFoot(w, h, scaleFactor, dc, xh)

	return img, nil
}

func (ie *ImageEnhancer) drawScale(w, h, sf float64, img *image.RGBA, dc *draw2dimg.GraphicContext) {

	// Scale Constants
	scaleHeight := sf * baseScaleHeight
	imgPad := sf * baseImgPad
	scaleTextFontSize := scaleFont(sf, baseScaleTextFontSize)
	scaleTotalHeight := imgPad + scaleHeight + scaleTextFontSize*fontSpaceFactor + imgPad
	footFontSize := scaleFont(sf, baseFootFontSize)
	footHeight := imgPad + footFontSize*fontSpaceFactor + imgPad

	dc.Save()
	dc.SetFillColor(color.White)
	dc.SetStrokeColor(color.Black)
	dc.SetFontSize(scaleTextFontSize)

	// Calculate Scale Box boundaries

	ex := int(w - scaleTotalHeight)
	tScale := ex / len(ie.colorMap)
	ex = len(ie.colorMap) * tScale
	pad := (w - float64(ex)) / 2

	ex = tScale * len(ie.colorMap)

	x0 := int(pad)
	x1 := int(w - pad)
	y0 := int(h - footHeight - scaleTotalHeight + imgPad)
	y1 := int(h - footHeight - (scaleTotalHeight - scaleHeight) + imgPad)

	// Reorder points
	if x0 > x1 {
		a := x0
		x0 = x1
		x1 = a
	}
	if y0 > y1 {
		a := y0
		y0 = y1
		y1 = a
	}

	scaleBox := image.Rect(x0, y0, x1, y1)

	// Draw LUT Color
	for i, v := range ie.originalColorMap {
		for y := scaleBox.Min.Y + 1; y < scaleBox.Max.Y-1; y++ {
			for x := scaleBox.Min.X + i*tScale; x < scaleBox.Min.X+(i+1)*tScale; x++ {
				img.Set(x, y, v)
			}
		}
	}

	// Draw border
	dc.MoveTo(float64(scaleBox.Min.X), float64(scaleBox.Min.Y))
	dc.LineTo(float64(scaleBox.Min.X), float64(scaleBox.Max.Y))
	dc.LineTo(float64(scaleBox.Max.X), float64(scaleBox.Max.Y))
	dc.LineTo(float64(scaleBox.Max.X), float64(scaleBox.Min.Y))
	dc.Close()
	dc.SetLineWidth(4 * sf)
	dc.SetStrokeColor(color.RGBA{
		R: 120,
		G: 120,
		B: 120,
		A: 255,
	})
	dc.Stroke()

	// Draw Texts

	scaleCount := 16
	scaleSpace := scaleBox.Dx() / scaleCount

	drawScale := func(data string, x, y float64) {
		w := float64(len(data)) * dc.GetFontSize()
		dc.StrokeStringAt(data, x-w/2, y+dc.GetFontSize())
		dc.FillStringAt(data, x-w/2, y+dc.GetFontSize())
	}

	drawScale("ºK", float64(scaleBox.Max.X), float64(scaleBox.Min.Y-int(imgPad*2))-dc.GetFontSize())

	for i := 0; i <= scaleCount; i++ {
		x := i * scaleSpace
		temp := ImageData.LutIndexToTemperature(ie.minTemp, ie.maxTemp, x/tScale)
		drawScale(fmt.Sprintf("%.0f", temp), float64(x+scaleBox.Min.X), float64(scaleBox.Max.Y+8))
	}

	dc.Restore()
}

func (ie *ImageEnhancer) drawTitle(title string, w, sf float64, dc *draw2dimg.GraphicContext) {
	dc.Save()
	imgPad := sf * baseImgPad
	titleFontSize := scaleFont(sf, baseTitleFontSize)

	dc.SetFontSize(titleFontSize)
	dc.SetFillColor(color.White)
	dc.SetStrokeColor(color.Black)
	dc.SetLineWidth(lineWidth * sf)

	dc.MoveTo(0, 0)

	x0, y0, x1, y1 := dc.GetStringBounds(title)
	titleX := (w - math.Abs(x1-x0)) / 2
	titleY := imgPad + math.Abs(y0-y1)

	dc.StrokeStringAt(title, titleX, titleY)
	dc.FillStringAt(title, titleX, titleY)

	dc.Restore()
}

func (ie *ImageEnhancer) drawFoot(w, h, sf float64, dc *draw2dimg.GraphicContext, xh *XRIT.Header) {
	dc.Save()

	imgPad := sf * baseImgPad
	footFontSize := scaleFont(sf, baseFootFontSize)

	footHeight := imgPad + footFontSize*fontSpaceFactor + imgPad

	x0 := int(0)
	x1 := int(w)
	y0 := int(h - footHeight + imgPad)
	y1 := int(h)

	// Reorder points
	if x0 > x1 {
		a := x0
		x0 = x1
		x1 = a
	}
	if y0 > y1 {
		a := y0
		y0 = y1
		y1 = a
	}

	dc.SetLineWidth(1)
	dc.SetStrokeColor(color.RGBA{
		R: 255,
		G: 255,
		B: 255,
		A: 255,
	})
	dc.MoveTo(float64(x0), float64(y0))
	dc.LineTo(float64(x1), float64(y0))
	dc.Stroke()

	versionString := "SatHelperApp"

	if SatHelperApp.GetVersion() != "<unknown>" {
		versionString += fmt.Sprintf(" %s.%s", SatHelperApp.GetVersion(), SatHelperApp.GetRevision())
	}

	dc.FillStringAt(versionString, imgPad, float64(y1)-imgPad)

	if xh.TimestampHeader != nil {
		dt := xh.TimestampHeader.GetDateTime().String()
		dc.MoveTo(0, 0)
		sx0, _, sx1, _ := dc.GetStringBounds(dt)
		sw := math.Abs(sx1 - sx0)
		dc.FillStringAt(dt, w-sw, float64(y1)-imgPad)
	}

	dc.Restore()
}
