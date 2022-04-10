package trace

import (
	"github.com/ariden83/blockchain/config"
	"github.com/ariden83/blockchain/internal/utils"
)

type State int

func (e State) String() string {
	return [...]string{"Minage", "Create", "Validate", "Done"}[e]
}

const (
	Minage State = iota
	Create
	Validate
	Done
)

type Message struct {
	ID    string
	State State
}

type Channel struct {
	channel chan Message
	id      string
	toClose bool
}

func (c *Channel) GetChan() chan Message {
	return c.channel
}

func (c *Channel) GetID() string {
	return c.id
}

func (c *Channel) Close() {
	c.toClose = true
	close(c.channel)
}

type Trace struct {
	channel     chan Message
	listChannel map[string]Channel
}

func New(cfg config.Trace) *Trace {
	if !cfg.Enabled {
		return nil
	}
	t := &Trace{
		channel:     make(chan Message),
		listChannel: map[string]Channel{},
	}

	go func() {
		t.setConcurrence()
	}()

	return t
}

func (t *Trace) setConcurrence() {
	for data := range t.channel {
		for _, c := range t.listChannel {
			if c.toClose {
				t.CloseReader(c)
			} else {
				c.channel <- data
			}
		}
	}
}

func (t *Trace) NewReader() *Channel {
	c := Channel{
		id:      utils.RandomString(5),
		channel: make(chan Message),
	}

	t.listChannel[c.id] = c
	return &c
}

func (t *Trace) CloseReader(ch Channel) {
	_, ok := t.listChannel[ch.GetID()]
	if ok {
		delete(t.listChannel, ch.GetID())
	}
}

func (t *Trace) Push(blockID string, state State) {
	if blockID == "" {
		return
	}
	t.channel <- Message{
		ID:    blockID,
		State: state,
	}
}
