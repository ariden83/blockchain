package event

// Stub event instance
type Stub struct{}

func (m *Stub) NewReader() chan Message { return make(chan Message) }
func (m *Stub) CloseReaders()           {}
func (m *Stub) Push(Message)            {}
