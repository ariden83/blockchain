package model

import (
	"context"
	"github.com/ariden83/blockchain/cmd/web/config"
	"github.com/ariden83/blockchain/pkg/api"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"io"
	"time"
)

type IModel interface {
	Post(string, PostInput) (io.ReadCloser, error)
}

type Model struct {
	log     *zap.Logger
	client  *grpc.ClientConn
	timeOut float64
}

type PostInput interface{}
type PostOutput interface{}

type Option func(e *Model)

func New(cfg config.BlockchainAPI, log *zap.Logger) (*Model, error) {

	opts := []grpc.DialOption{}
	opts = append(opts, grpc.WithInsecure())
	opts = append(opts, grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(cfg.MaxSizeCall)))
	conn, err := grpc.Dial(cfg.URL, opts...)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	return &Model{
		log:     log.With(zap.String("service", "grpc")),
		timeOut: cfg.TimeOut,
		client:  conn,
	}, nil
}

func (m *Model) ShutDown() {
	if m.client != nil {
		if err := m.client.Close(); err != nil {
			m.log.Error("fail to close connexion", zap.Error(err))
		}
	}
}

func (m *Model) GetWallet(ctx context.Context, mnemonic string) (*api.GetWalletOutput, error) {
	c := api.NewApiClient(m.client)

	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(m.timeOut)*time.Second)
	defer cancel()

	search := api.GetWalletInput{
		Mnemonic: mnemonic,
	}
	data, err := c.GetWallet(ctx, &search)
	if err != nil {
		m.log.Info("Cannot connect get user wallet", zap.Error(err))
		return data, err

	}
	return data, nil
}
