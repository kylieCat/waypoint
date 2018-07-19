package cmd

import (
	"fmt"
	"os"

	"github.com/im-auld/waypoint/waypoint"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "waypt",
	Short: "CLI for tracking versions of apps",
	Long:  "CLI for tracking versions of apps",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(cmd.Short)
	},
}

var db = waypoint.WaypointStore{
	DBFilePath: "/Users/iana/.waypt/waypt.db",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
