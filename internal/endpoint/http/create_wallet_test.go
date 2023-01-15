package http

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/go-chi/chi"
	"go.uber.org/zap"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ariden83/blockchain/config"
	"github.com/ariden83/blockchain/internal/event"
	transactionfactory "github.com/ariden83/blockchain/internal/transaction/factory"
	"github.com/ariden83/blockchain/internal/wallet"
	"github.com/ariden83/blockchain/pkg/api"
)

func Test_create_wallet(t *testing.T) {
	trans, err := transactionfactory.New(transactionfactory.Config{
		Implementation: transactionfactory.ImplementationStub,
	})
	assert.NoError(t, err)

	wallets, err := wallet.New(config.Wallet{}, zap.NewNop())
	assert.NoError(t, err)

	endpoint := New(WithTransactions(trans), WithEvents(event.New()), WithWallets(wallets))

	t.Run("nominal", func(t *testing.T) {
		providedRequest := api.CreateWalletInput{
			Password: []byte("password"),
		}
		validBodyRequest, err := json.Marshal(providedRequest)

		rw := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodPost, "/wallet", bytes.NewBuffer(validBodyRequest))

		rctx := chi.NewRouteContext()
		r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))

		endpoint.handleCreateWallet(rw, r)
		require.Equal(t, http.StatusCreated, rw.Code)
		require.NoError(t, err)

		var jsonResult api.CreateWalletOutput
		err = json.Unmarshal(rw.Body.Bytes(), &jsonResult)
		assert.NoError(t, err)
	})
}
