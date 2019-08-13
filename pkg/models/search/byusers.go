package search

import (
	"fmt"
	"github.com/jmoiron/sqlx"
)

func (f *SearchStore) GetByUser(email string) (*DispPet,error) {
	var pet DispPet
	err := f.DB.QueryRowx(`select name, description from pets p join users u on p.user_id=u.id where u.email=$1 `,email).StructScan(&pet)
	if err != nil {
		return &pet, err
	}
	return &pet, nil
}

func (f *SearchStore) GetAllPets() (pets []*DispPet, err error) {
	rows, err := f.DB.Query("select name, description from pets order by name")
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

