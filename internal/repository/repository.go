package repository

import "database/sql"

type UserRepository interface {
}


type Repository struct {
	UserRepository
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{}
}
