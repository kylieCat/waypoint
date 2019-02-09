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

	"crypto/tls"
	"net/http"

	"github.com/kylie-a/waypoint/pkg"
	"github.com/kylie-a/waypoint/pkg/docker"
	"github.com/kylie-a/waypoint/pkg/helm"
	"github.com/kylie-a/waypoint/pkg/k8s"
	"github.com/spf13/cobra"
)

var (
	target string
)

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
		deploy := conf.GetDeployment(target)

		tr := http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		}
		httpClient := &http.Client{Transport: &tr}
		helmClient := helm.NewClient(helm.HelmToken(os.Getenv("HELM_TOKEN")))
		k8sClient := k8s.NewClient(
			k8s.Context(deploy.Context),
			k8s.HTTPClient(httpClient),
		)
		dockerClient, err := docker.NewDockerClient()
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
		prevVer, err := ws.GetLatest(deploy.App)
		checkErr(err, true, false)
		newVer := prevVer.Bump(releaseType)
		if err := ws.Save(deploy.App, newVer); err != nil {
			checkErr(err, true, false)
		}
		release := pkg.NewRelease(
			conf,
			target,
			releaseType,
			pkg.Helm(helmClient),
			pkg.K8s(k8sClient),
			pkg.Docker(dockerClient),
			pkg.DB(ws),
			pkg.PrevVersion(prevVer),
			pkg.CurrentVersion(newVer),
		)
		release.Do(pkg.DefaultSteps)
	},
}

func init() {
	rootCmd.AddCommand(releaseCmd)
	releaseCmd.Flags().Bool("major", false, "Bump the major version up by one")
	releaseCmd.Flags().Bool("minor", false, "Bump the minor version up by one")
	releaseCmd.Flags().Bool("patch", false, "Bump the patch version up by one")
	releaseCmd.Flags().Bool("rebuild", false, "Reuse the latest version and redeploy")
	releaseCmd.Flags().StringVar(&target, "target", "", "The deployment to target in the conf file.")
}
