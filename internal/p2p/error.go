package p2p

import (
	"go.uber.org/zap"
)

var (
	failNegociateError     = "failed to negotiate security protocol: peer id mismatch"
	protocolError          = "protocol not supported"
	addressAMReadyUseError = "bind: address already in use"
	noGoodAddress          = "no good addresses"
)

func (e *EndPoint) Handle(err error) {
	if err != nil {
		e.log.Fatal("fatal error, can't continue", zap.Error(err))
	}
}
