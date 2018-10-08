package pkg

type Application struct {
	Name     string   `json:"name"`
	Versions Versions `json:"versions"`
}

func (app Application) GetKey() []byte {
	return []byte(app.Name)
}

type Applications []*Application

func (app Applications) Each(handler func(Record) error) error {
	for _, record := range app {
		err := handler(record)
		if err != nil {
			return err
		}
	}
	return nil
}

func (app Applications) Len() int           { return len(app) }
func (app Applications) Swap(i, j int)      { app[i], app[j] = app[j], app[i] }
func (app Applications) Less(i, j int) bool { return app[i].Name < app[j].Name }
