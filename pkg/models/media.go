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
	AddMediaPathDb(path string,userId int)
	GetExistedGallery(userId int) []string
}

func (m *MediaStore) AddLogoPathDb(path string,userId int){
	_, err := m.DB.Exec("insert into logos (logo_path,user_id)VALUES ($1,$2)", path, userId)
	if err != nil {
		log.Println(err)
	}
}
func (m *MediaStore) AddMediaPathDb(path string,userId int){
	_, err := m.DB.Exec("insert into gallery (file_path,user_id)VALUES ($1,$2)", path, userId)
	if err != nil {
		log.Println(err)
	}
}

func (m *MediaStore) GetLogo(userId int) string{
	var path string
	err := m.DB.QueryRow("select logo_path from logos where logo_path IS NOT NULL and user_id=$1 Order by created_time DESC LIMIT 1", userId).Scan(&path)
	if err != nil {
		log.Println(err)
	}
	return path
}

func (m *MediaStore) GetExistedLogo(userId int) []string{
	rows, err := m.DB.Query("select logo_path from logos where user_id =$1 order by created_time desc;", userId)
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
func (m *MediaStore) GetExistedGallery(userId int) []string{
	rows, err := m.DB.Query("select file_path from gallery where user_id =$1 order by created_time desc;", userId)
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