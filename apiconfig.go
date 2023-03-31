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
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	hitsString := fmt.Sprintf(`<html>

	<body>
		<h1>Welcome, Chirpy Admin</h1>
		<p>Chirpy has been visited %d times!</p>
	</body>
	
	</html>`, cfg.fileServes)
	io.WriteString(w, hitsString)
}
