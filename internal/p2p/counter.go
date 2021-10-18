package p2p

import (
	"encoding/binary"
	"fmt"
	net "github.com/libp2p/go-libp2p-core"
	"github.com/libp2p/go-libp2p-core/network"
	"go.uber.org/zap"
	"time"
)

func (e *EndPoint) handleCounter(s net.Stream) {
	go e.writeCounter(s)
	go e.readCounter(s)
}

func (e *EndPoint) writeCounter(s network.Stream) {
	var counter uint64

	for {
		<-time.After(time.Second)
		counter++

		err := binary.Write(s, binary.BigEndian, counter)
		if err != nil {
			e.log.Fatal("fail to write binary in counter", zap.Error(err))
		}
	}
}

func (e *EndPoint) readCounter(s network.Stream) {
	for {
		var counter uint64

		err := binary.Read(s, binary.BigEndian, &counter)
		if err != nil {
			e.log.Fatal("fail to read binary in counter", zap.Error(err))
		}

		e.log.Info(fmt.Sprintf("Received %d from %s", counter, s.ID()))
	}
}
