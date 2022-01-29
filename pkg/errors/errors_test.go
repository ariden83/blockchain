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
	}{
		{"bad request", New("bad request", WithStatus(http.StatusBadRequest)), http.StatusBadRequest},
		{"missing password", ErrMissingPassword, http.StatusBadRequest},
		{"missing fields", ErrMissingFields, http.StatusPreconditionFailed},
		{"password invalid", ErrInvalidPassword, http.StatusPreconditionFailed},
		{"internal error", ErrInternalError, http.StatusForbidden},
		{"internal error", ErrInternalDependencyError, http.StatusFailedDependency},
		{"seed not found", ErrSeedNotFound, http.StatusNotFound},
		{"invalid captcha", ErrInvalidCaptcha, http.StatusUnauthorized},
		{"error", errors.New("error"), http.StatusUpgradeRequired},
	} {
		t.Run(test.messageTest, func(t *testing.T) {
			assert.NotNil(t, test.error)
			assert.Equal(t, test.messageTest, test.error.Error())
			assert.Equal(t, test.statusWaiting, StatusCode(test.error))
		})
	}
}
