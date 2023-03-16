package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
)

type apiConfig struct {
	fileServes int
}

func (cfg *apiConfig) middlewareMetrics(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileServes++
	})
}

func main() {
	apiCfg := apiConfig{}
	mux := http.NewServeMux()
	mux.Handle("/", apiCfg.middlewareMetrics(http.FileServer(http.Dir("."))))
	mux.HandleFunc("/healthz", healthzHandler)
	mux.HandleFunc("/metrics", apiCfg.hitzHandler)
	corsMux := corsMiddleware(mux)

	server := &http.Server{
		Addr:    ":8080",
		Handler: corsMux,
	}
	log.Fatal(server.ListenAndServe())
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "*")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func healthzHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	io.WriteString(w, "OK")
}

func (cfg *apiConfig) hitzHandler(w http.ResponseWriter, r *http.Request) {
	hitsString := fmt.Sprintf("Hits: %d", cfg.fileServes)
	io.WriteString(w, hitsString)
}
