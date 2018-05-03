package state

import (
	"sort"
)

type Record interface {
	GetKey() []byte
}

type RecordList interface {
	Each(func(Record) error)
	sort.Interface
}
