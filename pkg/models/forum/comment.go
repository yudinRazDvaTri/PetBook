package forum

import (
	"time"
)

type Comment struct {
	CommendID   int       `json:"comment_id" db:"comment_id"`
	TopicID     int       `json:"topic_id" db:"topic_id"`
	UserID      int       `json:"user_id" db:"user_id"`
	CreatedTime time.Time `json:"created_time" db:"created_time"`
	Content     string    `json:"content" db:"content"`
}

//func (f *ForumStore) AddNewComment(userID int, content string) (err error) {
//	_, err = f.DB.Exec(
//		`insert into comments (user_id, content) values ($1, $2)`, userID, content)
//	if err != nil {
//		return fmt.Errorf("cannot affect rows in pets in db: %v", err)
//	}
//	return
//}
//
//func (f *ForumStore) GetTopicComments(topicID int) (comments []*Comment, err error) {
//	rows, err := f.DB.Query("select * from comments where topic_id = $1 order by created_time DESC", topicID)
//	if err != nil {
//		err = fmt.Errorf("Can't read comments-rows from db: %v", err)
//		return
//	}
//	defer rows.Close()
//	err = sqlx.StructScan(rows, &comments)
//	if err != nil {
//		err = fmt.Errorf("Can't scan comments-rows from db: %v", err)
//	}
//	return
//}