package genesis

import (
	"go.uber.org/zap"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/ariden83/blockchain/config"
	"github.com/ariden83/blockchain/internal/event"
	p2pfactory "github.com/ariden83/blockchain/internal/p2p/factory"
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

	p2p, err := p2pfactory.New(p2pfactory.Config{Implementation: p2pfactory.ImplementationStub}, nil, nil, nil, nil)
	assert.NoError(t, err)

	evt := &event.Event{}

	wallets, err := wallet.New(config.Wallet{}, zap.NewNop())
	assert.NoError(t, err)

	genesis := New(cfg, persistence, trans, p2p, evt, wallets)

	t.Run("New", func(t *testing.T) {
		assert.NotNil(t, genesis)
	})

	t.Run("Genesis without default target", func(t *testing.T) {
		assert.False(t, genesis.Genesis())
	})

	t.Run("Genesis with target set", func(t *testing.T) {
		p2p.SetTarget("/test")
		assert.True(t, genesis.Genesis())
	})
}

func Test_Genesis_Genesis(t *testing.T) {
	p2p, err := p2pfactory.New(p2pfactory.Config{Implementation: p2pfactory.ImplementationStub}, nil, nil, nil, nil)
	assert.NoError(t, err)

	// Test with p2p target
	gen := New(nil, nil, nil, p2p, nil, nil)
	assert.True(t, gen.Genesis())

	// Test without p2p target
	gen.p2p = p2p
	assert.False(t, gen.Genesis())
}

/*
func Test_Genesis_Load(t *testing.T) {
	stop := make(chan error)

	// Test case where persistence.DBExists returns false
	gen := New(nil, &MockPersistenceAdapter{}, nil, nil, nil, nil)
	assert.NotNil(t, gen)
	go func() {
		gen.Load(stop)
	}()

	err := <-stop
	assert.NotNil(t, err)
	assert.Equal(t, "fail to open DB files", err.Error())

	// Test case where GetLastHash returns an error
	gen = New(nil, &MockPersistenceAdapter{}, nil, nil, nil, nil)
	assert.NotNil(t, gen)
	go func() {
		gen.Load(stop)
	}()

	err = <-stop
	assert.NotNil(t, err)
	assert.Equal(t, "fail to get last hash", err.Error())

	// Test case where GetCurrentHashSerialize returns an error
	gen = New(nil, &MockPersistenceAdapter{}, nil, nil, nil, nil)
	assert.NotNil(t, gen)
	go func() {
		gen.Load(stop)
	}()

	err = <-stop
	assert.NotNil(t, err)
	assert.Equal(t, "fail to get current hash", err.Error())

	// Test case where serialization fails
	gen = New(nil, &MockPersistenceAdapter{}, nil, nil, nil, nil)
	assert.NotNil(t, gen)
	gen.persistence = &MockPersistenceAdapter{}
	go func() {
		gen.Load(stop)
	}()

	err = <-stop
	assert.NotNil(t, err)
	assert.Equal(t, "fail to serialize genesis: fail to deserialize hash serializes: unexpected end of JSON input", err.Error())

	// Test case where updating persistence fails
	gen = New(nil, &MockPersistenceAdapter{}, nil, nil, nil, nil)
	assert.NotNil(t, gen)
	gen.persistence = &MockPersistenceAdapter{}
	go func() {
		gen.Load(stop)
	}()

	err = <-stop
	assert.NotNil(t, err)
	assert.Equal(t, "fail to update persistence", err.Error())
}

func Test_Genesis_createGenesis(t *testing.T) {
	stop := make(chan error)

	// Test case where wallet creation fails
	gen := New(nil, nil, nil, nil, nil, &MockWallets{})
	assert.NotNil(t, gen)
	gen.wallets = &MockWallets{}
	go func() {
		gen.createGenesis(stop)
	}()

	err := <-stop
	assert.NotNil(t, err)
	assert.Equal(t, "fail to serialize genesis: seed is not set", err.Error())

	// Test case where serialization fails
	gen = New(nil, nil, nil, nil, nil, &MockWallets{})
	assert.NotNil(t, gen)
	go func() {
		gen.createGenesis(stop)
	}()

	err = <-stop
	assert.NotNil(t, err)
	assert.Equal(t, "fail to serialize genesis: fail to serialize genesis: json: unsupported type: chan struct {}", err.Error())

	// Test case where updating persistence fails
	gen = New(nil, &MockPersistenceAdapter{}, nil, nil, nil, &MockWallets{})
	assert.NotNil(t, gen)
	gen.persistence = &MockPersistenceAdapter{}
	gen.wallets = &MockWallets{}
	go func() {
		gen.createGenesis(stop)
	}()

	err = <-stop
	assert.NotNil(t, err)
	assert.Equal(t, "fail to update persistence", err.Error())
}
*/
