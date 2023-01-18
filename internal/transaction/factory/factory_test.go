package transactionfactory

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_New(t *testing.T) {
	t.Run("stub implementation", func(t *testing.T) {
		adapter, err := New(Config{Implementation: ImplementationStub})
		require.NotNil(t, adapter)
		require.NoError(t, err)
	})

	t.Run("unknown implementation", func(t *testing.T) {
		adapter, err := New(Config{Implementation: ImplementationUnknown})
		assert.Nil(t, adapter)
		assert.Error(t, err)
	})

	t.Run("invalid implementation", func(t *testing.T) {
		adapter, err := New(Config{Implementation: "invalid"})
		assert.Nil(t, adapter)
		assert.Error(t, err)
	})
}
