package MapDrawer

import (
	"bytes"
	"github.com/jonas-p/go-shp"
	"github.com/llgcode/draw2d/draw2dimg"
	"github.com/opensatelliteproject/SatHelperApp/XRIT/Geo"
	"image"
	"image/color"
	"math"
)

const defaultLineWidth = 3

var defaultLineColor = color.RGBA{R: 255, A: 255}

type MapDrawer struct {
	sections  []*MapSection
	cache     map[string][]*MapSection
	lineWidth float64
	lineColor color.Color
}

type MapSection struct {
	name     string
	polygons []shp.Polygon
	fields   map[string]string
}

func MakeMapDrawer(shapeFile string) (error, *MapDrawer) {
	shape, err := shp.Open(shapeFile)
	if err != nil {
		return err, nil
	}
	defer shape.Close()

	md := &MapDrawer{
		sections:  make([]*MapSection, 0),
		cache:     make(map[string][]*MapSection),
		lineWidth: defaultLineWidth,
		lineColor: defaultLineColor,
	}

	fields := shape.Fields()
	for shape.Next() {
		n, p := shape.Shape()

		var poly *shp.Polygon

		switch v := p.(type) {
		case *shp.Polygon:
			poly = v
		}

		if poly == nil {
			continue
		}

		// Let's split each polygon part into several polygons
		section := &MapSection{
			polygons: make([]shp.Polygon, 0),
			fields:   make(map[string]string),
		}

		if len(poly.Parts) > 1 {
			idx := poly.Parts[0]
			for i := 1; i < len(poly.Parts); i++ {
				p := shp.Polygon{
					Parts:    make([]int32, 1),
					Points:   poly.Points[idx:poly.Parts[i]],
					NumParts: 1,
				}
				p.NumPoints = int32(len(p.Points))

				idx = poly.Parts[i]
				section.polygons = append(section.polygons, p)
			}
			// Last Part
			p := shp.Polygon{
				Parts:    make([]int32, 1),
				Points:   poly.Points[idx:],
				NumParts: 1,
			}
			p.NumPoints = int32(len(p.Points))
			section.polygons = append(section.polygons, p)
		} else {
			section.polygons = append(section.polygons, *poly)
		}

		for k, f := range fields {
			fieldBytes := f.Name[:]

			a := bytes.Split(fieldBytes, []byte{0})

			field := string(a[0])

			val := shape.ReadAttribute(n, k)
			section.fields[field] = val
			if field == "name" {
				section.name = val
			}
		}

		md.sections = append(md.sections, section)
	}

	return nil, md
}

func (md *MapDrawer) SetLineColor(c color.Color) {
	md.lineColor = c
}

func (md *MapDrawer) SetLineWidth(w float64) {
	md.lineWidth = w
}

func (md *MapDrawer) DrawMap(img *image.RGBA, gc *Geo.Converter) {
	w := float64(img.Bounds().Dx())
	h := float64(img.Bounds().Dy())
	draw := false

	ctx := draw2dimg.NewGraphicContext(img)
	ctx.SetStrokeColor(md.lineColor)
	ctx.SetLineWidth(md.lineWidth)

	for _, v := range md.sections {
		for _, poly := range v.polygons {
			p0 := poly.Points[0]
			lastX, lastY := gc.LatLon2XYf(p0.Y, p0.X)

			ctx.BeginPath()
			ctx.MoveTo(lastX, lastY)

			for _, p := range poly.Points {
				lat := p.Y
				lon := p.X
				if lat < gc.MaxLatitude() && lat > gc.MinLatitude() && lon < gc.MaxLongitude() && lon > gc.MinLongitude() {
					x, y := gc.LatLon2XYf(lat, lon)

					cx := float64(x)
					cy := float64(y)

					if (!math.IsNaN(lastX) && !math.IsNaN(lastY)) &&
						(x > 0 && y > 0) &&
						(cx < w && cy < h) &&
						(lastX > 0 && lastY > 0) &&
						(lastX < w && lastY < h) {

						if cx != lastX && cy != lastY {
							ctx.LineTo(cx, cy)
						}
						draw = true
					}
					lastX = cx
					lastY = cy
				}
			}
			if draw {
				ctx.Close()
				ctx.Stroke()
			}
		}
	}
}
