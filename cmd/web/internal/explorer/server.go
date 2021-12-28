package explorer

import (
	"context"
	"github.com/ariden83/blockchain/cmd/web/config"
	"github.com/ariden83/blockchain/cmd/web/internal/auth"
	"github.com/ariden83/blockchain/cmd/web/internal/metrics"
	"github.com/ariden83/blockchain/cmd/web/internal/middleware"
	"github.com/ariden83/blockchain/cmd/web/internal/model"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/urfave/negroni"
	"go.uber.org/zap"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

type Explorer struct {
	log           *zap.Logger
	cfg           *config.Config
	baseURL       string
	server        *http.Server
	router        *mux.Router
	model         *model.Model
	auth          *auth.Auth
	metricsServer *http.Server
	metrics       *metrics.Metrics
}

type Healthz struct {
	Result   bool     `json:"result"`
	Messages []string `json:"messages"`
	Version  string   `json:"version"`
}

func New(options ...func(*Explorer)) *Explorer {
	e := &Explorer{
		router: mux.NewRouter(),
	}
	for _, o := range options {
		o(e)
	}
	return e
}

func WithConfig(cfg *config.Config) func(*Explorer) {
	return func(e *Explorer) {
		e.cfg = cfg
		e.baseURL = "http://localhost" + cfg.BuildPort(cfg.Port)
	}
}

func WithLogs(logs *zap.Logger) func(*Explorer) {
	return func(e *Explorer) {
		e.log = logs.With(zap.String("service", "http"))
	}
}

func WithMetrics(m *metrics.Metrics) func(*Explorer) {
	return func(e *Explorer) {
		e.metrics = m
	}
}

func WithModel(evt *model.Model) func(*Explorer) {
	return func(e *Explorer) {
		e.model = evt
	}
}

func WithAuth(a *auth.Auth) func(*Explorer) {
	return func(e *Explorer) {
		e.auth = a
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

	n := negroni.New(negroni.HandlerFunc(middleware.DefaultHeader))
	n.UseFunc(e.tokenHeader)
	n.UseFunc(e.requestIDHeader)

	n.Use(negroni.HandlerFunc(func(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		route := strings.ToLower(r.Method)
		route = strings.Replace(route, "/", "_", 0)

		jsonHandler := promhttp.InstrumentHandlerInFlight(
			e.metrics.InFlight,

			promhttp.InstrumentHandlerResponseSize(
				e.metrics.ResponseSize.MustCurryWith(prometheus.Labels{"service": route}),

				promhttp.InstrumentHandlerRequestSize(
					e.metrics.RequestSize.MustCurryWith(prometheus.Labels{"service": route}),

					promhttp.InstrumentHandlerCounter(
						e.metrics.RouteCountReqs.MustCurryWith(prometheus.Labels{"service": route}),

						promhttp.InstrumentHandlerDuration(
							e.metrics.ResponseDuration.MustCurryWith(prometheus.Labels{"service": route}),
							next)))))

		jsonHandler.ServeHTTP(rw, r)
	}))

	e.server = &http.Server{
		Addr:           ":" + strconv.Itoa(e.cfg.Port),
		Handler:        n,
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
