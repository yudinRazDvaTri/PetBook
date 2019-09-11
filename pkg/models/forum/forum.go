package forum

import "github.com/jmoiron/sqlx"

type ForumStorer interface {
	//Topic Methods
	GetAllTopics() (topics []Topic, err error)
	CreateNewTopic(userID int, title, description string) (err error)
	GetTopicByID(topicID int) (topic Topic, err error)
	NewViewTopic(userName string, topic Topic) (viewTopic ViewTopic, err error)
	//Comment Methods
	GetTopicComments(topicID int) (comments []Comment, err error)
	AddNewComment(topicID, userID, parentID int, content string) (err error)
	RateComment(commentID, userID int) (bool, error)
	NewViewComment(userName string, comment Comment) (viewComment ViewComment, err error)
}

type ForumStore struct {
	DB *sqlx.DB
}
