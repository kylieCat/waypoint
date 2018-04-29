package cmd

import (
    "github.com/spf13/cobra"
    "fmt"
    "os"
    "github.com/im-auld/waypoint/waypoint/state"
)

var rootCmd = &cobra.Command{
  Use:   "waypt",
  Short: "CLI for tracking versions of apps",
  Long: "CLI for tracking versions of apps",
  Run: func(cmd *cobra.Command, args []string) {
      fmt.Println(cmd.Short)
  },
}

var db = state.WaypointStore{
    DBFilePath: "./waypt.db",
}

func Execute() {
  if err := rootCmd.Execute(); err != nil {
    fmt.Println(err)
    os.Exit(1)
  }
}
