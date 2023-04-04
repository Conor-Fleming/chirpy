package database

import (
	"encoding/json"
	"errors"
	"os"
	"sync"
)

var idCounter int = 0

type DB struct {
	path string
	mux  *sync.RWMutex
}

type DBSchema struct {
	Chirps map[int]Chirp `json:"Chirps"`
}

type Chirp struct {
	ID   int    `json:"id"`
	Body string `json:"body"`
}

func NewDB(filepath string) (*DB, error) {
	db := DB{
		path: filepath,
	}

	err := db.ensureDB()
	if err != nil {
		return nil, err
	}

	return &db, nil
}

func (db *DB) ensureDB() error {
	_, err := os.ReadFile(db.path)
	if err != nil {
		dbErr := db.createDB()
		return dbErr
	}

	return nil
}

func (db *DB) createDB() error {
	data, err := json.Marshal(DBSchema{
		Chirps: make(map[int]Chirp),
	})
	if err != nil {
		return err
	}
	err = os.WriteFile(db.path, data, 0644)
	if err != nil {
		return err
	}

	return nil
}

func (db *DB) readDB() (DBSchema, error) {
	var schema DBSchema
	data, err := os.ReadFile(db.path)
	if err != nil {
		return schema, err
	}
	err = json.Unmarshal(data, &schema)
	if err != nil {
		return schema, errors.New("Could not unmarshal Data from DB")
	}

	return schema, nil
}

func (db *DB) wrtiteDB(data DBSchema) error {
	json, err := json.Marshal(data)
	if err != nil {
		return err
	}

	err = os.WriteFile(db.path, json, 0644)
	if err != nil {
		return err
	}

	return nil
}
