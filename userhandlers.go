package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
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
		Password   string `json:"password"`
		Email      string `json:"email"`
		Token_time int    `json:"expires_in_seconds"`
	}

	decoder := json.NewDecoder(r.Body)
	user := parameters{}
	err := decoder.Decode(&user)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, errors.New("error decoding email to user object"))
		return
	}

	cfg.createJWT(user.Token_time)

	result, err := cfg.dbClient.UserLogin(user.Email, user.Password)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, errors.New("error loging in"))
		return
	}

	respondWithJSON(w, http.StatusOK, result)
}

func (cfg apiConfig) createJWT(tokenTime int) {
	if tokenTime > 24 || tokenTime == 0 {
		tokenTime = 24
	}

	now := time.Now().UTC()
	expiration := time.Now().Add(time.Duration(tokenTime))
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer:    "chirpy",
		IssuedAt:  now,
		ExpiresAt: expiration,
		Subject:   string(user.ID),
	})

	token.SignedString(cfg.jwtSecret)
}
