// Geo
// Geographic Conversion Tools
// Based on: http://www.cgms-info.org/documents/pdf_cgms_03.pdf
package Geo

import "math"

var MAXLON = Deg2Rad(75)
var MINLON = Deg2Rad(-75)
var MAXLAT = Deg2Rad(79)
var MINLAT = Deg2Rad(-79)

const radiusPoles = 6356.7523
const radiusEquator = 6378.1370
const vehicleDistance = 42142.5833
const planetAspectRatio = (radiusPoles * radiusPoles) / (radiusEquator * radiusEquator)
const surfaceDistance = vehicleDistance*vehicleDistance - radiusEquator*radiusEquator

const deg2radc = math.Pi / 180
const rad2degc = 180 / math.Pi

func Deg2Rad(deg float64) float64 {
	return deg * deg2radc
}

func Rad2Deg(rad float64) float64 {
	return rad * rad2degc
}

func LonLat2XY(satLon, lon, lat float64, coff int, cfac float64, loff int, lfac float64) (c, l int) {
	cf, lf := LonLat2XYf(satLon, lon, lat, coff, cfac, loff, lfac)
	c = int(cf)
	l = int(lf)
	return
}

func LonLat2XYf(satLon, lon, lat float64, coff int, cfac float64, loff int, lfac float64) (c, l float64) {
	subLon := Deg2Rad(satLon)
	lon -= subLon

	lon = math.Min(math.Max(lon, MINLON), MAXLON)
	lat = math.Min(math.Max(lat, MINLAT), MAXLAT)

	psi := math.Atan(planetAspectRatio * math.Tan(lat))
	re := radiusPoles / (math.Sqrt(1 - (1-planetAspectRatio)*math.Cos(psi)*math.Cos(psi)))

	r1 := vehicleDistance - re*math.Cos(psi)*math.Cos(lon)
	r2 := -1 * re * math.Cos(psi) * math.Sin(lon)
	r3 := re * math.Sin(psi)

	rn := math.Sqrt(r1*r1 + r2*r2 + r3*r3)
	x := math.Atan(-1 * r2 / r1)
	y := math.Asin(-1 * r3 / rn)
	x = Rad2Deg(x)
	y = Rad2Deg(y)

	c = float64(coff) + (x*cfac)/0x10000
	l = float64(loff) + (y*lfac)/0x10000

	return
}

func XY2LonLat(satLon float64, c, l, coff int, cfac float64, loff int, lfac float64) (lat, lon float64) {
	subLon := Deg2Rad(satLon)

	x := float64(c-coff) * 0x10000 / cfac
	y := float64(l-loff) * 0x10000 / lfac

	sinx, cosx := math.Sincos(Deg2Rad(x))
	siny, cosy := math.Sincos(Deg2Rad(y))

	a1 := vehicleDistance * vehicleDistance * cosx * cosx * cosy * cosy
	a2 := (cosy*cosy + planetAspectRatio*siny*siny) * surfaceDistance

	if a1 < a2 {
		return 0, 0
	}

	sd := math.Sqrt(a1 - a2)
	sn := (vehicleDistance*cosx*cosy - sd) / (cosy*cosy + planetAspectRatio*siny*siny)
	s1 := vehicleDistance - sn*cosx*cosy
	s2 := sn * sinx * cosy
	s3 := -1 * sn * siny

	sxy := math.Sqrt(s1*s1 + s2*s2)

	lat = 0
	lon = 0

	if s1 == 0 {
		if s2 > 0 {
			lon = Deg2Rad(90)
		} else {
			lon = Deg2Rad(-90)
		}
		lon += subLon
	} else {
		lon = math.Atan(s2/s1) + subLon
	}

	if sxy == 0 {
		if planetAspectRatio*s3 > 0 {
			lat = Deg2Rad(90)
		} else {
			lat = Deg2Rad(-90)
		}
	} else {
		lat = math.Atan(planetAspectRatio * s3 / sxy)
	}

	return
}
