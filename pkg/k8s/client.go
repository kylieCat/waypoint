package k8s

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"os/exec"

	kx "github.com/kylie-a/kx/pkg"

	"github.com/kylie-a/requests"
)

// Routes
const (
	ListPodsTemplate  = "/api/v1/namespaces/%s/pods"
	GetPodTemplate    = "/api/v1/namespaces/%s/pods/%s"
	ListNodes          = "/api/v1/nodes/"
	defaultHostPort   = 50000
	defaultTargetPort = 44134
)

var conf *kx.KubeConfig
var clusterData *kx.ClusterUserData

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
	kubeConfig string
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
	var err error

	client := &Client{
		kubeConfig: "~/.kube/config",
		namespace:  "kube-system",
		labels:     []string{"app=helm", "name=tiller"},
		http:       requests.NewClient(),
		hostPort:   defaultHostPort,
		targetPort: defaultTargetPort,
	}
	client = client.ApplyOptions(opts...)

	if clusterData, err = getClusterData(client.context); err != nil {
		log.Fatalf(err.Error())
	}
	client.token = clusterData.User.UserConf.AuthProvider.Config.AccessToken
	client.endpoint = clusterData.Cluster.ClusterConf.Server
	return client
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

	if resp, err = c.http.Get(route, requests.WithBearerToken(c.token), requests.WithQueryParams(params)); err != nil {
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

func (c *Client) ListNodes() {
	var resp *requests.Response
	var err error

	route := c.formatURL(ListNodes)

	if resp, err = c.http.Get(route, requests.WithBearerToken(c.token)); err != nil {
		log.Fatalf(err.Error())
	}
	var nodes map[string]interface{}
	if err = resp.JSON(&nodes); err != nil {
		log.Fatalf(err.Error())
	}
	fmt.Println(nodes)
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
		//fmt.Printf("[debug] SERVER: localhost:%s\n", c.hostPort)
		//fmt.Printf("[debug] ARGS: %v\n", args)
		if err := exec.CommandContext(ctx, "kubectl", args...).Run(); err != nil {
			fmt.Printf("ERROR: %s\n", err.Error())
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

func getClusterData(ctxName string) (*kx.ClusterUserData, error) {
	var err error

	if conf == nil {
		if conf, err = kx.GetDefaultKubeConfig(); err != nil {
			log.Fatalf("couldn't get kube config: %s", err.Error())
		}

		if clusterData, err = conf.GetClusterUserData(ctxName); err != nil {
			log.Fatalf("no cluster with name %s: %s", ctxName, err.Error())
		}
	}
	return clusterData, nil
}
