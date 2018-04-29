package state

import (
	"strconv"
	"strings"
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
