package search

import (
	"fmt"
	"github.com/jmoiron/sqlx"
)

type DispPet struct{
	Name string `json:"name" db:"name"'`
	Description string `json:"description" db:"description"'`
}
func (f *SearchStore) GetFilterPets(m map[string]string) (pets []*DispPet, err error) {
	var (
		query string
		firstParam bool
	)
	parameters:=make([]string,5)
	if age, ok := m["age"]; ok {
		parameters[0]="age="+age
	}
	if animalType, ok := m["animal_type"]; ok {
		parameters[1]="animal_type='"+animalType+"'"
	}
	if breed, ok := m["breed"]; ok {
		parameters[2]="breed='"+breed+"'"
	}
	if weight, ok := m["weight"]; ok {
		parameters[3]="weight="+weight
	}
	if gender, ok := m["gender"]; ok {
		parameters[4]="gender='"+gender+"'"
	}
	if name,ok:=m["name"];ok{
		parameters[4]="name='"+name+"'"
	}
	for _, str := range parameters {
		if str!="" && !firstParam {
			query += str
			firstParam=true
		}else if str!="" && firstParam{
			query+=" and " + str
		}

	}
	query+=" order by name;"
	rows, err := f.DB.Query("select name, description from pets where " + query)
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

