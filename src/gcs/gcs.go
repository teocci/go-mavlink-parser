// Package gcs
// Created by RTT.
// Author: teocci@yandex.com on 2021-Sep-06
package gcs

import (
	"log"
	"math"
)

const (
	MinSouth = 35.61469485754825
	MinWest  = 127.26129834574632

	MaxNorth = 36.03345205933448
	MaxEast  = 127.6288468798199
)

// SCS represents a spherical coordinate system using latitude, longitude.
type SCS struct {
	Lat float64
	Lon float64
}

// Delta represents de difference between two SCS coordinates.
type Delta struct {
	Lat float64
	Lon float64
}

func (c SCS) LatInRange() bool {
	return c.Lat >= MinSouth && c.Lat <= MaxNorth
}

func (c SCS) LonInRange() bool {
	return c.Lon >= MinWest && c.Lon <= MaxEast
}

func (c SCS) InRange() bool {
	return c.LatInRange() && c.LonInRange()
}

func (c SCS) Delta(r SCS) Delta {
	return Delta{
		Lat: c.Lat - r.Lat,
		Lon: c.Lon - r.Lon,
	}
}

func (c SCS) CentiMetersTo(r SCS) MetricLength {
	return distance(c, r, CM)
}

func (c SCS) MetersTo(r SCS) MetricLength {
	return Distance(c, r)
}

func (c SCS) Equals(r SCS) bool {
	return c.Lat == r.Lat && c.Lon == r.Lon
}

func (c SCS) toRadians() SCS {
	return SCS{
		Lat: degreesToRadians(c.Lat),
		Lon: degreesToRadians(c.Lon),
	}
}

func Distance(c, r SCS) MetricLength {
	return distance(c, r, M)
}

func distance(orig, dest SCS, unit MetricLength) MetricLength {
	orig = orig.toRadians()
	dest = dest.toRadians()

	var c float64
	if orig.Equals(dest) {
		c = 0
	} else {
		delta := orig.Delta(dest)

		a := math.Pow(math.Sin(delta.Lat/2), 2) + math.Cos(orig.Lat)*math.Cos(dest.Lat)*math.Pow(math.Sin(delta.Lon/2), 2)
		c = 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	}

	var dist MetricLength
	switch unit {
	case CM:
		dist = MetricLength(c) * earthRadiusCm
	case M:
		dist = MetricLength(c) * earthRadiusM
	case KM:
		dist = MetricLength(c) * earthRadiusKm
	case Mi:
		dist = MetricLength(c) * earthRadiusMi
	case NM:
		dist = MetricLength(c) * earthRadiusNM
	default:
		log.Fatal("metric unit not defined")
	}

	return dist
}

// degreesToRadians converts from degrees to radians.
func degreesToRadians(d float64) float64 {
	return d * math.Pi / 180
}
