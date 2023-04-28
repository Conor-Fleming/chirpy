package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type parameters struct {
	Password   string `json:"password"`
	Email      string `json:"email"`
	Token_time int    `json:"expires_in_seconds"`
}

func (cfg apiConfig) postUserHandler(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	user := parameters{}
	err := decoder.Decode(&user)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, errors.New("error decoding email to user object"))
		return
	}

	result, err := cfg.dbClient.CreateUser(user.Email, user.Password, user.Token_time)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, errors.New("error creating user"))
		return
	}

	respondWithJSON(w, http.StatusCreated, result)
}

func (cfg apiConfig) userLoginHandler(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	user := parameters{}
	err := decoder.Decode(&user)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, errors.New("error decoding email to user object"))
		return
	}

	result, err := cfg.dbClient.UserLogin(user.Email, user.Password)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, errors.New("error loging in"))
		return
	}

	token := cfg.createJWT(user)

	respondWithJSON(w, http.StatusOK, result)
}

func (cfg apiConfig) updateUserHandler(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("Authorization")
	token = strings.TrimPrefix(token, "Bearer")

}

func (cfg apiConfig) createJWT(user parameters) *jwt.Token {
	if user.Token_time > 24 || user.Token_time == 0 {
		user.Token_time = 24
	}

	now := time.Now().UTC()
	expiration := time.Now().Add(time.Duration(user.Token_time))
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer:    "chirpy",
		IssuedAt:  now,
		ExpiresAt: expiration,
	})

	return token
}
