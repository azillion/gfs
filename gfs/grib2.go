package gfs

import (
	"time"
)

// GRIB2 is simplified GRIB2 file structure
type GRIB2 struct {
	RefTime     time.Time
	VerfTime    time.Time
	Name        string
	Description string
	Unit        string
	Level       string
	Values      []Value
}

// Value is data item of GRIB2 file
type Value struct {
	Longitude float64
	Latitude  float64
	Value     float32
}
