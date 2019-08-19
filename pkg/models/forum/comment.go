package forum

import (
	"fmt"
	"github.com/lib/pq"
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
	Rating int
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

func (f *ForumStore) RateComment(commentID, userID int) (bool, error) {
	rateOk := true

	_, err := f.DB.Exec(
		`insert into ratings (comment_id, user_id) values ($1, $2);`,
		commentID, userID)
	if err != nil {
		if pgerr, ok := err.(*pq.Error); ok {
			if pgerr.Code == "23505" {
				err = fmt.Errorf("User %d already liked comment %d: %v",
					userID, commentID, err)
			}
		}

		err = fmt.Errorf("Error while trying to rate %d comment by %d user in DB: %v",
			userID, commentID, err)
	}
	return rateOk, err
}

func (f *ForumStore) GetCommentRating(commentID int) (rating int, err error) {
	err = f.DB.QueryRowx(
		`SELECT COUNT(*) FROM ratings WHERE comment_id = $1;`, commentID).Scan(&rating)
	if err != nil {
		err = fmt.Errorf("Can't count rating of comment with %d id from DB: %v.\n",
			commentID, err)
	}
	return
}