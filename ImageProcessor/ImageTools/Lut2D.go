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

type Lut2D struct {
	lut [][]color.Color
}

func MakeLut2D(filename string) (*Lut2D, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}

	defer f.Close()

	return MakeLut2DFromReader(f)
}

func MakeLut2DFromMemory(data []byte) (*Lut2D, error) {
	return MakeLut2DFromReader(bytes.NewReader(data))
}

func MakeLut2DFromReader(r io.Reader) (*Lut2D, error) {
	img, _, err := image.Decode(r)

	if err != nil {
		return nil, err
	}

	b := img.Bounds()

	if b.Dx() != 256 || b.Dy() != 256 {
		return nil, fmt.Errorf("invalid dimensions %dx%d. Expected 256x256", b.Dx(), b.Dy())
	}

	lut2d := &Lut2D{
		lut: make([][]color.Color, 256),
	}

	for i := 0; i < 256; i++ {
		lut2d.lut[i] = make([]color.Color, 256)
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

	for y := 0; y < 256; y++ {
		for x := 0; x < 256; x++ {
			lut2d.lut[y][x] = cr.At(x, y)
		}
	}

	return lut2d, nil
}

func (l2d *Lut2D) Apply(a, b *image.Gray) (*image.RGBA, error) {
	if !a.Bounds().Eq(b.Bounds()) {
		return nil, fmt.Errorf("the images does not have same size")
	}

	out := image.NewRGBA(a.Bounds())

	s := a.Bounds()

	for y := 0; y < s.Dy(); y++ {
		for x := 0; x < s.Dx(); x++ {
			i := a.Pix[y*a.Stride+x]
			j := b.Pix[y*b.Stride+x]
			out.Set(x, y, l2d.lut[i][j])
		}
	}

	return out, nil
}
