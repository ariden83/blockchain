package errors

import (
	"net/http"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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

func (e *errorString) statusFromGRPC() codes.Code {
	for statusCode, num := range listStatus {
		if statusCode == e.status {
			return num
		}
	}
	return 0
}

func StatusCode(err error) int {
	if errWithStatus, ok := err.(*errorString); ok {
		return errWithStatus.Status()
	}
	return http.StatusUpgradeRequired
}

func GRPC(e error) error {
	if err, ok := e.(*errorString); ok {
		if listStatus[err.Status()] != 0 {
			return status.Errorf(
				listStatus[err.Status()],
				err.Error(),
			)
		}
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
	http.StatusResetContent:       codes.InvalidArgument,
}

var (
	ErrMissingPassword         = New("missing password", WithStatus(http.StatusBadRequest))
	ErrMissingFields           = New("missing fields", WithStatus(http.StatusPreconditionFailed))
	ErrEmptyField              = New("empty field", WithStatus(http.StatusPreconditionFailed))
	ErrInvalidPassword         = New("password invalid", WithStatus(http.StatusPreconditionFailed))
	ErrInternalError           = New("internal error", WithStatus(http.StatusForbidden))
	ErrInternalDependencyError = New("internal error", WithStatus(http.StatusFailedDependency))
	ErrSeedNotFound            = New("seed not found", WithStatus(http.StatusNotFound))
	ErrRecreatePassword        = New("recreate seed password", WithStatus(http.StatusResetContent))
	ErrInvalidCaptcha          = New("invalid captcha", WithStatus(http.StatusUnauthorized))
	ErrAlreadyConnected        = New("alreadyConnected", WithStatus(http.StatusFound))
	ErrNotEnoughFunds          = New("not enough funds", WithStatus(http.StatusUnauthorized))
	ErrorSeedPasswordInvalid   = New("invalid password", WithStatus(http.StatusResetContent))
	ErrCreatedBlockIsInvalid   = New("new block created is invalid", WithStatus(http.StatusRequestedRangeNotSatisfiable))
	ErrNotImplemented          = New("not implemented", WithStatus(http.StatusNotImplemented))
)
