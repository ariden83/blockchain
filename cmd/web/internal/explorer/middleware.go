package explorer

import (
	"context"
	"github.com/ariden83/blockchain/cmd/web/internal/logger"
	uuid "github.com/satori/go.uuid"
	"go.uber.org/zap"
	"net/http"
)

func (e *Explorer) tokenHeader(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	/*if err := e.auth.TokenValid(r); err != nil {
		next(rw, r)
		return
	}


	token := e.auth.ExtractToken(r)

	e.setTokenHeaders(rw, ts)*/

	next(rw, r)
}

func (e *Explorer) requestIDHeader(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	var reqID string
	reqID = r.Header.Get(RequestIDHeaderKey)
	if reqID == "" {
		reqID = uuid.NewV4().String()
	}
	w.Header().Set(RequestIDHeaderKey, reqID)
	ctx := context.WithValue(r.Context(), RequestIDKey, reqID)
	ctx = logger.ToContext(ctx, e.log.With(zap.String(RequestIDKey, reqID)))
	next(w, r)
}
