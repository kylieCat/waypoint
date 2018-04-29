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
	"os"
	"fmt"
	"github.com/im-auld/waypoint/waypoint/state"
)

// bumpCmd represents the bump command
var bumpCmd = &cobra.Command{
	Use:   "bump",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		version, err := db.GetMostRecent(args[0])
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(2)
		}
		var newVersion state.Version
		if cmd.Flag("major").Changed {
			newVersion = version.BumpMajor()
		}
		if cmd.Flag("minor").Changed {
			newVersion = version.BumpMinor()
		}
		if cmd.Flag("patch").Changed {
			newVersion = version.BumpPatch()
		}
		db.NewVersion(args[0], &newVersion)
		fmt.Println(newVersion.SemVer())
	},
}

func init() {
	rootCmd.AddCommand(bumpCmd)
	bumpCmd.Flags().Bool("major", false, "Bump the verisn up by one")
	bumpCmd.Flags().Bool("minor", false, "Bump the verisn up by one")
	bumpCmd.Flags().Bool("patch", false, "Bump the verisn up by one")
}
