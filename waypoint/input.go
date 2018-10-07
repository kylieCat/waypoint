package waypoint

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/mitchellh/go-homedir"
	"google.golang.org/api/option"
	"gopkg.in/yaml.v2"
)

type AuthKind string

const (
	ApiKey    AuthKind = "apiKey"
	CredsFile AuthKind = "credsFile"
	chartsAPI          = "/api/charts"
)

type TillerConf struct {
	Namespace string   `json:"namespace" yaml:"namespace"`
	Context   string   `json:"context" yaml:"context"`
	Endpoint  string   `json:"endpoint" yaml:"endpoint"`
	Labels    []string `json:"labels" yaml:"labels"`
}

type HelmConf struct {
	Name     string   `json:"name" yaml:"name"`
	ChartDir string   `json:"chartDir" yaml:"chartDir"`
	DestDir  string   `json:"destDir" yaml:"destDir"`
	Save     bool     `json:"save" yaml:"save"`
	Set      []string `json:"set" yaml:"set"`
}

type DockerConf struct {
	Repo    string `json:"repo" yaml:"repo"`
	Creds   string `json:"creds" yaml:"creds"`
	Context string `json:"context" yaml:"context"`
	File    string `json:"file" yaml:"file"`
}

type Deployment struct {
	App     string     `json:"app" yaml:"app"`
	Project string     `json:"project" yaml:"project"`
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
	return fmt.Sprintf("%s%s/%s/%s", d.Helm.Name, chartsAPI, d.App, version)
}

func (d Deployment) GetHelmPostURL() string {
	return fmt.Sprintf("%s%s", d.Helm.Name, chartsAPI)
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

type Deployments map[string]Deployment

type Auth struct {
	Kind    AuthKind `json:"kind" yaml:"kind"`
	Project string   `json:"project" yaml:"project"`
	Value   string   `json:"value" yaml:"value"`
}

type Config struct {
	Auth        Auth        `json:"auth" yaml:"auth"`
	Deployments Deployments `json:"deployments" yaml:"deployments"`
}

func (c Config) GetDeployment(dep string) Deployment {
	d, ok := c.Deployments[dep]
	if !ok {
		fmt.Printf("deploymetn %s not configured", dep)
		os.Exit(2)
	}
	return d
}

func (c Config) GetAuth() option.ClientOption {
	switch c.Auth.Kind {
	case ApiKey:
		return option.WithAPIKey(c.Auth.Value)
	case CredsFile:
		return option.WithCredentialsFile(c.Auth.Value)
	default:
		return nil
	}
}

func GetConf(fileName string) *Config {
	conf := new(Config)
	fileName, err := homedir.Expand(fileName)
	b, err := ioutil.ReadFile(fileName)
	if err != nil {
		fmt.Println(err.Error())
		return nil
	}
	err = yaml.Unmarshal(b, &conf)
	if err != nil {
		fmt.Println(err.Error())
		return nil
	}
	return conf
}
