package p2p

import (
	"fmt"
	"go.uber.org/zap"
)

var streamResetError = fmt.Errorf("stream reset")

func (e *EndPoint) Handle(err error) {
	if err != nil {
		e.log.Fatal("fatal error, can't continue", zap.Error(err))
	}
}
