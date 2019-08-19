package forum

import (
	"fmt"
	"github.com/dpgolang/PetBook/pkg/logger"
	"time"

	"github.com/jmoiron/sqlx"
)

type Topic struct {
	TopicID        int       `json:"topic_id" db:"topic_id"`
	UserID         int       `json:"user_id" db:"user_id"`
	CreatedTime    time.Time `json:"created_time" db:"created_time"`
	Title          string    `json:"title" db:"title"`
	Description    string    `json:"description" db:"description"`
	CommentsNumber int       `json:"comments_number" db:"comments_number"`
}

type ViewTopic struct {
	UserName string
	Topic
}

func (f *ForumStore) CreateNewTopic(userID int, title, description string) (err error) {
	_, err = f.DB.Exec(
		`insert into topics (user_id, title, description) values ($1, $2, $3)`,
		userID, title, description)
	if err != nil {
		err = fmt.Errorf("cannot affect rows in topics table of db: %v", err)
	}
	return
}

func (f *ForumStore) GetAllTopics() (topics []Topic, err error) {
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
	for i := range topics {
		err = f.SetNumberOfComments(topics[i].TopicID)
		if err != nil {
			logger.Error(err)
		}
	}
	return
}

func (f *ForumStore) GetTopicComments(topicID int) (comments []Comment, err error) {
	rows, err := f.DB.Query("select * from comments where topic_id = $1 order by created_time ASC", topicID)
	if err != nil {
		err = fmt.Errorf("Can't read comment-rows from db: %v", err)
		return
	}
	defer rows.Close()
	err = sqlx.StructScan(rows, &comments)
	if err != nil {
		err = fmt.Errorf("Can't scan comment-rows from db: %v", err)
	}
	return
}

func (f *ForumStore) SetNumberOfComments(topicID int) (err error) {
	_, err = f.DB.Exec(
		`UPDATE topics SET comments_number =
				(SELECT count(*) FROM comments WHERE topic_id = $1)
				WHERE topic_id = $1;`, topicID)
	if err != nil {
		err = fmt.Errorf("Can't number of comments in topic %d from DB: %v",
			topicID, err)
	}
	return
}

func (f *ForumStore) GetTopicByID(topicID int) (topic Topic, err error) {
	err = f.SetNumberOfComments(topicID)
	if err != nil {
		logger.Error(err)
	}
	err = f.DB.QueryRowx(
		`SELECT * FROM topics WHERE topic_id = $1`, topicID).StructScan(&topic)
	if err != nil {
		err = fmt.Errorf("Error occurred while trying read topic with $d id from DB: %v.\n", topicID, err)
	}
	return
}
