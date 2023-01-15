package metrics

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/ariden83/blockchain/config"
	"github.com/ariden83/blockchain/internal/logger"
)

func Test_Metrics(t *testing.T) {
	cfg := &config.Config{}
	logs := logger.InitLight(cfg.Log)
	defer logs.Sync()

	mtc := New(config.Metrics{}, logs)
	assert.NotNil(t, mtc)
}
