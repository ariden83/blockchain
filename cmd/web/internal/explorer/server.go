package explorer

import (
	"github.com/ariden83/blockchain/cmd/web/config"
	"github.com/ariden83/blockchain/cmd/web/internal/model"
	"github.com/ariden83/blockchain/cmd/web/internal/token"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

type Explorer struct {
	log     *zap.Logger
	cfg     *config.Config
	baseURL string
	server  *http.Server
	router  *mux.Router
	model   *model.Model
	token   *token.Token
}

func New(cfg *config.Config, log *zap.Logger, m *model.Model, t *token.Token) *Explorer {
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
	e.loadRoutes()
	e.loadMiddleware()
	e.listenOrDie()
}

func (Explorer) loadMiddleware() {}

func (e *Explorer) listenOrDie() {
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

	log.Fatal(e.server.ListenAndServe())
}
