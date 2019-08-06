package gfs

import (
	"fmt"
	"time"
)

const (
	baseURLFormat string = "https://nomads.ncep.noaa.gov/cgi-bin/filter_gfs_%s.pl" // resolution
	fileURIFormat string = "gfs.t%sz.pgrb2.%s.%s"                                  // time, resolution, anl/filesuffix
	dirURIFormat  string = "%%2Fgfs.%s%%2F%s"                                      // date of data, time frame of data

	maxFRange int = 128
)

type levelKey string

// Level is a region at a certain altitude
type Level struct {
	uriKey     string
	isIncluded bool
}

type climateVariableKey string

// ClimateVariable is a portion of GRIB2 climate data
type ClimateVariable struct {
	uriKey     string
	isIncluded bool
}

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

// FileSuffix is the final part of the filename
type FileSuffix string

// NCEPRepository holds the data that constructs the URL
type NCEPRepository struct {
	resolution Resolution
	dateRange  DateRange

	// params
	levels                   map[levelKey]Level
	levelsURICache           string
	climateVariables         map[climateVariableKey]ClimateVariable
	climateVariablesURICache string

	region Region

	URIs []string
}

// GetBaseURL gets the base URL of the repository
func (ncep *NCEPRepository) GetBaseURL() (string, error) {
	if ncep.resolution == "" {
		return "", fmt.Errorf("no resolution set")
	}
	return fmt.Sprintf(baseURLFormat, ncep.resolution), nil
}

// GetURIs get the URIs
func (ncep *NCEPRepository) GetURIs() ([]string, error) {
	if ncep.dateRange.End.Sub(time.Now()) > 0 {
		panic("end date can not be in the future")
	}

	d := ncep.dateRange.Start
	end := ncep.dateRange.End
	for end.Sub(d).Hours() > 24 {
		date := d.Format("20060102")
		_, err := ncep.GetURIsForDate(date)
		if err != nil {
			return nil, err
		}
		d = d.Add(time.Hour * 24)
	}

	return ncep.URIs, nil
}

// GetURIsForDate get the URIs for a specific date
func (ncep *NCEPRepository) GetURIsForDate(date string) ([]string, error) {
	return ncep.GetURIsForDateAndTime(date, AllTimeFrames)
}

// GetURIsForDateAndTime get the URIs for a specific date and time frame
func (ncep *NCEPRepository) GetURIsForDateAndTime(date string, timeFrame TimeFrame) ([]string, error) {
	if _, err := ncep.GetBaseURL(); err != nil {
		return nil, fmt.Errorf("failed to get the base url")
	}

	if len(ncep.URIs) > 0 {
		return ncep.URIs, nil
	}

	// build .anl URI since it's unique
	anlURI := ncep.buildURI(date, timeFrame, "anl")
	ncep.URIs = append(ncep.URIs, anlURI)

	// create range of time frames
	var timeFrames []TimeFrame
	if timeFrame == AllTimeFrames {
		timeFrames = append(timeFrames, Zulu)
		timeFrames = append(timeFrames, ZeroSixHundredHours)
		timeFrames = append(timeFrames, TwelveHundredHours)
		timeFrames = append(timeFrames, EighteenHundredHours)
	} else {
		timeFrames = append(timeFrames, timeFrame)
	}

	// loop through the time frames and build the URIs for each
	for _, tf := range timeFrames {
		// build the f URIs
		for i := 0; i <= maxFRange; i++ {
			f := i * 3
			// build the file suffix
			var suffix FileSuffix
			suffix = FileSuffix(fmt.Sprintf("f%03d", f))

			// build the URI
			fURI := ncep.buildURI(date, tf, suffix)

			// add the URI to the URIs slice
			ncep.URIs = append(ncep.URIs, fURI)
		}
	}

	return ncep.URIs, nil
}

func (ncep *NCEPRepository) getLevels() string {
	if ncep.levelsURICache != "" {
		return ncep.levelsURICache
	}

	if len(ncep.levels) == 0 {
		ncep.levelsURICache = "all_lev=on"
		return ncep.levelsURICache
	}
	// TODO: Read levels from csv file
	// TODO: Read config for which levels to enable
	return ""
}

func (ncep *NCEPRepository) getClimateVariables() string {
	if ncep.climateVariablesURICache != "" {
		return ncep.climateVariablesURICache
	}

	if len(ncep.climateVariables) == 0 {
		ncep.climateVariablesURICache = "all_var=on"
		return ncep.climateVariablesURICache
	}
	// TODO: Read vars from csv file
	// TODO: Read config for which vars to enable
	return ""
}

func (ncep *NCEPRepository) buildURI(date string, timeFrame TimeFrame, fs FileSuffix) string {
	fileURI := fmt.Sprintf(fileURIFormat, timeFrame, ncep.resolution, fs)
	levelURI := ncep.getLevels()
	climateVariableURI := ncep.getClimateVariables()
	regionURI := ncep.region.ToURI()
	dirURI := fmt.Sprintf(dirURIFormat, date, timeFrame)
	URI := fmt.Sprintf("?file=%s&%s&%s&%s&%s", fileURI, levelURI, climateVariableURI, regionURI, dirURI)
	fmt.Println(URI)
	return URI
}

// LoadParams reads the param object into the repository
func (ncep *NCEPRepository) LoadParams(p *Params) error {
	ncep.resolution = p.Resolution
	ncep.dateRange = *p.DateRange
	return nil
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
