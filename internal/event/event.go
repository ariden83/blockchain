package event

import (
	"github.com/ariden83/blockchain/internal/blockchain"
	"github.com/satori/go.uuid"
)

type Message struct {
	Type  EventType
	ID    string
	Value []byte
}

type EventType int

type Event struct {
	channel      chan Message
	channelBlock chan blockchain.Block
	listChannel  []chan Message
}

const (
	BlockChain EventType = iota
	BlockChainFull
	NewBlock
	Wallet
	Pool
	Files
	Address
)

func (e EventType) String() string {
	return [...]string{"Blockchain", "BlockChainFull", "Block", "Wallets", "Pool", "files", "Address"}[e]
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

func (e *Event) Push(m Message) {
	if m.ID == "" {
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
	e.channelBlock = make(chan blockchain.Block)
	return e.channelBlock
}
