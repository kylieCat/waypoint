package pkg

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/mitchellh/go-homedir"
	"gopkg.in/yaml.v2"
)

type GCPAuthKind string
type BackendKind string

const (
	DataStore BackendKind = "datastore"
	Bolt      BackendKind = "bolt"
	MongoDB   BackendKind = "mongo"
	Dynamo    BackendKind = "dynamo"
	ApiKey    GCPAuthKind = "apiKey"
	CredsFile GCPAuthKind = "credsFile"
	chartsAPI             = "/api/charts"
)

func toStringSlice(data []interface{}) []string {
	results := make([]string, 0, len(data))
	for _, label := range data {
		results = append(results, label.(string))
	}
	return results
}

type TillerConf struct {
	Namespace string   `json:"namespace" yaml:"namespace"`
	Context   string   `json:"context" yaml:"context"`
	Endpoint  string   `json:"endpoint" yaml:"endpoint"`
	Labels    []string `json:"labels" yaml:"labels"`
}

func (t TillerConf) auxConf() AuxConf {
	return newAuxTillerConf()
}

func (t *TillerConf) UnmarshalJSON(data []byte) error {
	var unmarshaled ConfUnmarshaler
	var err error

	if unmarshaled, err = defaultUnmarshalJSON(data, t); err != nil {
		return err
	}
	concrete := unmarshaled.(TillerConf)
	*t = concrete
	return nil
}

func (t *TillerConf) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var unmarshaled ConfUnmarshaler
	var err error

	if unmarshaled, err = defaultUnmarshalYAML(unmarshal, t); err != nil {
		return err
	}
	concrete := unmarshaled.(TillerConf)
	*t = concrete
	return nil
}

func (t *TillerConf) applyDefaults(defaults Deployment) {
	if t.Namespace == "" {
		t.Namespace = defaults.Tiller.Namespace
	}
	if t.Context == "" {
		t.Context = defaults.Tiller.Context
	}
	if t.Endpoint == "" {
		t.Endpoint = defaults.Tiller.Endpoint
	}
	if t.Labels == nil || len(t.Labels) == 0 {
		t.Labels = defaults.Tiller.Labels
	}
}

type auxTillerConf struct {
	Namespace string   `json:"namespace" yaml:"namespace"`
	Context   string   `json:"context" yaml:"context"`
	Endpoint  string   `json:"endpoint" yaml:"endpoint"`
	Labels    []string `json:"labels" yaml:"labels"`
}

func (t auxTillerConf) conf(data map[string]interface{}) ConfUnmarshaler {
	for key, value := range data {
		switch key {
		case "namespace":
			t.Namespace = value.(string)
		case "context":
			t.Context = value.(string)
		case "endpoint":
			t.Endpoint = value.(string)
		case "labels":
			labels := toStringSlice(value.([]interface{}))
			t.Labels = labels
		}
	}
	return TillerConf{
		Namespace: t.Namespace,
		Context:   t.Context,
		Endpoint:  t.Endpoint,
		Labels:    t.Labels,
	}
}

func newAuxTillerConf() AuxConf {
	return auxTillerConf{
		Namespace: "kube-system",
		Endpoint:  "http://localhost:50000",
		Labels:    []string{"app=helm", "name=tiller"},
	}
}

type HelmConf struct {
	Name     string   `json:"name" yaml:"name"`
	ChartDir string   `json:"chartDir" yaml:"chartDir"`
	DestDir  string   `json:"destDir" yaml:"destDir"`
	Save     bool     `json:"save" yaml:"save"`
	Args     []string `json:"args" yaml:"args"`
}

func (h HelmConf) auxConf() AuxConf {
	return newAuxHelmConf()
}

func (h *HelmConf) UnmarshalJSON(data []byte) error {
	var unmarshaled ConfUnmarshaler
	var err error

	if unmarshaled, err = defaultUnmarshalJSON(data, h); err != nil {
		return err
	}
	concrete := unmarshaled.(HelmConf)
	*h = concrete
	return nil
}

func (h *HelmConf) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var unmarshaled ConfUnmarshaler
	var err error

	if unmarshaled, err = defaultUnmarshalYAML(unmarshal, h); err != nil {
		return err
	}
	concrete := unmarshaled.(HelmConf)
	*h = concrete
	return nil
}

func (h *HelmConf) applyDefaults(defaults Deployment) {
	if h.Name == "" {
		h.Name = defaults.Helm.Name
	}
	if h.ChartDir == "" {
		h.ChartDir = defaults.Helm.ChartDir
	}
	if h.DestDir == "" {
		h.DestDir = defaults.Helm.DestDir
	}
	if !h.Save {
		h.Save = defaults.Helm.Save
	}
	if h.Args == nil || len(h.Args) == 0 {
		fmt.Println(defaults.Helm.Args)
		h.Args = defaults.Helm.Args
	}
}

type auxHelmConf struct {
	Name     string   `json:"name" yaml:"name"`
	ChartDir string   `json:"chartDir" yaml:"chartDir"`
	DestDir  string   `json:"destDir" yaml:"destDir"`
	Save     bool     `json:"save" yaml:"save"`
	Args     []string `json:"args" yaml:"args"`
}

