package database

import "errors"

func (db *DB) CreateUser(email string) (User, error) {
	db.mux.Lock()
	defer db.mux.Unlock()

	userData, err := db.readDB()
	if err != nil {
		return User{}, err
	}

	userid := len(userData.Users) + 1
	if _, ok := userData.Users[userid]; !ok {
		user := User{
			ID:    userid,
			Email: email,
		}

		userData.Users[user.ID] = user
		db.wrtiteDB(userData)
		return user, nil
	}

	return User{}, errors.New("could not create user")
}
