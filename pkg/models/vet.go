package models

import (
	//"github.com/dpgolang/PetBook/pkg/logger"
	"fmt"
	"github.com/jmoiron/sqlx"
)

type Vet struct {
	ID          int    `db:"user_id"`
	Name        string `db:"name"`
	Qualification     string `db:"qualification"`
	Surname       string `db:"surname"`
	Category         string `db:"category"`
	Certificates string `db:"certificates"`
}

type VetStore struct {
	DB *sqlx.DB
}

type VetStorer interface {
	RegisterVet(vet *Vet) error
	//DisplayName(userID int) (name string, err error)
	GetVetEnums() ([]string, error)
	UpdateVet(vet *Vet)error
}

func (c *VetStore) RegisterVet(vet *Vet) error {
	_, err := c.DB.Exec("insert into vets (user_id, name, qualification, surname, category, certificates) values ($1, $2, $3, $4, $5, $6);",
		vet.ID, vet.Name, vet.Qualification, vet.Surname, vet.Category, vet.Certificates)
	if err != nil {
		return fmt.Errorf("cannot affect rows in pets in db: %v", err)
	}
	return nil
}

//func (c *VetStore) DisplayVetName(userID int) (name string, err error) {
//	var email, vetName string
//	err = c.DB.QueryRow(
//		`SELECT email FROM users WHERE id = $1`, userID).Scan(&email)
//	if err != nil {
//		return "", fmt.Errorf("Error occurred while trying read email of user with $d id: %v.\n", userID, err)
//	}
//	err = c.DB.QueryRow(
//		`SELECT name FROM pets WHERE user_id = $1`, userID).Scan(&vetName)
//	if err != nil {
//		return "", fmt.Errorf("Error occurred while trying read petName of user with $d id: %v.\n", userID, err)
//	}
//	name = email + "/" + vetName
//
//	return
//}

func (p *VetStore) UpdateVet(vet *Vet) error{
	_, err := p.DB.Exec("update vets set name=$1, qualification=$2,surname=$3, category=$4,certificates=$5 where user_id = $6",
		vet.Name, vet.Qualification, vet.Surname, vet.Category, vet.Certificates, vet.ID)
	if err != nil {
		return fmt.Errorf("cannot affect rows in vets in db: %v", err)
	}
	return nil
}
func (p *VetStore) GetVetEnums() ([]string, error) {
	var vetType []string
	var vtype string
	rows, err := p.DB.Queryx("SELECT unnest(enum_range(NULL::class))::text")
	if err != nil {
		fmt.Println("Error in getting enums")
	}
	for rows.Next() {
		err = rows.Scan(&vtype)
		if err != nil {
			fmt.Errorf("cannot affect rows in pets in db: %v", err)
		}
		vetType = append(vetType, vtype)
	}
	return vetType,nil
}
