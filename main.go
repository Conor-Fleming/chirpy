package main

import (
	"log"
	"net/http"

	"github.com/Conor-Fleming/chirpy/database"
	"github.com/go-chi/chi"
)

func main() {
	apiCfg := apiConfig{}
	dbClient, err := database.NewDB("database.json")
	if err != nil {
		log.Fatal(err)
	}
	router := chi.NewRouter()
	apiRouter := chi.NewRouter()
	adminRouter := chi.NewRouter()

	apiRouter.Get("/healthz", healthzHandler)
	adminRouter.Get("/metrics", apiCfg.hitzHandler)

	apiRouter.Post("/chirps", validateHandler)

	router.Mount("/", apiCfg.middlewareMetrics(http.FileServer(http.Dir("."))))
	router.Mount("/api", apiRouter)
	router.Mount("/admin", adminRouter)

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

		next.ServeHTTP(w, r)
	})
}
