package backend

import (
	"context"

	"log"

	"cloud.google.com/go/datastore"
	"github.com/kylie-a/waypoint/pkg"
	"github.com/mitchellh/go-homedir"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

type WaypointStoreDS struct {
	client *datastore.Client
	auth   pkg.BackendAuthConf
}

func NewWaypointStoreDS(conf *pkg.Config) *WaypointStoreDS {
	client, err := datastore.NewClient(context.Background(), conf.Backend.Conf["project"], getAuth(conf))
	if err != nil {
		panic(err.Error())
	}
	return &WaypointStoreDS{client: client}
}

func (ds WaypointStoreDS) GetLatest(app string) (*pkg.Version, error) {
	parentKey := datastore.NameKey("application", app, nil)
	q := datastore.NewQuery("release").Ancestor(parentKey).Order("-Timestamp").Limit(1)
	iter := ds.client.Run(context.Background(), q)
	var version pkg.Version
	_, err := iter.Next(&version)
	if err == iterator.Done {
		return &version, nil
	}
	if err != nil {
		return nil, err
	}
	return &version, nil
}

func (ds WaypointStoreDS) All(app string) (pkg.Versions, error) {
	parentKey := datastore.NameKey("application", app, nil)
	q := datastore.NewQuery("release").Ancestor(parentKey)
	iter := ds.client.Run(context.Background(), q)
	versions := make(pkg.Versions, 0)
	var version pkg.Version
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

func (ds WaypointStoreDS) Save(app string, version *pkg.Version) error {
	parentKey := datastore.NameKey("application", app, nil)
	key := datastore.NameKey("release", version.SemVer(), parentKey)
	_, err := ds.client.Put(context.Background(), key, version)
	return err
}

func (ds WaypointStoreDS) AddApplication(name string, initialVersion string) error {
	key := datastore.NameKey("application", name, nil)
	app := &pkg.Application{Name: name}
	if _, err := ds.client.Put(context.Background(), key, app); err != nil {
		return err
	}
	parts, _ := pkg.GetPartsFromSemVer(initialVersion)
	version := pkg.NewVersion(parts[0], parts[1], parts[2])
	return ds.Save(app.Name, &version)
}

func getAuth(conf *pkg.Config) option.ClientOption {
	c := conf.Backend.Conf
	switch pkg.GCPAuthKind(c["kind"]) {
	case pkg.CredsFile:
		f, _ := homedir.Expand(c["value"])
		return option.WithCredentialsFile(f)
	case pkg.ApiKey:
		return option.WithAPIKey(c["value"])
	}
	return nil
}
