package explorer

import (
	"context"
	"github.com/ariden83/blockchain/cmd/web/internal/logger"
	uuid "github.com/satori/go.uuid"
	"go.uber.org/zap"
	"net/http"
)

func (e *Explorer) tokenHeader(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	userKey := r.Header.Get("User_Key")
	if userKey == "" {
		next(rw, r)
		return
	}

	ts, err := e.auth.CreateToken(userKey)
	if err != nil {
		e.fail(http.StatusUnprocessableEntity, err, rw)
		return
	}
	saveErr := e.auth.CreateAuth(r.Context(), userKey, ts)
	if saveErr != nil {
		e.fail(http.StatusUnprocessableEntity, err, rw)
		return
	}

	e.setTokenHeaders(rw, ts)
	e.setUserKeyHeaders(rw, userKey)

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
