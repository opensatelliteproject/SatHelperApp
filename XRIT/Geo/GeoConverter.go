package Geo

import (
	"crypto/sha256"
	"fmt"
	"github.com/opensatelliteproject/SatHelperApp/ImageProcessor/Projector"
	"github.com/opensatelliteproject/SatHelperApp/XRIT"
	"math"
	"regexp"
	"strconv"
)

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
func MakeGeoConverter(satelliteLongitude float64, coff, loff int, cfac, lfac float64, fixAspect bool, imageWidth int) Projector.ProjectionConverter {
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

// MakeGeoConverterFromXRIT Creates a new instance of GeoConverter from a XRIT File Header
func MakeGeoConverterFromXRIT(xh *XRIT.Header) (Projector.ProjectionConverter, error) {

	x := regexp.MustCompile(`.*\((.*)\)`)

	regMatch := x.FindStringSubmatch(xh.ImageNavigationHeader.ProjectionName)
	if len(regMatch) < 2 {
		return nil, fmt.Errorf("cannot find projection lon at %s", xh.ImageNavigationHeader.ProjectionName)
	}

	lon, err := strconv.ParseFloat(regMatch[1], 64)

	if err != nil {
		return nil, err
	}

	if xh.ImageNavigationHeader == nil {
		return nil, fmt.Errorf("no image navigation header")
	}

	inh := xh.ImageNavigationHeader

	if xh.SegmentIdentificationHeader != nil && xh.SegmentIdentificationHeader.COMS1 {
		xh.ImageNavigationHeader.LineScalingFactor >>= 9 // Not sure why is needed
		xh.ImageNavigationHeader.LineScalingFactor -= 180000
	}

	return MakeSimpleGeoConverter(lon, int(inh.ColumnOffset), int(inh.LineOffset), float64(inh.ColumnScalingFactor), float64(inh.LineScalingFactor)), nil

}

// MakeSimpleGeoConverter Creates a new instance of GeoConverter
//  Same as MakeGeoConverter but with fixAspect disabled and imageWidth = 0
//  satelliteLongitude => Satellite longitude.
//  coff => Column Offset
//  loff => Line Offset
//  cfac => Column Scaling Factor
//  lfac => Line Scaling Factor
func MakeSimpleGeoConverter(satelliteLongitude float64, coff, loff int, cfac, lfac float64) Projector.ProjectionConverter {
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

func (gc *Converter) Hash() string {
	s := fmt.Sprintf("%f%d%d%f%f%v%d%v", gc.satelliteLongitude, gc.coff, gc.loff, gc.lfac, gc.cfac, gc.fixAspect, gc.imageWidth, gc.cropLeft)
	h := sha256.New()
	_, _ = h.Write([]byte(s))

	return fmt.Sprintf("%x", h.Sum(nil))
}

// endregion
