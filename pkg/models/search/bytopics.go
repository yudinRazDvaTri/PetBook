package search

import (
	"fmt"

	"github.com/dpgolang/PetBook/pkg/models/forum"
	"github.com/jmoiron/sqlx"
)

//Get filter topics by topics description
func (f *SearchStore) GetTopicsBySearch(search string) (topics []forum.Topic, err error) {

	rows, err := f.DB.Query("select * from topics where description ilike '%' || $1 || '%'", search)
	if err != nil {
		err = fmt.Errorf("Can't read topics-rows from db: %v", err)
		return
	}
	defer rows.Close()
	err = sqlx.StructScan(rows, &topics)
	if err != nil {
		err = fmt.Errorf("Can't scan topics-rows from db: %v", err)
	}
	return
}
