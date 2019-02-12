package pkg

//type Args struct {
//	Name           *string   `json:"name,omitempty" yaml:"name,omitempty" `
//	Namespace      *string   `json:"namespace,omitempty" yaml:"namespace,omitempty" `
//	ValueFiles     []string  `json:"valueFiles,omitempty" yaml:"valueFiles,omitempty" `
//	ChartPath      *string   `json:"chartPath,omitempty" yaml:"chartPath,omitempty" `
//	DryRun         *bool     `json:"dryRun,omitempty" yaml:"dryRun,omitempty" `
//	DisableHooks   *bool     `json:"disableHooks,omitempty" yaml:"disableHooks,omitempty" `
//	DisableCRDHook *bool     `json:"disableCrdHook,omitempty" yaml:"disableCrdHook,omitempty" `
//	Replace        *bool     `json:"replace,omitempty" yaml:"replace,omitempty" `
//	Verify         *bool     `json:"verify,omitempty" yaml:"verify,omitempty" `
//	Keyring        *string   `json:"keyring,omitempty" yaml:"keyring,omitempty" `
//	Out            io.Writer `json:"out,omitempty" yaml:"out,omitempty" `
//	Values         []string  `json:"values,omitempty" yaml:"values,omitempty" `
//	StringValues   []string  `json:"stringValues,omitempty" yaml:"stringValues,omitempty" `
//	FileValues     []string  `json:"fileValues,omitempty" yaml:"fileValues,omitempty" `
//	NameTemplate   *string   `json:"nameTemplate,omitempty" yaml:"nameTemplate,omitempty" `
//	Version        *string   `json:"version,omitempty" yaml:"version,omitempty" `
//	Timeout        *int64    `json:"timeout,omitempty" yaml:"timeout,omitempty" `
//	Wait           *bool     `json:"wait,omitempty" yaml:"wait,omitempty" `
//	RepoURL        *string   `json:"repoUrl,omitempty" yaml:"repoUrl,omitempty" `
//	Username       *string   `json:"username,omitempty" yaml:"username,omitempty" `
//	Password       *string   `json:"password,omitempty" yaml:"password,omitempty" `
//	Devel          *bool     `json:"devel,omitempty" yaml:"devel,omitempty" `
//	DepUp          *bool     `json:"depUp,omitempty" yaml:"depUp,omitempty" `
//	Description    *string   `json:"description,omitempty" yaml:"description,omitempty" `
//	CertFile       *string   `json:"certFile,omitempty" yaml:"certFile,omitempty" `
//	KeyFile        *string   `json:"keyFile,omitempty" yaml:"keyFile,omitempty" `
//	CaFile         *string   `json:"caFile,omitempty" yaml:"caFile,omitempty" `
//}
//
//func (a Args) auxConf() AuxConf {
//	return newAuxArgs()
//}

//func (a *Args) applyDefaults(defaults Deployment) {
//	if a.Name == nil {
//		a.Name = defaults.Helm.Args.Name
//	}
//	if a.Namespace == nil {
//		a.Namespace = defaults.Helm.Args.Namespace
//	}
//	if a.ValueFiles == nil {
//		a.ValueFiles = defaults.Helm.Args.ValueFiles
//	}
//	if a.ChartPath == nil {
//		a.ChartPath = defaults.Helm.Args.ChartPath
//	}
//	//if a.DryRun == bool {
//	//	a.DryRun = defaults.Helm.Args.DryRun
//	//}
//	//if a.DisableHooks == bool {
//	//	a.DisableHooks = defaults.Helm.Args.DisableHooks
//	//}
//	//if a.DisableCRDHook == bool {
//	//	a.DisableCRDHook = defaults.Helm.Args.DisableCRDHook
//	//}
//	//if a.Replace == "" {
//	//	a.Replace = defaults.Helm.Args.Replace
//	//}
//	//if a.Verify == bool {
//	//	a.Verify = defaults.Helm.Args.Verify
//	//}
//	if a.Keyring == nil {
//		a.Keyring = defaults.Helm.Args.Keyring
//	}
//	//if a.Out == "" {
//	//	a.Out = defaults.Helm.Args.Out
//	//}
//	if a.Values == nil || len(a.Values) == 0 {
//		a.Values = defaults.Helm.Args.Values
//	}
//	if a.StringValues == nil || len(a.StringValues) == 0 {
//		a.StringValues = defaults.Helm.Args.StringValues
//	}
//	if a.FileValues == nil || len(a.FileValues) == 0 {
//		a.FileValues = defaults.Helm.Args.FileValues
//	}
//	if a.NameTemplate == nil {
//		a.NameTemplate = defaults.Helm.Args.NameTemplate
//	}
//	if a.Version == nil {
//		a.Version = defaults.Helm.Args.Version
//	}
//	if a.Timeout == nil {
//		a.Timeout = defaults.Helm.Args.Timeout
//	}
//	//if a.Wait == bool {
//	//	a.Wait = defaults.Helm.Args.Wait
//	//}
//	if a.RepoURL == nil {
//		a.RepoURL = defaults.Helm.Args.RepoURL
//	}
//	if a.Username == nil {
//		a.Username = defaults.Helm.Args.Username
//	}
//	if a.Password == nil {
//		a.Password = defaults.Helm.Args.Password
//	}
//	//if a.Devel == bool {
//	//	a.Devel = defaults.Helm.Args.Devel
//	//}
//	//if a.DepUp == bool {
//	//	a.DepUp = defaults.Helm.Args.DepUp
//	//}
//	if a.Description == nil {
//		a.Description = defaults.Helm.Args.Description
//	}
//	if a.CertFile == nil {
//		a.CertFile = defaults.Helm.Args.CertFile
//	}
//	if a.KeyFile == nil {
//		a.KeyFile = defaults.Helm.Args.KeyFile
//	}
//	if a.CaFile == nil {
//		a.CaFile = defaults.Helm.Args.CaFile
//	}
//}

