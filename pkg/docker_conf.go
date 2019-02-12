package pkg

type DockerConf struct {
	Repo    string `json:"repo" yaml:"repo"`
	Creds   string `json:"creds" yaml:"creds"`
	Context string `json:"context" yaml:"context"`
	File    string `json:"file" yaml:"file"`
}

func (d DockerConf) auxConf() AuxConf {
	return newAuxDockerConf()
}

func (d *DockerConf) UnmarshalJSON(data []byte) error {
	var unmarshaled ConfUnmarshaler
	var err error

	if unmarshaled, err = defaultUnmarshalJSON(data, d); err != nil {
		return err
	}
	concrete := unmarshaled.(DockerConf)
	*d = concrete
	return nil
}

func (d *DockerConf) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var unmarshaled ConfUnmarshaler
	var err error

	if unmarshaled, err = defaultUnmarshalYAML(unmarshal, d); err != nil {
		return err
	}
	concrete := unmarshaled.(DockerConf)
	*d = concrete
	return nil
}

func (d *DockerConf) applyDefaults(defaults Deployment) {
	if d.Repo == "" {
		d.Repo = defaults.Docker.Repo
	}
	if d.Creds == "" {
		d.Repo = defaults.Docker.Creds
	}
	if d.Context == "" {
		d.Repo = defaults.Docker.Context
	}
	if d.File == "" {
		d.Repo = defaults.Docker.File
	}
}

type auxDockerConf struct {
	Repo    string `json:"repo" yaml:"repo,omitempty"`
	Creds   string `json:"creds" yaml:"creds,omitempty"`
	Context string `json:"context" yaml:"context,omitempty"`
	File    string `json:"file" yaml:"file,omitempty"`
}

func newAuxDockerConf() auxDockerConf {
	return auxDockerConf{
		Repo:    "",
		Creds:   "docker-credential-gcr",
		Context: ".",
		File:    "Dockerfile",
	}
}

func (dc auxDockerConf) conf(data map[string]interface{}) ConfUnmarshaler {
	for key, value := range data {
		switch key {
		case "repo":
			dc.Repo = value.(string)
		case "creds":
			dc.Creds = value.(string)
		case "context":
			dc.Context = value.(string)
		case "file":
			dc.File = value.(string)
		}
	}
	return DockerConf{
		Repo:    dc.Repo,
		Creds:   dc.Creds,
		Context: dc.Context,
		File:    dc.File,
	}
}

func DefaultDockerConf() DockerConf {
	return DockerConf{
		Repo:    "gcr.io/{{.DockerRepo}}",
		Creds:   "docker-credential-gcr",
		Context: ".",
		File:    "Dockerfile",
	}
}
