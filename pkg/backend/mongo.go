package backend

import "github.com/kylie-a/waypoint/pkg"

type WaypointStoreMongo struct {
}

func NewWaypointStoreMongo(conf *pkg.Config) *WaypointStoreMongo {
	return &WaypointStoreMongo{}
}

func (ws *WaypointStoreMongo) GetLatest(app string) (*pkg.Version, error) {
	panic("implement me")
}

func (ws *WaypointStoreMongo) All(app string) (pkg.Versions, error) {
	panic("implement me")
}

func (ws *WaypointStoreMongo) Save(app string, version *pkg.Version) error {
	panic("implement me")
}

func (ws *WaypointStoreMongo) AddApplication(name string, initialVersion string) error {
	panic("implement me")
}
