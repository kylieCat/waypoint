package db

import (
	"errors"
	"fmt"
	"os"
	"sort"

	"github.com/boltdb/bolt"
	"github.com/kylie-a/waypoint/pkg"
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

type WaypointStoreBolt struct {
	DBFilePath string `json:"db_file_path"`
}

func NewWaypointStoreBolt(dbPath string) WaypointStoreBolt {
	return WaypointStoreBolt{
		DBFilePath: dbPath,
	}
}

func (wp WaypointStoreBolt) GetLatest(app string) (*pkg.Version, error) {
	versions, err := wp.All(app)
	if err != nil {
		return nil, err
	}
	versionCount := len(versions)
	if versionCount == 0 {
		return nil, errors.New("no versions found for app")
	}
	return &versions[versionCount-1], err
}

func (wp WaypointStoreBolt) All(app string) (pkg.Versions, error) {
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
		err = b.ForEach(func(k, v []byte) error {
			raw = append(raw, v)
			return nil
		})
		return err
	})
	versions := make(pkg.Versions, len(raw))
	for idx, r := range raw {
		var version pkg.Version
		err := msgpack.Unmarshal(r, &version)
		if err != nil {
			return nil, err
		}
		versions[idx] = version
	}
	sort.Sort(versions)
	return versions, err
}

func (wp WaypointStoreBolt) Save(app string, version *pkg.Version) error {
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

func (wp WaypointStoreBolt) AddApplication(name string, initialVersion string) error {
	db, err := getDB(wp.DBFilePath, 0600, nil)
	if err != nil {
		return err
	}
	defer func(){
		err := db.Close()
		if err != nil {
			fmt.Println(err.Error())
		}
	}()

	return db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucket([]byte(name))
		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}
		return nil
	})
}
