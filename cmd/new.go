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

// newCmd represents the new command
var newCmd = &cobra.Command{
	Use:   "new",
	Short: "Add a new application to waypoint",
	Long: "Add a new application to waypoint",
	Run: func(cmd *cobra.Command, args []string) {
		initial := cmd.Flag("initial").Value.String()
        err := db.AddApplication(args[0], initial)
        if err != nil {
           fmt.Errorf(err.Error())
        }
        fmt.Printf("Added app %s and set initial version to %s", args[0], initial)
	},
}

func init() {
	rootCmd.AddCommand(newCmd)
	newCmd.Flags().String("initial", "0.1.0", "Set the initial version for an app")
}
