package iterator

import (
	"testing"

	"github.com/stretchr/testify/assert"

	persistencefactory "github.com/ariden83/blockchain/internal/persistence/factory"
)

func Test_New(t *testing.T) {
	p, err := persistencefactory.New(persistencefactory.Config{
		Implementation: persistencefactory.ImplementationStub,
	})
	assert.NoError(t, err)

	iter := New(p)
	assert.NotNil(t, iter)

	t.Run("Int", func(t *testing.T) {
		block, err := iter.Next()
		assert.Error(t, err)
		assert.Nil(t, block)
	})
}
