package forum

import (
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
)

type Topic struct {
	TopicID     int       `json:"topic_id" db:"topic_id"`
	UserID      int       `json:"user_id" db:"user_id"`
	CreatedTime time.Time `json:"created_time" db:"created_time"`
	Title       string    `json:"title" db:"title"`
	Description string    `json:"description" db:"description"`
}

func (f *ForumStore) GetAllTopics() (topics []*Topic, err error) {
	rows, err := f.DB.Query("select * from topics order by created_time DESC")
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

func (f *ForumStore) CreateNewTopic(topic *Topic) (err error) {
	_, err = f.DB.Exec(
		`insert into topics (user_id, title, description) values ($1, $2, $3)`,
		topic.UserID, topic.Title, topic.Description)
	if err != nil {
		return fmt.Errorf("cannot affect rows in topics table of db: %v", err)
	}
	return
}