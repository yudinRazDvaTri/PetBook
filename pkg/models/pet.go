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
	DisplayName(userID int, role string) (name string, err error)
	GetPetEnums() ([]string, error)
	UpdatePet(pet *Pet) error
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

func (c *PetStore) DisplayName(userID int, role string) (name string, err error) {
	var email, petName, vetName string
	err = c.DB.QueryRow(
		`SELECT email FROM users WHERE id = $1`, userID).Scan(&email)
	if err != nil {
		return "", fmt.Errorf("Error occurred while trying read email of user with %b id: %v.\n", userID,  err)
	}

	if role == "pet" {
		err = c.DB.QueryRow(
			`SELECT name FROM pets WHERE user_id = $1`, userID).Scan(&petName)
		if err != nil {
			return "", fmt.Errorf("Error occurred while trying read petName of user with %b id: %v.\n",userID, err)
		}
		name = email + "/" + petName
	} else if role == "vet" {
		err = c.DB.QueryRow(
			`SELECT name FROM vets WHERE user_id = $1`, userID).Scan(&vetName)
		if err != nil {
			return "", fmt.Errorf("Error occurred while trying read vetName of user with %b id: %v.\n", userID, err)
		}
		name = email + "/" + vetName
	}

	return
}

func (p *PetStore) UpdatePet(pet *Pet) error {
	_, err := p.DB.Exec(`INSERT into pets(user_id, age, name, animal_type, breed, weight, gender, description) 
								values ($1, $2, $3, $4, $5, $6, $7, $8)
								ON CONFLICT (user_id) DO UPDATE 
								SET age = $2,
								name = $3,
								animal_type = $4,
								breed = $5,
								weight = $6,
								gender = $7,
								description = $8`, pet.ID, pet.Age, pet.Name, pet.PetType, pet.Breed, pet.Weight, pet.Gender, pet.Description)

	if err != nil {
		return fmt.Errorf("Error occurred while trying to update pet table: %v.\n", err)
	}
	return nil


	return nil
}

func (p *PetStore) GetPetEnums() ([]string, error) {
	var petType []string
	var ptype string
	rows, err := p.DB.Queryx("SELECT unnest(enum_range(NULL::kind_of_animal))::text")
	if err != nil {
		return nil, fmt.Errorf("cannot connect to database: %v", err)
	}
	for rows.Next() {
		err = rows.Scan(&ptype)
		if err != nil {
			return nil, fmt.Errorf("cannot insert in db: %v", err)
		}
		petType = append(petType, ptype)
	}
	return petType, nil
}
