package event

import (
	"github.com/satori/go.uuid"
	"go.uber.org/zap"

	"github.com/ariden83/blockchain/config"
	"github.com/ariden83/blockchain/internal/event/trace"
	"github.com/ariden83/blockchain/internal/p2p/validator"
)

// Message represent a message push on a channel.
type Message struct {
	Type  EventType
	ID    string
	Value []byte
}

// EventType represent a message event type.
type EventType int

// Event represent a new adapter Event.
type Event struct {
	channel      chan Message
	channelBlock chan validator.Validator
	listChannel  []chan Message
	trace        *trace.Trace
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

func New(options ...func(*Event)) *Event {
	c := make(chan Message)
	e := &Event{
		channel: c,
	}

	for _, o := range options {
		o(e)
	}

	go func() {
		e.setConcurrence()
	}()
	return e
}

func WithTraces(cfg config.Trace, logs *zap.Logger) func(*Event) {
	return func(e *Event) {
		if cfg.Enabled {
			e.trace = trace.New(cfg, logs.With(zap.String("service", "traces")))
		}
	}
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

func (e *Event) PushBlock(block validator.Validator) {
	e.channelBlock <- block
}

func (e *Event) PushTrace(blockID string, state trace.State) {
	if e.trace != nil {
		e.trace.Push(blockID, state)
	}
}

func (e *Event) NewBlockReader() chan validator.Validator {
	e.channelBlock = make(chan validator.Validator)
	return e.channelBlock
}

func (e *Event) NewTraceReader() *trace.Channel {
	return e.trace.NewReader()
}

// CloseTraceReader must close trace reader.
func (e *Event) CloseTraceReader(c trace.Channel) {
	e.trace.CloseReader(c)
}

// CloseReaders must close all readers.
func (e *Event) CloseReaders() {
	for _, c := range e.listChannel {
		close(c)
	}
}