//type auxArgs struct {
//	Name           *string   `json:"name" yaml:"name" `
//	Namespace      *string   `json:"namespace" yaml:"namespace" `
//	ValueFiles     []string  `json:"valueFiles" yaml:"valueFiles" `
//	ChartPath      *string   `json:"chartPath" yaml:"chartPath" `
//	DryRun         *bool     `json:"dryRun" yaml:"dryRun" `
//	DisableHooks   *bool     `json:"disableHooks" yaml:"disableHooks" `
//	DisableCRDHook *bool     `json:"disableCrdHook" yaml:"disableCrdHook" `
//	Replace        *bool     `json:"replace" yaml:"replace" `
//	Verify         *bool     `json:"verify" yaml:"verify" `
//	Keyring        *string   `json:"keyring" yaml:"keyring" `
//	Out            io.Writer `json:"out" yaml:"out" `
//	Values         []string  `json:"values" yaml:"values" `
//	StringValues   []string  `json:"stringValues" yaml:"stringValues" `
//	FileValues     []string  `json:"fileValues" yaml:"fileValues" `
//	NameTemplate   *string   `json:"nameTemplate" yaml:"nameTemplate" `
//	Version        *string   `json:"version" yaml:"version" `
//	Timeout        *int64    `json:"timeout" yaml:"timeout" `
//	Wait           *bool     `json:"wait" yaml:"wait" `
//	RepoURL        *string   `json:"repoUrl" yaml:"repoUrl" `
//	Username       *string   `json:"username" yaml:"username" `
//	Password       *string   `json:"password" yaml:"password" `
//	Devel          *bool     `json:"devel" yaml:"devel" `
//	DepUp          *bool     `json:"depUp" yaml:"depUp" `
//	Description    *string   `json:"description" yaml:"description" `
//	CertFile       *string   `json:"certFile" yaml:"certFile" `
//	KeyFile        *string   `json:"keyFile" yaml:"keyFile" `
//	CaFile         *string   `json:"caFile" yaml:"caFile" `
//}
//
//func newAuxArgs() auxArgs {
//	return auxArgs{}
//}
//
//func (a auxArgs) conf(data map[string]interface{}) ConfUnmarshaler {
//	for key, value := range data {
//		switch key {
//		case "name":
//			*a.Name = value.(string)
//		case "namespace":
//			*a.Namespace = value.(string)
//		case "valueFiles":
//			a.ValueFiles = value.([]string)
//		case "chartPath":
//			*a.ChartPath = value.(string)
//		case "dryRun":
//			*a.DryRun = value.(bool)
//		case "disableHooks":
//			*a.DisableHooks = value.(bool)
//		case "disableCrdHook":
//			*a.DisableCRDHook = value.(bool)
//		case "replace":
//			*a.Replace = value.(bool)
//		case "verify":
//			*a.Verify = value.(bool)
//		case "keyring":
//			*a.Keyring = value.(string)
//		case "out":
//			a.Out = value.(io.Writer)
//		case "values":
//			a.Values = value.([]string)
//		case "stringValues":
//			a.StringValues = value.([]string)
//		case "fileValues":
//			a.FileValues = value.([]string)
//		case "nameTemplate":
//			*a.NameTemplate = value.(string)
//		case "version":
//			a.Version = value.(*string)
//		case "timeout":
//			*a.Timeout = value.(int64)
//		case "wait":
//			*a.Wait = value.(bool)
//		case "repoUrl":
//			*a.RepoURL = value.(string)
//		case "username":
//			*a.Username = value.(string)
//		case "password":
//			*a.Password = value.(string)
//		case "devel":
//			*a.Devel = value.(bool)
//		case "depUp":
//			*a.DepUp = value.(bool)
//		case "description":
//			*a.Description = value.(string)
//		case "certFile":
//			*a.CertFile = value.(string)
//		case "keyFile":
//			*a.KeyFile = value.(string)
//		case "caFile":
//			*a.CaFile = value.(string)
//		}
//	}
//	return Args{
//		Name:           a.Name,
//		Namespace:      a.Namespace,
//		ValueFiles:     a.ValueFiles,
//		ChartPath:      a.ChartPath,
//		DryRun:         a.DryRun,
//		DisableHooks:   a.DisableHooks,
//		DisableCRDHook: a.DisableCRDHook,
//		Replace:        a.Replace,
//		Verify:         a.Verify,
//		Keyring:        a.Keyring,
//		Out:            a.Out,
//		Values:         a.Values,
//		StringValues:   a.StringValues,
//		FileValues:     a.FileValues,
//		NameTemplate:   a.NameTemplate,
//		Version:        a.Version,
//		Timeout:        a.Timeout,
//		Wait:           a.Wait,
//		RepoURL:        a.RepoURL,
//		Username:       a.Username,
//		Password:       a.Password,
//		Devel:          a.Devel,
//		DepUp:          a.DepUp,
//		Description:    a.Description,
//		CertFile:       a.CertFile,
//		KeyFile:        a.KeyFile,
//		CaFile:         a.CaFile,
//	}
//}

