package model

import (
	"context"
	"io"
	"time"

	"go.uber.org/zap"
	"google.golang.org/grpc"

	"github.com/ariden83/blockchain/cmd/web/config"
	"github.com/ariden83/blockchain/pkg/api"
)

type IModel interface {
	Post(string, PostInput) (io.ReadCloser, error)
}

type Model struct {
	log       *zap.Logger
	connexion *grpc.ClientConn
	client    *api.ApiClient
	timeOut   float64
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
	logCTX := log.With(zap.String("service", "grpc"))
	logCTX.Info("init grpc connexion", zap.String("url", cfg.URL))

	client := api.NewApiClient(conn)

	return &Model{
		log:       logCTX,
		timeOut:   cfg.TimeOut,
		connexion: conn,
		client:    &client,
	}, nil
}

func (m *Model) ShutDown() {
	if m.client != nil {
		if err := m.connexion.Close(); err != nil {
			m.log.Error("fail to close connexion", zap.Error(err))
		}
	}
}

func (m *Model) GetWallet(ctx context.Context, mnemonic, password []byte) (*api.GetWalletOutput, error) {
	var cancel context.CancelFunc
	ctx, cancel = context.WithTimeout(ctx, time.Duration(m.timeOut)*time.Second)
	defer cancel()

	data, err := (*m.client).GetWallet(ctx, &api.GetWalletInput{
		Password: password,
		Mnemonic: mnemonic,
	})
	if err != nil {
		m.log.Info("Cannot get user wallet", zap.Error(err))
		return data, err
	}
	return data, nil
}

func (m *Model) CreateWallet(ctx context.Context, password []byte) (*api.CreateWalletOutput, error) {
	var cancel context.CancelFunc
	ctx, cancel = context.WithTimeout(ctx, time.Duration(m.timeOut)*time.Second)
	defer cancel()

	data, err := (*m.client).CreateWallet(ctx, &api.CreateWalletInput{
		Password: password,
	})
	if err != nil {
		m.log.Info("Cannot create user wallet", zap.Error(err))
		return data, err
	}
	return data, nil
}

func (m *Model) GetBalance(ctx context.Context, privKey string) (*api.GetBalanceOutput, error) {
	var cancel context.CancelFunc
	ctx, cancel = context.WithTimeout(ctx, time.Duration(m.timeOut)*time.Second)
	defer cancel()

	data, err := (*m.client).GetBalance(ctx, &api.GetBalanceInput{
		PrivKey: []byte(privKey),
	})
	if err != nil {
		m.log.Info("Cannot connect get user wallet", zap.Error(err))
		return data, err

	}
	return data, nil
}

func (m *Model) ValidWallet(ctx context.Context, privKey []byte) (*api.ValidWalletOutput, error) {
	var cancel context.CancelFunc
	ctx, cancel = context.WithTimeout(ctx, time.Duration(m.timeOut)*time.Second)
	defer cancel()

	data, err := (*m.client).ValidWallet(ctx, &api.ValidWalletInput{
		PrivKey: privKey,
	})
	if err != nil {
		m.log.Info("Cannot valid inscription", zap.Error(err))
		return data, err
	}
	return data, nil
}
