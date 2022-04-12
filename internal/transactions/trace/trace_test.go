package trace

import (
	"github.com/ariden83/blockchain/internal/logger"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/ariden83/blockchain/config"
)

func (t *Trace) getListID() map[string]Channel {
	return t.listChannel
}

func Test_trace(t *testing.T) {
	cfg := config.Log{
		Path:     "./tmp/logs",
		CLILevel: "info",
		WithFile: false,
	}

	logs := logger.InitLight(cfg)

	n := New(config.Trace{false}, logs)
	assert.Nil(t, n)

	n = New(config.Trace{true}, logs)

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

	waitChannel := make(chan bool)
	waitChannel1 := make(chan bool)

	i := 0
	j := 0

	go func() {
		for {
			if result, more := <-channel1.GetChan(); more {
				assert.NotEmpty(t, result.ID)
				j++
				waitChannel <- true
			}
		}
	}()

	go func() {
		for {
			if result, more := <-channel.GetChan(); more {
				assert.NotEmpty(t, result.ID)
				i++
				waitChannel1 <- true
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

	<-waitChannel
	<-waitChannel1

	assert.Equal(t, 1, i)
	assert.Equal(t, 1, j)
}
