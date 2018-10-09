package backend

// UnkownBackendError is a custom error type.
type UnkownBackendError struct {
	Message string
}

func (e UnkownBackendError) Error() string {
	return e.Message
}

// NewUnkownBackendError creates a new `UnkownBackendError`.
func NewUnkownBackendError(backend string) UnkownBackendError {
	return UnkownBackendError{Message: "unknown bqackend: " + backend}
}
