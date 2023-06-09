package main

import (
	"log"
	"net/http"
	"os"

	"github.com/Conor-Fleming/chirpy/database"
	"github.com/go-chi/chi"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	dbClient, err := database.NewDB("database.json")
	if err != nil {
		log.Fatal()
	}
	apiCfg := apiConfig{
		dbClient:  *dbClient,
		jwtSecret: os.Getenv("JWT_SECRET"),
	}
	router := chi.NewRouter()
	apiRouter := chi.NewRouter()
	adminRouter := chi.NewRouter()

	apiRouter.Get("/healthz", healthzHandler)
	adminRouter.Get("/metrics", apiCfg.hitzHandler)

	apiRouter.Put("/users", apiCfg.updateUserHandler)
	apiRouter.Post("/users", apiCfg.postUserHandler)
	apiRouter.Post("/login", apiCfg.userLoginHandler)

	apiRouter.Post("/chirps", apiCfg.postChirpHandler)
	apiRouter.Get("/chirps", apiCfg.getChirpsHandler)
	apiRouter.Get("/chirps/{chirpID}", apiCfg.GetchirpByID)

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
