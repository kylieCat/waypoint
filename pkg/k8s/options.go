package k8s

import (
	"net/http"
	"github.com/kylie-a/requests"
)

type Option func(client *Client)

func Endpoint(value string) Option {
	return func(client *Client) {
		client.endpoint = value
	}
}

func Namespace(value string) Option {
	return func(client *Client) {
		client.namespace = value
	}
}

func Context(value string) Option {
	return func(client *Client) {
		client.context = value
	}
}

func Labels(value []string) Option {
	return func(client *Client) {
		client.labels = value
	}
}

func Token(value string) Option {
	return func(client *Client) {
		client.token = value
	}
}

func HostPort(value int) Option {
	return func(client *Client) {
		client.hostPort = value
	}
}

func TargetPort(value int) Option {
	return func(client *Client) {
		client.targetPort = value
	}
}

func HTTPClient(value *http.Client) Option {
	return func(client *Client) {
		client.http = requests.NewClient(requests.CustomClient(value))
	}
}
