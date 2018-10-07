package k8s

import "fmt"

// NoPodsFoundError is a custom error type.
type NoPodsFoundError struct {
	Message string
}

func (e NoPodsFoundError) Error() string {
	return e.Message
}

// NewNoPodsFoundError creates a new `NoPodsFoundError`.
func NewNoPodsFoundError(args ...interface{}) NoPodsFoundError {
	return NoPodsFoundError{Message: fmt.Sprintf("no pods found with params: %v", args...)}
}
