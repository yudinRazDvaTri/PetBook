package search

import (
	"github.com/dpgolang/PetBook/pkg/models/forum"
	"github.com/jmoiron/sqlx"
)

type SearchStorer interface {
	GetAllPets()([]*DispPet, error)
	GetByUser(email string)(*DispPet,error)
	GetFilterPets(m map[string]interface{})([]*DispPet,error)
	GetTopicsBySearch(search string)([]forum.Topic, error)
}

type SearchStore struct {
	DB *sqlx.DB
}
