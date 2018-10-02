// Copyright © 2018 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"github.com/spf13/cobra"
	"fmt"
)

var (
	appName, releaseType, target string
)

func printError()     {}
func printWarning()   {}
func checkError()     {}
func getVersion()     {}
func buildImage()     {}
func cleanup()        {}
func packageChart()   {}
func updateHelmRepo() {}

// deployCmd represents the deploy command
var releaseCmd = &cobra.Command{
	Use:   "release",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		app := args[0]
		prevVer, err := db.GetMostRecent(app)
		checkErr(err)
		releaseType := getReleaseType(cmd)
		newVer := bumpVersion(appName, *prevVer, releaseType)
		fmt.Println(newVer)
	},
}

func init() {
	rootCmd.AddCommand(releaseCmd)
	releaseCmd.Flags().StringVar(&appName, "app", "", "Application to release")
	releaseCmd.Flags().StringVar(&releaseType, "type", "", "Release type")
	releaseCmd.Flags().StringVar(&target, "target", "", "Target env")
}
