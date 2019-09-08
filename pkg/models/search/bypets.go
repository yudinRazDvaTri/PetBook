package search

import (
	"fmt"
	"strings"

	"github.com/dpgolang/PetBook/pkg/models"
	"github.com/jmoiron/sqlx"
)

type DispPet struct {
	Name        string `json:"name" db:"name"'`
	Description string `json:"description" db:"description"'`
	UserID      int    `json:"user_id" db:"user_id"`
}

//Function that is called in the template. It returns a boolean value,
//hether the user can subscribe to this user
func (f *DispPet) CanFollow(userID int, petsFollowing []*models.FollowerPets) bool {
	if f.UserID == userID {
		return false
	}
	for _, val := range petsFollowing {
		if int(val.UserID) == f.UserID {
			return false
		}
	}
	return true
}

//Get filter pets by different parameters from db
func (f *SearchStore) GetFilterPets(m map[string]interface{}) (pets []*DispPet, err error) {
	var where []string
	var values []interface{}
	var count = 1
	for _, k := range []string{"age", "animal_type", "breed", "weight", "gender", "name"} {
		if v, ok := m[k]; ok {
			values = append(values, v)
			where = append(where, fmt.Sprintf("%s = $%d", k, count))
			count++
		}
	}
	rows, err := f.DB.Query("SELECT name, description, user_id FROM pets WHERE "+strings.Join(where, " AND "), values...)
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
