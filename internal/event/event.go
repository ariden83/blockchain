package event

import "github.com/ariden83/blockchain/internal/blockchain"

type EventType int

type Event struct {
	channel      chan EventType
	channelBlock chan blockchain.Block
	listChannel  []chan EventType
}

const (
	BlockChain EventType = iota
	NewBlock
	Wallet
	Pool
	Files
)

func (e EventType) String() string {
	return [...]string{"Blockchain", "Block", "Wallets", "Pool", "files"}[e]
}

func New() *Event {
	c := make(chan EventType)
	e := &Event{
		channel: c,
	}
	go func() {
		e.setConcurrence()
	}()
	return e
}

func (e *Event) Push(evt EventType) {
	e.channel <- evt
}

func (e *Event) setConcurrence() {
	for data := range e.channel {
		for _, c := range e.listChannel {
			c <- data
		}
	}
}

func (e *Event) NewReader() chan EventType {
	newChan := make(chan EventType)
	e.listChannel = append(e.listChannel, newChan)
	return newChan
}

func (e *Event) PushBlock(block blockchain.Block) {
	e.channelBlock <- block
}

func (e *Event) NewBlockReader() chan blockchain.Block {
	newChan := make(chan blockchain.Block)
	return newChan
}
