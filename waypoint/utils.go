package waypoint

import (
	"strconv"
	"strings"
)

type ReleaseType string

const (
	Major ReleaseType = "major"
	Minor ReleaseType = "minor"
	Patch ReleaseType = "patch"
)

func GetPartsFromSemVer(semver string) ([]int, error) {
	parts := make([]int, 0)
	for _, part := range strings.Split(semver, ".") {
		p, err := strconv.Atoi(part)
		if err != nil {
			return nil, err
		}
		parts = append(parts, p)
	}
	return parts, nil
}
