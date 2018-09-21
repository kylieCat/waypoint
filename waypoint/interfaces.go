package waypoint

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

type DataBase interface {
	GetMostRecent(app string) (*Version, error)
	ListAll(app string) (Versions, error)
	NewVersion(app string, version *Version) error
	AddApplication(name string, initialVersion string) error
}
