package search

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"strings"
)

type DispPet struct{
	Name string `json:"name" db:"name"'`
	Description string `json:"description" db:"description"'`
}
func (f *SearchStore) GetFilterPets(m map[string]interface{}) (pets []*DispPet, err error) {
	var where []string
	var values []interface{}
	var count = 1
	for _, k := range []string{"age","animal_type","breed","weight","gender","name"}{
		if v, ok := m[k]; ok {
			values = append(values, v)
			where = append(where, fmt.Sprintf("%s = $%d", k, count))
			count++
		}
	}
	rows, err := f.DB.Query("SELECT name, description FROM pets WHERE " + strings.Join(where, " AND "), values... )
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

