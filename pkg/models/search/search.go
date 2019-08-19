package search

import "github.com/jmoiron/sqlx"

type SearchStorer interface {

	GetAllPets()([]*DispPet, error)
	GetByUser(email string)(*DispPet,error)
	GetFilterPets(m map[string]string)([]*DispPet,error)

}

type SearchStore struct {
	DB *sqlx.DB
}
