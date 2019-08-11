package forum

import "github.com/jmoiron/sqlx"

type ForumStorer interface {
	GetAllTopics() (topics []*Topic, err error)
	CreateNewTopic(topic *Topic) (err error)
}

type ForumStore struct {
	DB *sqlx.DB
}

