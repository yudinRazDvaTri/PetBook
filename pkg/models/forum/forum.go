package forum

import "github.com/jmoiron/sqlx"

type ForumStorer interface {
	GetAllTopics() (topics []*Topic, err error)
	CreateNewTopic(userID int, title, description string) (err error)
	AddNewComment(topicID, userID int, content string) (err error)
	GetTopicComments(topicID int) (comments []Comment, err error)
}

type ForumStore struct {
	DB *sqlx.DB
}
