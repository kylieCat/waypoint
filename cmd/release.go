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
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
	"github.com/kylie-a/waypoint/waypoint"
	"github.com/spf13/cobra"
)

var (
	target string
)

type Release struct {
	conf        *waypoint.Config
	deploy      waypoint.Deployment
	typ         waypoint.ReleaseType
	prevVersion *waypoint.Version
	newVersion  *waypoint.Version
	docker      *client.Client
}

func NewRelease(conf *waypoint.Config, target string, typ waypoint.ReleaseType) Release {
	var newVer *waypoint.Version

	deploy := conf.GetDeployment(target)
	prevVer, err := db.GetMostRecent(deploy.App)
	checkErr(err, true, false)
	if typ != waypoint.Rebuild {
		newVer = bumpVersion(deploy.App, *prevVer, typ)
	} else {
		newVer = prevVer
	}
	cli, err := client.NewEnvClient()
	checkErr(err, true, false)
	return Release{
		conf:        conf,
		deploy:      deploy,
		typ:         typ,
		prevVersion: prevVer,
		newVersion:  newVer,
		docker:      cli,
	}
}

func (r Release) Do() {
	for _, step := range steps {
		step.Execute(r)
	}
}

func (r Release) getDockerListOpts(image string) types.ImageListOptions {
	a := filters.NewArgs()
	a.Add("reference", image)
	return types.ImageListOptions{Filters: a}
}

// deployCmd represents the deploy command
var releaseCmd = &cobra.Command{
	Use:   "release",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		releaseType := getReleaseType(cmd)
		release := NewRelease(conf, target, releaseType)
		release.Do()
		//release.cleanup()
		//release.buildImage()
		//packageChart()
		//updateHelmRepo()
	},
}

func init() {
	rootCmd.AddCommand(releaseCmd)
	releaseCmd.Flags().Bool("major", false, "Bump the major version up by one")
	releaseCmd.Flags().Bool("minor", false, "Bump the minor version up by one")
	releaseCmd.Flags().Bool("patch", false, "Bump the patch version up by one")
	releaseCmd.Flags().Bool("rebuild", false, "Reuse the latest version up by one")
	releaseCmd.Flags().StringVar(&target, "target", "", "Target env")
}
