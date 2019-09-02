package models

import (
	"github.com/jmoiron/sqlx"
	"log"
	"time"
)

type Media struct {
	Id int
	UserId int
	LogoPath string
	GalleryPath string
	DocsPath string
	Time time.Time
}

type MediaStore struct {
	DB *sqlx.DB
}

type MediaStorer interface {
	AddLogoPathDb(path string,userId int)
	GetLogo(userId int) string
	GetExistedLogo(userId int) []string
}

func (m *MediaStore) AddLogoPathDb(path string,userId int){
	_, err := m.DB.Exec("insert into media (logo_path,user_id)VALUES ($1,$2)", path, userId)
	if err != nil {
		log.Println(err)
	}
}

func (m *MediaStore) GetLogo(userId int) string{
	var path string
	err := m.DB.QueryRow("select logo_path from media where logo_path IS NOT NULL and media_id=(select max(media_id)from media) and user_id=$1", userId).Scan(&path)
	if err != nil {
		log.Println(err)
	}
	return path
}

func (m *MediaStore) GetExistedLogo(userId int) []string{
	rows, err := m.DB.Query("select logo_path from media where user_id =$1 order by created_time desc;", userId)
	if err != nil {
		log.Println(err)
	}
	var results []string
	for rows.Next() {
		var p string
		err = rows.Scan(&p)
		if err != nil {
			log.Println(err)
		}
		results = append(results, p)
	}
	return results
}