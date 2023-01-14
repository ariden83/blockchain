package persistenceadapter

// Adapter is an interface which describes all the methods to interact with persistence.
type Adapter interface {
	Close() error
	DBExists() bool
	GetLastHash() ([]byte, error)
	GetCurrentHashSerialize(hash []byte) ([]byte, error)
	LastHash() []byte
	SetLastHash(lastHash []byte)
	Update(lastHash []byte, hashSerialize []byte) error
}
