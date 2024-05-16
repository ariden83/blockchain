package http

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ariden83/blockchain/internal/p2p/address"
)

func Test_get_Servers_address(t *testing.T) {
	endpoint := New()
	address.IAM.SetAddress("http://127.0.0.1")

	t.Run("nominal", func(t *testing.T) {
		expectedResp := `[
  "http://127.0.0.1"
]`

		rw := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/address", nil)
		r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, chi.NewRouteContext()))

		endpoint.handleGetServersAddress(rw, r)
		require.Equal(t, http.StatusOK, rw.Code)

		bodyBytes, err := io.ReadAll(rw.Body)
		assert.Equal(t, expectedResp, string(bodyBytes))
		assert.NoError(t, err)
	})
}
