package transactionfactory

import (
	"fmt"

	transactionadapter "github.com/ariden83/blockchain/internal/transaction"
	"github.com/ariden83/blockchain/internal/transaction/impl/stub"
	"github.com/ariden83/blockchain/internal/transaction/impl/transaction"
)

const (
	// ImplementationUnknown is an unknown and invalid adapter implementation.
	ImplementationUnknown = ""
	// ImplementationStub is an implementation that does nothing and return no error.
	ImplementationStub = "stub"
	// ImplementationTransaction uses a transaction impl.
	ImplementationTransaction = "transaction"
)

// Config struct which describe how build an ad adapter instance.
type Config struct {
	Implementation string `mapstructure:"impl"`
}

// New creates a new Adapter instance based on a Config.
func New(config Config, options ...func(transactions *transaction.Transactions)) (transactionadapter.Adapter, error) {
	var (
		adapter transactionadapter.Adapter
		err     error
	)

	switch config.Implementation {
	case ImplementationTransaction:
		adapter = transaction.New(options...)
	case ImplementationStub:
		adapter = stub.New()
	case ImplementationUnknown:
		fallthrough
	default:
		return adapter, fmt.Errorf("unknown implementation %s", config.Implementation)
	}

	return adapter, err
}
