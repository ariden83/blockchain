package explorer

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/ariden83/blockchain/cmd/web/config"
	"github.com/ariden83/blockchain/cmd/web/internal/metrics"
	"github.com/ariden83/blockchain/cmd/web/internal/model"
	"github.com/ariden83/blockchain/cmd/web/internal/token"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
	"net/http"
	"net/http/pprof"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

type Explorer struct {
	log           *zap.Logger
	cfg           *config.Config
	baseURL       string
	server        *http.Server
	router        *mux.Router
	model         *model.Model
	token         *token.Token
	metricsServer *http.Server
	metrics       *metrics.Metrics
}

type Healthz struct {
	Result   bool     `json:"result"`
	Messages []string `json:"messages"`
	Version  string   `json:"version"`
}

func New(cfg *config.Config, log *zap.Logger, m *model.Model, t *token.Token, mtc *metrics.Metrics) *Explorer {
	return &Explorer{
		log:     log,
		cfg:     cfg,
		baseURL: "http://localhost" + cfg.BuildPort(cfg.Port),
		router:  mux.NewRouter(),
		model:   m,
		token:   t,
		metrics: mtc,
	}
}

func (e *Explorer) Start(stop chan error) {
	e.log.Info("start web server")
	e.loadTemplates()
	e.loadRoutes()
	e.loadMiddleware()
	e.listenOrDie(stop)
}

func (Explorer) loadMiddleware() {}

func (e *Explorer) listenOrDie(stop chan error) {
	e.log.Info("Start listening HTTP Server", zap.Int("port", e.cfg.Port))

	_, err := os.Stat(filepath.Join(e.cfg.StaticDir, "index.css"))
	if err != nil {
		e.log.Fatal("fail to read index.css")
	}

	mux := http.NewServeMux()

	fileServer := http.FileServer(http.Dir(e.cfg.StaticDir))

	mux.Handle("/static/", http.StripPrefix("/static", fileServer))
	mux.Handle("/", e.router)

	e.server = &http.Server{
		Addr:           ":" + strconv.Itoa(e.cfg.Port),
		Handler:        mux,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	if err = e.server.ListenAndServe(); err != nil {
		stop <- err
	}
}

func (e *Explorer) StartMetricsServer(stop chan error) {
	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		message := "The service " + e.cfg.Name + " responds correctly"
		res := Healthz{Result: true, Messages: []string{message}, Version: e.cfg.Version}
		js, err := json.Marshal(res)
		if err != nil {
			e.log.Fatal("Fail to jsonify", zap.Error(err))
		}
		if _, err := w.Write(js); err != nil {
			e.log.Fatal("Fail to Write response in http.ResponseWriter", zap.Error(err))
			return
		}
	})

	mux.HandleFunc("/readiness", func(w http.ResponseWriter, r *http.Request) {
		result := true
		message := "The service " + e.cfg.Name + " responds correctly"

		res := Healthz{Result: result, Messages: []string{message}, Version: e.cfg.Version}
		js, err := json.Marshal(res)
		if err != nil {
			e.log.Fatal("Fail to jsonify", zap.Error(err))
		}
		if _, err := w.Write(js); err != nil {
			e.log.Fatal("Fail to Write response in http.ResponseWriter", zap.Error(err))
			return
		}
	})

	mux.Handle("/metrics", promhttp.Handler())
	e.PProf(mux)

	addr := fmt.Sprintf("%s:%d", e.cfg.Metrics.Host, e.cfg.Metrics.Port)
	e.metricsServer = &http.Server{
		Addr:           addr,
		Handler:        mux,
		ReadTimeout:    time.Duration(e.cfg.Healthz.ReadTimeout) * time.Second,
		WriteTimeout:   time.Duration(e.cfg.Healthz.WriteTimeout) * time.Second,
		IdleTimeout:    120 * time.Second,
		MaxHeaderBytes: 1 << 12,
	}
	go func() {
		e.log.Info("Listening HTTP for healthz route", zap.String("address", addr))
		if err := e.metricsServer.ListenAndServe(); err != nil {
			stop <- err
		}
	}()
}

func (e *Explorer) PProf(mux *http.ServeMux) {
	mux.HandleFunc("/debug/pprof/", pprof.Index)
	mux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
	mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	mux.Handle("/debug/pprof/allocs", pprof.Handler("allocs"))
	mux.Handle("/debug/pprof/heap", pprof.Handler("heap"))
	mux.Handle("/debug/pprof/mutex", pprof.Handler("mutex"))
	mux.Handle("/debug/pprof/goroutine", pprof.Handler("goroutine"))
	mux.Handle("/debug/pprof/block", pprof.Handler("block"))
	mux.Handle("/debug/pprof/threadcreate", pprof.Handler("threadcreate"))
	mux.Handle("/debug/pprof/trace", pprof.Handler("trace"))
}

func (e *Explorer) Shutdown(ctx context.Context) {
	e.server.Shutdown(ctx)
	e.metricsServer.Shutdown(ctx)
}
