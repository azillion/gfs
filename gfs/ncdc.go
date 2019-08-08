package gfs

import (
	"fmt"
	"time"
)

const (
	baseURLFormat string = "https://nomads.ncdc.noaa.gov/data/"
	fileURIFormat string = "gfs.t%sz.pgrb2.%s.%s" // time, resolution, anl/filesuffix
	dirURIFormat  string = "%%2Fgfs.%s%%2F%s"     // date of data, time frame of data
)

// FileSuffix is the final part of the filename
type FileSuffix string

// NCDCRepository holds the data that constructs the URL
type NCDCRepository struct {
	dateRange                  DateRange
	isAdditionalPrecipIncluded bool

	URIs []string
}

// GetBaseURL gets the base URL of the repository
func (ncdc *NCDCRepository) GetBaseURL() (string, error) {
	return fmt.Sprintf(baseURLFormat, ncdc.resolution), nil
}

// GetURIs get the URIs
func (ncdc *NCDCRepository) GetURIs() ([]string, error) {
	if ncdc.dateRange.End.Sub(time.Now()) > 0 {
		panic("end date can not be in the future")
	}

	d := ncdc.dateRange.Start
	end := ncdc.dateRange.End
	for end.Sub(d).Hours() > 24 {
		date := d.Format("20060102")
		_, err := ncdc.GetURIsForDate(date)
		if err != nil {
			return nil, err
		}
		d = d.Add(time.Hour * 24)
	}

	return ncdc.URIs, nil
}

// GetURIsForDate get the URIs for a specific date
func (ncdc *NCDCRepository) GetURIsForDate(date string) ([]string, error) {
	return ncdc.GetURIsForDateAndTime(date, AllTimeFrames)
}

// GetURIsForDateAndTime get the URIs for a specific date and time frame
func (ncdc *NCDCRepository) GetURIsForDateAndTime(date string, timeFrame TimeFrame) ([]string, error) {
	if _, err := ncdc.GetBaseURL(); err != nil {
		return nil, fmt.Errorf("failed to get the base url")
	}

	if len(ncdc.URIs) > 0 {
		return ncdc.URIs, nil
	}

	// build .anl URI since it's unique
	anlURI := ncdc.buildURI(date, timeFrame, "anl")
	ncdc.URIs = append(ncdc.URIs, anlURI)

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
			fURI := ncdc.buildURI(date, tf, suffix)

			// add the URI to the URIs slice
			ncdc.URIs = append(ncdc.URIs, fURI)
		}
	}

	return ncdc.URIs, nil
}

func (ncdc *NCDCRepository) buildURI(date string, timeFrame TimeFrame, fs FileSuffix) string {
	fileURI := fmt.Sprintf(fileURIFormat, timeFrame, ncdc.resolution, fs)
	levelURI := ncdc.getLevels()
	climateVariableURI := ncdc.getClimateVariables()
	regionURI := ncdc.region.ToURI()
	dirURI := fmt.Sprintf(dirURIFormat, date, timeFrame)
	URI := fmt.Sprintf("?file=%s&%s&%s&%s&%s", fileURI, levelURI, climateVariableURI, regionURI, dirURI)
	// fmt.Println(URI)
	return URI
}

// LoadParams reads the param object into the repository
func (ncdc *NCDCRepository) LoadParams(p *Params) error {
	ncdc.dateRange = p.DateRange

	return nil
}
