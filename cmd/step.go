package cmd

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/docker/docker/api/types"
	"github.com/kylie-a/requests"
)

type Step struct {
	StartMesg     func(r Release) string
	step          func(r Release) error
	exitOnErr     bool
	ShouldExecute func() bool
}

func (s Step) Execute(r Release) {
	fmt.Printf(s.StartMesg(r))
	if err := s.step(r); err != nil {
		if s.exitOnErr {
			printError(err.Error())
			os.Exit(2)
		} else {
			printWarning(err.Error())
			return
		}
	}
	fmt.Printf("%s\n", DONE)
}

var steps = []Step{
	{
		StartMesg: func(r Release) string {
			return fmt.Sprintf("Deleting previous image %s...", r.deploy.TaggedImageName(r.prevVersion.SemVer()))
		},
		step: func(r Release) error {
			var list []types.ImageSummary
			var err error

			listOpts := r.getDockerListOpts(r.deploy.TaggedImageName(r.prevVersion.SemVer()))
			if list, err = r.docker.ImageList(context.Background(), listOpts); err != nil {
				return err
			}

			var img types.ImageSummary
			if len(list) == 1 {
				img = list[0]
			} else {
				return errors.New(fmt.Sprintf("%d images found; skipping", len(list)))
			}
			_, err = r.docker.ImageRemove(context.Background(), img.ID, types.ImageRemoveOptions{})
			return err
		},
		exitOnErr: false,
	},
	{
		StartMesg: func(r Release) string{
			return fmt.Sprintf("Removing previous Helm chart  %s/%s...", r.deploy.App, r.prevVersion.SemVer())
		},
		step: func(r Release) error {
			var resp *requests.Response
			var err error

			if resp, err = requests.Delete(r.deploy.GetHelmURL(r.prevVersion.SemVer()), requests.WithBasicAuth(token)); err != nil {
				return err
			}
			if resp.Code != 200 {
				return errors.New("error deleting previous chart: " + resp.Content())
			}
			return nil
		},
		exitOnErr: false,
	},
	{
		StartMesg: func(r Release) string{
			return fmt.Sprintf("Building image %s...", r.deploy.TaggedImageName(r.newVersion.SemVer()))
		},
		step: func(r Release) error {
			var ctx io.Reader
			var err error

			if ctx, err = r.deploy.GetContext(); err != nil {
				return err
			}
			opts := types.ImageBuildOptions{
				Tags: []string{r.deploy.TaggedImageName(r.newVersion.SemVer())},
			}
			if _, err = r.docker.ImageBuild(context.Background(), ctx, opts); err != nil {
				return err
			}
			return nil
		},
		exitOnErr: true,
	},
	{
		StartMesg: func(r Release) string{
			return fmt.Sprintf("Pushing image %s...", r.deploy.TaggedImageName(r.newVersion.SemVer()))
		},
		step: func(r Release) error {
			var err error

			opts := types.ImagePushOptions{RegistryAuth: r.deploy.GetDockerAuth()}
			ref := r.deploy.TaggedImageName(r.newVersion.SemVer())
			if _, err = r.docker.ImagePush(context.Background(), ref, opts); err != nil {
				return err
			}
			return nil
		},
		exitOnErr: true,
	},
	{
		StartMesg: func(r Release) string{
			return fmt.Sprintf("Creating Helm chart %s:%s...", r.deploy.App, r.newVersion.SemVer())
		},
		step: func(r Release) error {
			// helm package --version ${version} manifests/${app} -d ./manifests/ &> /tmp/build
			return nil
		},
		exitOnErr: true,
	},
	{
		StartMesg: func(r Release) string{
			return fmt.Sprintf("Uploading Helm chart to %s...", r.deploy.GetHelmURL(r.newVersion.SemVer()))
		},
		step: func(r Release) error {
			//curl -s --data-binary "@./manifests/${app}-${version}.tgz" ${repo_url} &> /tmp/build
			return nil
		},
		exitOnErr: true,
	},
	{
		StartMesg: func(r Release) string{
			return "Updating Helm index file..."
		},
		step: func(r Release) error {
			//helm repo index manifests/${app} --url ${repo_url}
			return nil
		},
		exitOnErr: false,
	},
	{
		StartMesg: func(r Release) string{
			return "Updating Helm chart repos..."
		},
		step: func(r Release) error {
			// helm repo update
			return nil
		},
		exitOnErr: false,
	},
}
