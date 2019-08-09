package gfs

import (
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
)

const (
	ncepBaseURLFormat string = "https://nomads.ncep.noaa.gov/cgi-bin/filter_gfs_%s.pl" // resolution
	ncepFileURIFormat string = "gfs.t%sz.pgrb2.%s.%s"                                  // time, resolution, anl/filesuffix
	ncepDirURIFormat  string = "%%2Fgfs.%s%%2F%s"                                      // date of data, time frame of data

	maxFRange int = 128
)

// Level is a region at a certain altitude
type Level struct {
	uriKey     string
	isIncluded bool
}

// ClimateVariable is a portion of GRIB2 climate data
type ClimateVariable struct {
	uriKey     string
	isIncluded bool
}

// NCEPRepository holds the data that constructs the URL
type NCEPRepository struct {
	resolution Resolution
	dateRange  DateRange
	timeFrames []TimeFrame

	// params
	levels                   map[string]Level
	levelsURICache           string
	climateVariables         map[string]ClimateVariable
	climateVariablesURICache string

	region Region

	URIs []string
}

// LoadParams reads the param object into the repository
func (ncep *NCEPRepository) LoadParams(p *Params) error {
	ncep.resolution = p.Resolution
	ncep.dateRange = p.DateRange
	ncep.timeFrames = getTimeFrames(p.TimeFrame)
	return nil
}

// GetBaseURL gets the base URL of the repository
func (ncep *NCEPRepository) GetBaseURL() (string, error) {
	if ncep.resolution == "" {
		return "", fmt.Errorf("no resolution set")
	}
	return fmt.Sprintf(ncepBaseURLFormat, ncep.resolution), nil
}

// GetURIs get the URIs
func (ncep *NCEPRepository) GetURIs() ([]string, error) {
	ncep.URIs = ncep.URIs[:0]

	if ncep.dateRange.End.Sub(time.Now()) > 0 {
		return nil, fmt.Errorf("end date can not be in the future")
	}

	loops := getNumberOfLoops(ncep.dateRange.Start, ncep.dateRange.End)
	d := ncep.dateRange.Start
	for i := 0; i < loops; i++ {
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
	ncep.URIs = ncep.URIs[:0]
	// loop through the time frames and build the URIs for each
	for _, tf := range ncep.timeFrames {
		_, err := ncep.GetURIsForDateAndTime(date, tf)
		if err != nil {
			return nil, err
		}
	}
	return ncep.URIs, nil
}

// GetURIsForDateAndTime get the URIs for a specific date and time frame
func (ncep *NCEPRepository) GetURIsForDateAndTime(date string, timeFrame TimeFrame) ([]string, error) {
	ncep.URIs = ncep.URIs[:0]

	// build .anl URI since it's unique
	anlURI := ncep.buildURI(date, timeFrame, "anl")
	ncep.URIs = append(ncep.URIs, anlURI)

	// build the f URIs
	for i := 0; i <= maxFRange; i++ {
		f := i * 3
		// build the file suffix
		suffix := FileSuffix(fmt.Sprintf("f%03d", f))

		// build the URI
		fURI := ncep.buildURI(date, timeFrame, suffix)

		// add the URI to the URIs slice
		ncep.URIs = append(ncep.URIs, fURI)
	}

	return ncep.URIs, nil
}

func (ncep *NCEPRepository) buildURI(date string, timeFrame TimeFrame, fs FileSuffix) string {
	fileURI := fmt.Sprintf(ncepFileURIFormat, timeFrame, ncep.resolution, fs)
	levelURI := ncep.getLevels()
	climateVariableURI := ncep.getClimateVariables()
	regionURI := ncep.region.ToURI()
	dirURI := fmt.Sprintf(ncepDirURIFormat, date, timeFrame)

	URI := fmt.Sprintf("?file=%s&%s&%s&%s&%s", fileURI, levelURI, climateVariableURI, regionURI, dirURI)
	logrus.Debug(URI)
	return URI
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