type Args map[string]interface{}

func (a Args) applyDefaults(defaults Deployment) Args {
	if a == nil {
		a = make(Args)
	}
	def := Args{
		"values": []string{"image.tag=0.7.1"},
		"name": "wayex",
		"version": "0.7.1",
	}
	for key, value := range def {
		if _, ok := a[key]; !ok {
			a[key] = value
		}
	}
	return a
}

type HelmConf struct {
	Name     string `json:"name" yaml:"name"`
	ChartDir string `json:"chartDir" yaml:"chartDir"`
	DestDir  string `json:"destDir" yaml:"destDir"`
	Save     bool   `json:"save" yaml:"save"`
	Args     Args   `json:"args" yaml:"args"`
}

func (h HelmConf) auxConf() AuxConf {
	return newAuxHelmConf()
}

func (h *HelmConf) UnmarshalJSON(data []byte) error {
	var unmarshaled ConfUnmarshaler
	var err error

	if unmarshaled, err = defaultUnmarshalJSON(data, h); err != nil {
		return err
	}
	concrete := unmarshaled.(HelmConf)
	*h = concrete
	return nil
}

func (h *HelmConf) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var unmarshaled ConfUnmarshaler
	var err error

	if unmarshaled, err = defaultUnmarshalYAML(unmarshal, h); err != nil {
		return err
	}
	concrete := unmarshaled.(HelmConf)
	*h = concrete
	return nil
}

func (h *HelmConf) applyDefaults(defaults Deployment) {
	if h.Name == "" {
		h.Name = defaults.Helm.Name
	}
	if h.ChartDir == "" {
		h.ChartDir = defaults.Helm.ChartDir
	}
	if h.DestDir == "" {
		h.DestDir = defaults.Helm.DestDir
	}
	if !h.Save {
		h.Save = defaults.Helm.Save
	}
	h.Args = h.Args.applyDefaults(defaults)
}

type auxHelmConf struct {
	Name     string `json:"name" yaml:"name"`
	ChartDir string `json:"chartDir" yaml:"chartDir"`
	DestDir  string `json:"destDir" yaml:"destDir"`
	Save     bool   `json:"save" yaml:"save"`
	Args     Args   `json:"args" yaml:"args"`
}

func (h auxHelmConf) conf(data map[string]interface{}) ConfUnmarshaler {
	for key, value := range data {
		switch key {
		case "name":
			h.Name = value.(string)
		case "chartDir":
			h.ChartDir = value.(string)
		case "destDir":
			h.DestDir = value.(string)
		case "save":
			h.Save = value.(bool)
			//case "args":
			//	var args Args
			//	yaml.Unmarshal([]byte(value.(map[interface{}]interface{})), &args)
			//	h.Args = value.(Args)
		}
	}
	return HelmConf{
		Name:     h.Name,
		ChartDir: h.ChartDir,
		DestDir:  h.DestDir,
		Save:     h.Save,
		Args:     h.Args,
	}
}

func newAuxHelmConf() AuxConf {
	return auxHelmConf{
		ChartDir: "./deploy/{{.App}}",
		DestDir:  "./deploy",
		Save:     true,
		Args:     Args{},
	}
}

func DefaultHelmConf() HelmConf {
	return HelmConf{
		Name:     "test",
		ChartDir: "deploy/{{.App}}",
		DestDir:  "deploy/",
		Save:     true,
		Args:     Args{},
	}
}

func toStringSlice(data []interface{}) []string {
	results := make([]string, 0, len(data))
	for _, label := range data {
		results = append(results, label.(string))
	}
	return results
}
