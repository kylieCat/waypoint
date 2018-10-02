package waypoint

import (
	"io/ioutil"
	"gopkg.in/yaml.v2"
	"fmt"
	"google.golang.org/api/option"
)

type AuthKind string

const (
	ApiKey    AuthKind = "apiKey"
	CredsFile AuthKind = "credsFile"
)

type Auth struct {
	Kind  AuthKind `json:"kind" yaml:"kind"`
	Value string   `json:"value" yaml:"value"`
}

type ConfigFile struct {
	Project string `json:"project" yaml:"project"`
	Auth    Auth   `json:"auth" yaml:"auth"`
}

func (c ConfigFile) GetAuth() option.ClientOption {
	switch c.Auth.Kind {
	case ApiKey:
		return option.WithAPIKey(c.Auth.Value)
	case CredsFile:
		return option.WithCredentialsFile(c.Auth.Value)
	default:
		return nil
	}
}

func GetConf(fileName string) *ConfigFile {
	conf := new(ConfigFile)
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
