package grpc

import (
	"io"

	"github.com/ariden83/blockchain/pkg/api"
)

func (e *EndPoint) GetTraces(_ *api.TraceInput, stream api.Api_GetTracesServer) error {
	channel := e.transaction.Trace()
	if channel == nil {
		return nil
	}

	defer func() {
		e.transaction.CloseTrace(*channel)
		e.log.Info("close get traces")
	}()

	go func() {
		select {
		case <-stream.Context().Done():
			channel.Close()
			return
		case <-e.stop:
			channel.Close()
		}
	}()

	for {
		if result, more := <-channel.GetChan(); more {
			err := stream.Send(&api.TraceOutput{
				Id:    result.ID,
				State: result.State.String(),
			})

			if err == io.EOF {
				return nil
			} else if err != nil {
				return err
			}
		} else {
			break
		}
	}

	return nil
}
