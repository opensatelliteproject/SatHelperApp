package ImageTools

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"io"
	"os"
	"reflect"
)

type Lut1D struct {
	lut []color.Color
}

func MakeLut1D(filename string) (*Lut1D, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}

	defer f.Close()

	return MakeLut1DFromReader(f)
}

func MakeLut1DFromMemory(data []byte) (*Lut1D, error) {
	return MakeLut1DFromReader(bytes.NewReader(data))
}

func MakeLut1DFromColors(colors []color.Color) (*Lut1D, error) {
	if len(colors) != 256 {
		return nil, fmt.Errorf("invalid dimension %d. Expected 256", len(colors))
	}

	l := &Lut1D{
		lut: make([]color.Color, 256),
	}

	copy(l.lut, colors)

	return l, nil
}

func MakeLut1DFromReader(r io.Reader) (*Lut1D, error) {
	img, _, err := image.Decode(r)

	if err != nil {
		return nil, err
	}

	b := img.Bounds()

	if b.Dx() != 256 || b.Dy() != 1 {
		return nil, fmt.Errorf("invalid dimensions %dx%d. Expected 256x1", b.Dx(), b.Dy())
	}

	lut2d := &Lut1D{
		lut: make([]color.Color, 256),
	}

	var cr ColorReader

	switch v := img.(type) {
	case *image.NRGBA:
		o := image.NewRGBA(v.Bounds())
		draw.Draw(o, v.Bounds(), v, v.Bounds().Min, draw.Src)
		cr = o
	case *image.RGBA:
		cr = v
	case *image.Gray:
		cr = v
	default:
		return nil, fmt.Errorf("invalid image type: %s", reflect.TypeOf(img))
	}

	for x := 0; x < 256; x++ {
		lut2d.lut[x] = cr.At(x, 0)
	}

	return lut2d, nil
}

func (l1d *Lut1D) ApplyFromGray(a *image.Gray) (*image.RGBA, error) {
	out := image.NewRGBA(a.Bounds())

	s := a.Bounds()

	for y := 0; y < s.Dy(); y++ {
		for x := 0; x < s.Dx(); x++ {
			i := a.Pix[y*a.Stride+x]
			out.Set(x, y, l1d.lut[i])
		}
	}

	return out, nil
}

func (l1d *Lut1D) ApplyFromRGBA(a *image.RGBA) (*image.RGBA, error) {
	out := image.NewRGBA(a.Bounds())

	s := a.Bounds()

	for y := 0; y < s.Dy(); y++ {
		for x := 0; x < s.Dx(); x++ {
			i := int(a.Pix[y*a.Stride+x*4])
			out.Set(x, y, l1d.lut[int(i)])
		}
	}

	return out, nil
}
