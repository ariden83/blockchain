package errors

import (
	"errors"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Error(t *testing.T) {
	for _, test := range []struct {
		messageTest   string
		error         error
		statusWaiting int
		messageGRPC error
	}{
		{
			"bad request",
			New("bad request", WithStatus(http.StatusBadRequest)),
			http.StatusBadRequest,
			errors.New("rpc error: code = InvalidArgument desc = bad request"),
		},
		{
			"missing password",
			ErrMissingPassword,
			http.StatusBadRequest,
			errors.New("rpc error: code = InvalidArgument desc = missing password"),
		},
		{
			"missing fields",
			ErrMissingFields,
			http.StatusPreconditionFailed,
			errors.New("rpc error: code = FailedPrecondition desc = missing fields"),
		},
		{
			"password invalid",
			ErrInvalidPassword,
			http.StatusPreconditionFailed,
			errors.New("rpc error: code = FailedPrecondition desc = password invalid"),
		},
		{
			"internal error",
			ErrInternalError,
			http.StatusForbidden,
			errors.New("rpc error: code = Internal desc = internal error"),
		},
		{
			"internal error",
			ErrInternalDependencyError,
			http.StatusFailedDependency,
			errors.New("rpc error: code = Internal desc = internal error"),
		},
		{
			"seed not found",
			ErrSeedNotFound,
			http.StatusNotFound,
			errors.New("rpc error: code = NotFound desc = seed not found"),
		},
		{
			"invalid captcha",
			ErrInvalidCaptcha,
			http.StatusUnauthorized,
			errors.New("rpc error: code = PermissionDenied desc = invalid captcha"),
		},
		{
			"error",
			errors.New("error"),
			http.StatusUpgradeRequired,
			errors.New("rpc error: code = InvalidArgument desc = error"),
		},
	} {
		t.Run(test.messageTest, func(t *testing.T) {
			assert.NotNil(t, test.error)
			assert.Equal(t, test.messageTest, test.error.Error())
			assert.Equal(t, test.statusWaiting, StatusCode(test.error))
			assert.Equal(t, test.messageGRPC.Error(), GRPC(test.error).Error())
		})
	}
}
