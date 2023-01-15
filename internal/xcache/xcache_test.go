package xcache

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_New_Cache(t *testing.T) {
	t.Run("new cache adapter", func(t *testing.T) {
		cacheAdapter, err := New()
		assert.NoError(t, err)
		assert.NotNil(t, cacheAdapter)
	})

	t.Run("cache with WithSize option", func(t *testing.T) {
		cacheAdapter, err := New(WithSize(5))
		assert.NoError(t, err)
		assert.NotNil(t, cacheAdapter)
		assert.Equal(t, cacheAdapter.posSize, int32(5))
	})

	t.Run("cache with WithPruneSize option", func(t *testing.T) {
		cacheAdapter, err := New(WithPruneSize(5))
		assert.NoError(t, err)
		assert.NotNil(t, cacheAdapter)
		assert.Equal(t, cacheAdapter.posPruneSize, int32(5))
	})

	t.Run("cache with  WithTTL  option", func(t *testing.T) {
		cacheAdapter, err := New(WithTTL(5))
		assert.NoError(t, err)
		assert.NotNil(t, cacheAdapter)
		assert.Equal(t, cacheAdapter.posTTL, time.Duration(5))
	})

	t.Run("cache with  WithNegSize option", func(t *testing.T) {
		cacheAdapter, err := New(WithNegSize(5))
		assert.NoError(t, err)
		assert.NotNil(t, cacheAdapter)
		assert.Equal(t, cacheAdapter.negSize, int32(5))
	})

	t.Run("cache with  WithNegPruneSize option", func(t *testing.T) {
		cacheAdapter, err := New(WithNegPruneSize(5))
		assert.NoError(t, err)
		assert.NotNil(t, cacheAdapter)
		assert.Equal(t, cacheAdapter.negPruneSize, int32(5))
	})

	t.Run("cache with  WithNegTTL option", func(t *testing.T) {
		cacheAdapter, err := New(WithNegTTL(5))
		assert.NoError(t, err)
		assert.NotNil(t, cacheAdapter)
		assert.Equal(t, cacheAdapter.negTTL, time.Duration(5))
	})

	t.Run("cache with  WithStaleFetchers option", func(t *testing.T) {
		cacheAdapter, err := New(WithStaleFetchers(5))
		assert.NoError(t, err)
		assert.NotNil(t, cacheAdapter)
		assert.Equal(t, cacheAdapter.staleFetchers, 5)
	})

	t.Run("cache with  WithStaleQueueSize option", func(t *testing.T) {
		cacheAdapter, err := New(WithStaleQueueSize(5))
		assert.NoError(t, err)
		assert.NotNil(t, cacheAdapter)
		assert.Equal(t, cacheAdapter.staleQueueSize, 5)
	})

	t.Run("cache with  WithFetchers option", func(t *testing.T) {
		cacheAdapter, err := New(WithFetchers(5))
		assert.NoError(t, err)
		assert.NotNil(t, cacheAdapter)
		assert.Equal(t, cacheAdapter.maxFetchers, 5)
	})

	t.Run("cache with  WithStale option", func(t *testing.T) {
		t.Run("with true value", func(t *testing.T) {
			cacheAdapter, err := New(WithStale(true))
			assert.NoError(t, err)
			assert.NotNil(t, cacheAdapter)
			assert.True(t, cacheAdapter.canUseStale)
		})

		t.Run("with false value", func(t *testing.T) {
			cacheAdapter, err := New(WithStale(false))
			assert.NoError(t, err)
			assert.NotNil(t, cacheAdapter)
			assert.False(t, cacheAdapter.canUseStale)
		})
	})

	t.Run("cache with  WithStaleValidator option", func(t *testing.T) {
		cacheAdapter, err := New(WithStaleValidator(func(interface{}, time.Duration) bool {
			return true
		}))
		assert.NoError(t, err)
		assert.NotNil(t, cacheAdapter)
	})
}
