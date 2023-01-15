package metrics

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/ariden83/blockchain/config"
)

func Test_Metrics(t *testing.T) {
	mtc := New(config.Metrics{}, nil)
	assert.NotNil(t, mtc)
}
