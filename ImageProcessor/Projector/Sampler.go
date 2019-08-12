package Projector

import (
	"image"
	"image/color"
)

type Sampler2D struct {
	numChannels int
	stride      int
	pixData     []byte
}

func MakeSampler2D(img image.Image) *Sampler2D {
	o := &Sampler2D{}

	switch v := img.(type) {
	case *image.RGBA:
		o.numChannels = 4
		o.pixData = v.Pix
		o.stride = v.Stride
	case *image.Gray:
		o.numChannels = 1
		o.pixData = v.Pix
		o.stride = v.Stride
	default:
		panic("image type not supported")
	}

	return o
}

func (s *Sampler2D) GetPixelGray(x, y float64) color.Gray {
	return color.Gray{
		Y: BilinearInterp(s.pixData, x, y, s.stride, s.numChannels, 0),
	}
}

func (s *Sampler2D) GetPixel(x, y float64) color.Color {
	var r, g, b, a byte

	if x < 0 || y < 0 {
		return color.Black
	}

	r = BilinearInterp(s.pixData, x, y, s.stride, s.numChannels, 0)
	a = 255

	switch s.numChannels {
	case 1:
		g = r
		b = r
	case 3:
		g = BilinearInterp(s.pixData, x, y, s.stride, s.numChannels, 1)
		b = BilinearInterp(s.pixData, x, y, s.stride, s.numChannels, 2)
	case 4:
		g = BilinearInterp(s.pixData, x, y, s.stride, s.numChannels, 1)
		b = BilinearInterp(s.pixData, x, y, s.stride, s.numChannels, 2)
		a = BilinearInterp(s.pixData, x, y, s.stride, s.numChannels, 3)
	}

	return color.RGBA{
		R: r,
		G: g,
		B: b,
		A: a,
	}
}

func BilinearInterp(data []byte, x, y float64, mw int, numChannels, colorChannel int) byte {
	rx := int(x)
	ry := int(y)
	fracX := x - float64(rx)
	fracY := y - float64(ry)

	if fracX == 0 && fracY == 0 { // Integer amount
		return valueAtImage(data, rx, ry, mw, numChannels, colorChannel)
	}

	invFracX := 1 - fracX
	invFracY := 1 - fracY

	a := valueAtImageF(data, rx, ry, mw, numChannels, colorChannel)
	b := valueAtImageF(data, rx+1, ry, mw, numChannels, colorChannel)
	c := valueAtImageF(data, rx, ry+1, mw, numChannels, colorChannel)
	d := valueAtImageF(data, rx+1, ry+1, mw, numChannels, colorChannel)

	v := (a*invFracX+b*fracX)*invFracY + (c*invFracX+d*fracX)*fracY

	return byte(v)
}

func valueAtImageF(data []byte, x, y, mw, numChannels, colorChannel int) float64 {
	return float64(valueAtImage(data, x, y, mw, numChannels, colorChannel))
}

func valueAtImage(data []byte, x, y, mw, numChannels, colorChannel int) byte {
	px := y*mw + x*numChannels + colorChannel
	if px >= len(data) {
		return 0
	}
	return data[px]
}
