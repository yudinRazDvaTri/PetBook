package models

import (
	"database/sql"
	"fmt"
	"github.com/jmoiron/sqlx"
	"log"
)

func logErr(err error) {
	if err != nil {
		log.Println(err)
	}
}

type User struct {
	ID        int    `json:"id" db:"id"`
	Email     string `json:"email" db:"email"`
	Login     string `json:"login" db:"login"`
	UserType  string `db:"pet_or_vet"`
	Firstname string `json:"firstname" db:"firstname"`
	Lastname  string `json:"lastname" db:"lastname"`
	Password  string `json:"password" db:"password"`
}

type UserStorer interface {
	GetUsers() ([]User, error)
	GetUser(user *User) error
	Register(user *User) error
	ChangePassword(user *User, newPassword string) error
	Login(userСhecking *User) error
	GetPet(user *User) (Pet, error)
	ReadUserID(user *User) error
}

type UserStore struct {
	DB *sqlx.DB
}

func (c *UserStore) GetUsers() ([]User, error) {
	rows, err := c.DB.Query("select * from users")
	logErr(err)
	defer rows.Close()
	users := []User{}
	err = sqlx.StructScan(rows, &users)
	if err != nil {
		logErr(err)
		return users, fmt.Errorf("cannot scan users from db: %v", err)
	}
	return users, nil
}

func (c *UserStore) GetUser(user *User) error {
	err := c.DB.QueryRowx("select * from users where email=$1", user.Email).StructScan(user)
	if err != nil {
		return fmt.Errorf("cannot scan user from db: %v", err)
	}
	return nil
}

func (c *UserStore) ReadUserID(user *User) error {
	err := c.DB.QueryRow("select id from users where email=$1", user.Email).Scan(user.ID)
	if err != nil {
		return fmt.Errorf("cannot scan userID from db: %v", err)
	}
	return nil
}

func (c *UserStore) Register(user *User) error {
	_, err := c.DB.Exec("insert into users (email,firstname, lastname, login ,password) values ($1,$2,$3, $4, $5)",
		user.Email, user.Firstname, user.Lastname, user.Login, user.Password)
	if err != nil {
		return fmt.Errorf("cannot affect rows in users in db: %v", err)
	}
	return nil
}

func (c *UserStore) ChangePassword(user *User, newPassword string) error {
	res, err := c.DB.Exec("UPDATE users SET password=$1 WHERE email = $2",
		newPassword, user.Email)
	if err != nil {
		return fmt.Errorf("cannot update users in db: %v", err)
	}
	num, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("cannot affect rows in users in db: %v", err)
	}
	if num != 1 {
		return fmt.Errorf("cannot find this user")
	}
	user.Password = newPassword
	return nil
}

//TODO: create custom error
func (c *UserStore) Login(userСhecking *User) error {
	var passwordFromBase string
	err := c.DB.QueryRow("select password from users where email=$1", userСhecking.Email).Scan(&passwordFromBase)
	if userСhecking.Password != passwordFromBase || err == sql.ErrNoRows {
		return fmt.Errorf("wrong login data")
	}
	if err != nil {
		return fmt.Errorf("cannot login this user: %v", err)
	}
	return nil
}

func (c *UserStore) GetPet(user *User) (Pet, error) {
	pet := Pet{}
	err := c.DB.QueryRowx(
		`SELECT user_id,name,animal_type, breed, age, weight, gender 
		FROM pets p, users u 
		WHERE p.user_id = u.id  
		AND u.email = $1 `, user.Email).StructScan(&pet)
	if err != nil {
		return pet, err
	}
	return pet, nil
}
