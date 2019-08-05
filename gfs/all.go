package gfs

import (
	"flag"
	"fmt"
	"io/ioutil"
	"math"
	"net/http"
	"os"
	"strings"
	"time"
)

const (
	inputDateLayout = "2006-01-02"
	urlDateLayout   = "20060102"
)

var (
	startDate    string
	endDate      string
	outputFolder string
	baseURL1     = "https://nomads.ncep.noaa.gov/cgi-bin/filter_gfs_1p00.pl?file=gfs.t"
	baseURL2     = "z.pgrb2.1p00."
	baseURL3     = "&all_lev=on&all_var=on&leftlon=0&rightlon=360&toplat=90&bottomlat=-90&dir=%2Fgfs."
)

func init() {
	flag.StringVar(&startDate, "b", "2006-01-02", "begin date <YYYY-MM-DD>")
	flag.StringVar(&endDate, "e", "2014-01-02", "end date <YYYY-MM-DD>")
	flag.StringVar(&outputFolder, "o", "./", "output folder")
}

func main() {
	flag.Parse()

	start, err := time.Parse(inputDateLayout, startDate)
	if err != nil {
		panic(err)
	}
	dataTime := start

	end, err := time.Parse(inputDateLayout, endDate)
	if err != nil {
		panic(err)
	}
	if end.Sub(time.Now()) > 0 {
		panic("end date can not be in the future")
	}

	err = os.Chdir(outputFolder)
	if err != nil {
		panic(err)
	}

	diff := end.Sub(start)
	days := diff.Hours() / 24
	if days < 1.0 {
		days = 1.0
	}
	days = math.Floor(days)
	loops := int(days) * 4

	regions := []string{"anl",
		"f000", "f003", "f006", "f009",
		"f012", "f015", "f018", "f021",
		"f024", "f027", "f030", "f033",
		"f036", "f039", "f042", "f045",
		"f048", "f051", "f054", "f057"}

	for i := 0; i < loops; i++ {
		for j := 0; j < len(regions); j++ {
			data := getData(dataTime, regions[j])
			fileName := formatFileName(dataTime, regions[j])
			saveData(fileName, data)
			dataTime = increTime(dataTime)
		}
	}
}

func saveData(fileName string, data []byte) {
	saveFile, err := os.Create(fileName)
	if err != nil {
		panic(err)
	}

	_, err = saveFile.Write(data)
	if err != nil {
		panic(err)
	}

	err = saveFile.Sync()
	if err != nil {
		panic(err)
	}

	err = saveFile.Close()
	if err != nil {
		panic(err)
	}
	fmt.Printf("saved %s\n", saveFile.Name())
}

func getData(t time.Time, s string) []byte {
	reqURL := formatURL(t, s)
	resp, err := http.Get(reqURL)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	return body
}

func formatURL(t time.Time, s string) string {
	var urlParts []string
	urlParts = append(urlParts, baseURL1)
	urlParts = append(urlParts, t.Format("15"))
	urlParts = append(urlParts, baseURL2)
	urlParts = append(urlParts, s)
	urlParts = append(urlParts, baseURL3)
	urlParts = append(urlParts, t.Format(urlDateLayout))
	urlParts = append(urlParts, "%2F")
	urlParts = append(urlParts, t.Format("15"))
	return strings.Join(urlParts, "")
}

func formatFileName(t time.Time, s string) string {
	var urlParts []string
	urlParts = append(urlParts, "gfs")
	urlParts = append(urlParts, t.Format("2006010215"))
	urlParts = append(urlParts, s)
	return strings.Join(urlParts, ".")
}

func increTime(t time.Time) time.Time {
	return t.Add(time.Hour * 6)
}
