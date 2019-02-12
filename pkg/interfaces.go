package pkg

import (
	"context"
	"encoding/json"
	"io"
)

type IDocker interface {
	RemoveImage(taggedImageName string) error
	BuildImage(taggedImageName, buildCtx string) error
	PushImage(ref, repo, credHelper string) error
	Auth(repo, credHelper string) (string, error)
	GetContext(buildCtx string) (io.Reader, error)
}
type IStorage interface {
	GetLatest(app string) (*Version, error)
	All(app string) (Versions, error)
	Save(app string, version *Version) error
	AddApplication(name string, initialVersion string) error
}

type IHelm interface {
	RemoveChart(app, repoName, version string) error
	UploadChart(ch []byte, repoName string) error
	HasChart(app, repoName, version string) bool
	Package(src, version, dest string, saveLocal bool) error
	IsInstalled(name string) bool
	Deploy(name, src, ns string, opts Args) error
	Install(src, ns string, opts Args) error
	Upgrade(app, src string, opts Args) error
	UpdateRepo(repoName string) error
	UpdateRepos() error
	UpdateIndex(chartSrc, baseURL string) error
}

type IK8s interface {
	GetTillerPod() (string, error)
	ListNodes()
	StartForwarder() (context.CancelFunc, error)
}

type BackendConf interface {
	GetKind() BackendKind
	GetAuth() BackendAuthConf
}

type BackendAuthConf interface {
	GetKind() GCPAuthKind
}

type ConfUnmarshaler interface {
	auxConf() AuxConf
}

type AuxConf interface {
	conf(map[string]interface{}) ConfUnmarshaler
}

func defaultUnmarshalJSON(data []byte, conf ConfUnmarshaler) (ConfUnmarshaler, error) {
	var raw map[string]interface{}
	jd := conf.auxConf()

	if err := json.Unmarshal(data, &raw); err != nil {
		return nil, err
	}
	return jd.conf(raw), nil
}

func defaultUnmarshalYAML(unmarshal func(interface{}) error, conf ConfUnmarshaler) (ConfUnmarshaler, error) {
	var raw map[string]interface{}
	jd := conf.auxConf()
	if err := unmarshal(&raw); err != nil {
		return nil, err
	}

	//fmt.Printf("conf: %v\n", conf)
	//fmt.Printf("raw: %v\n", raw)
	//fmt.Printf("jd: %v\n", jd)
	//fmt.Println("===============================================")
	return jd.conf(raw), nil
}
