package pkg

import (
	"encoding/json"
	"sort"
)

type Record interface {
	GetKey() []byte
}

type RecordList interface {
	Each(func(Record) error) error
	sort.Interface
}

type BackendConf interface {
	GetKind() BackendKind
	GetAuth() BackendAuthConf
}

type BackendAuthConf interface {
	GetKind() GCPAuthKind
}

type BackendService interface {
	GetLatest(app string) (*Version, error)
	All(app string) (Versions, error)
	Save(app string, version *Version) error
	AddApplication(name string, initialVersion string) error
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

	return jd.conf(raw), nil
}