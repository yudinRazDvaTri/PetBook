package models

import (
	"database/sql"
	"fmt"
	"github.com/dpgolang/PetBook/pkg/utilerr"
	"github.com/jmoiron/sqlx"
	"os"
	"time"
)

type Media struct {
	Id          int
	UserId      int
	LogoPath    string
	GalleryPath string
	DocsPath    string
	Time        time.Time
}

type MediaStore struct {
	DB *sqlx.DB
}

type MediaStorer interface {
	AddLogoPathDb(path string, userId int) error
	GetLogo(userId int) (string, error)
	AddMediaPathDb(path string, userId int) error
	GetExistedGallery(userId int) ([]string, error)
	DeleteFile(path string) error
}
//Add user`s logo path in DB by getting path as parameter
func (m *MediaStore) AddLogoPathDb(path string, userId int) error {
	_, err := m.DB.Exec("insert into logos (logo_path,user_id)VALUES ($1,$2)", path, userId)
	if err != nil {
		return fmt.Errorf("cannot execute database query: %v", err)
	}
	return nil
}
//Add user`s photo path in DB by 1 photo in moment
func (m *MediaStore) AddMediaPathDb(path string, userId int) error {
	_, err := m.DB.Exec("insert into gallery (file_path,user_id)VALUES ($1,$2)", path, userId)
	if err != nil {
		return fmt.Errorf("cannot execute database query: %v", err)
	}
	return nil
}
//Select user`s logo path in DB
func (m *MediaStore) GetLogo(userId int) (string, error) {
	var path string
	err := m.DB.QueryRow("select logo_path from logos where logo_path IS NOT NULL and user_id=$1 Order by created_time DESC LIMIT 1", userId).Scan(&path)
	if err != nil {
		if err == sql.ErrNoRows {
			return path, &utilerr.LogoDoesNotExist{Description: "User doesn't have a logo."}
		}
		return path, fmt.Errorf("cannot connect to database: %v", err)
	}
	return path, nil
}
//Select user`s photos paths from DB
func (m *MediaStore) GetExistedGallery(userId int) ([]string, error) {
	var results []string
	rows, err := m.DB.Query("select file_path from gallery where user_id =$1 order by created_time desc;", userId)
	if err != nil {
		return results, fmt.Errorf("cannot connect to database: %v", err)
	}
	for rows.Next() {
		var p string
		err = rows.Scan(&p)
		if err != nil {
			return nil, fmt.Errorf("cannot get photo from media in db: %v", err)
		}
		results = append(results, p)
	}
	return results, nil
}
//Func delete file from Folder and DB
func (m *MediaStore) DeleteFile(path string) error {
	err := os.Remove(path)
	if err != nil {
		return fmt.Errorf("cannot delete from folder: %v", err)
	}
	_, err2 := m.DB.Exec("delete from gallery where file_path = $1", path)
	if err2 != nil {
		return fmt.Errorf("cannot execute database query: %v", err2)
	}
	fmt.Println("==> done deleting file")
	return nil
}
