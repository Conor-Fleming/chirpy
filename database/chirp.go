package database

import (
	"errors"
	"sort"
)

var idCounter int = 0

func (db *DB) CreateChirp(body string) (Chirp, error) {
	//read DB to get map of "Chirps"
	db.mux.Lock()
	defer db.mux.Unlock()
	chirpData, _ := db.readDB()

	//Create chirp obj with body and unique ID
	//update map with new chirp
	idCounter += 1
	if _, ok := chirpData.Chirps[idCounter]; !ok {
		chirp := Chirp{
			ID:   idCounter,
			Body: body,
		}
		//write updated map to db and return new Chirp
		chirpData.Chirps[idCounter] = chirp
		db.wrtiteDB(chirpData)
		return chirp, nil
	}

	return Chirp{}, errors.New("error occured when creating chirp")
}

func (db *DB) GetChirps() ([]Chirp, error) {
	db.mux.Lock()
	defer db.mux.Unlock()
	chirps, err := db.readDB()
	if err != nil {
		return nil, errors.New("error reading DB")
	}

	keys := make([]int, 0)
	for k := range chirps.Chirps {
		keys = append(keys, k)
	}
	sort.Slice(keys, func(i, j int) bool {
		return keys[i] < keys[j]
	})

	chirpSlice := make([]Chirp, 0)
	for _, v := range keys {
		chirpSlice = append(chirpSlice, chirps.Chirps[v])
	}

	return chirpSlice, nil
}
