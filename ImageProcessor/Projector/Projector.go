package Projector

import (
	"image"
	"image/color"
	"math"
	"runtime"
	"sync"
)

type Projector struct {
	gc ProjectionConverter
}

func MakeProjector(gc ProjectionConverter) *Projector {
	return &Projector{
		gc: gc,
	}
}

func (p *Projector) DrawLatLonLines(src *image.RGBA, thickness int, c color.Color) {
	t0 := -thickness / 2
	t1 := thickness - t0

	for lat := p.gc.MinLatitude(); lat < p.gc.MaxLatitude(); lat += 10 {
		for lon := p.gc.MinLongitude(); lon < p.gc.MaxLongitude(); lon += 0.01 {
			x, y := p.gc.LatLon2XY(lat, lon)
			for i := -t0; i < t1; i++ {
				src.Set(x, y+i, c)
			}
		}
	}

	for lon := p.gc.MinLongitude(); lon < p.gc.MaxLongitude(); lon += 10 {
		for lat := p.gc.MinLatitude(); lat < p.gc.MaxLatitude(); lat += 0.01 {
			x, y := p.gc.LatLon2XY(lat, lon)
			for i := -t0; i < t1; i++ {
				src.Set(x, y+i, c)
			}
		}
	}
}

func (p *Projector) ReprojectLinearGray(src *image.Gray) *image.Gray {
	dst := image.NewGray(src.Bounds())
	sampler := MakeSampler2D(src)
	width := src.Bounds().Dx()
	height := src.Bounds().Dy()

	p.reprojectChunkGray(dst, sampler, 0, 0, width, height)

	return dst
}

func (p *Projector) ReprojectLinear(src *image.RGBA) *image.RGBA {
	dst := image.NewRGBA(src.Bounds())
	sampler := MakeSampler2D(src)
	width := src.Bounds().Dx()
	height := src.Bounds().Dy()

	p.reprojectChunk(dst, sampler, 0, 0, width, height)

	return dst
}

func (p *Projector) ReprojectLinearMultiThread(src image.Image) *image.RGBA {
	dst := image.NewRGBA(src.Bounds())
	sampler := MakeSampler2D(src)
	width := src.Bounds().Dx()
	height := src.Bounds().Dy()

	n := runtime.NumCPU() / 3 // Let's not use ALL available stuff since that can break other stuff
	n = int(math.Round(float64(n)))
	if n < 1 {
		n = 1
	}

	n2 := int(math.Sqrt(float64(n)))

	nX := n2
	nY := n - n2

	if nY == 0 {
		nY = 1 // Single Thread
	}

	wg := sync.WaitGroup{}

	deltaY := height / nY
	deltaX := width / nX

	for i := 0; i < nX; i++ {
		sx := i * deltaX
		ex := (i + 1) * deltaX
		if ex > width || (i+1 == nX && ex != width) {
			ex = width
		}

		for j := 0; j < nY; j++ {
			sy := j * deltaY
			ey := (j + 1) * deltaY

			if ey > height || (j+1 == nY && ey != height) {
				ey = height
			}

			wg.Add(1)
			go func() {
				defer wg.Done()
				p.reprojectChunk(dst, sampler, sx, sy, ex, ey)
			}()
		}
	}
	wg.Wait()

	return dst
}

func (p *Projector) ReprojectLinearMultiThreadGray(src image.Image) *image.Gray {
	dst := image.NewGray(src.Bounds())
	sampler := MakeSampler2D(src)
	width := src.Bounds().Dx()
	height := src.Bounds().Dy()

	n := runtime.NumCPU()
	n2 := int(math.Sqrt(float64(n)))

	nX := n2
	nY := n - n2

	if nY == 0 {
		nY = 1 // Single Thread
	}

	wg := sync.WaitGroup{}

	deltaY := height / nY
	deltaX := width / nX

	for i := 0; i < nX; i++ {
		sx := i * deltaX
		ex := (i + 1) * deltaX
		if ex > width || (i+1 == nX && ex != width) {
			ex = width
		}

		for j := 0; j < nY; j++ {
			sy := j * deltaY
			ey := (j + 1) * deltaY

			if ey > height || (j+1 == nY && ey != height) {
				ey = height
			}

			wg.Add(1)
			go func() {
				defer wg.Done()
				p.reprojectChunkGray(dst, sampler, sx, sy, ex, ey)
			}()
		}
	}
	wg.Wait()

	return dst
}

func (p *Projector) reprojectChunk(dst *image.RGBA, sampler *Sampler2D, sx, sy, ex, ey int) {
	dPtr := dst.Pix
	stride := dst.Stride
	width := dst.Bounds().Dx()
	height := dst.Bounds().Dy()

	for y := sy; y < ey; y++ {
		for x := sx; x < ex; x++ {
			lat := (p.gc.MaxLatitude() - p.gc.TrimLatitude()) - ((float64(y) * (p.gc.LatitudeCoverage() - p.gc.TrimLatitude()*2)) / float64(height))
			lon := ((float64(x) * (p.gc.LongitudeCoverage() - p.gc.TrimLongitude()*2)) / float64(width)) + (p.gc.MinLongitude() + p.gc.TrimLongitude())

			if lat > p.gc.MaxLatitude() || lat < p.gc.MinLatitude() || lon > p.gc.MaxLongitude() || lon < p.gc.MinLongitude() {
				dPtr[y*stride+x] = 0
				dPtr[y*stride+x+1] = 0
				dPtr[y*stride+x+2] = 0
				dPtr[y*stride+x+3] = 255
			} else {
				tx, ty := p.gc.LatLon2XYf(lat, lon)
				c := sampler.GetPixel(tx, ty)

				r, g, b, a := c.RGBA()

				dPtr[y*stride+x*4+0] = uint8(r)
				dPtr[y*stride+x*4+1] = uint8(g)
				dPtr[y*stride+x*4+2] = uint8(b)
				dPtr[y*stride+x*4+3] = uint8(a)
			}
		}
	}
}

func (p *Projector) reprojectChunkGray(dst *image.Gray, sampler *Sampler2D, sx, sy, ex, ey int) {
	dPtr := dst.Pix
	stride := dst.Stride
	width := dst.Bounds().Dx()
	height := dst.Bounds().Dy()

	for y := sy; y < ey; y++ {
		for x := sx; x < ex; x++ {
			lat := (p.gc.MaxLatitude() - p.gc.TrimLatitude()) - ((float64(y) * (p.gc.LatitudeCoverage() - p.gc.TrimLatitude()*2)) / float64(height))
			lon := ((float64(x) * (p.gc.LongitudeCoverage() - p.gc.TrimLongitude()*2)) / float64(width)) + (p.gc.MinLongitude() + p.gc.TrimLongitude())

			if lat > p.gc.MaxLatitude() || lat < p.gc.MinLatitude() || lon > p.gc.MaxLongitude() || lon < p.gc.MinLongitude() {
				dPtr[y*stride+x] = 0
			} else {
				tx, ty := p.gc.LatLon2XYf(lat, lon)
				c := sampler.GetPixelGray(tx, ty)

				dPtr[y*stride+x] = uint8(c.Y)
			}
		}
	}
}
