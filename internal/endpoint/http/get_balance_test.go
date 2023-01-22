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

func Test_get_Balance(t *testing.T) {
	trans, err := transactionfactory.New(transactionfactory.Config{
		Implementation: transactionfactory.ImplementationStub,
	})
	assert.NoError(t, err)

	wallets, err := wallet.New(config.Wallet{}, zap.NewNop())
	assert.NoError(t, err)

	endpoint := New(WithTransactions(trans), WithEvents(event.New()), WithWallets(wallets))

	t.Run("nominal", func(t *testing.T) {
		providedRequest := api.GetBalanceInput{
			PrivKey: []byte("private-key"),
		}
		validBodyRequest, err := json.Marshal(providedRequest)

		rw := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodPost, "/balance", bytes.NewBuffer(validBodyRequest))
		r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, chi.NewRouteContext()))

		endpoint.handleGetBalance(rw, r)
		require.Equal(t, http.StatusOK, rw.Code)
		require.NoError(t, err)

		var jsonResult api.GetBalanceOutput
		err = json.Unmarshal(rw.Body.Bytes(), &jsonResult)
		assert.NoError(t, err)
	})
}
