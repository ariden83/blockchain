// Package trace implements a tracing mechanism to track the state of blockchain events.
package trace

import (
	"fmt"

	"go.uber.org/zap"

	"github.com/ariden83/blockchain/internal/utils"
)

type Config struct {
	Enabled bool `config:"memory_enabled"`
}

// State represents the state of a blockchain event.
type State int

// String returns the string representation of the State.
func (e State) String() string {
	return [...]string{"Mining", "Creating", "Validating", "Done", "Failed"}[e]
}

// Constants representing different states of a blockchain event.
const (
	Mining State = iota
	Creating
	Validating
	Done
	Failed
)

// Message represents a message to be pushed on a trace channel.
type Message struct {
	ID    string
	State State
}

// Channel represents a trace channel.
type Channel struct {
	channel chan Message
	id      string
	toClose bool
}

// GetChan returns the underlying channel of the trace channel.
func (c *Channel) GetChan() chan Message {
	return c.channel
}

// GetID returns the ID of the trace channel.
func (c *Channel) GetID() string {
	return c.id
}

// Close closes the trace channel.
func (c *Channel) Close() {
	c.toClose = true
	if c.channel != nil {
		close(c.channel)
	}
}

// Trace represents the trace mechanism to track blockchain events.
type Trace struct {
	channel     chan Message
	listChannel map[string]Channel
	log         *zap.Logger
}

// New creates a new Trace instance with the given configuration and logger.
func New(cfg Config, log *zap.Logger) *Trace {
	if !cfg.Enabled {
		return nil
	}
	t := &Trace{
		channel:     make(chan Message),
		listChannel: map[string]Channel{},
		log:         log,
	}

	go func() {
		t.setConcurrence()
	}()

	return t
}

// setConcurrence sets concurrency for message processing.
func (t *Trace) setConcurrence() {
	for data := range t.channel {
		for _, c := range t.listChannel {
			if c.channel != nil {
				c.channel <- data
			}
		}
	}
}

// NewReader creates a new trace reader channel.
func (t *Trace) NewReader() *Channel {
	c := Channel{
		id:      utils.RandomString(uint8(5)),
		channel: make(chan Message),
	}

	t.listChannel[c.id] = c
	return &c
}

// CloseReader closes the trace reader channel.
func (t *Trace) CloseReader(ch Channel) {
	_, ok := t.listChannel[ch.GetID()]
	if ok {
		delete(t.listChannel, ch.GetID())
	}
}

// Push pushes a message onto the trace channel.
func (t *Trace) Push(blockID string, state State) {
	if blockID == "" {
		return
	}
	t.log.Info(fmt.Sprintf("send message in trace channel %s %+v", blockID, state))
	t.channel <- Message{
		ID:    blockID,
		State: state,
	}
}
