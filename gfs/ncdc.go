package gfs

import (
	"fmt"
	"time"
)

const (
	ncdcBaseURLFormat string = "https://nomads.ncdc.noaa.gov/data/"
	// fileURIFormat     string = "gfs.t%sz.pgrb2.%s.%s" // time, resolution, anl/filesuffix
	// dirURIFormat      string = "%%2Fgfs.%s%%2F%s"     // date of data, time frame of data
)

// NCDCRepository holds the data that constructs the URL
type NCDCRepository struct {
	dateRange                  DateRange
	timeFrames                 []TimeFrame
	isAdditionalPrecipIncluded bool

	URIs []string
}

// LoadParams reads the param object into the repository
func (ncdc *NCDCRepository) LoadParams(p *Params) error {
	ncdc.dateRange = p.DateRange
	ncdc.timeFrames = getTimeFrames(p.TimeFrame)
	ncdc.isAdditionalPrecipIncluded = p.IsAdditionalPrecipIncluded

	return nil
}

// GetBaseURL gets the base URL of the repository
func (ncdc *NCDCRepository) GetBaseURL() (string, error) {
	return fmt.Sprintf(ncdcBaseURLFormat), nil
}

// GetURIs get the URIs
func (ncdc *NCDCRepository) GetURIs() ([]string, error) {
	ncdc.URIs = ncdc.URIs[:0]

	if ncdc.dateRange.End.Sub(time.Now()) > 0 {
		panic("end date can not be in the future")
	}

	loops := getNumberOfLoops(ncdc.dateRange.Start, ncdc.dateRange.End)
	d := ncdc.dateRange.Start
	for i := 0; i < loops; i++ {
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
	ncdc.URIs = ncdc.URIs[:0]
	// loop through the time frames and build the URIs for each
	for _, tf := range ncdc.timeFrames {
		_, err := ncdc.GetURIsForDateAndTime(date, tf)
		if err != nil {
			return nil, err
		}
	}
	return ncdc.URIs, nil
}

// GetURIsForDateAndTime get the URIs for a specific date and time frame
func (ncdc *NCDCRepository) GetURIsForDateAndTime(date string, timeFrame TimeFrame) ([]string, error) {
	ncdc.URIs = ncdc.URIs[:0]
	if _, err := ncdc.GetBaseURL(); err != nil {
		return nil, fmt.Errorf("failed to get the base url")
	}

	// build .anl URI since it's unique
	anlURI := ncdc.buildURI(date, timeFrame, "anl")
	ncdc.URIs = append(ncdc.URIs, anlURI)

	// build the f URIs
	for i := 0; i <= maxFRange; i++ {
		f := i * 3
		// build the file suffix
		suffix := FileSuffix(fmt.Sprintf("f%03d", f))

		// build the URI
		fURI := ncdc.buildURI(date, timeFrame, suffix)

		// add the URI to the URIs slice
		ncdc.URIs = append(ncdc.URIs, fURI)
	}

	return ncdc.URIs, nil
}

func (ncdc *NCDCRepository) buildURI(date string, timeFrame TimeFrame, fs FileSuffix) string {
	// fileURI := fmt.Sprintf(fileURIFormat, timeFrame, ncdc.resolution, fs)
	// dirURI := fmt.Sprintf(dirURIFormat, date, timeFrame)
	// URI := fmt.Sprintf("?file=%s&%s&%s&%s&%s", fileURI, dirURI)
	// fmt.Println(URI)
	// return URI
	return ""
}
