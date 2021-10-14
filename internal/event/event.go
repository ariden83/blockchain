package event

type EventType int

type Event struct {
	Chan chan EventType
}

const (
	Blockchain EventType = iota
	Wallets
	Pool
)

func (e EventType) String() string {
	return [...]string{"Blockchain", "Wallets", "Pool"}[e]
}

func New() *Event {
	c := make(chan EventType)
	return &Event{
		Chan: c,
	}
}
