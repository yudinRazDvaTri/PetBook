package models

import "time"

type Comment struct {
	ID        int       `json:"id" db:"id"`
	TopicID   int       `json:"topic_id" db:"topic_id"`
	CreatorID int       `json:"creator_id" db:"creator_id"`
	Created   time.Time `json:"created" db:"created"`
	Edited    time.Time `json:"edited" db:"edited"`
	Content   string    `json:"content" db:"content"`
}
