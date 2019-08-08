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
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/azillion/nimbus/gfs"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// getCmd represents the get command
var getCmd = &cobra.Command{
	Use:   "get [data source]",
	Short: "Get files from NOMADS.",
	Long:  `Download files from NOMADS`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("requires a data source argument")
		}
		if strings.EqualFold("gfs", args[0]) {
			return nil
		}
		return fmt.Errorf("invalid data source specified: %s", args[0])
	},
	Run: func(cmd *cobra.Command, args []string) {
		if strings.EqualFold("gfs", args[0]) {
			handleGFSDataSource()
			return
		}

		logrus.Infof("no valid data source provided: %s", args[0])
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
	getCmd.Flags().StringVarP(&cfgFile, "config", "c", "", "config file (default is ./get-config.nimbus.yaml)")
}

func parseConfigFile() (*gfs.Params, error) {
	params := gfs.Params{
		RepositoryType: gfs.NCEPRepoType,
		Resolution:     gfs.OneDegree,
		DateRange: gfs.DateRange{
			Start: time.Now().AddDate(0, 0, -8),
			End:   time.Now(),
		},
		TimeFrame: gfs.AllTimeFrames,
		IsAdditionalPrecipIncluded: false,
	}
	err := viper.Unmarshal(&params)
	if err != nil {
		return nil, err
	}

	var dateRangeStrings gfs.DateRangeStrings
	err = viper.UnmarshalKey("date_range", &dateRangeStrings)
	if err != nil {
		return nil, err
	}

	params.DateRange.LoadFromStrings(dateRangeStrings.Start, dateRangeStrings.End)
	logrus.Debug(params)

	return &params, nil
}

func handleGFSDataSource() error {
	params, err := parseConfigFile()
	if err != nil {
		return err
	}
	gfsService := gfs.NewService(params)
	gfsService.GetFiles()

	return nil
}
