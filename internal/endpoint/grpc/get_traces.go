package grpc

import (
	"fmt"
	"io"

	"github.com/ariden83/blockchain/pkg/api"
)

func (e *EndPoint) GetTraces(_ *api.TraceInput, stream api.Api_GetTracesServer) error {
	fmt.Println(fmt.Sprintf("******************************************************* GetTraces start"))
	channel := e.transaction.Trace()
	if channel == nil {
		return nil
	}

	defer channel.Close()

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

		select {
		case <-stream.Context().Done():
			fmt.Println(fmt.Sprintf("******************************************************* stream.Context().Done()"))
			channel.Close()
			return nil
		case <-e.stop:
			fmt.Println(fmt.Sprintf("******************************************************* e stop"))
			channel.Close()
			return nil
		default:
		}
	}

	fmt.Println(fmt.Sprintf("******************************************************* e stop"))
	return nil
}
