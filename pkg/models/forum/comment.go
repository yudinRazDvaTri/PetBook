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
	Likes []int64
	Comment
}

type ByRating []ViewComment

// Methods to sort ViewComments by Rating
func (v ByRating) Len() int {
	return len(v)
}
func (v ByRating) Swap(i, j int) {
	v[i], v[j] = v[j], v[i]
}
func (v ByRating) Less(i, j int) bool {
	return len(v[i].Likes) < len(v[j].Likes)
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

func (f *ForumStore) GetCommentRatings(commentID int) (likes []int64, err error) {
	rows, err := f.DB.Query(`SELECT user_id FROM ratings WHERE comment_id = $1;`, commentID)
	if err != nil {
		err = fmt.Errorf("Can't read rating-rows from db: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var like int64
		err = rows.Scan(&like)
		if err != nil {
			err = fmt.Errorf("Can't scan rating-row from db: %v", err)
		}
		likes = append(likes, like)
	}

	return
}
