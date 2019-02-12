package pkg

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"text/template"

	"github.com/mitchellh/go-homedir"
)

type Release struct {
	deploy      *Deployment
	typ         ReleaseType
	prevVersion *Version
	newVersion  *Version
	docker      IDocker
	helm        IHelm
	k8s         IK8s
	db          IStorage
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
	return fmt.Sprintf("%s%s/%s/%s", r.deploy.Helm.Name, ChartsAPI, r.deploy.App, version)
}

func (r Release) GetHelmPostURL() string {
	return fmt.Sprintf("%s%s", r.deploy.Helm.Name, ChartsAPI)
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

func (r Release) Format(tmpl string) (string, error) {
	var buf bytes.Buffer

	f := fmtData{
		App:        r.App(),
		NewVersion: r.newVersion.SemVer(),
		OldVersion: r.prevVersion.SemVer(),
	}
	t := template.Must(template.New("letter").Parse(tmpl))
	err := t.Execute(&buf, f)
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}

func (r *Release) ApplyOptions(opts ...ReleaseOption) *Release {
	for _, opt := range opts {
		opt(r)
	}
	return r
}

func NewRelease(conf *Config, target string, typ ReleaseType, opts ...ReleaseOption) *Release {
	deploy := conf.GetDeployment(target)
	release := &Release{
		deploy: deploy,
		typ:    typ,
	}
	release.ApplyOptions(opts...)
	return release
}
