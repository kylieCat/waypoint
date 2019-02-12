package pkg

import (
	"fmt"
	"os"

	"context"
)

type Step struct {
	StartMesg     func(r Release) string
	Action        func(r Release) error
	ExitOnErr     bool
	ShouldExecute func(r Release) bool
}

func (s Step) Execute(r Release) {
	fmt.Printf(s.StartMesg(r))
	if s.ShouldExecute(r) {
		if err := s.Action(r); err != nil {
			if s.ExitOnErr {
				printError(err.Error())
				os.Exit(2)
			} else {
				printWarning(err.Error())
				return
			}
		}
		fmt.Println(green("✅ DONE!"))
	} else {
		printSkipping()
	}
}

const shouldBuild = true

var DeletePrevImage = Step{
	StartMesg: func(r Release) string {
		return fmt.Sprintf("Deleting previous image %s...", r.deploy.TaggedImageName(r.prevVersion.SemVer()))
	},
	Action: func(r Release) error {
		imageName := r.deploy.TaggedImageName(r.prevVersion.SemVer())
		return r.docker.RemoveImage(imageName)
	},
	ExitOnErr: false,
	ShouldExecute: func(r Release) bool {
		return shouldBuild
	},
}

var DeletePrevChart = Step{
	StartMesg: func(r Release) string {
		return fmt.Sprintf("Removing previous Helm chart  %s...", r.deploy.GetHelmDeleteURL(r.prevVersion.SemVer()))
	},
	Action: func(r Release) error {
		return r.helm.RemoveChart(r.App(), r.HelmRepoName(), r.newVersion.SemVer())
	},
	ExitOnErr: false,
	ShouldExecute: func(r Release) bool {
		return r.helm.HasChart(r.App(), r.HelmRepoName(), r.newVersion.SemVer())
	},
}

var BuildImage = Step{
	StartMesg: func(r Release) string {
		return fmt.Sprintf("Building image %s...", r.deploy.TaggedImageName(r.newVersion.SemVer()))
	},
	Action: func(r Release) error {
		imageName := r.deploy.TaggedImageName(r.newVersion.SemVer())
		return r.docker.BuildImage(imageName, r.deploy.DockerContext())
	},
	ExitOnErr: true,
	ShouldExecute: func(r Release) bool {
		return shouldBuild
	},
}

var PushImage = Step{
	StartMesg: func(r Release) string {
		return fmt.Sprintf("Pushing image %s...", r.deploy.TaggedImageName(r.newVersion.SemVer()))
	},
	Action: func(r Release) error {
		ref := r.deploy.TaggedImageName(r.newVersion.SemVer())
		return r.docker.PushImage(ref, r.deploy.GetDockerRepo(), r.deploy.DockerCredHelper())
	},
	ExitOnErr: true,
	ShouldExecute: func(r Release) bool {
		return shouldBuild
	},
}

var CreateChart = Step{
	StartMesg: func(r Release) string {
		return fmt.Sprintf("Creating Helm chart %s:%s...", r.deploy.App, r.newVersion.SemVer())
	},
	Action: func(r Release) error {
		ver := r.newVersion.SemVer()
		src, _ := r.Format(r.GetHelmChartSrc())
		dest, _ := r.Format(r.GetHelmChartDest())
		shouldSaveLocal := r.SaveHelmLocal()

		return r.helm.Package(src, ver, dest, shouldSaveLocal)
	},
	ExitOnErr: true,
	ShouldExecute: func(r Release) bool {
		return true
	},
}

var UploadChart = Step{
	StartMesg: func(r Release) string {
		return fmt.Sprintf("Uploading Helm chart to %s...", r.deploy.GetHelmPostURL())
	},
	Action: func(r Release) error {
		ch := r.GetHelmPackage(r.newVersion.SemVer())
		return r.helm.UploadChart(ch, r.HelmRepoName())
	},
	ExitOnErr: true,
	ShouldExecute: func(r Release) bool {
		return true
	},
}

var UpdateHelmIndex = Step{
	StartMesg: func(r Release) string {
		return "Updating Helm index file..."
	},
	Action: func(r Release) error {
		//helm repo index manifests/${cmd} --url ${repo_url}
		src, _ := r.Format(r.deploy.GetHelmChartDir())
		baseURL := r.deploy.GetHelmPostURL()

		return r.helm.UpdateIndex(src, baseURL)
	},
	ExitOnErr: false,
	ShouldExecute: func(r Release) bool {
		return true
	},
}

var UpdateHelmRepos = Step{
	StartMesg: func(r Release) string {
		return "Updating Helm chart repos..."
	},
	Action: func(r Release) error {
		return r.helm.UpdateRepo(r.HelmRepoName())
	},
	ExitOnErr: false,
	ShouldExecute: func(r Release) bool {
		return true
	},
}

var DeployChart = Step{
	StartMesg: func(r Release) string {
		appVer := fmt.Sprintf("%s%s", r.deploy.App, r.newVersion.SemVer())
		return fmt.Sprintf("Deploying chart %s...", green(appVer))
	},
	Action: func(r Release) error {
		var cancel context.CancelFunc
		var err error

		src := r.GetHelmPackagePath(r.newVersion.SemVer())
		if cancel, err = r.k8s.StartForwarder(); err != nil {
			return err
		}
		defer cancel()
		return r.helm.Deploy(r.App(), src,"sre", r.deploy.Helm.Args)
	},
	ExitOnErr: true,
	ShouldExecute: func(r Release) bool {
		return shouldBuild
	},
}

var InstallChart = Step{
	StartMesg: func(r Release) string {
		appVer := fmt.Sprintf("%s%s", r.deploy.App, r.newVersion.SemVer())
		return fmt.Sprintf("Installing chart %s...", green(appVer))
	},
	Action: func(r Release) error {
		var cancel context.CancelFunc
		var err error

		src := r.GetHelmPackagePath(r.newVersion.SemVer())
		if cancel, err = r.k8s.StartForwarder(); err != nil {
			return err
		}
		defer cancel()
		return r.helm.Install(src, "sre", Args{})
	},
	ExitOnErr: true,
	ShouldExecute: func(r Release) bool {
		return shouldBuild
	},
}

var UpgradeChart = Step{
	StartMesg: func(r Release) string {
		appVer := fmt.Sprintf("%s%s", r.deploy.App, r.newVersion.SemVer())
		return fmt.Sprintf("Updating release %s...", green(appVer))
	},
	Action: func(r Release) error {
		var cancel context.CancelFunc
		var err error

		src := r.deploy.GetHelmPackagePath(r.newVersion.SemVer())
		if cancel, err = r.k8s.StartForwarder(); err != nil {
			return err
		}
		defer cancel()

		app := r.App()
		return r.helm.Upgrade(app, src, Args{})
	},
	ExitOnErr: true,
	ShouldExecute: func(r Release) bool {
		return shouldBuild
	},
}

var DefaultSteps = []Step{
	DeletePrevImage,
	DeletePrevChart,
	BuildImage,
	PushImage,
	CreateChart,
	UploadChart,
	UpdateHelmIndex,
	UpdateHelmRepos,
	InstallChart,
	UpgradeChart,
}
