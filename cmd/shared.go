package cmd

import (
	"fmt"
	"github.com/kylie-a/waypoint/waypoint"
	"github.com/spf13/cobra"
	"os"
)

const (
	GREEN     = "\033[0;38;5;2m"
	YELLOW    = "\033[0;38;5;11m"
	RED       = "\033[0;38;5;9m"
	COLOR_OFF = "\033[0m"
	SUCCESS   = "SUCCESS!\n"
	ERROR_MSG = "${RED}ERROR: %s ${COLOR_OFF}\n"
)

var (
	DONE = fmt.Sprintf("%sDONE!%s", GREEN, COLOR_OFF)
)

var db waypoint.DataBase

func InitDB(conf *waypoint.Config) {
	db = waypoint.NewWaypointStoreDS(conf.Auth.Project, conf.GetAuth())
}

func printError(msg string) {
	fmt.Printf("%sERROR: %s.%s\n", RED, msg, COLOR_OFF)
}

func printWarning(msg string) {
	fmt.Printf("%sWARNING: %s.%s\n", YELLOW, msg, COLOR_OFF)
}

func checkErr(err error, exitOnError, done bool) {
	if err != nil {
		if exitOnError {
			printError(err.Error())
			os.Exit(2)
		} else {
			printWarning(err.Error())
		}
	}
	if done {
		fmt.Printf("%s\n", DONE)
	}
}

func getReleaseType(cmd *cobra.Command) waypoint.ReleaseType {
	if cmd.Flag("major").Changed {
		return waypoint.Major
	}
	if cmd.Flag("minor").Changed {
		return waypoint.Minor
	}
	if cmd.Flag("patch").Changed {
		return waypoint.Patch
	}
	if cmd.Flag("rebuild").Changed {
		return waypoint.Rebuild
	}
	return waypoint.Minor
}

func bumpVersion(appName string, version waypoint.Version, releaseType waypoint.ReleaseType) *waypoint.Version {
	var newVersion waypoint.Version
	switch releaseType {
	case waypoint.Major:
		newVersion = version.BumpMajor()
	case waypoint.Minor:
		newVersion = version.BumpMinor()
	case waypoint.Patch:
		newVersion = version.BumpPatch()
	}
	checkErr(db.NewVersion(appName, &newVersion), true, false)
	return &newVersion
}
