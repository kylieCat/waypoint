package pkg

type TillerConf struct {
	Namespace string   `json:"namespace" yaml:"namespace"`
	Context   string   `json:"context" yaml:"context"`
	Endpoint  string   `json:"endpoint" yaml:"endpoint"`
	Labels    []string `json:"labels" yaml:"labels"`
}

func (t TillerConf) auxConf() AuxConf {
	return newAuxTillerConf()
}

func (t *TillerConf) UnmarshalJSON(data []byte) error {
	var unmarshaled ConfUnmarshaler
	var err error

	if unmarshaled, err = defaultUnmarshalJSON(data, t); err != nil {
		return err
	}
	concrete := unmarshaled.(TillerConf)
	*t = concrete
	return nil
}

func (t *TillerConf) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var unmarshaled ConfUnmarshaler
	var err error

	if unmarshaled, err = defaultUnmarshalYAML(unmarshal, t); err != nil {
		return err
	}
	concrete := unmarshaled.(TillerConf)
	*t = concrete
	return nil
}

func (t *TillerConf) applyDefaults(defaults Deployment) {
	if t.Namespace == "" {
		t.Namespace = defaults.Tiller.Namespace
	}
	if t.Context == "" {
		t.Context = defaults.Tiller.Context
	}
	if t.Endpoint == "" {
		t.Endpoint = defaults.Tiller.Endpoint
	}
	if t.Labels == nil || len(t.Labels) == 0 {
		t.Labels = defaults.Tiller.Labels
	}
}

type auxTillerConf struct {
	Namespace string   `json:"namespace" yaml:"namespace"`
	Context   string   `json:"context" yaml:"context"`
	Endpoint  string   `json:"endpoint" yaml:"endpoint"`
	Labels    []string `json:"labels" yaml:"labels"`
}

func (t auxTillerConf) conf(data map[string]interface{}) ConfUnmarshaler {
	for key, value := range data {
		switch key {
		case "namespace":
			t.Namespace = value.(string)
		case "context":
			t.Context = value.(string)
		case "endpoint":
			t.Endpoint = value.(string)
		case "labels":
			labels := toStringSlice(value.([]interface{}))
			t.Labels = labels
		}
	}
	return TillerConf{
		Namespace: t.Namespace,
		Context:   t.Context,
		Endpoint:  t.Endpoint,
		Labels:    t.Labels,
	}
}

func newAuxTillerConf() AuxConf {
	return auxTillerConf{
		Namespace: "kube-system",
		Endpoint:  "http://localhost:50000",
		Labels:    []string{"cmd=helm", "name=tiller"},
	}
}

func DefaultTillerConf() TillerConf {
	return TillerConf{
		Namespace: "kube-system",
		Endpoint:  "http://localhost:50000",
		Labels:    []string{"app=helm", "name=tiller"},
	}
}


