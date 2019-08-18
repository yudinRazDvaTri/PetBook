package forum

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"time"
)

type Comment struct {
	CommentID   int       `json:"comment_id" db:"comment_id"`
	TopicID     int       `json:"topic_id" db:"topic_id"`
	UserID      int       `json:"user_id" db:"user_id"`
	CreatedTime time.Time `json:"created_time" db:"created_time"`
	Content     string    `json:"content" db:"content"`
}

type ViewComment struct {
	UserName string
	Comment
}

func (f *ForumStore) AddNewComment(topicID, userID int, content string) (err error) {
	_, err = f.DB.Exec(
		`insert into comments (topic_id, user_id, content) values ($1, $2, $3);`,
		topicID, userID, content)
	if err != nil {
		err = fmt.Errorf("cannot affect rows in comments table of db: %v", err)
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