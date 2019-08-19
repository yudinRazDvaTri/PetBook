package forum

import "github.com/jmoiron/sqlx"

type ForumStorer interface {
	//Topic Methods
	GetAllTopics() (topics []Topic, err error)
	CreateNewTopic(userID int, title, description string) (err error)
	GetTopicByID(topicID int) (topic Topic, err error)
	//Comment Methods
	AddNewComment(topicID, userID int, content string) (err error)
	GetTopicComments(topicID int) (comments []Comment, err error)
}

type ForumStore struct {
	DB *sqlx.DB
}
