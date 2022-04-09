package grpc

import (
	"io"

	"go.uber.org/zap"

	"github.com/ariden83/blockchain/pkg/api"
	pkgErr "github.com/ariden83/blockchain/pkg/errors"
)

func (e *EndPoint) GetTraces(_ *api.TraceInput, stream api.Api_GetTracesServer) error {

	channel := e.transaction.Trace()
	if channel == nil {
		return nil
	}

	go func() {
		select {
		case <-stream.Context().Done():
			channel.Close()
			break
		case <-e.stop:
			channel.Close()
			break
		}
		return
	}()

	for {
		if result, more := <-channel.GetChan(); more {
			if result.ID == "" {
				continue
			}
			err := stream.Send(&api.TraceOutput{
				Id:    result.ID,
				State: result.State.String(),
			})

			if err == io.EOF {
				e.log.Info("stream Reached EOF")
				return nil
			}

			if err != nil {
				e.log.Error("error on sending trace to stream", zap.Error(err), zap.String("id", result.ID))
				return pkgErr.ErrInternalError
			}

		} else {
			break
		}
	}

	return nil
}
