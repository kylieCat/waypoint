package cmd

import (
	"fmt"
	"os"

	"github.com/kylie-a/waypoint/pkg"
	"github.com/spf13/cobra"
)

var (
	cfgFile string
	conf    *pkg.Config
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "pkg",
	Short: "A command line interface for deploying SRE services",
	Long:  `A command line interface for deploying SRE services`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	var err error
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", ".waypoint.yaml", "config file (default is .waypoint.yaml)")
	if conf, err = pkg.GetConf(cfgFile); err != nil {
		fmt.Println(err.Error())
		os.Exit(2)
	}
	InitDB(conf)
}
