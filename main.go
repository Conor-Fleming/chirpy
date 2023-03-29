package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi"
)

func main() {
	apiCfg := apiConfig{}
	router := chi.NewRouter()
	apiRouter := chi.NewRouter()
	adminRouter := chi.NewRouter()
	router.Mount("/", apiCfg.middlewareMetrics(http.FileServer(http.Dir("."))))
	apiRouter.Mount("/api", apiCfg.middlewareMetrics(http.FileServer(http.Dir("."))))
	adminRouter.Mount("/admin", apiCfg.middlewareMetrics(http.FileServer(http.Dir("."))))
	apiRouter.Get("/healthz", healthzHandler)
	adminRouter.Get("/metrics", apiCfg.hitzHandler)
	corsMux := corsMiddleware(router)

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
		if r.Method != "GET" {
			w.WriteHeader(405)
			return
		}
		next.ServeHTTP(w, r)
	})
}
