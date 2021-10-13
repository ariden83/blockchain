package endpoint

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"net/http"
	"net/http/pprof"
	"strconv"
	"time"
)

func (e *EndPoint) ListenMetrics(stop chan error) {
	mux := e.makeHealthzRouter()

	e.metricsServer = &http.Server{
		Addr:           ":" + strconv.Itoa(e.config.Metrics.Port),
		Handler:        mux,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		IdleTimeout:    120 * time.Second,
		MaxHeaderBytes: 1 << 12,
	}
	go func() {
		log.Println("Metrics Server Listening on port :", strconv.Itoa(e.config.Metrics.Port))
		if err := e.metricsServer.ListenAndServe(); err != nil {
			stop <- fmt.Errorf("cannot start healthz server %s", err)
		}
	}()
}

func (e *EndPoint) makeHealthzRouter() http.Handler {
	muxRouter := mux.NewRouter()
	muxRouter.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		message := "The service responds correctly"
		res := Healthz{Result: true, Messages: []string{message}, Version: e.config.Version}
		js, err := json.Marshal(res)
		if err != nil {
			log.Fatalf(fmt.Sprintf("Fail to jsonify %s", err))
		}
		if _, err := w.Write(js); err != nil {
			log.Fatalf(fmt.Sprintf("Fail to Write response in http.ResponseWriter %s", err))
			return
		}
	})

	muxRouter.HandleFunc("/readiness", func(w http.ResponseWriter, r *http.Request) {
		result := true
		message := "Service responds correctly"

		res := Healthz{Result: result, Messages: []string{message}, Version: e.config.Version}
		js, err := json.Marshal(res)
		if err != nil {
			log.Fatalf(fmt.Sprintf("Fail to jsonify %s", err))
		}
		if _, err := w.Write(js); err != nil {
			log.Fatalf(fmt.Sprintf("Fail to Write response in http.ResponseWriter %s", err))
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
