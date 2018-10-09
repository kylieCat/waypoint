package pkg

import (
	"github.com/kylie-a/waypoint/pkg/docker"
	"github.com/kylie-a/waypoint/pkg/helm"
	"github.com/kylie-a/waypoint/pkg/k8s"
)

type ReleaseOption func(release *Release)

func Conf(value *Config) ReleaseOption {
	return func(release *Release) {
		release.conf = value
	}
}

func Deploy(value Deployment) ReleaseOption {
	return func(release *Release) {
		release.deploy = value
	}
}

func Type(value ReleaseType) ReleaseOption {
	return func(release *Release) {
		release.typ = value
	}
}

func PrevVersion(value *Version) ReleaseOption {
	return func(release *Release) {
		release.prevVersion = value
	}
}

func CurrentVersion(value *Version) ReleaseOption {
	return func(release *Release) {
		release.newVersion = value
	}
}

func Docker(value *docker.Client) ReleaseOption {
	return func(release *Release) {
		release.docker = value
	}
}

func Helm(value *helm.Client) ReleaseOption {
	return func(release *Release) {
		release.helm = value
	}
}

func K8s(value *k8s.Client) ReleaseOption {
	return func(release *Release) {
		release.k8s = value
	}
}

func DB(value BackendService) ReleaseOption {
	return func(release *Release) {
		release.ws = value
	}
}