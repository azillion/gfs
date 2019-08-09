package gfs

import "fmt"

// Region contains the region of the data
type Region struct {
	LeftLon   float32
	RightLon  float32
	TopLat    float32
	BottomLat float32
}

// ToURI returns a URI string of the region
func (r *Region) ToURI() string {
	return fmt.Sprintf("leftlon=%1.2f&rightlon=%1.2f&toplat=%1.2f&bottomlat=%1.2f", r.LeftLon, r.RightLon, r.TopLat, r.BottomLat)
}

// FullEarthRegion returns a region the full size of Earth
func FullEarthRegion() *Region {
	return &Region{
		LeftLon:   0.0,
		RightLon:  360.0,
		TopLat:    90.0,
		BottomLat: -90.0,
	}
}
