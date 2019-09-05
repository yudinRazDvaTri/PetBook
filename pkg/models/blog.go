package models

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"log"
	"time"
)

type Blog struct {
	Id      int
	UserId  string
	PetName string
	Message string
	Date    time.Time
	LogoPath string
}

type BlogStore struct {
	DB *sqlx.DB
}

type BlogStorer interface {
	GetBlog(userid int) ([]Blog,error)
	CreateBlog(form string, idUser int) error
	DeleteBlog(blogid string) error
}

func logFatal(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func (b *BlogStore) GetBlog(userID int) ([]Blog ,error){
	rows, err := b.DB.Query("select blog_id, content, created_time, name from blog,pets where blog.user_id =$1 and pets.user_id=blog.user_id order by created_time desc;", userID)
	if err != nil {
		return nil,fmt.Errorf("cannot connect to database: %v", err)
	}
	tRes := Blog{}
	var results []Blog
	for rows.Next() {
		var blogid int
		var message, name string
		var time time.Time
		err = rows.Scan(&blogid, &message, &time, &name)
		tRes.Id = blogid
		tRes.PetName = name
		tRes.Date = time
		tRes.Message = message
		results = append(results, tRes)
		if err != nil {
			return nil,fmt.Errorf("cannot insert message to messages in db: %v", err)
		}
	}
	return results,nil
}

//func (b *BlogStore) GetBlog() []Blog{
//	rows, err:= b.DB.Query("select blog_id, content from blog")
//	if err != nil{
//		logFatal(err)
//	}
//	tRes:=Blog{}
//	var results []Blog
//	for rows.Next(){
//		var id int
//		var message string
//		err = rows.Scan(&id,&message)
//		tRes.Id=id
//		tRes.Message=message
//		results = append(results,tRes)
//		if err != nil{
//			logFatal(err)
//		}
//	}
//	return results
//}

func (b *BlogStore) CreateBlog(form string, idUser int) error{
	result, err := b.DB.Exec("insert into blog (content,user_id) values ($1,$2);", form, idUser)
	if err != nil {
		return fmt.Errorf("cannot execute database query: %v", err)
	}
	_, err = result.RowsAffected()
	if err != nil {
		return fmt.Errorf("cannot create post: %v", err)
	}
	return nil
}

func (b *BlogStore) DeleteBlog(blogid string) error{
	result, err := b.DB.Exec("delete from blog where blog_id = $1", blogid)
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
