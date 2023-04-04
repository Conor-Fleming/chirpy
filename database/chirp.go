package database

import "errors"

func (db *DB) CreateChirp(body string) (Chirp, error) {
	//read DB to get map of "Chirps"
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
	//chir
	return nil, nil
}
