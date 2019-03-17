package Projector

import (
	"crypto/sha256"
	"fmt"
)

// LinearConverter
// This is used for pixel map to coordinates from a reprojected image.
type LinearConverter struct {
	minLat float64
	minLon float64
	maxLat float64
	maxLon float64

	imgWidth  int
	imgHeight int
}

// MakeLinearConverter Creates a new instance of LinearConverter
//  This is used for pixel map to coordinates from a reprojected image.
//  imgWidth => Image Width in pixels
//  imgHeight => Image Height in pixels
//  pc => Previous Projection Converter (before reprojection)
func MakeLinearConverter(imgWidth, imgHeight int, pc ProjectionConverter) ProjectionConverter {
	return &LinearConverter{
		minLat:    pc.MinLatitude(),
		maxLat:    pc.MaxLatitude(),
		minLon:    pc.MinLongitude(),
		maxLon:    pc.MaxLongitude(),
		imgWidth:  imgWidth,
		imgHeight: imgHeight,
	}
}

// LatLon2XY Converts Latitude/Longitude to Pixel X/Y
// lat => Latitude in Degrees
// lon => Longitude in Degrees
func (gc *LinearConverter) LatLon2XY(lat, lon float64) (x, y int) {
	xf, yf := gc.LatLon2XYf(lat, lon)

	x = int(xf)
	y = int(yf)

	return
}

// LatLon2XYf Converts Latitude/Longitude to Pixel X/Y (float64)
// lat => Latitude in Degrees
// lon => Longitude in Degrees
func (gc *LinearConverter) LatLon2XYf(lat, lon float64) (x, y float64) {
	nLat := -((lat - gc.maxLat + gc.TrimLatitude()) / (gc.LatitudeCoverage() - gc.TrimLatitude()*2))
	nLon := (lon - gc.minLon - gc.TrimLongitude()) / (gc.LongitudeCoverage() - gc.TrimLongitude()*2)

	x = nLon * float64(gc.imgWidth)
	y = nLat * float64(gc.imgHeight)

	if x < 0 {
		x = 0
	}
	if x > float64(gc.imgWidth) {
		x = float64(gc.imgWidth)
	}

	if y < 0 {
		y = 0
	}

	if y > float64(gc.imgHeight) {
		y = float64(gc.imgHeight)
	}

	return
}

// XY2LatLon Converts Pixel X/Y to Latitude/Longitude
// lat => Latitude in Degrees
// lon => Longitude in Degrees
func (gc *LinearConverter) XY2LatLon(x, y int) (lat, lon float64) {
	nX := float64(x) / float64(gc.imgWidth)
	nY := float64(y) / float64(gc.imgHeight)

	lat = (gc.MaxLatitude() - gc.TrimLatitude()) - (nY * (gc.LatitudeCoverage() - gc.TrimLatitude()*2))
	lon = (nX * (gc.LongitudeCoverage() - gc.TrimLongitude()*2)) + (gc.MinLongitude() + gc.TrimLongitude())

	return
}

// region Getters

// ColumnOffset returns the number of pixels that the image is offset from left
func (gc *LinearConverter) ColumnOffset() int {
	return 0
}

// LineOffset returns the number of pixels that the image is offset from top
func (gc *LinearConverter) LineOffset() int {
	return 0
}

// CropLeft returns the number of pixels that should be cropped
func (gc *LinearConverter) CropLeft() int {
	return 0
}

// MaxLatitude returns the Maximum Visible Latitude
func (gc *LinearConverter) MaxLatitude() float64 {
	return gc.maxLat
}

// MinLatitude returns Minimum Visible Latitude
func (gc *LinearConverter) MinLatitude() float64 {
	return gc.minLat
}

// MaxLongitude returns Maximum visible Longitude
func (gc *LinearConverter) MaxLongitude() float64 {
	return gc.maxLon
}

// MinLongitude returns Minimum visible latitude
func (gc *LinearConverter) MinLongitude() float64 {
	return gc.minLon
}

// LatitudeCoverage returns Coverage of the view in Latitude Degrees
func (gc *LinearConverter) LatitudeCoverage() float64 {
	return gc.MaxLatitude() - gc.MinLatitude()
}

// LongitudeCoverage returns Coverage of the view in Longitude Degrees
func (gc *LinearConverter) LongitudeCoverage() float64 {
	return gc.MaxLongitude() - gc.MinLongitude()
}

// TrimLongitude returns Longitude Trim parameter for removing artifacts on Reprojection (in degrees)
func (gc *LinearConverter) TrimLongitude() float64 {
	return 16
}

// TrimLatitude returns Latitude Trim parameter for removing artifacts on Reprojection (in degrees)
func (gc *LinearConverter) TrimLatitude() float64 {
	return 16
}

func (gc *LinearConverter) Hash() string {
	s := fmt.Sprintf("LinearConverter%f%f%f%f%d%d", gc.maxLon, gc.minLon, gc.maxLat, gc.minLat, gc.imgWidth, gc.imgHeight)
	h := sha256.New()
	_, _ = h.Write([]byte(s))

	return fmt.Sprintf("%x", h.Sum(nil))
}

// endregion
