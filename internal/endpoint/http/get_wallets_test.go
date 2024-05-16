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
	"go.uber.org/zap"

	"github.com/ariden83/blockchain/internal/event"
	transactionfactory "github.com/ariden83/blockchain/internal/transaction/factory"
	"github.com/ariden83/blockchain/internal/wallet"
)

func Test_get_Wallets(t *testing.T) {
	trans, err := transactionfactory.New(transactionfactory.Config{
		Implementation: transactionfactory.ImplementationStub,
	})
	assert.NoError(t, err)

	wallets, err := wallet.New(wallet.Config{}, zap.NewNop())
	assert.NoError(t, err)

	_, err = wallets.Create([]byte("test"))
	assert.NoError(t, err)

	seeds, err := wallets.GetSeeds()
	assert.NoError(t, err)

	wallets.UpdateSeeds(seeds)

	endpoint := New(WithTransactions(trans), WithEvents(event.New()), WithWallets(wallets))

	t.Run("nominal", func(t *testing.T) {
		rw := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/wallets", nil)
		r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, chi.NewRouteContext()))

		endpoint.handleGetWallets(rw, r)
		require.Equal(t, http.StatusOK, rw.Code)

		bodyBytes, err := io.ReadAll(rw.Body)
		assert.Equal(t, "null", string(bodyBytes))
		assert.NoError(t, err)
	})
}
