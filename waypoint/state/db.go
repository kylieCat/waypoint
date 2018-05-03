package state

import (
	"errors"
	"fmt"
	"os"

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

func (wp WaypointStore) AddApplication(name string, initialVersion string) error {
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
