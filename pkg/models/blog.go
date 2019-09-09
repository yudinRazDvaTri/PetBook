package models

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"log"
	"time"
)

type Blog struct {
	Id       int
	UserId   string
	Name     string
	Message  string
	Date     time.Time
	LogoPath string
}

type BlogStore struct {
	DB *sqlx.DB
}

type BlogStorer interface {
	GetPetBlog(userid int) ([]Blog,error)
	CreateBlog(form string, idUser int) error
	DeleteBlog(blogid string) error
	GetVetBlog(userID int) ([]Blog ,error)
}

func logFatal(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func (b *BlogStore) GetPetBlog(userID int) ([]Blog ,error){
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
		tRes.Name = name
		tRes.Date = time
		tRes.Message = message
		results = append(results, tRes)
		if err != nil {
			return nil,fmt.Errorf("cannot insert message to messages in db: %v", err)
		}
	}
	return results,nil
}
func (b *BlogStore) GetVetBlog(userID int) ([]Blog ,error){
	rows, err := b.DB.Query("select blog_id, content, created_time, name from blog,vets where blog.user_id =$1 and vets.user_id=blog.user_id order by created_time desc;", userID)
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
		tRes.Name = name
		tRes.Date = time
		tRes.Message = message
		results = append(results, tRes)
		if err != nil {
			return nil,fmt.Errorf("cannot insert message to messages in db: %v", err)
		}
	}
	return results,nil
}




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
