package database

import (
	"errors"
	"fmt"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

func (db *DB) CreateUser(email, password string, token_time int) (User, error) {
	email = strings.ToLower(email)
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return User{}, fmt.Errorf("Create User hash: %v", err)
	}

	db.mux.Lock()
	defer db.mux.Unlock()

	userData, err := db.readDB()
	if err != nil {
		return User{}, err
	}

	userid := len(userData.Users) + 1
	if _, ok := userData.Users[email]; !ok {
		user := User{
			ID:           userid,
			Email:        email,
			PasswordHash: string(hashed),
		}
		authenticated := User{
			ID:    userid,
			Email: email,
		}

		userData.Users[email] = user
		db.wrtiteDB(userData)
		return authenticated, nil
	}

	return User{}, errors.New("User already exists with that email")
}

func (db *DB) UserLogin(email, password string) (User, error) {
	email = strings.ToLower(email)
	db.mux.Lock()
	defer db.mux.Unlock()

	userData, err := db.readDB()
	if err != nil {
		return User{}, err
	}

	if _, ok := userData.Users[email]; !ok {
		return User{}, errors.New("user does not exist")
	}

	//get stored hash and compare to given
	user := userData.Users[email]
	hash := userData.Users[email].PasswordHash
	err = bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		return User{}, errors.New("incorrect password")
	}
	authenticated := User{
		ID:    user.ID,
		Email: user.Email,
	}
	return authenticated, nil
}
