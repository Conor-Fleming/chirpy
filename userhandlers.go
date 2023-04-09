package main

import (
	"encoding/json"
	"errors"
	"net/http"
)

func (cfg apiConfig) postUserHandler(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email string `json:"email"`
	}

	decoder := json.NewDecoder(r.Body)
	user := parameters{}
	err := decoder.Decode(&user)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, errors.New("error decoding email to user object"))
		return
	}

	result, err := cfg.dbClient.CreateUser(user.Email)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, errors.New("error creating user"))
		return
	}

	respondWithJSON(w, http.StatusCreated, result)
}

//implement function to normalize emails potentially error if not correct format
