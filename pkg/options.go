package pkg

type ReleaseOption func(release *Release)

func Deploy(value *Deployment) ReleaseOption {
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

func Docker(value IDocker) ReleaseOption {
	return func(release *Release) {
		release.docker = value
	}
}

func Helm(value IHelm) ReleaseOption {
	return func(release *Release) {
		release.helm = value
	}
}

func K8s(value IK8s) ReleaseOption {
	return func(release *Release) {
		release.k8s = value
	}
}

func DB(value IStorage) ReleaseOption {
	return func(release *Release) {
		release.db = value
	}
}
