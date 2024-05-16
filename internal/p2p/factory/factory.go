package factory

import (
	"fmt"

	"go.uber.org/zap"

	"github.com/ariden83/blockchain/internal/event"
	"github.com/ariden83/blockchain/internal/p2p"
	p2padapter "github.com/ariden83/blockchain/internal/p2p/impl/p2p"
	"github.com/ariden83/blockchain/internal/p2p/impl/stub"
	persistenceadapter "github.com/ariden83/blockchain/internal/persistence"
	"github.com/ariden83/blockchain/internal/wallet"
)

const (
	// ImplementationUnknown is an unknown and invalid adapter implementation.
	ImplementationUnknown = ""
	// ImplementationStub is an implementation that does nothing and return no error.
	ImplementationStub = "stub"
	// ImplementationP2P uses a P2P impl.
	ImplementationP2P = "p2p"
)

// Config struct which describe how build an ad adapter instance.
type Config struct {
	Implementation string `mapstructure:"impl"`

	Config p2padapter.Config
}

// New creates a new Adapter instance based on a Config.
func New(cfg Config,
	per persistenceadapter.Adapter,
	wallets wallet.IWallets,
	logs *zap.Logger,
	evt *event.Event,
	opts ...p2padapter.Option) (p2p.Adapter, error) {

	var (
		adapter p2p.Adapter
		err     error
	)

	switch cfg.Implementation {
	case ImplementationP2P:
		adapter = p2padapter.New(cfg.Config, per, wallets, logs, evt, opts...)
	case ImplementationStub:
		adapter = stub.New()
	case ImplementationUnknown:
		fallthrough
	default:
		return adapter, fmt.Errorf("unknown implementation %s", cfg.Implementation)
	}

	return adapter, err
}
