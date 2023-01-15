package event

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"

	"github.com/ariden83/blockchain/config"
)

func Test_New_Event(t *testing.T) {
	t.Run("new event adapter", func(t *testing.T) {
		eventAdapter := New()
		assert.NotNil(t, eventAdapter)
	})

	t.Run("new event with WithTraces option", func(t *testing.T) {
		t.Run("new event with WithTraces option: enabled", func(t *testing.T) {
			eventAdapter := New(WithTraces(config.Trace{Enabled: true}, zap.NewNop()))
			assert.NotNil(t, eventAdapter)
			assert.NotNil(t, eventAdapter.trace)
		})

		t.Run("new event with WithTraces option: not enabled", func(t *testing.T) {
			eventAdapter := New(WithTraces(config.Trace{Enabled: false}, zap.NewNop()))
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
