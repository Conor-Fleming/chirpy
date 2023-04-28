package database

import (
	"encoding/json"
	"errors"
	"os"
	"sync"
)

type DB struct {
	path string
	mux  *sync.RWMutex
}

type DBSchema struct {
	Chirps map[int]Chirp   `json:"Chirps"`
	Users  map[string]User `json:"Users"`
}

type Chirp struct {
	ID   int    `json:"id"`
	Body string `json:"body"`
}

type User struct {
	ID           int    `json:"id"`
	Email        string `json:"email"`
	PasswordHash string `json:"password"`
	Token        string `json:"token"`
}

func NewDB(filepath string) (*DB, error) {
	db := DB{
		path: filepath,
		mux:  &sync.RWMutex{},
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
		Users:  make(map[string]User),
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
