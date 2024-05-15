package trace

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"

	"github.com/ariden83/blockchain/config"
	"github.com/ariden83/blockchain/internal/logger"
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

func Test_PushAndRead(t *testing.T) {
	// Create a new Trace instance
	trace := New(config.Trace{Enabled: true}, zap.NewNop())
	assert.NotNil(t, trace)

	// Create a new trace reader channel
	reader := trace.NewReader()
	assert.NotNil(t, reader)

	// Push a message onto the trace channel
	trace.Push("block123", Mining)

	// Read the message from the trace reader channel
	select {
	case msg := <-reader.GetChan():
		assert.Equal(t, "block123", msg.ID)
		assert.Equal(t, Mining, msg.State)
	case <-time.After(1 * time.Second):
		t.Error("Timeout while waiting for message")
	}

	// Close the trace reader channel
	trace.CloseReader(*reader)
}

func Test_PushWithoutTrace(t *testing.T) {
	// Create a new Trace instance with tracing disabled
	trace := New(config.Trace{Enabled: false}, zap.NewNop())
	assert.Nil(t, trace)
}

func Test_CloseReader(t *testing.T) {
	// Create a new Trace instance
	trace := New(config.Trace{Enabled: true}, zap.NewNop())
	assert.NotNil(t, trace)

	// Create a new trace reader channel
	reader := trace.NewReader()
	assert.NotNil(t, reader)

	// Close the trace reader channel
	trace.CloseReader(*reader)

	// Try to read from the closed trace reader channel (should not block)
	select {
	case <-reader.GetChan():
		t.Error("Received unexpected message from closed channel")
	default:
		// No message received, as expected
	}
}
