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

type Backend struct {
	Kind BackendKind   `json:"kind" yaml:"kind"`
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
