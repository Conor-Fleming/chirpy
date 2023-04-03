package database

import (
	"sync"
)

type DB struct {
	path string
	mux  sync.Mutex
}

func NewDB(filepath string) (*DB, error) {
	db := DB{
		path: filepath,
	}

	return &db, nil
}
