package mondialrelayfactory

import (
	"fmt"

	persistenceadapter "github.com/ariden83/blockchain/internal/persistence"
	"github.com/ariden83/blockchain/internal/persistence/impl/badger"
	"github.com/ariden83/blockchain/internal/persistence/impl/stub"
)

const (
	// ImplementationUnknown is an unknown and invalid adapter implementation.
	ImplementationUnknown = ""
	// ImplementationStub is an implementation that does nothing and return no error.
	ImplementationStub = "stub"
	// ImplementationBadger uses a badger impl.
	ImplementationBadger = "badger"
)

// Config struct which describe how build an ad adapter instance.
type Config struct {
	Implementation string `mapstructure:"impl"`
	Badger         badger.Config
}

// New creates a new Adapter instance based on a Config.
func New(config Config) (persistenceadapter.Adapter, error) {
	var (
		adapter persistenceadapter.Adapter
		err     error
	)

	switch config.Implementation {
	case ImplementationBadger:
		adapter, err = badger.New(config.Badger)
	case ImplementationStub:
		adapter = stub.New()
	case ImplementationUnknown:
		fallthrough
	default:
		return adapter, fmt.Errorf("unknown implementation %s", config.Implementation)
	}

	return adapter, err
}
