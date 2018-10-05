package waypoint

import (
	"os"

	"github.com/kylie-a/waypoint/waypoint/docker"
	helm2 "github.com/kylie-a/waypoint/waypoint/helm"
	"github.com/mitchellh/go-homedir"
	"fmt"
	"path/filepath"
	"io/ioutil"
)

var db DataBase

type Release struct {
	conf        *Config
	deploy      Deployment
	typ         ReleaseType
	prevVersion *Version
	newVersion  *Version
	docker      *docker.Client
	helm        *helm2.Client
}

func (r Release) Do(steps []Step) {
	for _, step := range steps {
		step.Execute(r)
	}
}

func (r Release) App() string {
	return r.deploy.App
}

func (r Release) HelmRepoName() string {
	return r.deploy.Helm.Name
}

func (r Release) GetDockerRepo() string {
	return r.deploy.Docker.Repo
}

func (r Release) GetDockerfile() string {
	filePath, _ := homedir.Expand(r.deploy.Docker.File)
	return filePath
}

func (r Release) DockerCredHelper() string {
	return r.deploy.Docker.Creds
}

func (r Release) DockerContext() string {
	return r.deploy.Docker.Context
}

func (r Release) ImageName() string {
	dkr := r.deploy.Docker
	return fmt.Sprintf("%s/%s", dkr.Repo, r.deploy.App)
}

func (r Release) TaggedImageName(version string) string {
	return fmt.Sprintf("%s:%s", r.deploy.ImageName(), version)
}

func (r Release) GetHelmDeleteURL(version string) string {
	return fmt.Sprintf("%s%s/%s/%s", r.deploy.Helm.Name, chartsAPI, r.deploy.App, version)
}

func (r Release) GetHelmPostURL() string {
	return fmt.Sprintf("%s%s", r.deploy.Helm.Name, chartsAPI)
}

func (r Release) GetHelmChartSrc() string {
	filePath, _ := homedir.Expand(r.deploy.Helm.ChartDir)
	return filePath
}

func (r Release) GetHelmChartDest() string {
	filePath, _ := homedir.Expand(r.deploy.Helm.DestDir)
	return filePath
}

func (r Release) GetHelmPackagePath(version string) string {
	fileName := fmt.Sprintf("%s-%s.tgz", r.deploy.App, version)
	return filepath.Join(r.GetHelmChartDest(), fileName)
}

func (r Release) GetHelmPackage(version string) []byte {
	var data []byte
	var err error

	filePath := r.GetHelmPackagePath(version)
	if data, err = ioutil.ReadFile(filePath); err != nil {
		return nil
	}
	return data
}

func (r Release) SaveHelmLocal() bool {
	return r.deploy.Helm.Save
}

func NewRelease(conf *Config, target string, typ ReleaseType) Release {
	var newVer *Version

	db = NewWaypointStoreDS(conf.Auth.Project, conf.GetAuth())
	deploy := conf.GetDeployment(target)
	prevVer, err := db.GetMostRecent(deploy.App)
	checkErr(err, true, false)
	newVer = prevVer.Bump(typ)
	dkr, err := docker.NewDockerClient()
	checkErr(err, true, false)
	helm := helm2.NewHelmClient(helm2.HelmToken(os.Getenv("HELM_TOKEN")))
	return Release{
		conf:        conf,
		deploy:      deploy,
		typ:         typ,
		prevVersion: prevVer,
		newVersion:  newVer,
		docker:      dkr,
		helm:        helm,
	}
}
