package main

import (
	"fmt"
	"io"
	"net/http"
)

type apiConfig struct {
	fileServes int
}

func (cfg *apiConfig) middlewareMetrics(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileServes++
		next.ServeHTTP(w, r)
	})
}

func (cfg *apiConfig) hitzHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	hitsString := fmt.Sprintf("Hits: %d", cfg.fileServes)
	io.WriteString(w, hitsString)
}
