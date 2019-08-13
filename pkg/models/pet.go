package models

import (
	//"github.com/dpgolang/PetBook/pkg/logger"
	"fmt"
	"github.com/jmoiron/sqlx"
)

type Pet struct {
	ID          int    `db:"user_id"`
	Name        string `db:"name"`
	PetType     string `db:"animal_type"`
	Breed       string `db:"breed"`
	Age         string `db:"age"`
	Weight      string `db:"weight"`
	Gender      string `db:"gender"`
	Description string `db:"description"`
}

type PetStore struct {
	DB *sqlx.DB
}

type PetStorer interface {
	RegisterPet(pet *Pet) error
}

// TODO: rewrite to update into
func (c *PetStore) RegisterPet(pet *Pet) error {
	_, err := c.DB.Exec("insert into pets (user_id, name, animal_type, breed, age, weight, gender, description) values ($1, $2, $3, $4, $5, $6, $7, $8)",
		pet.ID, pet.Name, pet.PetType, pet.Breed, pet.Age, pet.Weight, pet.Gender, pet.Description)
	if err != nil {
		return fmt.Errorf("cannot affect rows in pets in db: %v", err)
	}
	return nil
}
