package p2p

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"

	"github.com/ariden83/blockchain/internal/event"
	persistencefactory "github.com/ariden83/blockchain/internal/persistence/factory"
	"github.com/ariden83/blockchain/internal/wallet"
)

func Test_New_Cache(t *testing.T) {
	cfg := Config{}
	persistence, err := persistencefactory.New(persistencefactory.Config{
		Implementation: persistencefactory.ImplementationStub,
	})
	assert.NoError(t, err)

	wallets, err := wallet.New(wallet.Config{}, zap.NewNop())
	assert.NoError(t, err)

	evt := event.New()

	t.Run("new cache adapter", func(t *testing.T) {
		cacheAdapter := New(cfg, persistence, wallets, zap.NewNop(), evt)
		assert.NotNil(t, cacheAdapter)
	})

	t.Run("cache with WithSize option", func(t *testing.T) {
		cacheAdapter := New(cfg, persistence, wallets, zap.NewNop(), evt, WithXCache(XCache{}))
		assert.NoError(t, err)
		assert.NotNil(t, cacheAdapter)
		assert.NotNil(t, cacheAdapter.xCache, int32(5))
	})
}
