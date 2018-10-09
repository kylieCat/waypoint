package pkg

import (
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
	//yaml.Unmarshaler
	//json.Unmarshaler
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
