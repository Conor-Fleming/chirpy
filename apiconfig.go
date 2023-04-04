package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/Conor-Fleming/chirpy/database"
)

type apiConfig struct {
	fileServes int
	dbClient   database.DB
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

func (cfg *apiConfig) getChirpsHandler(w http.ResponseWriter, r *http.Request) {

}

func (cfg apiConfig) postChirpHandler(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}

	decoder := json.NewDecoder(r.Body)
	chirp := parameters{}
	err := decoder.Decode(&chirp)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, errors.New("something went wrong"))
		return
	}
	//ensure chirp is under limit of 140 chars
	if len(chirp.Body) > 140 {
		respondWithError(w, 400, errors.New("Chirp is too long"))
		return
	}

	chirp.Body = cleanChirp(chirp.Body)

	//write chirp to DB and respond with ok
	cfg.dbClient.CreateChirp(chirp.Body)
	respondWithJSON(w, http.StatusOK, validBody{
		Cleaned_Body: chirp.Body,
	})
}

func cleanChirp(body string) string {
	words := strings.Split(body, " ")
	for i, v := range words {
		if strings.ToLower(v) == "kerfuffle" || strings.ToLower(v) == "sharbert" || strings.ToLower(v) == "fornax" {
			words[i] = "****"
		}
	}
	return strings.Join(words, " ")
}
