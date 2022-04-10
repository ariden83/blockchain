package trace

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/ariden83/blockchain/config"
)

func (t *Trace) getListID() map[string]Channel {
	return t.listChannel
}

func Test_trace(t *testing.T) {
	n := New(config.Trace{false})
	assert.Nil(t, n)

	n = New(config.Trace{true})

	list := n.getListID()
	assert.Equal(t, 0, len(list))

	channel := n.NewReader()
	assert.NotNil(t, channel)
	defer channel.Close()

	channel1 := n.NewReader()
	assert.NotNil(t, channel1)
	defer channel1.Close()

	c := channel.GetChan()
	assert.NotNil(t, c)
	id := channel.GetID()
	assert.NotEmpty(t, id)

	go func() {
		for {
			if result, more := <-channel.GetChan(); more {
				assert.NotEmpty(t, result.ID)
			}
		}
	}()

	list = n.getListID()
	assert.Equal(t, 2, len(list))
	assert.NotEmpty(t, list[id])

	n.Push("id", 1)
	n.CloseReader(*channel)

	list = n.getListID()
	assert.Equal(t, 1, len(list))
}
