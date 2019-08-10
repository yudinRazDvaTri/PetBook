package models

import (
	"time"

	"github.com/jmoiron/sqlx"
)

type Comment struct {
	CommendID   int       `json:"comment_id" db:"comment_id"`
	TopicID     int       `json:"topic_id" db:"topic_id"`
	UserID      int       `json:"user_id" db:"user_id"`
	CreatedTime time.Time `json:"created_time" db:"created_time"`
	Content     string    `json:"content" db:"content"`
}

type CommentStore struct {
	DB *sqlx.DB
}

