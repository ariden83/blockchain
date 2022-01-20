package explorer

import (
	"context"
	"github.com/ariden83/blockchain/cmd/web/config"
	"github.com/ariden83/blockchain/cmd/web/internal/auth"
	"github.com/ariden83/blockchain/cmd/web/internal/auth/classic"
	"github.com/ariden83/blockchain/cmd/web/internal/locales"
	"github.com/ariden83/blockchain/cmd/web/internal/metrics"
	"github.com/ariden83/blockchain/cmd/web/internal/model"
	"github.com/ariden83/blockchain/cmd/web/internal/recaptcha"
	"github.com/go-oauth2/oauth2/v4"
	"github.com/go-oauth2/oauth2/v4/errors"
	"github.com/go-oauth2/oauth2/v4/generates"
	"github.com/go-oauth2/oauth2/v4/manage"
	"github.com/go-oauth2/oauth2/v4/models"
	"github.com/go-oauth2/oauth2/v4/server"
	"github.com/go-oauth2/oauth2/v4/store"
	"github.com/go-session/session"
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
	metadata      config.Metadata
	baseURL       string
	server        *http.Server
	router        *mux.Router
	model         *model.Model
	auth          *auth.Auth
	metricsServer *http.Server
	metrics       *metrics.Metrics
	authServer    *server.Server
	reCaptcha     *recaptcha.Captcha
	locales       *locales.Locales
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

func WithMetadata(metadata config.Metadata) func(*Explorer) {
	return func(e *Explorer) {
		e.metadata = metadata
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

func WithLocales(cfg config.Locales) func(*Explorer) {
	return func(e *Explorer) {
		e.locales = locales.New(cfg)
	}
}

func WithRecaptcha(cfg config.ReCaptcha, log *zap.Logger) func(*Explorer) {
	c := recaptcha.New(cfg, log)
	if c == nil {
		return func(e *Explorer) {}
	}
	return func(e *Explorer) {
		e.reCaptcha = c
	}
}

func (e *Explorer) Start(stop chan error) {
	e.log.Info("start web server")
	e.loadTemplates()
	e.loadRoutes()
	e.loadConnectedRoutes()
	e.loadNonConnectedRoutes()

	e.loadAPIRoutes()
	e.loadAPIConnectedRoutes()
	e.loadAPINonConnectedRoutes()

	e.listenOrDie(stop)
}

func (e *Explorer) listenOrDie(stop chan error) {
	e.log.Info("Start listening HTTP Server", zap.Int("port", e.cfg.Port))

	_, err := os.Stat(filepath.Join(e.cfg.StaticDir, "index.css"))
	if err != nil {
		e.log.Fatal("fail to read index.css")
	}
	mux := http.NewServeMux()
	fileServer := http.FileServer(http.Dir(e.cfg.StaticDir))
	e.manageAuth()
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))
	mux.Handle("/", e.router)
	n := negroni.New()
	// n.UseFunc(e.tokenHeader)
	n.UseFunc(e.requestIDHeader)
	n.UseFunc(e.dumpRequest)
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

	n.UseHandler(mux)
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

func (e *Explorer) manageAuth() {
	classicAuth, ok := e.auth.API[classic.Name]
	if !ok {
		return
	}

	manager := manage.NewDefaultManager()

	manager.SetAuthorizeCodeExp(time.Minute * 10)
	manager.SetPasswordTokenCfg(manage.DefaultPasswordTokenCfg)
	manager.SetAuthorizeCodeTokenCfg(manage.DefaultAuthorizeCodeTokenCfg)
	manager.SetRefreshTokenCfg(manage.DefaultRefreshTokenCfg)
	manager.SetClientTokenCfg(manage.DefaultClientTokenCfg)
	manager.SetImplicitTokenCfg(&manage.Config{AccessTokenExp: time.Hour * 2, RefreshTokenExp: time.Hour * 24 * 7, IsGenerateRefresh: true})

	manager.MustTokenStorage(store.NewMemoryTokenStore())
	// generate jwt access token
	// manager.MapAccessGenerate(generates.NewJWTAccessGenerate("", []byte("00000000"), jwt.SigningMethodHS512))
	manager.MapAccessGenerate(generates.NewAccessGenerate())

	clientStore := store.NewClientStore()
	clientStore.Set(classicAuth.Config().ClientID, &models.Client{
		ID:     classicAuth.Config().ClientID,
		Secret: classicAuth.Config().ClientSecret,
	})

	manager.MapClientStorage(clientStore)
	e.authServer = server.NewServer(server.NewConfig(), manager)
	e.authServer.SetAllowGetAccessRequest(true)
	e.authServer.SetClientInfoHandler(server.ClientFormHandler)
	e.authServer.SetRefreshingValidationHandler(func(ti oauth2.TokenInfo) (allowed bool, err error) {
		return true, nil
	})
	e.authServer.SetUserAuthorizationHandler(func(w http.ResponseWriter, r *http.Request) (userID string, err error) {
		store, err := session.Start(r.Context(), w, r)
		if err != nil {
			return
		}

		data, ok := store.Get(sessionLabelUserID)
		if !ok {
			w.Header().Set("Location", defaultPageLogin)
			w.WriteHeader(http.StatusFound)
			return
		}

		userID = data.(string)
		store.Delete(sessionLabelUserID)
		store.Save()
		return
	})
	e.authServer.SetInternalErrorHandler(func(err error) (re *errors.Response) {
		e.log.Error("Internal Error", zap.Error(err))
		return
	})
	e.authServer.SetResponseErrorHandler(func(re *errors.Response) {
		e.log.Error("Response Error", zap.Error(re.Error))
	})
}

func (e *Explorer) Shutdown(ctx context.Context) {
	e.server.Shutdown(ctx)
	e.metricsServer.Shutdown(ctx)
}
