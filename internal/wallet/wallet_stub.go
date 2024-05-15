package wallet

// Stub wallet instance
type Stub struct{}

func (m *Stub) Create([]byte) (*Seed, error) { return nil, nil }
