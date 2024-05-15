package p2p

import (
	"go.uber.org/zap"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/ariden83/blockchain/config"
	"github.com/ariden83/blockchain/internal/event"
	persistencefactory "github.com/ariden83/blockchain/internal/persistence/factory"
	"github.com/ariden83/blockchain/internal/wallet"
)

func Test_New_Cache(t *testing.T) {
	cfg := config.P2P{}
	persistence, err := persistencefactory.New(persistencefactory.Config{
		Implementation: persistencefactory.ImplementationStub,
	})
	assert.NoError(t, err)

	wallets, err := wallet.New(config.Wallet{}, zap.NewNop())
	assert.NoError(t, err)

	evt := event.New()

	t.Run("new cache adapter", func(t *testing.T) {
		cacheAdapter := New(cfg, persistence, wallets, zap.NewNop(), evt)
		assert.NotNil(t, cacheAdapter)
	})

	t.Run("cache with WithSize option", func(t *testing.T) {
		cacheAdapter := New(cfg, persistence, wallets, zap.NewNop(), evt, WithXCache(config.XCache{}))
		assert.NoError(t, err)
		assert.NotNil(t, cacheAdapter)
		assert.NotNil(t, cacheAdapter.xCache, int32(5))
	})
}
