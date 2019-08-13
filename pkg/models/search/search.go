package search

import "github.com/jmoiron/sqlx"

type SearchStorer interface {
	GetAllPets()([]*DispPet, error)
	GetByUser(email string)(*DispPet,error)
}

type SearchStore struct {
	DB *sqlx.DB
}
