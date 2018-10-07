package helm

import (
	"k8s.io/helm/pkg/helm"
	"github.com/pkg/errors"
	"crypto/tls"
	"github.com/golang/protobuf/proto"
	"github.com/mitchellh/go-homedir"
	"k8s.io/helm/pkg/helm/helmpath"
	"context"
)

type HelmOption func(client *Client)

func HelmHost(host string) HelmOption {
	return func(client *Client) {
		client.tillerOpts = append(client.tillerOpts, helm.Host(host))
	}
}

func HelmWithTLS(cfg *tls.Config) HelmOption {
	return func(client *Client) {
		client.tillerOpts = append(client.tillerOpts, helm.WithTLS(cfg))
	}
}

func HelmBeforeCall(fn func(context.Context, proto.Message) error) HelmOption {
	return func(client *Client) {
		client.tillerOpts = append(client.tillerOpts, helm.BeforeCall(fn))
	}
}

func HelmConnectTimeout(timeout int64) HelmOption {
	return func(client *Client) {
		client.tillerOpts = append(client.tillerOpts, helm.ConnectTimeout(timeout))
	}
}

func HelmHome(value string) HelmOption {
	return func(client *Client) {
		path, _ := homedir.Expand(value)
		client.env.Home = helmpath.Home(path)
	}
}

func HelmToken(value string) HelmOption {
	return func(client *Client) {
		client.token = value
	}
}

type InstallOption interface {
	Get() helm.InstallOption
	Set(value interface{}) error
}

type ValueOverrides struct {
	opt helm.InstallOption
}

func (vo ValueOverrides) Get() helm.InstallOption {
	return vo.opt
}

func (vo ValueOverrides) Set(raw interface{}) error {
	if value, ok := raw.([]byte); ok {
		vo.opt = helm.ValueOverrides(value)
	}
	return errors.New("wrong type")
}

type ReleaseName struct {
	opt helm.InstallOption
}

func (r ReleaseName) Get() helm.InstallOption {
	return r.opt
}

func (r ReleaseName) Set(raw interface{}) error {
	if value, ok := raw.(string); ok {
		r.opt = helm.ReleaseName(value)
	}
	return errors.New("can't use value")
}

type InstallTimeout struct {
	opt helm.InstallOption
}

func (i InstallTimeout) Get() helm.InstallOption {
	return i.opt
}

func (i InstallTimeout) Set(raw interface{}) error {
	if value, ok := raw.(int64); ok {
		i.opt = helm.InstallTimeout(value)
	}
	return errors.New("can't use value")
}

type InstallWait struct {
	opt helm.InstallOption
}

func (i InstallWait) Get() helm.InstallOption {
	return i.opt
}

func (i InstallWait) Set(raw interface{}) error {
	if value, ok := raw.(bool); ok {
		i.opt = helm.InstallWait(value)
	}
	return errors.New("can't use value")
}

type InstallDescription struct {
	opt helm.InstallOption
}

func (i InstallDescription) Get() helm.InstallOption {
	return i.opt
}

func (i InstallDescription) Set(raw interface{}) error {
	if value, ok := raw.(string); ok {
		i.opt = helm.InstallDescription(value)
	}
	return errors.New("can't use value")
}

type InstallDryRun struct {
	opt helm.InstallOption
}

func (i InstallDryRun) Get() helm.InstallOption {
	return i.opt
}

func (i InstallDryRun) Set(raw interface{}) error {
	if value, ok := raw.(bool); ok {
		i.opt = helm.InstallDryRun(value)
	}
	return errors.New("can't use value")
}

type InstallDisableHooks struct {
	opt helm.InstallOption
}

func (i InstallDisableHooks) Get() helm.InstallOption {
	return i.opt
}

func (i InstallDisableHooks) Set(raw interface{}) error {
	if value, ok := raw.(bool); ok {
		i.opt = helm.InstallDisableHooks(value)
	}
	return errors.New("can't use value")
}

type InstallDisableCRDHook struct {
	opt helm.InstallOption
}

func (i InstallDisableCRDHook) Get() helm.InstallOption {
	return i.opt
}

func (i InstallDisableCRDHook) Set(raw interface{}) error {
	if value, ok := raw.(bool); ok {
		i.opt = helm.InstallDisableCRDHook(value)
	}
	return errors.New("can't use value")
}

type InstallReuseName struct {
	opt helm.InstallOption
}

func (i InstallReuseName) Get() helm.InstallOption {
	return i.opt
}

func (i InstallReuseName) Set(raw interface{}) error {
	if value, ok := raw.(bool); ok {
		i.opt = helm.InstallReuseName(value)
	}
	return errors.New("can't use value")
}

type optionsMap map[string]InstallOption

func (o optionsMap) Get(optName string) InstallOption {
	return o[optName]
}

var optMap = optionsMap{
	"valueOverrides": ValueOverrides{},
	"releaseName": ReleaseName{},
	"installTimeout": InstallTimeout{},
	"installWait": InstallWait{},
	"installDescription": InstallDescription{},
	"installDryRun": InstallDryRun{},
	"installDisableHooks": InstallDisableHooks{},
	"installDisableCRDHook": InstallDisableCRDHook{},
	"installReuseName": InstallReuseName{},
}

func(o optionsMap) getOptions(args map[string]interface{}) []helm.InstallOption {
	out := make([]helm.InstallOption, 0)
	for key, value := range args {
		option, ok := o[key]
		if ok {
			if err := option.Set(value); err == nil {
				out = append(out, option.Get())
			}
		}
	}
	return out
}
