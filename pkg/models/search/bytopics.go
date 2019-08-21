package search

import (
	"fmt"
	"github.com/dpgolang/PetBook/pkg/models/forum"
	"github.com/jmoiron/sqlx"
)

func (f *SearchStore) GetTopicsBySearch(search string)(topics []forum.Topic,err error){
	query:=" '%"+search+"%'"+"order by created_time DESC;"
	rows, err := f.DB.Query("select * from topics where description ilike "+ query)
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
