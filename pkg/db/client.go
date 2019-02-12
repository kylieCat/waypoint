package db

import (
	"github.com/kylie-a/waypoint/pkg"
)

type Client struct {
	pkg.IStorage
}

type Auth struct {
	pkg.BackendAuthConf
}

func NewClient(conf *pkg.Config) (*Client, error) {
	client := &Client{}
	switch conf.Backend.Kind {
	case pkg.DataStore:
		client.IStorage = NewWaypointStoreDS(conf)
	case pkg.Bolt:
		client.IStorage = NewWaypointStoreBolt(conf.Backend.Conf["dbPath"])
	case pkg.MongoDB:
		client.IStorage = NewWaypointStoreMongo(conf)
	case pkg.Dynamo:
		client.IStorage = NewWaypointStoreDDB(conf)
	default:
		return nil, NewUnkownBackendError(string(conf.Backend.Kind))
	}
	return client, nil
}
