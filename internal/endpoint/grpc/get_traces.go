package grpc

import (
	"go.uber.org/zap"

	"github.com/ariden83/blockchain/pkg/api"
	pkgErr "github.com/ariden83/blockchain/pkg/errors"
)

func (e *EndPoint) GetTraces(_ *api.TraceInput, stream api.Api_GetTracesServer) error {

	channel := e.transaction.Trace()
	if channel == nil {
		return nil
	}

	defer channel.Close()

	for {
		if result, more := <-channel.GetChan(); more {

			if err := stream.Send(&api.TraceOutput{
				Id:    result.ID,
				State: result.State.String(),
			}); err != nil {
				e.log.Error("error on sending trace to stream", zap.Error(err), zap.String("id", result.ID))
				return pkgErr.ErrInternalError
			}

		} else {
			break
		}
	}

	return nil
}
