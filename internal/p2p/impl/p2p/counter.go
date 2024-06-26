package p2p

import (
	"encoding/binary"
	"time"

	"github.com/libp2p/go-libp2p/core/network"
	"go.uber.org/zap"
)

func (e *Adapter) handleCounter(s network.Stream) {
	go e.writeCounter(s)
	go e.readCounter(s)
}

func (e *Adapter) writeCounter(s network.Stream) {
	var counter uint64

	for {
		<-time.After(time.Second)
		counter++

		err := binary.Write(s, binary.BigEndian, counter)
		if err != nil {
			e.log.Error("fail to write binary in counter", zap.Uint64("counter", counter), zap.Error(err))
			break
		}
	}
}

func (e *Adapter) readCounter(s network.Stream) {
	for {
		var counter uint64

		err := binary.Read(s, binary.BigEndian, &counter)
		if err != nil {
			e.log.Error("fail to read binary in counter", zap.Error(err))
			break
		}
	}
}
