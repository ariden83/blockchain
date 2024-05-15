// Package P2P represent a peer to peer linked on go-libp2p.
package stub

func New() *Stub {
	return &Stub{}
}

// Stub
type Stub struct{}

func (m *Stub) Enabled() bool   { return true }
func (m *Stub) HasTarget() bool { return true }

func (m *Stub) Listen(stop chan error)     {}
func (m *Stub) Target() string             { return "" }
func (m *Stub) PushMsgForFiles(chan error) {}

func (m *Stub) SetTarget(target string) {}

func (m *Stub) Shutdown() {}
