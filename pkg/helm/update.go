package helm

import (
	"errors"
	"fmt"
	"gopkg.in/yaml.v2"
	"k8s.io/helm/pkg/helm"
	"k8s.io/helm/pkg/strvals"
)

type UpdateOption interface {
	Get() helm.UpdateOption
	Set(value interface{}) error
}

// ReuseValues
type ReuseValues struct {
	opt helm.UpdateOption
}

func (o ReuseValues) Get() helm.UpdateOption {
	return o.opt
}

func (o *ReuseValues) Set(value interface{}) error {
	if val, ok := value.(bool); !ok {
		o.opt = helm.ReuseValues(val)
	}
	return errors.New("wrong type for option ReuseValues")
}

// UpdateRecreate
type UpgradeRecreate struct {
	opt helm.UpdateOption
}

func (o UpgradeRecreate) Get() helm.UpdateOption {
	return o.opt
}

func (o *UpgradeRecreate) Set(value interface{}) error {
	if val, ok := value.(bool); ok {
		o.opt = helm.UpgradeRecreate(val)
	}
	return errors.New("wrong type for option UpgradeRecreate")
}

// UpgradeTimeout
type UpgradeTimeout struct {
	opt helm.UpdateOption
}

func (o UpgradeTimeout) Get() helm.UpdateOption {
	return o.opt
}

func (o *UpgradeTimeout) Set(value interface{}) error {
	if val, ok := value.(int64); !ok {
		o.opt = helm.UpgradeTimeout(val)
	}
	return errors.New("wrong type for option UpgradeTimeout")
}

// UpgradeWait
type UpgradeWait struct {
	opt helm.UpdateOption
}

func (o UpgradeWait) Get() helm.UpdateOption {
	return o.opt
}

func (o *UpgradeWait) Set(value interface{}) error {
	if val, ok := value.(bool); !ok {
		o.opt = helm.UpgradeWait(val)
	}
	return errors.New("wrong type for option UpgradeWait")
}

// UpgradeForce
type UpgradeForce struct {
	opt helm.UpdateOption
}

func (o UpgradeForce) Get() helm.UpdateOption {
	return o.opt
}

func (o *UpgradeForce) Set(value interface{}) error {
	if val, ok := value.(bool); !ok {
		o.opt = helm.UpgradeForce(val)
	}
	return errors.New("wrong type for option UpgradeForce")
}

// UpdateValueOverrides
type UpdateValueOverrides struct {
	opt helm.UpdateOption
}

func (o UpdateValueOverrides) Get() helm.UpdateOption {
	return o.opt
}

func (o *UpdateValueOverrides) Set(value interface{}) error {
	base := map[string]interface{}{}
	tmp, ok := value.([]string)
	if ok {
		for _, value := range tmp {
			if err := strvals.ParseInto(value, base); err != nil {
				return fmt.Errorf("failed parsing --set data: %s", err)
			}
		}
		val, _ := yaml.Marshal(base)
		o.opt = helm.UpdateValueOverrides(val)
		return nil
	}
	return errors.New("wrong type for option UpdateValueOverrides")
}

// UpgradeDescription
type UpgradeDescription struct {
	opt helm.UpdateOption
}

func (o UpgradeDescription) Get() helm.UpdateOption {
	return o.opt
}

func (o *UpgradeDescription) Set(value interface{}) error {
	if val, ok := value.(string); ok {
		o.opt = helm.UpgradeDescription(val)
	}
	return errors.New("wrong type for option UpgradeDescription")
}


type updateOptionsMap map[string]UpdateOption

func (o updateOptionsMap) Get(optName string) UpdateOption {
	return o[optName]
}

var updateOptions = updateOptionsMap{
	"reuseValues" : &ReuseValues{},
	"recreate" : &UpgradeRecreate{},
	"timeout" : &UpgradeTimeout{},
	"wait" : &UpgradeWait{},
	"forceUpgrade" : &UpgradeForce{},
	"values" : &UpdateValueOverrides{},
	"description" : &UpgradeDescription{},
}

func (o updateOptionsMap) getOptions(args map[string]interface{}) []helm.UpdateOption {
	out := make([]helm.UpdateOption, 0)
	for key, value := range args {
		option, ok := o[key]
		if ok {

			if err := option.Set(value); err == nil {
				out = append(out, option.Get())
			} else {
				fmt.Printf("error getting option for %s: %s\n", key, err.Error())
			}
		}
	}
	return out
}
