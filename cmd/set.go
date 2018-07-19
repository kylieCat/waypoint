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

	"github.com/im-auld/waypoint/waypoint"
	"github.com/spf13/cobra"
)

// setCmd represents the set command
var setCmd = &cobra.Command{
	Use:   "set",
	Short: "Set the most recent semantic version for an application.",
	Long:  `Set the most recent semantic version for an application. Defaults to 1.0.0`,
	Run: func(cmd *cobra.Command, args []string) {
		parts, err := waypoint.GetPartsFromSemVer(cmd.Flag("semver").Value.String())
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(2)
		}
		version := waypoint.NewVersion(parts[waypoint.MAJOR], parts[waypoint.MINOR], parts[waypoint.PATCH])
		err = db.NewVersion(args[0], &version)
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(2)
		}
		fmt.Printf("Added version: %s for app %s\n", version.SemVer(), args[0])
	},
}

func init() {
	rootCmd.AddCommand(setCmd)
	setCmd.Flags().String("semver", "1.0.0", "The semver for the app. Defaults to 1.0.0")
}
