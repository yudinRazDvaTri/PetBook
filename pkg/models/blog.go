package models

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"log"
	"time"
)

type Blog struct{
	Id int
	UserId string
	Message string
	Date time.Time
}

type BlogStore struct {
	DB *sqlx.DB
}

type BlogStorer interface {
	GetBlog() []Blog
	CreateBlog (form string,idUser int)
	DeleteBlog (blogid string)
}

func logFatal(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func (b *BlogStore) GetBlog() []Blog{
	rows, err:= b.DB.Query("select blog_id, content from blog")
	if err != nil{
		logFatal(err)
	}
	tRes:=Blog{}
	var results []Blog
	for rows.Next(){
		var id int
		var message string
		err = rows.Scan(&id,&message)
		tRes.Id=id
		tRes.Message=message
		results = append(results,tRes)
		if err != nil{
			logFatal(err)
		}
	}
	return results
}

func (b *BlogStore) CreateBlog (form string,idUser int){

	result, err:= b.DB.Exec("insert into blog (content,user_id) values ($1,$2);",form,idUser)
	if err != nil {
		log.Fatal(err)
	}
	_,err = result.RowsAffected()
	if err != nil {
		log.Println(err)
		return
	}
}

func (b *BlogStore) DeleteBlog (blogid string){
	result, err := b.DB.Exec("delete from blog where blog_id = $1;",blogid)
	if err != nil {
		log.Println("didn't delete blodid ",blogid)
		return
	}
	n,err := result.RowsAffected()
	if err != nil{
		log.Println("didn't delete blodid ",blogid)
		return
	}
	fmt.Println("rows affected - ",n)
}

