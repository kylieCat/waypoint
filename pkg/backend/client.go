package backend

import (
	"github.com/kylie-a/waypoint/pkg"
)

type Client struct {
	pkg.BackendService
}

type Auth struct {
	pkg.BackendAuthConf
}

func NewClient(conf *pkg.Config) (*Client, error) {
	client := &Client{}
	switch conf.Backend.Kind {
	case pkg.DataStore:
		client.BackendService = NewWaypointStoreDS(conf)
	case pkg.Bolt:
		client.BackendService = NewWaypointStoreBolt(conf)
	case pkg.MongoDB:
		client.BackendService = NewWaypointStoreMongo(conf)
	case pkg.Dynamo:
		client.BackendService = NewWaypointStoreDDB(conf)
	default:
		return nil, NewUnkownBackendError(string(conf.Backend.Kind))
	}
	return client, nil
}
