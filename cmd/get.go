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

	"github.com/azillion/nimbus/gfs"
	"github.com/spf13/cobra"
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
		if strings.EqualFold("ncep", args[0]) {
			return nil
		}
		return fmt.Errorf("invalid data source specified: %s", args[0])
	},
	Run: func(cmd *cobra.Command, args []string) {
		if strings.EqualFold("ncep", args[0]) {
			fmt.Println("ncep called")
			defaultParams := gfs.Params{
				Resolution: gfs.OneDegree,
				DateRange: &gfs.DateRange{
					Start: time.Now().AddDate(0, 0, -8),
					End:   time.Now(),
				},
				TimeFrame: gfs.AllTimeFrames,
			}
			gfs.GetFiles(gfs.NCEPRepoType, &defaultParams)
		} else {
			fmt.Println("get called")
		}
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
	// getCmd.Flags().StringVarP(&configFilePath, "file", "f", "", "path to data request parameters file")
}
