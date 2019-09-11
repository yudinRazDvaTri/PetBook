package search

import (
	"github.com/dpgolang/PetBook/pkg/models/forum"
	"github.com/jmoiron/sqlx"
)

type SearchStorer interface {
	GetAllPets(userID int) ([]*DispPet, error)
	GetByUser(userID int, email string) *DispPet
	GetFilterPets(m map[string]interface{}) ([]*DispPet, error)
	GetTopicsBySearch(search string) ([]forum.Topic, error)
}

type SearchStore struct {
	DB *sqlx.DB
}
