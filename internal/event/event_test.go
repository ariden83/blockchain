package event

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"

	"github.com/ariden83/blockchain/internal/event/trace"
)

func Test_New_Event(t *testing.T) {
	t.Run("new event adapter", func(t *testing.T) {
		eventAdapter := New()
		assert.NotNil(t, eventAdapter)
	})

	t.Run("new event with WithTraces option", func(t *testing.T) {
		t.Run("new event with WithTraces option: enabled", func(t *testing.T) {
			eventAdapter := New(WithTraces(trace.Config{Enabled: true}, zap.NewNop()))
			assert.NotNil(t, eventAdapter)
			assert.NotNil(t, eventAdapter.trace)
		})

		t.Run("new event with WithTraces option: not enabled", func(t *testing.T) {
			eventAdapter := New(WithTraces(trace.Config{Enabled: false}, zap.NewNop()))
			assert.NotNil(t, eventAdapter)
			assert.Nil(t, eventAdapter.trace)
		})
	})

	t.Run("push message and received it", func(t *testing.T) {
		typeProvided := NewBlock
		eventAdapter := New()
		done := make(chan bool, 2)
		ready := make(chan bool, 2)

		go func() {
			newChan := eventAdapter.NewReader()
			ready <- true
			for {
				mess := <-newChan
				assert.Equal(t, typeProvided, mess.Type)
				done <- true
				return
			}
		}()

		go func() {
			newChan := eventAdapter.NewReader()
			ready <- true
			for {
				mess := <-newChan
				assert.Equal(t, typeProvided, mess.Type)
				done <- true
				return
			}
		}()

		<-ready
		<-ready

		eventAdapter.Push(Message{
			Type:  typeProvided,
			ID:    "ID-2",
			Value: []byte("message-value"),
		})

		assert.Len(t, eventAdapter.listChannel, 2)
		<-done
		<-done

		eventAdapter.CloseReaders()
	})
}

func Test_PushAndRead(t *testing.T) {
	// Create a new Event instance
	event := New(WithTraces(trace.Config{Enabled: true}, zap.NewNop()))
	assert.NotNil(t, event)

	// Create a new message reader channel
	reader := event.NewReader()
	assert.NotNil(t, reader)

	// Push a message onto the event channel
	event.Push(Message{Type: BlockChain, ID: "block123", Value: []byte("data")})

	// Read the message from the event reader channel
	select {
	case msg := <-reader:
		assert.Equal(t, BlockChain, msg.Type)
		assert.Equal(t, "block123", msg.ID)
		assert.Equal(t, []byte("data"), msg.Value)
		event.CloseReaders()
	case <-time.After(1 * time.Second):
		t.Error("Timeout while waiting for message")
	}
}

func Test_PushWithoutTrace(t *testing.T) {
	// Create a new Event instance with tracing disabled
	event := New(WithTraces(trace.Config{Enabled: true}, zap.NewNop()))
	assert.NotNil(t, event)
	defer event.CloseReaders()

	// Try to push a message (should not panic)
	assert.NotPanics(t, func() {
		event.Push(Message{Type: BlockChain, ID: "block123", Value: []byte("data")})
	})
}

func Test_CloseReaders(t *testing.T) {
	// Create a new Event instance
	event := New(WithTraces(trace.Config{Enabled: true}, zap.NewNop()))
	assert.NotNil(t, event)

	// Create a new message reader channel
	reader := event.NewReader()
	assert.NotNil(t, reader)

	// Close all message reader channels
	event.CloseReaders()

	// Try to read from the closed message reader channel and ensure no value is received
	select {
	case msg, ok := <-reader:
		if ok {
			t.Errorf("Received unexpected message from closed channel: %v", msg)
		}
	default:
		// No message received, as expected
	}
}
