package gfs

import (
	"math"
	"time"
)

const (
	// NCEPRepoType get files from NCEP
	NCEPRepoType RepositoryType = "NCEP"
	// NCDCRepoType get files from NCDC
	NCDCRepoType RepositoryType = "NCDC"
)

// RepositoryType the type of repository being accessed
type RepositoryType string

// Repository interface for different NOMADS file servers
type Repository interface {
	GetBaseURL() (string, error)
	GetURIs() ([]string, error)
	GetURIsForDate(date string) ([]string, error)
	GetURIsForDateAndTime(date string, timeFrame TimeFrame) ([]string, error)
	LoadParams(*Params) error
}

// NewRepository NewRepository
func NewRepository(rt RepositoryType) Repository {
	if rt == NCEPRepoType {
		return new(NCEPRepository)
	} else if rt == NCDCRepoType {
		return nil
	}
	return nil
}

// allTimeFrames returns a slice of all time frames
func allTimeFrames() []TimeFrame {
	timeFrames := []TimeFrame{Zulu, ZeroSixHundredHours, TwelveHundredHours, EighteenHundredHours}
	return timeFrames
}

func getTimeFrames(tf TimeFrame) []TimeFrame {
	// create range of time frames
	var timeFrames []TimeFrame
	if tf == AllTimeFrames {
		timeFrames = allTimeFrames()
	} else {
		timeFrames = []TimeFrame{tf}
	}
	return timeFrames
}

func getNumberOfLoops(start, end time.Time) int {
	loops := end.Sub(start).Hours() / 24
	loops = math.Floor(loops)
	return int(loops)
}
