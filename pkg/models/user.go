package models

import (
	"database/sql"
	"fmt"
	"github.com/dpgolang/PetBook/pkg/utilerr"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
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

type UserStore struct {
	DB *sqlx.DB
}

type UserStorer interface {
	GetUsers() ([]User, error)
	GetUser(userID int) (User, error)
	Register(user *User) error
	ChangePassword(user *User, newPassword string) error
	Login(email, userPassword string) (int, error)
	GetPet(userID int) (Pet, error)
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

func (c *UserStore) GetUser(userID int) (User, error) {
	var user User
	err := c.DB.QueryRowx("select * from users where id=$1", userID).StructScan(&user)
	if err != nil {
		return user, fmt.Errorf("cannot scan user from db: %v", err)
	}
	return user, nil
}

func (c *UserStore) Register(user *User) error {
	_, err := c.DB.Exec("insert into users (email, firstname, lastname, login, password) values ($1,$2,$3, $4, $5)",
		user.Email, user.Firstname, user.Lastname, user.Login, user.Password)

	if err != nil {
		if _, ok := err.(*pq.Error); ok {
			return &utilerr.UniqueTaken{Description: "Id or login has already been taken!"}
		}
		return fmt.Errorf("Error occurred while trying to add new user: %v.\n", err)
	}

	return nil
}

func (c *UserStore) ChangePassword(user *User, newPassword string) error {
	res, err := c.DB.Exec("UPDATE users SET password=$1 WHERE email = $2",
		newPassword, user.Email)
	if err != nil {
		return fmt.Errorf("Error occurred while trying to update user in db: %v.\n", err)
	}

	num, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("Error occurred while trying to count affected rows in users in db: %v.\n", err)
	}

	if num != 1 {
		return &utilerr.WrongCredentials{Description: "Wrong email or password."}
	}

	user.Password = newPassword
	return nil
}

func (c *UserStore) Login(email, userPassword string) (int, error) {
	var passwordFromBase string
	var idFromBase int
	err := c.DB.QueryRow("select password, id from users where email=$1", email).Scan(&passwordFromBase, &idFromBase)

	if err != nil {
		if err == sql.ErrNoRows {
			return 0, &utilerr.WrongCredentials{Description: "Wrong email or password."}
		}
		return 0, fmt.Errorf("Error occurred while trying to login user: %v.\n", err)
	}

	if err = bcrypt.CompareHashAndPassword([]byte(passwordFromBase), []byte(userPassword)); err != nil {
		return 0, &utilerr.WrongCredentials{Description: "Wrong email or password."}
	}
	return idFromBase, nil
}

func (c *UserStore) GetPet(userID int) (Pet, error) {
	var pet Pet
	err := c.DB.QueryRowx(
		`SELECT user_id,name,animal_type, breed, age, weight, gender 
		FROM pets p, users u 
		WHERE p.user_id = u.id  
		AND u.id = $1 `, userID).StructScan(&pet)

	if err != nil {
		return pet, err
	}

	return pet, nil
}
