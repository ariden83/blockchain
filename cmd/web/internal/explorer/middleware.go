package explorer

import (
	"context"
	"errors"
	"github.com/go-session/session"
	"net/http"
	"os"
	"time"

	uuid "github.com/satori/go.uuid"
	"go.uber.org/zap"

	"github.com/ariden83/blockchain/cmd/web/internal/logger"
)

func (e *Explorer) requestIDHeader(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	var reqID string
	reqID = r.Header.Get(XRequestIDHeaderKey)
	if reqID == "" {
		reqID = uuid.NewV4().String()
	}
	w.Header().Set(XRequestIDHeaderKey, reqID)
	ctx := context.WithValue(r.Context(), RequestIDHeader, reqID)
	ctx = logger.ToContext(ctx, e.log.With(zap.String(RequestIDHeader, reqID)))
	next(w, r)
}

func jsonHeader(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json;charset=UTF-8")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, HEAD")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Accept-ranges", "items")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.Header().Set("Cache-Control", "no-store")
		w.Header().Set("Pragma", "no-cache")
		now := time.Now()
		w.Header().Set("Date", now.String())
		// Call the next handler, which can be another middleware in the chain, or the final handler.
		next.ServeHTTP(w, r)
	})
}

func (e *Explorer) dumpRequest(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	if e.cfg.DumpVar {
		_ = dumpRequest(os.Stdout, "login", r) // Ignore the error
	}
	errors.New("")
	next(w, r)
}

func (e *Explorer) validateToken(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		store, err := session.Start(r.Context(), rw, r)
		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}
		accessToken, ok := store.Get(sessionLabelAccessToken)
		if !ok {
			rw.Header().Set("Location", defaultPageLogin)
			rw.WriteHeader(http.StatusFound)
			return
		}
		if _, err := e.authServer.Manager.LoadAccessToken(r.Context(), accessToken.(string)); err != nil {
			// try to refresh token
			if err = e.refreshToken(rw, r); err != nil {
				rw.Header().Set("Location", defaultPageLogin)
				rw.WriteHeader(http.StatusFound)
				return
			}
		}
		next.ServeHTTP(rw, r)
	})
}

func (e *Explorer) hasValidToken(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		store, err := session.Start(r.Context(), rw, r)
		if err != nil {
			return
		}
		accessToken, ok := store.Get(sessionLabelAccessToken)
		// si pas de token, on reste sur la page
		if !ok {
			next.ServeHTTP(rw, r)
			return
		}
		// si token valide, on change de page
		if _, err := e.authServer.Manager.LoadAccessToken(r.Context(), accessToken.(string)); err == nil {
			rw.Header().Set("Location", defaultPageLogged)
			rw.WriteHeader(http.StatusFound)
			return
		}
		// on reste sur la page
		next.ServeHTTP(rw, r)
	})
}
