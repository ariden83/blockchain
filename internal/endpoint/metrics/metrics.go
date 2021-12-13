package metrics

import (
	"context"
	"encoding/json"
	"github.com/ariden83/blockchain/config"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
	"net/http"
	"net/http/pprof"
	"strconv"
	"time"
)

type EndPoint struct {
	cfg    config.Metrics
	server *http.Server
	log    *zap.Logger
}

type Healthz struct {
	Result   bool     `json:"result"`
	Messages []string `json:"messages"`
	Version  string   `json:"version"`
}

func New(
	cfg config.Metrics,
	logs *zap.Logger,
) *EndPoint {
	e := &EndPoint{
		cfg: cfg,
		log: logs.With(zap.String("service", "metrics")),
	}

	return e
}

func (e *EndPoint) Listen() error {
	mux := e.router()

	e.server = &http.Server{
		Addr:           ":" + strconv.Itoa(e.cfg.Port),
		Handler:        mux,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		IdleTimeout:    120 * time.Second,
		MaxHeaderBytes: 1 << 12,
	}

	e.log.Info("Metrics Server start", zap.Int("port", e.cfg.Port))
	return e.server.ListenAndServe()
}

func (e *EndPoint) router() http.Handler {
	muxRouter := mux.NewRouter()
	muxRouter.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		message := "The service responds correctly"
		res := Healthz{Result: true, Messages: []string{message}}
		js, err := json.Marshal(res)
		if err != nil {
			e.log.Fatal("Fail to jsonify healthz response", zap.Error(err))
		}
		if _, err := w.Write(js); err != nil {
			e.log.Fatal("Fail to Write response in http.ResponseWriter", zap.Error(err))
			return
		}
	})

	muxRouter.HandleFunc("/readiness", func(w http.ResponseWriter, r *http.Request) {
		result := true
		message := "Service responds correctly"

		res := Healthz{Result: result, Messages: []string{message}}
		js, err := json.Marshal(res)
		if err != nil {
			e.log.Fatal("Fail to jsonify", zap.Error(err))
		}
		if _, err := w.Write(js); err != nil {
			e.log.Fatal("Fail to Write response in http.ResponseWriter", zap.Error(err))
			return
		}
	})

	muxRouter.Handle("/metrics", promhttp.Handler())

	muxRouter.HandleFunc("/debug/pprof/", pprof.Index)
	muxRouter.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	muxRouter.HandleFunc("/debug/pprof/profile", pprof.Profile)
	muxRouter.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	muxRouter.Handle("/debug/pprof/allocs", pprof.Handler("allocs"))
	muxRouter.Handle("/debug/pprof/heap", pprof.Handler("heap"))
	muxRouter.Handle("/debug/pprof/mutex", pprof.Handler("mutex"))
	muxRouter.Handle("/debug/pprof/goroutine", pprof.Handler("goroutine"))
	muxRouter.Handle("/debug/pprof/block", pprof.Handler("block"))
	muxRouter.Handle("/debug/pprof/threadcreate", pprof.Handler("threadcreate"))
	muxRouter.Handle("/debug/pprof/trace", pprof.Handler("trace"))

	return muxRouter
}

func (e *EndPoint) Shutdown(ctx context.Context) {
	err := e.server.Shutdown(ctx)
	if err != nil {
		e.log.Error("fail to shutdown server", zap.Error(err))
	}
}
