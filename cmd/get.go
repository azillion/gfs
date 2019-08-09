/*
Package cmd commands for nimbus
Copyright © 2019 Alexander Zillion <alex@alexzillion.com>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"strings"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/azillion/nimbus/gfs"
	"github.com/azillion/nimbus/util"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var outputFolder string

var defaultParams = &gfs.Params{
	RepositoryType: gfs.NCEPRepoType,
	Resolution:     gfs.OneDegree,
	DateRange: gfs.DateRange{
		Start: time.Now().AddDate(0, 0, -8),
		End:   time.Now(),
	},
	TimeFrame:                  gfs.AllTimeFrames,
	IsAdditionalPrecipIncluded: false,
}

// getCmd represents the get command
var getCmd = &cobra.Command{
	Use:   "get [config file]",
	Short: "Get files from NOMADS.",
	Long:  `Download files from NOMADS`,
	Args:  cobra.ExactArgs(1),
	PreRun: func(cmd *cobra.Command, args []string) {
		if debug {
			logrus.SetLevel(logrus.DebugLevel)
		}
		if len(args) < 1 {
			logrus.Fatal("requires a config file")
		}
		if util.FileExists(args[0]) {
			cfgFile = args[0]
			initConfig()
		} else {
			logrus.Fatal("config file does not exist or is unreadable: %s", args[0])
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		dataSource := viper.GetString("data_source")
		logrus.Debug(dataSource)
		if strings.EqualFold(dataSource, "gfs") {
			err := handleGFSDataSource()
			if err != nil {
				logrus.Fatal(err)
			}
			return
		}
		logrus.Infof("no valid config file provided: %s", cfgFile)
	},
}

func init() {
	rootCmd.AddCommand(getCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// getCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// getCmd.Flags().StringVarP(&cfgFile, "config", "c", "", "config file (required)")
	// getCmd.MarkFlagRequired("config")

	getCmd.Flags().StringVarP(&outputFolder, "output-folder", "o", "", "output folder (default is working directory)")
	viper.BindPFlag("output_folder", getCmd.Flags().Lookup("output-folder"))
}

func parseConfigFile() (*gfs.Params, error) {
	// unmarshal the config file over the default params
	params := *defaultParams
	err := viper.Unmarshal(&params)
	if err != nil {
		return nil, err
	}

	// parse the date values to strings
	var dateRangeStrings gfs.DateRangeStrings
	err = viper.UnmarshalKey("date_range", &dateRangeStrings)
	if err != nil {
		return nil, err
	}

	// convert the date range strings to time.Time and load them into the params
	params.DateRange.LoadFromStrings(dateRangeStrings.Start, dateRangeStrings.End)
	logrus.Debug(params)

	return &params, nil
}

func handleGFSDataSource() error {
	// parse the config file
	params, err := parseConfigFile()
	if err != nil {
		return err
	}
	logrus.Debug("parsed the config file")

	// create a new gfs service
	gfsService := gfs.NewService(params)
	logrus.Debug("created a new GFS service")

	gfsService.GetFiles()

	return nil
}
