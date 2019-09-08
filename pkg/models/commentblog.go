package models

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	//"log"
	"time"
)

type CommentBlog struct {
	CommentId       int
	CommentBlogId	 int
	CommentUserId   string
	CommentName     string
	CommentMessage  string
	CommentDate     time.Time
	//LogoPath string
}

type CommentBlogStore struct {
	DB *sqlx.DB
}

type CommentBlogStorer interface {
	GetPetCommentBlog(userid int) ([]CommentBlog,error)
	CreateCommentBlog(form string, idUser int,idBlog int) error
	DeleteCommentBlog(blogid string) error
	GetVetCommentBlog(userID int) ([]CommentBlog ,error)
}

//func logFatal(err error) {
//	if err != nil {
//		log.Fatal(err)
//	}
//}

func (b *CommentBlogStore) GetPetCommentBlog(userID int) ([]CommentBlog ,error){
	rows, err := b.DB.Query("select commentblog_id, blog_id, content, created_time, name from blog,pets,commentblog where blog.user_id =$1 and pets.user_id=blog.user_id order by created_time desc;", userID)
	if err != nil {
		return nil,fmt.Errorf("cannot connect to database: %v", err)
	}
	tRes := CommentBlog{}
	var results []CommentBlog
	for rows.Next() {
		var blogid int
		var message, name string
		var time time.Time
		err = rows.Scan(&blogid, &message, &time, &name)
		tRes.CommentId = blogid
		tRes.CommentName = name
		tRes.CommentDate = time
		tRes.CommentMessage = message
		results = append(results, tRes)
		if err != nil {
			return nil,fmt.Errorf("cannot insert message to messages in db: %v", err)
		}
	}
	return results,nil
}
func (b *CommentBlogStore) GetVetCommentBlog(userID int) ([]CommentBlog ,error){
	rows, err := b.DB.Query("select blog_id, content, created_time, name from blog,vets where blog.user_id =$1 and vets.user_id=blog.user_id order by created_time desc;", userID)
	if err != nil {
		return nil,fmt.Errorf("cannot connect to database: %v", err)
	}
	tRes := CommentBlog{}
	var results []CommentBlog
	for rows.Next() {
		var blogid int
		var message, name string
		var time time.Time
		err = rows.Scan(&blogid, &message, &time, &name)
		tRes.CommentId = blogid
		tRes.CommentName = name
		tRes.CommentDate = time
		tRes.CommentMessage = message
		results = append(results, tRes)
		if err != nil {
			return nil,fmt.Errorf("cannot insert message to messages in db: %v", err)
		}
	}
	return results,nil
}

func (b *CommentBlogStore) CreateCommentBlog(form string, idUser int,idBlog int) error{
	result, err := b.DB.Exec("insert into commentblog (content,user_id,blog_id) values ($1,$2,$3) ;", form, idUser,idBlog)
	if err != nil {
		return fmt.Errorf("cannot execute database query: %v", err)
	}
	_, err = result.RowsAffected()
	if err != nil {
		return fmt.Errorf("cannot create post: %v", err)
	}
	return nil
}

func (b *CommentBlogStore) DeleteCommentBlog(commentblogid string) error{
	result, err := b.DB.Exec("delete from commentblog where commentblog_id = $1", commentblogid)
	if err != nil {
		return fmt.Errorf("cannot execute database query: %v", err)
	}
	n, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("cannot delete post: %v", err)
	}
	fmt.Println("rows affected - ", n)
	return nil
}

