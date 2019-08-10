package models

import (
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
)

type Topic struct {
	TopicID     int       `json:"topic_id" db:"topic_id"`
	UserID      int       `json:"user_id" db:"user_id"`
	AnswerID    int       `json:"asnwer_id" db:"answer_id"`
	CreatedTime time.Time `json:"created_time" db:"created_time"`
	Title       string    `json:"title" db:"title"`
	Description string    `json:"description" db:"description"`
}

type TopicStore struct {
	DB *sqlx.DB
}

func (t *TopicStore) CreateTopic(topic *Topic) (err error) {
	_, err = t.DB.Exec(
		`insert into topics (user_id, title, description) values ($1, $2, $3)`,
		topic.UserID, topic.Title, topic.Description)
	if err != nil {
		return fmt.Errorf("cannot affect rows in pets in db: %v", err)
	}
	return
}
