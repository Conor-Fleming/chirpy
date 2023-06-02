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

	token, err := cfg.createJWT(user)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, errors.New("could not create token"))
	}

	result.Token = token

	respondWithJSON(w, http.StatusOK, result)
}

func (cfg apiConfig) updateUserHandler(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	user := parameters{}
	err := decoder.Decode(&user)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, errors.New("error decoding email to user object"))
		return
	}

	id, err := cfg.authorizeJWT(w, r)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, errors.New("could not authorize Token"))
	}

	updated, err := cfg.dbClient.UpdateUser(user.Email, user.Password, id)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, errors.New("could not update user"))
		return
	}

	respondWithJSON(w, http.StatusOK, updated)
}

func (cfg apiConfig) authorizeJWT(w http.ResponseWriter, r *http.Request) (string, error) {
	tokenString := r.Header.Get("Authorization")
	tokenString = strings.TrimPrefix(tokenString, "Bearer ")
	token, err := jwt.ParseWithClaims(tokenString, jwt.RegisteredClaims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte("test"), nil
	})
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, errors.New("could not parse jwt"))
		return "", err
	}

	expiration, err := token.Claims.GetExpirationTime()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, errors.New("could not get expiration time"))
		return "", err
	}

	if !token.Valid || expiration.Before(time.Now().UTC()) {
		respondWithError(w, http.StatusUnauthorized, errors.New("token has expired"))
	}
	userID, err := token.Claims.GetSubject()

	return userID, nil
}

func (cfg apiConfig) createJWT(user parameters) (string, error) {
	if user.Token_time > 24 || user.Token_time == 0 {
		user.Token_time = 24
	}

	expiration := jwt.NewNumericDate(time.Now().Add(time.Duration(user.Token_time)))
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer:    "chirpy",
		IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		ExpiresAt: expiration,
	})

	signedToken, err := token.SignedString(cfg.jwtSecret)
	if err != nil {
		return "", err
	}

	return signedToken, nil
}
