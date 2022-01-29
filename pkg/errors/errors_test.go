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
		{"missing password", New("missing password", WithStatus(http.StatusBadRequest)), http.StatusBadRequest},
		{"password invalid", New("password invalid", WithStatus(http.StatusUnauthorized)), http.StatusUnauthorized},
		{"internal error", New("internal error", WithStatus(http.StatusForbidden)), http.StatusForbidden},
		{"internal error", New("internal error", WithStatus(http.StatusFailedDependency)), http.StatusFailedDependency},
		{"seed not found", New("seed not found", WithStatus(http.StatusNotFound)), http.StatusNotFound},
		{"error", errors.New("error"), http.StatusUpgradeRequired},
	} {
		t.Run(test.messageTest, func(t *testing.T) {
			assert.NotNil(t, test.error)
			assert.Equal(t, test.messageTest, test.error.Error())
			assert.Equal(t, test.statusWaiting, StatusCode(test.error))
		})
	}
}
