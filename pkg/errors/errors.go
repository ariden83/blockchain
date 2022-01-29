package errors

import (
	"net/http"
)

// New returns an error that formats as the given text.
// Each call to New returns a distinct error value even if the text is identical.
func New(text string, options ...func(*errorString)) error {
	e := &errorString{s: text}
	for _, o := range options {
		o(e)
	}
	return e
}

// errorString is a trivial implementation of error.
type errorString struct {
	s      string
	status int
}

func (e *errorString) Error() string {
	return e.s
}

func (e *errorString) Status() int {
	return e.status
}

func WithStatus(status int) func(*errorString) {
	return func(e *errorString) {
		e.status = status
	}
}

func StatusCode(err error) int {
	if errWithStatus, ok := err.(*errorString); ok {
		return errWithStatus.Status()
	}
	return http.StatusUpgradeRequired
}

var (
	ErrMissingPassword         = New("missing password", WithStatus(http.StatusBadRequest))
	ErrInvalidPassword         = New("password invalid", WithStatus(http.StatusUnauthorized))
	ErrInternalError           = New("internal error", WithStatus(http.StatusForbidden))
	ErrInternalDependencyError = New("internal error", WithStatus(http.StatusFailedDependency))
	ErrSeedNotFound            = New("seed not found", WithStatus(http.StatusNotFound))
)
