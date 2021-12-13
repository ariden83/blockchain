package http

import "go.uber.org/zap"

func (e *EndPoint) Handle(err error) {
	if err != nil {
		e.log.Fatal("fatal error, can't continue", zap.Error(err))
	}
}
