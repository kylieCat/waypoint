package state

import (
	"errors"
	"fmt"
	"os"
	"sort"

	"github.com/boltdb/bolt"
	"github.com/vmihailenco/msgpack"
)

func getDB(path string, mode os.FileMode, options *bolt.Options) (*bolt.DB, error) {
	db, err := bolt.Open(path, mode, options)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func closeDb(db *bolt.DB) {
	err := db.Close()
	if err != nil {
		fmt.Println(err.Error())
	}
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
	versions, err := wp.ListAll(app)
	if err != nil {
		return nil, err
	}
	versionCount := len(versions)
	if versionCount == 0 {
		return nil, errors.New("no versions found for app")
	}
	return versions[versionCount-1], err
}

func (wp WaypointStore) ListAll(app string) (Versions, error) {
	db, err := getDB(wp.DBFilePath, 0600, nil)
	if err != nil {
		return nil, err
	}
	defer closeDb(db)
	var raw [][]byte
	err = db.View(func(tx *bolt.Tx) error {
		b, err := getBucket(tx, app)
		if err != nil {
			return err
		}
		b.ForEach(func(k, v []byte) error {
			raw = append(raw, v)
			return nil
		})
		return nil
	})
	versions := make(Versions, len(raw))
	for idx, r := range raw {
		var version Version
		err := msgpack.Unmarshal(r, &version)
		if err != nil {
			return nil, err
		}
		versions[idx] = &version
	}
	sort.Sort(versions)
	return versions, err
}

func (wp WaypointStore) NewVersion(app string, version *Version) error {
	db, err := getDB(wp.DBFilePath, 0600, nil)
	if err != nil {
		return err
	}
	defer closeDb(db)

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

func (wp WaypointStore) AddApplication(name string, initialVersion string) error {
	db, err := getDB(wp.DBFilePath, 0600, nil)
	if err != nil {
		return err
	}
	defer closeDb(db)

	err = db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucket([]byte(name))
		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}
		return nil
	})
	closeDb(db)
	if err != nil {
		return err
	}
	parts, err := GetPartsFromSemVer(initialVersion)
	if err != nil {
		return err
	}
	newVersion := NewVersion(parts[MAJOR], parts[MINOR], parts[PATCH])
	return wp.NewVersion(name, &newVersion)
}
