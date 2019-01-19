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
	"fmt"

	"os"

	"github.com/spf13/cobra"
)

// latestCmd represents the latest command
var latestCmd = &cobra.Command{
	Use:   "latest",
	Short: "Get the latest version for an app",
	Long:  `Gets the latest version for the provided appplication name. This gets the latest version by date.`,
	Run: func(cmd *cobra.Command, args []string) {
		latest, err := ws.GetLatest(conf.App)
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(2)
		}
		if cmd.Flag("commit").Changed {
			fmt.Println(latest.CommitHash)
			os.Exit(0)
		}
		fmt.Println(latest.SemVer())
	},
}

func init() {
	latestCmd.Flags().Bool("commit", false, "get commit instead")
	rootCmd.AddCommand(latestCmd)
}
