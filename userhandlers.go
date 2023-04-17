package main

import (
	"encoding/json"
	"errors"
	"net/http"
)

func (cfg apiConfig) postUserHandler(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Password string `json:"password"`
		Email    string `json:"email"`
	}

	decoder := json.NewDecoder(r.Body)
	user := parameters{}
	err := decoder.Decode(&user)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, errors.New("error decoding email to user object"))
		return
	}

	result, err := cfg.dbClient.CreateUser(user.Email, user.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, errors.New("error creating user"))
		return
	}

	respondWithJSON(w, http.StatusCreated, result)
}

func (cfg apiConfig) userLoginHandler(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Password string `json:"password"`
		Email    string `json:"email"`
	}

	decoder := json.NewDecoder(r.Body)
	user := parameters{}
	err := decoder.Decode(&user)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, errors.New("error decoding email to user object"))
		return
	}

	result, err := cfg.dbClient.UserLogin(user.Email, user.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, errors.New("error loging in"))
		return
	}

	respondWithJSON(w, http.StatusOK, result)
}
