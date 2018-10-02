package cmd

import (
	"fmt"
	"github.com/kylie-a/waypoint/waypoint"
	"github.com/spf13/cobra"
	"os"
)

var db waypoint.DataBase

func InitDB(dbType string) {
	if dbType == "datastore" {
		db = waypoint.NewWaypointStoreDS()
	}
	if dbType == "bolt" {
		db = waypoint.NewWaypointStoreBolt()
	}
}

func checkErr(err error) {
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(2)
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
	return waypoint.Minor
}

func bumpVersion(appName string, version waypoint.Version, releaseType waypoint.ReleaseType) waypoint.Version {
	var newVersion waypoint.Version
	switch releaseType {
	case waypoint.Major:
		newVersion = version.BumpMajor()
	case waypoint.Minor:
		newVersion = version.BumpMinor()
	case waypoint.Patch:
		newVersion = version.BumpPatch()
	}
	checkErr(db.NewVersion(appName, &newVersion))
	return newVersion
}
