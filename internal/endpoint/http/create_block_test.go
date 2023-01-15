package http

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/go-chi/chi"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	transactionfactory "github.com/ariden83/blockchain/internal/transaction/factory"
	"github.com/ariden83/blockchain/pkg/api"
)

func Test_create_block(t *testing.T) {
	trans, err := transactionfactory.New(transactionfactory.Config{
		Implementation: transactionfactory.ImplementationStub,
	})
	assert.NoError(t, err)

	endpoint := New(WithTransactions(trans))

	t.Run("nominal", func(t *testing.T) {
		providedRequest := api.CreateBlockInput{
			PrivKey: []byte("priv-key"),
		}
		validBodyRequest, err := json.Marshal(providedRequest)

		rw := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodPost, "/block", bytes.NewBuffer(validBodyRequest))

		rctx := chi.NewRouteContext()
		r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))

		endpoint.handleCreateBlock(rw, r)
		require.Equal(t, http.StatusProcessing, rw.Code)
		require.NoError(t, err)

		var jsonResult api.CreateBlockOutput
		err = json.Unmarshal(rw.Body.Bytes(), &jsonResult)
		assert.NoError(t, err)
	})
}
