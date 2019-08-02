package models

import (
//"database/sql"
//"github.com/jmoiron/sqlx"
//"petbook/store"
)

type User struct {
	ID        int    `json:"id" db:"id"`
	Email     string `json:"email" db:"email"`
	Login     string `json:"login" db:"login"`
	UserType  string `db:"pet_or_vet"`
	Firstname string `json:"firstname" db:"firstname"`
	Lastname  string `json:"lastname" db:"lastname"`
	Password  string `json:"password" db:"password"`
}

// func (u *User) CreateUser(db *sql.DB) err {
// 	u, err := modelLoading.LoadUser(db, u.Login)
// 	if err == sql.ErrNoRows {
// 		if err := modelLoading.CreateUser(db, u); err != nil {
// 			return err //something gone wrong while inserting to db
// 		}
// 		return nil // we added successfully
// 	}
// 	return err //if there we have user with such login
// }

// func (u *User) GetUser(db *sql.DB, login string) err {
// 	u, err := modelLoading.LoadUser(db, login)
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }

// func (u *User) ChangePassword(db *sql.DB, newPassword string) err {
// 	err := modelLoading.ChangePassword(db, u, newPassword)
// 	if err != nil {
// 		return err
// 	}
// 	u.Password = newPassword
// 	return nil
// }
