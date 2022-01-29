package errors

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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
	s          string
	status     int
	grpcStatus codes.Code
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
		e.grpcStatus = listStatus[status]
	}
}

func StatusCode(err error) int {
	if errWithStatus, ok := err.(*errorString); ok {
		return errWithStatus.Status()
	}
	return http.StatusUpgradeRequired
}

func GRPC(e error) error {
	if err, ok := e.(*errorString); ok {
		return status.Errorf(
			err.grpcStatus,
			err.Error(),
		)
	}
	return status.Errorf(
		codes.InvalidArgument,
		e.Error(),
	)
}

var listStatus = map[int]codes.Code{
	http.StatusBadRequest:         codes.InvalidArgument,
	http.StatusPreconditionFailed: codes.FailedPrecondition,
	http.StatusForbidden:          codes.Internal,
	http.StatusFailedDependency:   codes.Internal,
	http.StatusNotFound:           codes.NotFound,
	http.StatusUnauthorized:       codes.PermissionDenied,
}

var (
	ErrMissingPassword         = New("missing password", WithStatus(http.StatusBadRequest))
	ErrMissingFields           = New("missing fields", WithStatus(http.StatusPreconditionFailed))
	ErrInvalidPassword         = New("password invalid", WithStatus(http.StatusPreconditionFailed))
	ErrInternalError           = New("internal error", WithStatus(http.StatusForbidden))
	ErrInternalDependencyError = New("internal error", WithStatus(http.StatusFailedDependency))
	ErrSeedNotFound            = New("seed not found", WithStatus(http.StatusNotFound))
	ErrInvalidCaptcha          = New("invalid captcha", WithStatus(http.StatusUnauthorized))
)