func (h auxHelmConf) conf(data map[string]interface{}) ConfUnmarshaler {
	for key, value := range data {
		switch key {
		case "name":
			h.Name = value.(string)
		case "chartDir":
			h.ChartDir = value.(string)
		case "destDir":
			h.DestDir = value.(string)
		case "save":
			h.Save = value.(bool)
		case "args":
			vals := toStringSlice(value.([]interface{}))
			h.Args = vals
		}
	}
	return HelmConf{
		Name:     h.Name,
		ChartDir: h.ChartDir,
		DestDir:  h.DestDir,
		Save:     h.Save,
		Args:     h.Args,
	}
}

func newAuxHelmConf() AuxConf {
	return auxHelmConf{
		ChartDir: "./deploy",
		DestDir:  "./deploy",
		Save:     true,
		Args: []string{"--version={{.Version}}"},
	}
}

type DockerConf struct {
	Repo    string `json:"repo" yaml:"repo"`
	Creds   string `json:"creds" yaml:"creds"`
	Context string `json:"context" yaml:"context"`
	File    string `json:"file" yaml:"file"`
}

func (d DockerConf) auxConf() AuxConf {
	return newAuxDockerConf()
}

func (d *DockerConf) UnmarshalJSON(data []byte) error {
	var unmarshaled ConfUnmarshaler
	var err error

	if unmarshaled, err = defaultUnmarshalJSON(data, d); err != nil {
		return err
	}
	concrete := unmarshaled.(DockerConf)
	*d = concrete
	return nil
}

func (d *DockerConf) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var unmarshaled ConfUnmarshaler
	var err error

	if unmarshaled, err = defaultUnmarshalYAML(unmarshal, d); err != nil {
		return err
	}
	concrete := unmarshaled.(DockerConf)
	*d = concrete
	return nil
}

func (d *DockerConf) applyDefaults(defaults Deployment) {
	if d.Repo == "" {
		d.Repo = defaults.Docker.Repo
	}
	if d.Creds == "" {
		d.Repo = defaults.Docker.Creds
	}
	if d.Context == "" {
		d.Repo = defaults.Docker.Context
	}
	if d.File == "" {
		d.Repo = defaults.Docker.File
	}
}

type auxDockerConf struct {
	Repo    string `json:"repo" yaml:"repo,omitempty"`
	Creds   string `json:"creds" yaml:"creds,omitempty"`
	Context string `json:"context" yaml:"context,omitempty"`
	File    string `json:"file" yaml:"file,omitempty"`
}

func newAuxDockerConf() auxDockerConf {
	return auxDockerConf{
		Repo:    "",
		Creds:   "docker-credential-gcr",
		Context: ".",
		File:    "Dockerfile",
	}
}

func (dc auxDockerConf) conf(data map[string]interface{}) ConfUnmarshaler {
	for key, value := range data {
		switch key {
		case "repo":
			dc.Repo = value.(string)
		case "creds":
			dc.Creds = value.(string)
		case "context":
			dc.Context = value.(string)
		case "file":
			dc.File = value.(string)
		}
	}
	return DockerConf{
		Repo:    dc.Repo,
		Creds:   dc.Creds,
		Context: dc.Context,
		File:    dc.File,
	}
}

func DefaultTillerConf() TillerConf {
	return TillerConf{
		Namespace: "kube-system",
		Endpoint:  "http://localhost:50000",
		Labels:    []string{"app=helm", "name=tiller"},
	}
}

func DefaultHelmConf() HelmConf {
	return HelmConf{
		Name:     "test",
		ChartDir: "deploy/{appName}",
		DestDir:  "deploy/",
		Save:     true,
		Args:     []string{},
	}
}

func DefaultDockerConf() DockerConf {
	return DockerConf{
		Repo:    "gcr.io/{project}",
		Creds:   "docker-credential-gcr",
		Context: ".",
		File:    "Dockerfile",
	}
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

type Deployments map[string]*Deployment

type Backend struct {
	Kind BackendKind       `json:"kind" yaml:"kind"`
	Conf map[string]string `json:"conf" yaml:"conf"`
}

type Config struct {
	App         string      `json:"app" yaml:"app"`
	Backend     *Backend    `json:"backend" yaml:"backend"`
	Defaults    Deployment  `json:"defaults" yaml:"defaults" `
	Deployments Deployments `json:"deployments" yaml:"deployments"`
}

func (c Config) GetDeployment(dep string) *Deployment {
	d, ok := c.Deployments[dep]
	if !ok {
		fmt.Printf("deployment %s not configured", dep)
		os.Exit(2)
	}
	return d
}

func NewConfig() *Config {
	return &Config{
		App: "",
		Backend: &Backend{
			Kind: "datastore",
			Conf: make(map[string]string),
		},
		Deployments: make(Deployments),
	}
}

func GetConf(fileName string) (*Config, error) {
	extension := filepath.Ext(fileName)
	conf := NewConfig()
	fileName, err := homedir.Expand(fileName)
	b, err := ioutil.ReadFile(fileName)
	if err != nil {
		return nil, err
	}
	switch extension {
	case ".yaml":
		err = yaml.Unmarshal(b, &conf)
		if err != nil {
			return nil, err
		}
	case ".json":
		err = json.Unmarshal(b, &conf)
		if err != nil {
			return nil, err
		}
	}
	for _, deployment := range conf.Deployments {
		deployment.applyDefaults(conf.Defaults)
	}
	return conf, nil
}
