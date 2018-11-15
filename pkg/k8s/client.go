package k8s

import (
	"context"
	"fmt"
	"net/url"
	"os/exec"

	"io/ioutil"

	"github.com/kylie-a/requests"
	"github.com/mitchellh/go-homedir"
	"gopkg.in/yaml.v2"
)

// Routes
const (
	ListPodsTemplate = "/api/v1/namespaces/%s/pods"
	GetPodTemplate   = "/api/v1/namespaces/%s/pods/%s"
)

type Metadata struct {
	Name string `json:"name" yaml:"name"`
}

type Pod struct {
	Metadata Metadata `json:"metadata" yaml:"metadata"`
}

type ListPodsResponse struct {
	Items []Pod `json:"items" yaml:"items"`
}

type Client struct {
	endpoint   string
	namespace  string
	context    string
	labels     []string
	token      string
	http       *requests.Client
	hostPort   int
	targetPort int
}

func NewClient(opts ...Option) *Client {
	client := &Client{
		namespace:  "kube-system",
		labels:     []string{"app=helm", "name=tiller"},
		http:       requests.NewClient(),
		hostPort:   8081,
		targetPort: 44134,
	}
	return client.ApplyOptions(opts...)
}

func (c *Client) ApplyOptions(opts ...Option) *Client {
	for _, opt := range opts {
		opt(c)
	}
	return c
}

func (c *Client) GetTillerPod() (string, error) {
	var resp *requests.Response
	var err error

	route := c.formatURL(ListPodsTemplate, c.namespace)
	params := c.getParamsFromLabels()
	fmt.Println(getAccessToken(c.context))
	if resp, err = c.http.Get(route, requests.WithBearerToken(getAccessToken(c.context)), requests.WithQueryParams(params)); err != nil {
		return "", err
	}
	var listResp ListPodsResponse
	if err = resp.JSON(&listResp); err != nil {
		return "", err
	}
	if len(listResp.Items) != 1 {
		return "", NewNoPodsFoundError(params)
	}
	return listResp.Items[0].Metadata.Name, nil
}

func (c *Client) StartForwarder() (context.CancelFunc, error) {
	var podName string
	var err error

	ctx, cancel := context.WithCancel(context.Background())

	if podName, err = c.GetTillerPod(); err != nil {
		return nil, err
	}
	args := []string{
		"-n",
		c.namespace,
		"--context",
		c.context,
		"port-forward",
		podName,
		fmt.Sprintf("%d:%d", c.hostPort, c.targetPort),
	}
	go func() {
		//fmt.Println("[debug] SERVER: localhost:8081")
		//fmt.Printf("[debug] ARGS: %v\n", args)
		if err := exec.CommandContext(ctx, "kubectl", args...).Run(); err != nil {
			fmt.Println(err.Error())
			cancel()
		}
	}()
	return cancel, nil
}

func (c *Client) getParamsFromLabels() url.Values {
	vals := make(url.Values)
	for _, value := range c.labels {
		vals["labelSelector"] = append(vals["labelSelector"], value)
	}
	return vals
}

func (c *Client) formatURL(url string, args ...interface{}) string {
	url = fmt.Sprintf(url, args...)
	return fmt.Sprintf("%s%s", c.endpoint, url)
}

var conf *k8sConf

type clusterConf struct {
	CertificateAuthorityData string `yaml:"certificate-authority-data"`
	Server                   string `yaml:"server"`
}

type cluster struct {
	Cluster clusterConf `yaml:"cluster"`
	Name    string      `yaml:"name"`
}

type k8sContextConf struct {
	Cluster   string `yaml:"cluster"`
	Namespace string `yaml:"namespace"`
	User      string `yaml:"user"`
}

type k8sContext struct {
	Context k8sContextConf `yaml:"context"`
	Name    string         `yaml:"name"`
}

type authProviderConfig struct {
	AccessToken string `yaml:"access-token"`
	CmdArgs     string `yaml:"cmd-args"`
	CmdPath     string `yaml:"cmd-path"`
	Expiry      string `yaml:"expiry"`
	ExpiryKey   string `yaml:"expiry-key"`
	TokenKey    string `yaml:"token-key"`
}

type authProvider struct {
	Config authProviderConfig `yaml:"config"`
	Name   string             `yaml:"name"`
}

type k8sUser struct {
	AuthProvider authProvider `yaml:"auth-provider"`
}

type userConf struct {
	User k8sUser `yaml:"user"`
	Name string  `yaml:"name"`
}

type k8sConf struct {
	ApiVersion     string       `yaml:"apiVersion"`
	Clusters       []cluster    `yaml:"clusters"`
	Contexts       []k8sContext `yaml:"contexts"`
	CurrentContext string       `yaml:"current-context"`
	Kind           string       `yaml:"kind"`
	Preferences    interface{}  `yaml:"preferences"`
	Users          []userConf   `yaml:"users"`
}

func getConf() *k8sConf {
	if conf == nil {
		fileName, err := homedir.Expand("~/.kube/config")
		b, _ := ioutil.ReadFile(fileName)

		err = yaml.Unmarshal(b, &conf)
		if err != nil {
			fmt.Printf("ERROR: %s", err.Error())
			return conf
		}
	}
	return conf
}

func GetCurrentContext() string {
	conf := getConf()
	return conf.CurrentContext
}

func getAccessToken(ctxName string) string {
	conf := getConf()
	for _, ctx := range conf.Contexts {
		if ctx.Name == ctxName {
			userKey := ctx.Context.User
			for _, user := range conf.Users {
				if user.Name == userKey {
					return user.User.AuthProvider.Config.AccessToken
				}
			}
		}
	}
	return ""
}
