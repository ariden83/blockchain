package http

import (
	"encoding/json"
	"go.uber.org/zap"
	"io"
	"net/http"
)

type BodyReceived interface{}

func (e *EndPoint) decodeBody(rw http.ResponseWriter, logCTX *zap.Logger, body io.ReadCloser, req BodyReceived) error {
	var err error
	body = http.MaxBytesReader(rw, body, 1048)

	dec := json.NewDecoder(body)
	dec.DisallowUnknownFields()

	if err = dec.Decode(&req); err != nil {
		logCTX.Error("fail to decode input", zap.Error(err))
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return err
	}
	if err = dec.Decode(&struct{}{}); err != io.EOF {
		msg := "Request body must only contain a single JSON object"
		logCTX.Error(msg, zap.Error(err))
		http.Error(rw, msg, http.StatusBadRequest)
		return err
	}
	return nil
}

func (e *EndPoint) JSON(w http.ResponseWriter, code int, payload interface{}) {
	response, err := json.MarshalIndent(payload, "", "  ")
	if err != nil {
		e.log.Error("HTTP 500: Internal Server Error", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		if _, err = w.Write([]byte("HTTP 500: Internal Server Error")); err != nil {
			e.log.Error("fail to write response", zap.Error(err))
		}
		return
	}
	w.WriteHeader(code)
	if _, err = w.Write(response); err != nil {
		e.log.Error("fail to write response", zap.Error(err))
	}
}
