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
	chirpid := len(chirpData.Chirps) + 1
	if _, ok := chirpData.Chirps[chirpid]; !ok {
		chirp := Chirp{
			ID:   chirpid,
			Body: body,
		}
		//write updated map to db and return new Chirp
		chirpData.Chirps[chirpid] = chirp
		db.wrtiteDB(chirpData)
		return chirp, nil
	}

	return Chirp{}, errors.New("error occured when creating chirp")
}

func (db *DB) GetChirp(id int) (Chirp, error) {
	db.mux.Lock()
	defer db.mux.Unlock()
	chirps, err := db.readDB()
	if err != nil {
		return Chirp{}, errors.New("error reading DB")
	}

	if body, ok := chirps.Chirps[id]; ok {
		return body, nil
	}

	return Chirp{}, errors.New("The given ID doesnt exist")
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
