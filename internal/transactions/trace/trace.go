package trace

import (
	"fmt"
	"github.com/ariden83/blockchain/config"
	"github.com/ariden83/blockchain/internal/utils"
	"go.uber.org/zap"
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
	fmt.Println(fmt.Sprintf("******************************************************* close with id %s", c.id))
	if c.channel != nil {
		close(c.channel)
	}
}

type Trace struct {
	channel     chan Message
	listChannel map[string]Channel
	log         *zap.Logger
}

func New(cfg config.Trace, log *zap.Logger) *Trace {
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

func (t *Trace) setConcurrence() {
	for data := range t.channel {
		for _, c := range t.listChannel {
			if c.toClose || c.channel == nil {
				fmt.Println(fmt.Sprintf("******************************************** setConcurrence close"))
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
	fmt.Println(fmt.Sprintf("******************************************** send message in channel %s %+v", blockID, state))
	t.channel <- Message{
		ID:    blockID,
		State: state,
	}
}
