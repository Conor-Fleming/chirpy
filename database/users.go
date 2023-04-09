package database

import "errors"

func (db *DB) CreateUser(email string) (User, error) {
	db.mux.Lock()
	defer db.mux.Unlock()

	//how to deal with ID assignment

	userData, _ := db.readDB()
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
