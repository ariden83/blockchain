// Package event implements an event streaming platform to communicate and allow sharing of all events to different servers.
package event

import (
	uuid "github.com/satori/go.uuid"
	"go.uber.org/zap"

	"github.com/ariden83/blockchain/internal/event/trace"
	"github.com/ariden83/blockchain/internal/p2p/validator"
)

// Message represents a message pushed on a channel.
type Message struct {
	Type  EventType
	ID    string
	Value []byte
}

// EventType represents a message event type.
type EventType int

// Event represents a new adapter Event.
type Event struct {
	channel      chan Message
	channelBlock chan validator.Validator
	listChannel  []chan Message
	trace        *trace.Trace
}

// Constants for different event types.
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

// New creates a new Event instance with optional configurations.
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

// WithTraces configures Event instance with tracing options.
func WithTraces(cfg trace.Config, logs *zap.Logger) func(*Event) {
	return func(e *Event) {
		if cfg.Enabled {
			e.trace = trace.New(cfg, logs.With(zap.String("service", "traces")))
		}
	}
}

// Push pushes a message onto the channel.
func (e *Event) Push(m Message) {
	if m.ID == "" {
		m.ID = uuid.NewV4().String()
	}
	e.channel <- m
}

// setConcurrence sets concurrency for message processing.
func (e *Event) setConcurrence() {
	for data := range e.channel {
		for _, c := range e.listChannel {
			c <- data
		}
	}
}

// NewReader creates a new message reader channel.
func (e *Event) NewReader() chan Message {
	newChan := make(chan Message)
	e.listChannel = append(e.listChannel, newChan)
	return newChan
}

// PushBlock pushes a block onto the channel.
func (e *Event) PushBlock(block validator.Validator) {
	e.channelBlock <- block
}

// PushTrace pushes a trace onto the channel.
func (e *Event) PushTrace(blockID string, state trace.State) {
	if e.trace != nil {
		e.trace.Push(blockID, state)
	}
}

// NewBlockReader creates a new block reader channel.
func (e *Event) NewBlockReader() chan validator.Validator {
	e.channelBlock = make(chan validator.Validator)
	return e.channelBlock
}

// NewTraceReader creates a new trace reader channel.
func (e *Event) NewTraceReader() *trace.Channel {
	return e.trace.NewReader()
}

// CloseTraceReader closes the trace reader channel.
func (e *Event) CloseTraceReader(c trace.Channel) {
	e.trace.CloseReader(c)
}

// CloseReaders closes all message reader channels.
func (e *Event) CloseReaders() {
	for _, c := range e.listChannel {
		close(c)
	}
}
