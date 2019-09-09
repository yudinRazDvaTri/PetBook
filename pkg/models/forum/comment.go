package forum

import (
	"fmt"
	"sort"
	"time"

	"github.com/lib/pq"
)

type Comment struct {
	CommentID   int       `json:"comment_id" db:"comment_id"`
	TopicID     int       `json:"topic_id" db:"topic_id"`
	UserID      int       `json:"user_id" db:"user_id"`
	CreatedTime time.Time `json:"created_time" db:"created_time"`
	Content     string    `json:"content" db:"content"`
	ParentID    int       `json:"parent_id" db:"parent_id"`
	HaveKids    bool      `json:"have_kids" db:"have_kids"`
}

// View Alias-Struct to layout comment properly
type ViewComment struct {
	UserName  string
	LikersIDs []int64
	NestedComments []*ViewComment
	Comment
}

// Nesting comments
func TreeViewComments(vComments []ViewComment) (treeVComments []*ViewComment, err error) {

	for i := range vComments {

		if vComments[i].ParentID == 0 {
			treeVComments = append(treeVComments, &vComments[i])
		}
	}

	sort.Sort(sort.Reverse(ByRating(treeVComments)))

	for i := range treeVComments {
		for j := range vComments {
			if vComments[j].ParentID == treeVComments[i].CommentID {
				treeVComments[i].NestedComments = append(treeVComments[i].NestedComments, &vComments[j])
			}
		}
		sort.Sort(sort.Reverse(ByRating(treeVComments[i].NestedComments)))
	}

	return treeVComments, nil
}

// Method to discover if user can like comment
func (v *ViewComment) CanLike(userID int) bool {
	if v.UserID == userID {
		return false
	}
	for i := range v.LikersIDs {
		if int(v.LikersIDs[i]) == userID {
			return false
		}
	}
	return true
}

func (v *ViewComment) CanReply(userID int) bool {
	if v.UserID == userID {
		return false
	}
	return true
}

// ViewComment Constructor
func (f *ForumStore) NewViewComment(userName string, comment Comment) (viewComment ViewComment, err error) {
	likerIDs, err := f.getCommentLikers(comment.CommentID)
	if err != nil {
		err = fmt.Errorf("Can't read %d comment's likersIDs from DB: %v", comment.TopicID, err)
		return
	}
	var emptySlice []*ViewComment
	viewComment = ViewComment{userName, likerIDs,emptySlice, comment}
	return
}

type ByRating []*ViewComment

// Methods to sort ViewComments by Rating
func (v ByRating) Len() int {
	return len(v)
}
func (v ByRating) Swap(i, j int) {
	v[i], v[j] = v[j], v[i]
}
func (v ByRating) Less(i, j int) bool {
	return len(v[i].LikersIDs) < len(v[j].LikersIDs)
}

func (f *ForumStore) AddNewComment(topicID, userID, parentID int, content string) (err error) {
	_, err = f.DB.Exec(
		`insert into comments (topic_id, user_id, parent_id, content) values ($1, $2, $3, $4);`,
		topicID, userID, parentID, content)
	if err != nil {
		err = fmt.Errorf("can not add new comment to db: %v", err)
		return
	}
	return
}

func (f *ForumStore) RateComment(commentID, userID int) (bool, error) {
	rateOk := true

	_, err := f.DB.Exec(
		`insert into ratings (comment_id, user_id) values ($1, $2);`,
		commentID, userID)
	if err != nil {
		rateOk = false
		if pgerr, ok := err.(*pq.Error); ok {
			if pgerr.Code == "23505" {
				err = fmt.Errorf("User %d already liked comment %d: %v",
					userID, commentID, err)
				return rateOk, err
			}
		}
		err = fmt.Errorf("Error while trying to rate %d comment by %d user in DB: %v",
			userID, commentID, err)
	}
	return rateOk, err
}

func (f *ForumStore) getCommentLikers(commentID int) (likersIDs []int64, err error) {
	rows, err := f.DB.Query(`SELECT user_id FROM ratings WHERE comment_id = $1;`, commentID)
	if err != nil {
		err = fmt.Errorf("Can't read rating-rows from db: %v", err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var likerID int64
		err = rows.Scan(&likerID)
		if err != nil {
			err = fmt.Errorf("Can't scan rating-row from db: %v", err)
			return
		}
		likersIDs = append(likersIDs, likerID)
	}
	return
}

func (f *ForumStore) getCommentsIDs(topicID int) (commentsIDs []int64, err error) {
	rows, err := f.DB.Query(
		`SELECT comment_id FROM comments WHERE topic_id = $1;`, topicID)
	if err != nil {
		err = fmt.Errorf("Can't number of comments in topic %d from DB: %v", topicID, err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var commentID int64
		err = rows.Scan(&commentID)
		if err != nil {
			err = fmt.Errorf("Can't scan comment_id-row from db: %v", err)
			return
		}
		commentsIDs = append(commentsIDs, commentID)
	}
	return
}
