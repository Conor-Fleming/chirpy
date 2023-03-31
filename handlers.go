package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
)

type errorBody struct {
	Error string `json:"error"`
}

type validBody struct {
	Cleaned_Body string `json:"cleaned_body"`
}

var profaneWords = []string{
	"kerfuffle",
	"sharbert",
	"fornax",
}

func healthzHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	io.WriteString(w, "OK")
}

func validateHandler(w http.ResponseWriter, r *http.Request) {
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
	if len(chirp.Body) > 140 {
		respondWithError(w, 400, errors.New("Chirp is too long"))
		return
	}

	chirp.Body = cleanChirp(chirp.Body)
	fmt.Println(chirp.Body)

	respondWithJSON(w, http.StatusOK, validBody{
		Cleaned_Body: chirp.Body,
	})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	if payload != nil {
		response, err := json.Marshal(payload)
		if err != nil {
			log.Println("error marshalling", err)
			w.WriteHeader(500)
			response, _ := json.Marshal(errorBody{
				Error: "error marshalling",
			})
			w.Write(response)
			return
		}
		w.WriteHeader(code)
		w.Write(response)
	}
}

func respondWithError(w http.ResponseWriter, code int, err error) {
	respondWithJSON(w, code, errorBody{
		Error: err.Error(),
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
