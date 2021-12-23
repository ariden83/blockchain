package explorer

import (
	"encoding/json"
	"fmt"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
	"net/http"
	"net/http/pprof"
	"time"
)

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
