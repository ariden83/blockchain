package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/ariden83/blockchain/internal/endpoint"
	"github.com/ariden83/blockchain/internal/persistence"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

type server struct {
	endpoint *endpoint.EndPoint
}

func main() {
	defer os.Exit(0)

	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	pers := persistence.Init(os.Getenv("DB_PATH"))

	e := endpoint.Init(pers)
	defer e.Close()
	go func() {
		e.Genesis()
	}()

	s := server{
		endpoint: e,
	}

	log.Fatal(s.run())
}

// web server
func (s *server) run() error {
	mux := s.makeMuxRouter()
	httpPort := os.Getenv("PORT")
	log.Println("HTTP Server Listening on port :", httpPort)
	server := &http.Server{
		Addr:           ":" + httpPort,
		Handler:        mux,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	if err := server.ListenAndServe(); err != nil {
		return err
	}

	return nil
}

// create handlers
func (s *server) makeMuxRouter() http.Handler {
	muxRouter := mux.NewRouter()
	muxRouter.HandleFunc("/", s.handleGetBlockChain).Methods("GET")
	muxRouter.HandleFunc("/", s.handleWriteBlock).Methods("POST")
	return muxRouter
}

// write blockchain when we receive an http request
func (s *server) handleGetBlockChain(w http.ResponseWriter, r *http.Request) {
	s.endpoint.PrintBlockChain(w)
}

// takes JSON payload as an input for heart rate (BPM)
func (s *server) handleWriteBlock(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var m endpoint.Message

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&m); err != nil {
		respondWithJSON(w, r, http.StatusBadRequest, r.Body)
		return
	}
	defer r.Body.Close()

	newBlock := s.endpoint.GenerateBlock(m)

	respondWithJSON(w, r, http.StatusCreated, newBlock)

}

func respondWithJSON(w http.ResponseWriter, r *http.Request, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	response, err := json.MarshalIndent(payload, "", "  ")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("HTTP 500: Internal Server Error"))
		return
	}
	w.WriteHeader(code)
	w.Write(response)
}
