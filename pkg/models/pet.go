package models

import (
	"fmt"
	"github.com/jmoiron/sqlx"
)

type Pet struct {
	ID          int    `db:"user_id"`
	Name        string `db:"name"`
	PetType     string `db:"animal_type"`
	Breed       string `db:"breed"`
	Age         string `db:"age"`
	Weight      string `db: "weight"`
	Gender      string `db:"gender"`
	Description string `db:"description"`
}

type PetStore struct {
	DB *sqlx.DB
}

type PetStorer interface {
	GetPet(pet *Pet) error
	RegisterPet(pet *Pet) error
}

func (c *PetStore) GetPet(pet *Pet) error {
	err := c.DB.QueryRowx("select * from pets where user_id=$1", pet.ID).StructScan(pet)
	if err != nil {
		logErr(err)
		return fmt.Errorf("cannot scan pet from db: %v", err)
	}
	return nil
}

func (c *PetStore) RegisterPet(pet *Pet) error {
	_, err := c.DB.Exec("insert into pets (user_id, name, animal_type,breed,age,weight, gender) values ($1, $2, $3, $4, $5, $6, $7)",
		pet.ID, pet.Name, pet.PetType, pet.Breed, pet.Age, pet.Weight, pet.Gender)
	if err != nil {
		return fmt.Errorf("cannot affect rows in pets in db: %v", err)
	}
	return nil
}
