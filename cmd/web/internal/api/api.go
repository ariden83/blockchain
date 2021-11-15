package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/ariden83/blockchain/cmd/web/config"
	"go.uber.org/zap"
	"io"
	"net/http"
	"os"
	"time"
)

type Model struct {
	log     *zap.Logger
	baseURL string
	client  *http.Client
	TimeOut float64
}

type PostInput interface{}
type PostOutput interface{}

type Option func(e *Model)

func New(cfg *config.Config, log *zap.Logger) *Model {
	return &Model{
		log:     log,
		baseURL: "http://localhost" + cfg.BuildPort(cfg.Api.Port),
		client: &http.Client{
			Timeout: time.Duration(float64(time.Second) * cfg.Api.TimeOut),
		},
	}
}

func (m *Model) Post(path string, p PostInput) (io.ReadCloser, error) {
	postBody, err := json.Marshal(p)
	if err != nil {
		m.log.Error("An Error Occured when marshal post body", zap.Error(err), zap.String("path", path))
		return nil, err
	}
	reqBody := bytes.NewBuffer(postBody)
	//Leverage Go's HTTP Post function to make request
	req, err := http.NewRequest("POST", m.baseURL+path, reqBody)

	req.Header.Add("Accept", `application/json`)
	// add header for authentication
	req.Header.Add("Authorization", fmt.Sprintf("token %s", os.Getenv("TOKEN")))

	resp, err := m.client.Do(req)
	if err != nil {
		m.log.Error("An Error Occured when call api", zap.Error(err), zap.String("path", path))
		return nil, err
	}
	defer resp.Body.Close()
	return resp.Body, nil
}
