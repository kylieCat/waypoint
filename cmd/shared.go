package cmd

import (
	"fmt"
	"os"

	"github.com/kylie-a/waypoint/pkg"
	"github.com/kylie-a/waypoint/pkg/db"
	"github.com/spf13/cobra"
)

const (
	GREEN     = "\033[0;38;5;2m"
	YELLOW    = "\033[0;38;5;11m"
	RED       = "\033[0;38;5;9m"
	COLOR_OFF = "\033[0m"
	SUCCESS   = "Success!\n"
	ERROR_MSG = "${Red}ERROR: %s ${colorOff}\n"
)

var (
	DONE = fmt.Sprintf("%sDONE!%s", GREEN, COLOR_OFF)
)

var storage pkg.IStorage

func InitDB(conf *pkg.Config) {
	var err error

	if storage, err = db.NewClient(conf); err != nil {
		fmt.Println(err.Error())
	}
}

func colorPrint(color, msg string) string {
	return fmt.Sprintf("%s%s%s", color, msg, COLOR_OFF)
}

func green(msg string) string {
	return colorPrint(GREEN, msg)
}

func yellow(msg string) string {
	return colorPrint(YELLOW, msg)
}

func red(msg string) string {
	return colorPrint(RED, msg)
}

func printError(msg string) {
	msg = "ERROR: " + msg
	fmt.Println(red(msg))
}

func printWarning(msg string) {
	msg = "WARNING: " + msg
	fmt.Println(yellow(msg))
}

func printSkipping() {
	fmt.Println(yellow("SKIPPING: step ShouldExecute returned false"))
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
		fmt.Println(green("DONE!"))
	}
}

func getReleaseType(cmd *cobra.Command) pkg.ReleaseType {
	if cmd.Flag("major").Changed {
		return pkg.Major
	}
	if cmd.Flag("minor").Changed {
		return pkg.Minor
	}
	if cmd.Flag("patch").Changed {
		return pkg.Patch
	}
	if cmd.Flag("rebuild").Changed {
		return pkg.Rebuild
	}
	return pkg.Minor
}

func bumpVersion(appName string, version pkg.Version, releaseType pkg.ReleaseType) *pkg.Version {
	var newVersion pkg.Version
	switch releaseType {
	case pkg.Major:
		newVersion = version.BumpMajor()
	case pkg.Minor:
		newVersion = version.BumpMinor()
	case pkg.Patch:
		newVersion = version.BumpPatch()
	}
	checkErr(storage.Save(appName, &newVersion), true, false)
	return &newVersion
}
