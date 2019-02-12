package pkg

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"

	"github.com/mitchellh/go-homedir"
	"gopkg.in/yaml.v2"
)

type Deployment struct {
	App     string     `json:"app" yaml:"app"`
	Project string     `json:"project" yaml:"project"`
	Context string     `json:"context" yaml:"context"`
	Docker  DockerConf `json:"docker" yaml:"docker"`
	Helm    HelmConf   `json:"helm" yaml:"helm"`
	Tiller  TillerConf `json:"tiller" yaml:"tiller"`
}

func (d Deployment) GetDockerRepo() string {
	return d.Docker.Repo
}

func (d Deployment) GetDockerfile() string {
	filePath, _ := homedir.Expand(d.Docker.File)
	return filePath
}

func (d Deployment) DockerCredHelper() string {
	return d.Docker.Creds
}

func (d Deployment) DockerContext() string {
	return d.Docker.Context
}

func (d Deployment) ImageName() string {
	docker := d.Docker
	return fmt.Sprintf("%s/%s", docker.Repo, d.App)
}

func (d Deployment) TaggedImageName(version string) string {
	return fmt.Sprintf("%s:%s", d.ImageName(), version)
}

func (d Deployment) GetHelmDeleteURL(version string) string {
	return fmt.Sprintf("%s%s/%s/%s", d.Helm.Name, ChartsAPI, d.App, version)
}

func (d Deployment) GetHelmPostURL() string {
	return fmt.Sprintf("%s%s", d.Helm.Name, ChartsAPI)
}

func (d Deployment) GetHelmChartDir() string {
	filePath, _ := homedir.Expand(d.Helm.ChartDir)
	return filePath
}

func (d Deployment) GetHelmDestDir() string {
	filePath, _ := homedir.Expand(d.Helm.DestDir)
	return filePath
}

func (d Deployment) GetHelmPackagePath(version string) string {
	fileName := fmt.Sprintf("%s-%s.tgz", d.App, version)
	return filepath.Join(d.GetHelmDestDir(), fileName)
}

func (d Deployment) GetHelmPackage(version string) []byte {
	var data []byte
	var err error

	filePath := d.GetHelmPackagePath(version)
	if data, err = ioutil.ReadFile(filePath); err != nil {
		return nil
	}
	return data
}

func (d Deployment) SaveHelmLocal() bool {
	return d.Helm.Save
}

func (d Deployment) YAML() []byte {
	yml, _ := yaml.Marshal(d)
	return yml
}

func (d Deployment) JSON() []byte {
	js, _ := json.Marshal(d)
	return js
}

func (d Deployment) YAMLString() string {
	yml, _ := yaml.Marshal(d)
	return string(yml)
}

func (d Deployment) JSONString() string {
	js, _ := json.Marshal(d)
	return string(js)
}

func (d *Deployment) applyDefaults(defaults Deployment) {

	if d.App == "" {
		d.App = defaults.App
	}
	if d.Project == "" {
		d.Project = defaults.Project
	}
	if d.Context == "" {
		d.Context = defaults.Context
	}
	d.Docker.applyDefaults(defaults)
	d.Helm.applyDefaults(defaults)
	d.Tiller.applyDefaults(defaults)
}

func NewDeployment() Deployment {
	return Deployment{
		App:     "",
		Project: "",
		Docker:  DefaultDockerConf(),
		Helm:    DefaultHelmConf(),
		Tiller:  DefaultTillerConf(),
	}
}

type Deployments map[string]*Deployment
