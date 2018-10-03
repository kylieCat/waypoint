package waypoint

import (
	"encoding/base64"
	"fmt"
	"github.com/docker/docker-credential-helpers/client"
	"github.com/docker/docker/pkg/archive"
	"github.com/mitchellh/go-homedir"
	"google.golang.org/api/option"
	"gopkg.in/yaml.v2"
	"io"
	"io/ioutil"
	"os"
)

type AuthKind string

const (
	ApiKey    AuthKind = "apiKey"
	CredsFile AuthKind = "credsFile"
	chartsAPI          = "/api/charts"
)

type HelmConf struct {
	Name string `json:"name" yaml:"name"`
}

type DockerConf struct {
	Repo    string `json:"repo" yaml:"repo"`
	Context string `json:"context" yaml:"context"`
	File    string `json:"file" yaml:"file"`
}

type Deployment struct {
	App     string     `json:"app" yaml:"app"`
	Project string     `json:"project" yaml:"project"`
	Docker  DockerConf `json:"docker" yaml:"docker"`
	Helm    HelmConf   `json:"helm" yaml:"helm"`
}

func (d Deployment) ImageName() string {
	docker := d.Docker
	return fmt.Sprintf("%s/%s", docker.Repo, d.App)
}

func (d Deployment) TaggedImageName(version string) string {
	return fmt.Sprintf("%s:%s", d.ImageName(), version)
}

func (d Deployment) GetHelmURL(version string) string {
	return fmt.Sprintf("%s%s/%s/%s", d.Helm.Name, chartsAPI, d.App, version)
}

func (d Deployment) GetDockerfile() string {
	filePath, _ := homedir.Expand(d.Docker.File)
	return filePath
}

func (d Deployment) GetContext() (io.Reader, error) {
	ctxDir, err := homedir.Expand(d.Docker.Context)
	if err != nil {
		return nil, err
	}

	return archive.TarWithOptions(ctxDir, &archive.TarOptions{})
}

func (d Deployment) GetDockerAuth() string {
	creds, err := client.Get(client.NewShellProgramFunc("docker-credential-gcloud"), "https://gcr.io")
	if err != nil {
		fmt.Println(err)
		return ""
	}
	return base64.StdEncoding.EncodeToString([]byte(creds.Secret))
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
