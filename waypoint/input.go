package waypoint

import (
	"encoding/base64"
	"fmt"
	_ "github.com/docker/docker-credential-helpers/client"
	"github.com/docker/docker/pkg/archive"
	"github.com/mitchellh/go-homedir"
	"google.golang.org/api/option"
	"gopkg.in/yaml.v2"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"github.com/docker/docker-credential-helpers/client"
	"github.com/docker/docker/api/types"
	"encoding/json"
)

type AuthKind string

const (
	ApiKey    AuthKind = "apiKey"
	CredsFile AuthKind = "credsFile"
	chartsAPI          = "/api/charts"
)

type HelmConf struct {
	Name string `json:"name" yaml:"name"`
	ChartDir string `json:"chartDir" yaml:"chartDir"`
	DestDir string `json:"destDir" yaml:"destDir"`
	Save bool `json:"save" yaml:"save"`
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

func (d Deployment) GetHelmDeleteURL(version string) string {
	return fmt.Sprintf("%s%s/%s/%s", d.Helm.Name, chartsAPI, d.App, version)
}

func (d Deployment) GetHelmPostURL() string {
	return fmt.Sprintf("%s%s", d.Helm.Name, chartsAPI)
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
	creds, err := client.Get(client.NewShellProgramFunc("docker-credential-gcr"), "https://gcr.io")
	if err != nil {
		fmt.Println(err)
		return ""
	}
	authConfig := types.AuthConfig{
		Username: creds.Username,
		Password: creds.Secret,
	}
	encodedJSON, err := json.Marshal(authConfig)
	if err != nil {
		panic(err)
	}
	return base64.URLEncoding.EncodeToString(encodedJSON)
}

func (d Deployment) GetHelmChartDir () string {
	filePath, _ := homedir.Expand(d.Helm.ChartDir)
	return filePath
}

func (d Deployment) GetHelmDestDir () string {
	filePath, _ := homedir.Expand(d.Helm.DestDir)
	return filePath
}

func (d Deployment) GetHelmPackage(version string) []byte {
	var data []byte
	var err error

	fileName := fmt.Sprintf("%s-%s.tgz", d.App, version)
	filePath := filepath.Join(d.GetHelmDestDir(),  fileName)
	if data, err = ioutil.ReadFile(filePath); err != nil {
		return nil
	}
	return data
}

func (d Deployment) SaveHelmLocal () bool {
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
