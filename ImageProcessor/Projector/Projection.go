package Projector

type ProjectionConverter interface {
	// LatLon2XY Converts Latitude/Longitude to Pixel X/Y
	// lat => Latitude in Degrees
	// lon => Longitude in Degrees
	LatLon2XY(lat, lon float64) (x, y int)

	// LatLon2XYf Converts Latitude/Longitude to Pixel X/Y (float64)
	// lat => Latitude in Degrees
	// lon => Longitude in Degrees
	LatLon2XYf(lat, lon float64) (x, y float64)

	// XY2LatLon Converts Pixel X/Y to Latitude/Longitude
	// lat => Latitude in Degrees
	// lon => Longitude in Degrees
	XY2LatLon(x, y int) (lat, lon float64)

	Hash() string

	// ColumnOffset returns the number of pixels that the image is offset from left
	ColumnOffset() int

	// LineOffset returns the number of pixels that the image is offset from top
	LineOffset() int

	// CropLeft returns the number of pixels that should be cropped
	CropLeft() int

	// MaxLatitude returns the Maximum Visible Latitude
	MaxLatitude() float64

	// MinLatitude returns Minimum Visible Latitude
	MinLatitude() float64

	// MaxLongitude returns Maximum visible Longitude
	MaxLongitude() float64

	// MinLongitude returns Minimum visible latitude
	MinLongitude() float64

	// LatitudeCoverage returns Coverage of the view in Latitude Degrees
	LatitudeCoverage() float64

	// LongitudeCoverage returns Coverage of the view in Longitude Degrees
	LongitudeCoverage() float64

	// TrimLongitude returns Longitude Trim parameter for removing artifacts on Reprojection (in degrees)
	TrimLongitude() float64

	// TrimLatitude returns Latitude Trim parameter for removing artifacts on Reprojection (in degrees)
	TrimLatitude() float64
}
