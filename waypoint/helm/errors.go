package helm

// HelmRepoNotFoundError is a custom error type.
type HelmRepoNotFoundError struct {
	Message string
}

func (e HelmRepoNotFoundError) Error() string {
	return e.Message
}

// NewHelmRepoNotFoundError creates a new `HelmRepoNotFoundError`.
func NewHelmRepoNotFoundError(name string) HelmRepoNotFoundError {
	message := "no repo found for name: " + name
	return HelmRepoNotFoundError{Message: message}
}

// HelmRepoFileLoadError is a custom error type.
type HelmRepoFileLoadError struct {
	Message string
}

func (e HelmRepoFileLoadError) Error() string {
	return e.Message
}

// NewHelmRepoFileLoadError creates a new `HelmRepoFileLoadError`.
func NewHelmRepoFileLoadError(message string) HelmRepoFileLoadError {
	return HelmRepoFileLoadError{Message: message}
}

// NoHelmReposError is a custom error type.
type NoHelmReposError struct {
	Message string
}

func (e NoHelmReposError) Error() string {
	return e.Message
}

// NewNoHelmReposError creates a new `NoHelmReposError`.
func NewNoHelmReposError() NoHelmReposError {
	return NoHelmReposError{Message: "no repositories to show"}
}

// HelmUploadError is a custom error type.
type HelmUploadError struct {
	Message string
}

func (e HelmUploadError) Error() string {
	return e.Message
}

// NewHelmUploadError creates a new `HelmUploadError`.
func NewHelmUploadError(message string) HelmUploadError {
	return HelmUploadError{Message: "error adding chart: " + message}
}

// HelmDeleteError is a custom error type.
type HelmDeleteError struct {
	Message string
}

func (e HelmDeleteError) Error() string {
	return e.Message
}

// NewHelmDeleteError creates a new `HelmDeleteError`.
func NewHelmDeleteError(message string) HelmDeleteError {
	return HelmDeleteError{Message: "error deleting previous chart: " + message}
}