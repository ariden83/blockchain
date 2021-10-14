package event

type EventType int

type Event struct {
	channel chan EventType
}

const (
	BlockChain EventType = iota
	Wallet
	Pool
)

func (e EventType) String() string {
	return [...]string{"Blockchain", "Wallets", "Pool"}[e]
}

func New() *Event {
	c := make(chan EventType)
	return &Event{
		channel: c,
	}
}

func (e *Event) Push(evt EventType) {
	e.channel <- evt
}

func (e *Event) Get() chan EventType {
	return e.channel
}
