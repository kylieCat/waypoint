package state

import (
	"fmt"
	"time"
)

type VersionPart int

const (
	MAJOR VersionPart = iota
	MINOR
	PATCH
)

type Version struct {
	Major      int    `json:"major"`
	Minor      int    `json:"minor"`
	Patch      int    `json:"patch"`
	CommitHash string `json:"commit_hash"`
	Timestamp  int64  `json:"date"`
	parts      []int
}

func (v Version) GetKey() []byte {
	return []byte(v.SemVer())
}

func (v Version) SemVer() string {
	return fmt.Sprintf("%d.%d.%d", v.Major, v.Minor, v.Patch)
}

func (v Version) BumpMajor() Version {
	v.Major++
	return NewVersion(v.Major, 0, 0)
}

func (v Version) BumpMinor() Version {
	v.Minor++
	return NewVersion(v.Major, v.Minor, 0)
}

func (v Version) BumpPatch() Version {
	v.Patch++
	return NewVersion(v.Major, v.Minor, v.Patch)
}

func NewVersion(major, minor, patch int) Version {
	parts := []int{major, minor, patch}
	return Version{Major: major, Minor: minor, Patch: patch, parts: parts, Timestamp: time.Now().Unix()}
}

type Versions []*Version

func (v Versions) Each(handler func(Record) error) error {
	for _, record := range v {
		err := handler(record)
		if err != nil {
			return err
		}
	}
	return nil
}

func (vs Versions) Len() int      { return len(vs) }
func (vs Versions) Swap(i, j int) { vs[i], vs[j] = vs[j], vs[i] }
func (vs Versions) Less(i, j int) bool {
	return vs[i].Timestamp < vs[j].Timestamp
}
