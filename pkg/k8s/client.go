package k8s

import (
	"context"
	"fmt"
	"net/url"
	"os/exec"

	"github.com/kylie-a/requests"
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
	if resp, err = c.http.Get(route, requests.WithBasicAuth(c.token), requests.WithQueryParams(params)); err != nil {
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
