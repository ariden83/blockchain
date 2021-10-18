package event

import (
	"github.com/ariden83/blockchain/internal/blockchain"
	"github.com/satori/go.uuid"
)

type Message struct {
	Type EventType
	ID   string
}

type EventType int

type Event struct {
	channel      chan Message
	channelBlock chan blockchain.Block
	listChannel  []chan Message
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
	c := make(chan Message)
	e := &Event{
		channel: c,
	}
	go func() {
		e.setConcurrence()
	}()
	return e
}

func (e *Event) Push(evt EventType, ID string) {
	m := Message{
		Type: evt,
		ID:   ID,
	}
	if ID == "" {
		m.ID = uuid.NewV4().String()
	}
	e.channel <- m
}

func (e *Event) setConcurrence() {
	for data := range e.channel {
		for _, c := range e.listChannel {
			c <- data
		}
	}
}

func (e *Event) NewReader() chan Message {
	newChan := make(chan Message)
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
