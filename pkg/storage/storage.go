package storage

import "database/sql"

type Storage struct {
}

func NewPostgreStorage(db *sql.DB) *Storage{
	return &Storage{}
}