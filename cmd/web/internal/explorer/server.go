package explorer

import (
	"github.com/ariden83/blockchain/cmd/web/internal/api"
	"github.com/ariden83/blockchain/cmd/web/internal/config"
	"github.com/ariden83/blockchain/cmd/web/internal/token"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
	"log"
	"net/http"
	"strconv"
	"time"
)

type Explorer struct {
	log     *zap.Logger
	cfg     *config.Config
	baseURL string
	server  *http.Server
	router  *mux.Router
	model   *api.Model
	token   *token.Token
}

func New(cfg *config.Config, log *zap.Logger, m *api.Model, t *token.Token) *Explorer {
	return &Explorer{
		log:     log,
		cfg:     cfg,
		baseURL: "http://localhost" + cfg.BuildPort(cfg.Port),
		router:  mux.NewRouter(),
		model:   m,
		token:   t,
	}
}

func (e *Explorer) Start() {
	e.log.Info("start")
	e.loadTemplates()
	e.loadFileServer()
	e.loadRoutes()
	e.loadMiddleware()
	e.listenOrDie()
}

func (Explorer) loadMiddleware() {}

func (e *Explorer) listenOrDie() {
	e.log.Info("Start listening HTTP Server", zap.Int("port", e.cfg.Port))

	e.server = &http.Server{
		Addr:           ":" + strconv.Itoa(e.cfg.Port),
		Handler:        e.router,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	log.Fatal(e.server.ListenAndServe())
}
