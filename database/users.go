package database

import (
	"errors"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func (db *DB) CreateUser(email, password string) (User, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return User{}, fmt.Errorf("Create User hash: %v", err)
	}
	hashedPass := string(hashed)

	db.mux.Lock()
	defer db.mux.Unlock()

	userData, err := db.readDB()
	if err != nil {
		return User{}, err
	}

	userid := len(userData.Users) + 1
	if _, ok := userData.Users[userid]; !ok {
		user := User{
			ID:           userid,
			Email:        email,
			PasswordHash: hashedPass,
		}

		userData.Users[user.ID] = user
		db.wrtiteDB(userData)
		return user, nil
	}

	return User{}, errors.New("could not create user")
}

func (db *DB) UserLogin(email, password string) (User, error) {
	db.mux.Lock()
	defer db.mux.Unlock()

	userData, err := db.readDB()
	if err != nil {
		return User{}, err
	}

	//err = bcrypt.CompareHashAndPassword([]byte(userData.Users), )

	return User{}, errors.New("could not authenticate")
}
