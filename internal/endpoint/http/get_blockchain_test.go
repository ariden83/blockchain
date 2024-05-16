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

	persistencefactory "github.com/ariden83/blockchain/internal/persistence/factory"
)

func Test_get_Blockchain(t *testing.T) {
	per, err := persistencefactory.New(persistencefactory.Config{
		Implementation: persistencefactory.ImplementationStub,
	})
	assert.NoError(t, err)

	endpoint := New(WithPersistence(per))

	t.Run("nominal", func(t *testing.T) {
		rw := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/blockchain", nil)
		r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, chi.NewRouteContext()))

		endpoint.handleGetBlockChain(rw, r)
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, rw.Code)

		if rw.Body != nil {
			bodyBytes, err := io.ReadAll(rw.Body)
			require.NoError(t, err)
			require.Equal(t, []byte(""), bodyBytes)
		}
	})
}
