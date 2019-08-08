package models

import "time"

type Topic struct {
	ID          int       `json:"id" db:"id"`
	CreatorID   int       `json:"creator_id" db:"creator_id"`
	Created     time.Time `json:"created" db:"created"`
	Title       string    `json:"title" db:"title"`
	Description string    `json:"description" db:"description"`
}
