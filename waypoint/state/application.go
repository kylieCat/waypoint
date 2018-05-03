package state

type Application struct {
	Name     string   `json:"name"`
	Versions Versions `json:"versions"`
}

type Applications []*Application

func (app Applications) Len() int           { return len(app) }
func (app Applications) Swap(i, j int)      { app[i], app[j] = app[j], app[i] }
func (app Applications) Less(i, j int) bool { return app[i].Name < app[j].Name }
