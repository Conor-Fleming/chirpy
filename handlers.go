package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
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
