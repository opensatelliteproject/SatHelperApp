package Geo

import "math"

type Converter struct {
	satelliteLongitude float64 // Satellite Longitude
	coff               int     // Column Offset
	loff               int     // Line Offset
	cfac               float64 // Column Scaling Factor
	lfac               float64 // Line Scaling Factor
	fixAspect          bool    // Fix Aspect Ratio if needed
	imageWidth         int     // Image Width

	aspectRatio float64
	cropLeft    int
}

// MakeGeoConverter Creates a new instance of GeoConverter
//  satelliteLongitude => Satellite longitude.
//  coff => Column Offset
//  loff => Line Offset
//  cfac => Column Scaling Factor
//  lfac => Line Scaling Factor
//  fixAspect => If the aspect ratio should be fixed for cutting image
//  imageWidth => Image Width in pixels
func MakeGeoConverter(satelliteLongitude float64, coff, loff int, cfac, lfac float64, fixAspect bool, imageWidth int) *Converter {
	return &Converter{
		satelliteLongitude: satelliteLongitude,
		coff:               coff,
		loff:               loff,
		cfac:               cfac,
		lfac:               lfac,
		fixAspect:          fixAspect,
		imageWidth:         imageWidth,
		aspectRatio:        cfac / lfac,
		cropLeft:           coff - int(math.Min(float64(imageWidth-coff), float64(coff))),
	}
}

// MakeSimpleGeoConverter Creates a new instance of GeoConverter
//  Same as MakeGeoConverter but with fixAspect disabled and imageWidth = 0
//  satelliteLongitude => Satellite longitude.
//  coff => Column Offset
//  loff => Line Offset
//  cfac => Column Scaling Factor
//  lfac => Line Scaling Factor
func MakeSimpleGeoConverter(satelliteLongitude float64, coff, loff int, cfac, lfac float64) *Converter {
	return MakeGeoConverter(satelliteLongitude, coff, loff, cfac, lfac, false, 0)
}

// LatLon2XY Converts Latitude/Longitude to Pixel X/Y
// lat => Latitude in Degrees
// lon => Longitude in Degrees
func (gc *Converter) LatLon2XY(lat, lon float64) (x, y int) {
	x, y = LonLat2XY(gc.satelliteLongitude, Deg2Rad(lon), Deg2Rad(lat), gc.coff, gc.cfac, gc.loff, gc.lfac)

	if gc.fixAspect {
		y = int(float64(y) * gc.aspectRatio)
	}

	return
}

// LatLon2XYf Converts Latitude/Longitude to Pixel X/Y (float64)
// lat => Latitude in Degrees
// lon => Longitude in Degrees
func (gc *Converter) LatLon2XYf(lat, lon float64) (x, y float64) {
	x, y = LonLat2XYf(gc.satelliteLongitude, Deg2Rad(lon), Deg2Rad(lat), gc.coff, gc.cfac, gc.loff, gc.lfac)

	if gc.fixAspect {
		y *= gc.aspectRatio
	}

	return
}

// XY2LatLon Converts Pixel X/Y to Latitude/Longitude
// lat => Latitude in Degrees
// lon => Longitude in Degrees
func (gc *Converter) XY2LatLon(x, y int) (lat, lon float64) {
	lat, lon = XY2LonLat(gc.satelliteLongitude, x, y, gc.coff, gc.cfac, gc.loff, gc.lfac)
	lat = Rad2Deg(lat)
	lon = Rad2Deg(lon)
	return
}

// region Getters

// ColumnOffset returns the number of pixels that the image is offset from left
func (gc *Converter) ColumnOffset() int {
	return gc.coff
}

// LineOffset returns the number of pixels that the image is offset from top
func (gc *Converter) LineOffset() int {
	return gc.loff
}

// CropLeft returns the number of pixels that should be cropped
func (gc *Converter) CropLeft() int {
	return gc.cropLeft
}

// MaxLatitude returns the Maximum Visible Latitude
func (gc *Converter) MaxLatitude() float64 {
	return 79
}

// MinLatitude returns Minimum Visible Latitude
func (gc *Converter) MinLatitude() float64 {
	return -79
}

// MaxLongitude returns Maximum visible Longitude
func (gc *Converter) MaxLongitude() float64 {
	return gc.satelliteLongitude + 79
}

// MinLongitude returns Minimum visible latitude
func (gc *Converter) MinLongitude() float64 {
	return gc.satelliteLongitude - 79
}

// LatitudeCoverage returns Coverage of the view in Latitude Degrees
func (gc *Converter) LatitudeCoverage() float64 {
	return gc.MaxLatitude() - gc.MinLatitude()
}

// LongitudeCoverage returns Coverage of the view in Longitude Degrees
func (gc *Converter) LongitudeCoverage() float64 {
	return gc.MaxLongitude() - gc.MinLongitude()
}

// TrimLongitude returns Longitude Trim parameter for removing artifacts on Reprojection (in degrees)
func (gc *Converter) TrimLongitude() float64 {
	return 16
}

// TrimLatitude returns Latitude Trim parameter for removing artifacts on Reprojection (in degrees)
func (gc *Converter) TrimLatitude() float64 {
	return 16
}

// endregion
