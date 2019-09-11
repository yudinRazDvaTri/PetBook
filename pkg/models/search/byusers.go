package search

import (
	"fmt"

	"github.com/jmoiron/sqlx"
)

//Get sort animals by user
func (f *SearchStore) GetByUser(userID int, email string) *DispPet {
	var pet DispPet
	err := f.DB.QueryRowx(`select name, description, user_id from pets p join users u on p.user_id=u.id where u.email=$1 and user_id!=$2; `, email, userID).StructScan(&pet)
	if err != nil {
		return &pet
	}
	return &pet
}

//Get all pets from data base
func (f *SearchStore) GetAllPets(userID int) (pets []*DispPet, err error) {
	rows, err := f.DB.Query("select name, description,user_id from pets where user_id!=$1 order by name", userID)
	if err != nil {
		err = fmt.Errorf("Can't read pets from db: %v", err)
		return
	}
	defer rows.Close()
	err = sqlx.StructScan(rows, &pets)
	if err != nil {
		err = fmt.Errorf("Can't scan topics-rows from db: %v", err)
	}
	return
}
