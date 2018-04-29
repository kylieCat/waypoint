package state

import (
	"errors"
	"fmt"
	"github.com/boltdb/bolt"
	"github.com/vmihailenco/msgpack"
	"os"
	"time"
)

type VersionPart int

const (
	MAJOR VersionPart = iota
	MINOR
	PATCH
)

func getDB(path string, mode os.FileMode, options *bolt.Options) (*bolt.DB, error) {
	db, err := bolt.Open(path, mode, options)
	if err != nil {
		return nil, err
	}
	return db, nil
}

type Application struct {
	Name     string   `json:"name"`
	Versions Versions `json:"versions"`
}

type Applications []*Application

func (app Applications) Len() int           { return len(app) }
func (app Applications) Swap(i, j int)      { app[i], app[j] = app[j], app[i] }
func (app Applications) Less(i, j int) bool { return app[i].Name < app[j].Name }

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
	return NewVersion(v.Major, 0,0)
}

func (v Version) BumpMinor() Version {
	v.Minor++
	return NewVersion(v.Major, v.Minor, v.Patch)
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

func (vs Versions) Len() int      { return len(vs) }
func (vs Versions) Swap(i, j int) { vs[i], vs[j] = vs[j], vs[i] }
func (vs Versions) Less(i, j int) bool {
	return vs[i].Timestamp < vs[j].Timestamp
}

func getBucket(tx *bolt.Tx, bucketName string) (*bolt.Bucket, error) {
	bucket := tx.Bucket([]byte(bucketName))
	if bucket == nil {
		return nil, errors.New("no bucket found for bucketName")
	}
	return bucket, nil
}

type WaypointStore struct {
	DBFilePath string `json:"db_file_path"`
}

func (wp WaypointStore) GetMostRecent(app string) (*Version, error) {
	db, err := getDB(wp.DBFilePath, 0600, nil)
	defer db.Close()
	if err != nil {
		return nil, err
	}
	version := &Version{}
	err = db.View(func(tx *bolt.Tx) error {
		bucket, err := getBucket(tx, app)
		if err != nil {
			return err
		}
		notFound := errors.New("no versions found for app")
		if key, val := bucket.Cursor().Last(); key != nil {
			err := msgpack.Unmarshal(val, version)
			if err != nil {
				return err
			}
			notFound = nil
		}

		return notFound
	})
	return version, err
}

func (wp WaypointStore) ListAll(app string) ([]*Version, error) {
	db, err := getDB(wp.DBFilePath, 0600, nil)
	defer db.Close()
	if err != nil {
		return nil, err
	}

	versions := make([]*Version, 0)
	err = db.View(func(tx *bolt.Tx) error {
		b, err := getBucket(tx, app)
		if err != nil {
			return err
		}
		b.ForEach(func(k, v []byte) error {
			version := &Version{}
			if err := msgpack.Unmarshal(v, version); err != nil {
				return err
			}
			versions = append(versions, version)
			return nil
		})
		return nil
	})
	return versions, err
}

func (wp WaypointStore) NewVersion(app string, version *Version) error {
	db, err := getDB(wp.DBFilePath, 0600, nil)
	defer db.Close()
	if err != nil {
		return err
	}

	return db.Update(func(tx *bolt.Tx) error {
		data, err := msgpack.Marshal(version)
		if err != nil {
			return err
		}
		bucket, err := getBucket(tx, app)
		if err != nil {
			return err
		}
		return bucket.Put(version.GetKey(), data)
	})
}

func (wp WaypointStore) AddApplication(name string,initialVersion string) error {
	db, err := getDB(wp.DBFilePath, 0600, nil)
	defer db.Close()
	if err != nil {
		return err
	}
	err = db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucket([]byte(name))
		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}
        parts, err := GetPartsFromSemVer(initialVersion)
        if err != nil {
            return err
        }
        newVersion := NewVersion(parts[MAJOR], parts[MINOR], parts[PATCH])
		return wp.NewVersion(name, &newVersion)
	})
    if err != nil {
        return nil
    }
    db.Close()
	db, err = getDB(wp.DBFilePath, 0600, nil)
	defer db.Close()
	return db.Update(func(tx *bolt.Tx) error {
        parts, err := GetPartsFromSemVer(initialVersion)
        if err != nil {
            return err
        }
        newVersion := NewVersion(parts[MAJOR], parts[MINOR], parts[PATCH])
		return wp.NewVersion(name, &newVersion)
    })
}
