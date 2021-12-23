package explorer

import (
	"context"
	"github.com/ariden83/blockchain/cmd/web/config"
	"github.com/ariden83/blockchain/cmd/web/internal/metrics"
	"github.com/ariden83/blockchain/cmd/web/internal/model"
	"github.com/ariden83/blockchain/cmd/web/internal/token"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
	"net/http"
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
	e.loadAPIRoutes()
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

func (e *Explorer) Shutdown(ctx context.Context) {
	e.server.Shutdown(ctx)
	e.metricsServer.Shutdown(ctx)
}
