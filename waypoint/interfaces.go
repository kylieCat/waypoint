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
