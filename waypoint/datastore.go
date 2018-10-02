package waypoint

import (
	"context"

	"cloud.google.com/go/datastore"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
	"log"
)

type WaypointStoreDS struct {
	client *datastore.Client
}

func NewWaypointStoreDS(projectID string, opts ...option.ClientOption) *WaypointStoreDS {
	client, err := datastore.NewClient(context.Background(), projectID, opts...)
	if err != nil {
		panic(err.Error())
	}
	return &WaypointStoreDS{client: client}
}

func (ds WaypointStoreDS) GetMostRecent(app string) (*Version, error) {
	q := datastore.NewQuery("release").Order("-Timestamp").Limit(1)
	iter := ds.client.Run(context.Background(), q)
	var version Version
	_, err := iter.Next(&version)
	if err == iterator.Done {
		return &version, nil
	}
	if err != nil {
		return nil, err
	}
	return &version, nil
}

func (ds WaypointStoreDS) ListAll(app string) (Versions, error) {
	parentKey := datastore.NameKey("application", app, nil)
	q := datastore.NewQuery("release").Ancestor(parentKey)
	iter := ds.client.Run(context.Background(), q)
	versions := make(Versions, 0)
	var version Version
	_, err := iter.Next(&version)
	for err == nil {
		versions = append(versions, version)
		_, err = iter.Next(&version)
	}
	if err != iterator.Done && err != nil {
		log.Fatalf("Failed fetching results: %v", err)
	}
	return versions, nil
}

func (ds WaypointStoreDS) NewVersion(app string, version *Version) error {
	parentKey := datastore.NameKey("application", app, nil)
	key := datastore.NameKey("release", version.SemVer(), parentKey)
	_, err := ds.client.Put(context.Background(), key, version)
	return err
}

func (ds WaypointStoreDS) AddApplication(name string, initialVersion string) error {
	key := datastore.NameKey("application", name, nil)
	app := &Application{Name: name}
	if _, err := ds.client.Put(context.Background(), key, app); err != nil {
		return err
	}
	parts, _ := GetPartsFromSemVer(initialVersion)
	version := NewVersion(parts[0], parts[1], parts[2])
	return ds.NewVersion(app.Name, &version)
}
