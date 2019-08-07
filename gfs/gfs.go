package gfs

import (
	"fmt"
	"time"
)

const (
	// NCEPRepoType get files from NCEP
	NCEPRepoType RepositoryType = "NCEP"
	// NCDCRepoType get files from NCDC
	NCDCRepoType RepositoryType = "NCDC"

	// OneDegree 1.0 Degree of Longitudinal Resolution
	OneDegree Resolution = "1p00"
	// ZeroPointFiveDegree 0.5 Degrees of Longitudinal Resolution
	ZeroPointFiveDegree Resolution = "0p50"
	// ZeroPointTwoFiveDegree 0.25 Degrees of Longitudinal Resolution
	ZeroPointTwoFiveDegree Resolution = "0p25"

	// Zulu midnight - 0000
	Zulu TimeFrame = "00"
	// ZeroSixHundredHours 6 AM - 0600
	ZeroSixHundredHours TimeFrame = "06"
	// TwelveHundredHours noon - 1200
	TwelveHundredHours TimeFrame = "12"
	// EighteenHundredHours 6 PM - 1800
	EighteenHundredHours TimeFrame = "18"
	// AllTimeFrames self explanitory
	AllTimeFrames TimeFrame = "99"
)

// RepositoryType the type of repository being accessed
type RepositoryType string

// Resolution is the degree of resolution for the GFS data
type Resolution string

// DateRange is the range of files to download
type DateRange struct {
	Start time.Time
	End   time.Time
}

// LoadFromStrings LoadFromStrings
func (dr *DateRange) LoadFromStrings(s, e string) error {
	ss, err := time.Parse("2006-01-02", s)
	if err != nil {
		return err
	}
	dr.Start = ss
	ee, err := time.Parse("2006-01-02", e)
	if err != nil {
		return err
	}
	dr.End = ee
	return nil
}

// DateRangeStrings unmarshall the config file into strings
type DateRangeStrings struct {
	Start string `mapstructure:"start"`
	End   string `mapstructure:"end"`
}

// TimeFrame is the time the data was recorded
type TimeFrame string

// Params used when downloading grib2 GFS files
type Params struct {
	RepositoryType RepositoryType `mapstructure:"repository_type"`
	Resolution     Resolution     `mapstructure:"resolution"`
	DateRange      DateRange
	TimeFrame      TimeFrame `mapstructure:"time_frame"`
}

// Repository interface for different NOMADS file servers
type Repository interface {
	GetBaseURL() (string, error)
	GetURIs() ([]string, error)
	GetURIsForDate(date string) ([]string, error)
	GetURIsForDateAndTime(date string, timeFrame TimeFrame) ([]string, error)
	LoadParams(*Params) error
}

// Service holds the repository and params of the service
type Service struct {
	repository Repository
	params     *Params
}

// GetFiles get files from NOMADS
func (s *Service) GetFiles() error {
	if s.params.TimeFrame == AllTimeFrames {
		_, err := s.repository.GetURIs()
		if err != nil {
			panic(err)
		}
	}
	return nil
}

// NewService creates a new gfs service
func NewService(p *Params) *Service {
	r := getRepository(p.RepositoryType)
	if r == nil {
		panic(fmt.Errorf("no repository of that type"))
	}
	err := r.LoadParams(p)
	if err != nil {
		panic(fmt.Errorf("error loading params: %v", err))
	}
	return &Service{
		repository: r,
		params:     p,
	}
}

func getRepository(rt RepositoryType) Repository {
	if rt == NCEPRepoType {
		return new(NCEPRepository)
	} else if rt == NCDCRepoType {
		return nil
	}
	return nil
}

// TODO: Reimplement below

// func init() {
// flag.StringVar(&startDate, "b", "2006-01-02", "begin date <YYYY-MM-DD>")
// flag.StringVar(&endDate, "e", "2014-01-02", "end date <YYYY-MM-DD>")
// flag.StringVar(&outputFolder, "o", "./", "output folder")
// }

// func main() {
// flag.Parse()

// start, err := time.Parse(inputDateLayout, startDate)
// if err != nil {
// 	panic(err)
// }
// dataTime := start

// end, err := time.Parse(inputDateLayout, endDate)
// if err != nil {
// 	panic(err)
// }
// if end.Sub(time.Now()) > 0 {
// 	panic("end date can not be in the future")
// }

// err = os.Chdir(outputFolder)
// if err != nil {
// 	panic(err)
// }

// diff := end.Sub(start)
// days := diff.Hours() / 24
// if days < 1.0 {
// 	days = 1.0
// }
// days = math.Floor(days)
// loops := int(days) * 4

// regions := []string{"anl",
// 	"f000", "f003", "f006", "f009",
// 	"f012", "f015", "f018", "f021",
// 	"f024", "f027", "f030", "f033",
// 	"f036", "f039", "f042", "f045",
// 	"f048", "f051", "f054", "f057"}

// for i := 0; i < loops; i++ {
// 	for j := 0; j < len(regions); j++ {
// 		data := getData(dataTime, regions[j])
// 		fileName := formatFileName(dataTime, regions[j])
// 		saveData(fileName, data)
// 		dataTime = increTime(dataTime)
// 	}
// }
// }

// func saveData(fileName string, data []byte) {
// 	saveFile, err := os.Create(fileName)
// 	if err != nil {
// 		panic(err)
// 	}

// 	_, err = saveFile.Write(data)
// 	if err != nil {
// 		panic(err)
// 	}

// 	err = saveFile.Sync()
// 	if err != nil {
// 		panic(err)
// 	}

// 	err = saveFile.Close()
// 	if err != nil {
// 		panic(err)
// 	}
// 	fmt.Printf("saved %s\n", saveFile.Name())
// }

// func getData(t time.Time, s string) []byte {
// 	reqURL := formatURL(t, s)
// 	resp, err := http.Get(reqURL)
// 	if err != nil {
// 		panic(err)
// 	}
// 	defer resp.Body.Close()

// 	body, err := ioutil.ReadAll(resp.Body)
// 	if err != nil {
// 		panic(err)
// 	}

// 	return body
// }

// func formatURL(t time.Time, s string) string {
// 	var urlParts []string
// 	urlParts = append(urlParts, baseURL1)
// 	urlParts = append(urlParts, t.Format("15"))
// 	urlParts = append(urlParts, baseURL2)
// 	urlParts = append(urlParts, s)
// 	urlParts = append(urlParts, baseURL3)
// 	urlParts = append(urlParts, t.Format(urlDateLayout))
// 	urlParts = append(urlParts, "%2F")
// 	urlParts = append(urlParts, t.Format("15"))
// 	return strings.Join(urlParts, "")
// }

// func formatFileName(t time.Time, s string) string {
// 	var urlParts []string
// 	urlParts = append(urlParts, "gfs")
// 	urlParts = append(urlParts, t.Format("2006010215"))
// 	urlParts = append(urlParts, s)
// 	return strings.Join(urlParts, ".")
// }
