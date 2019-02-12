package db

import (
	"github.com/guregu/dynamo"
	"github.com/kylie-a/waypoint/pkg"
)

type WaypointStoreDDB struct {
	client *dynamo.DB
}

func NewWaypointStoreDDB(conf *pkg.Config) *WaypointStoreDDB {
	return &WaypointStoreDDB{}
}

func (ws *WaypointStoreDDB) GetLatest(app string) (*pkg.Version, error) {
	panic("implement me")
}

func (ws *WaypointStoreDDB) All(app string) (pkg.Versions, error) {
	panic("implement me")
}

func (ws *WaypointStoreDDB) Save(app string, version *pkg.Version) error {
	panic("implement me")
}

func (ws *WaypointStoreDDB) AddApplication(name string, initialVersion string) error {
	panic("implement me")
}
