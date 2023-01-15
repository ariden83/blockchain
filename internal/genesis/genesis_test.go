package genesis

import (
	"go.uber.org/zap"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/ariden83/blockchain/config"
	"github.com/ariden83/blockchain/internal/event"
	"github.com/ariden83/blockchain/internal/p2p"
	persistencefactory "github.com/ariden83/blockchain/internal/persistence/factory"
	transactionfactory "github.com/ariden83/blockchain/internal/transaction/factory"
	"github.com/ariden83/blockchain/internal/wallet"
)

func Test_Genesis(t *testing.T) {
	cfg := &config.Config{}
	persistence, err := persistencefactory.New(persistencefactory.Config{
		Implementation: persistencefactory.ImplementationStub,
	})
	assert.NoError(t, err)

	trans, err := transactionfactory.New(transactionfactory.Config{Implementation: transactionfactory.ImplementationStub})
	assert.NoError(t, err)

	endPoint := &p2p.EndPoint{}
	evt := &event.Event{}

	wallets, err := wallet.New(config.Wallet{}, zap.NewNop())
	assert.NoError(t, err)

	genesis := New(cfg, persistence, trans, endPoint, evt, wallets)

	t.Run("New", func(t *testing.T) {
		assert.NotNil(t, genesis)
	})

	t.Run("Genesis without default target", func(t *testing.T) {
		assert.False(t, genesis.Genesis())
	})

	t.Run("Genesis with target set", func(t *testing.T) {
		endPoint.SetTarget("/test")
		assert.True(t, genesis.Genesis())
	})
}
