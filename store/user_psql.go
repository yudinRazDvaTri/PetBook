package store

import (
	"database/sql"
	"test/models"

	"fmt"
	"github.com/jmoiron/sqlx"
	"log"
)

func logErr(err error) {
	if err != nil {
		log.Println(err)
	}
}

type UserStore struct {
	DB *sqlx.DB
}

func (c *UserStore) GetUsers() ([]models.User, error) {
	rows, err := c.DB.Query("select * from users")
	logErr(err)
	defer rows.Close()
	users := []models.User{}
	err = sqlx.StructScan(rows, &users)
	if err != nil {
		logErr(err)
		return users, fmt.Errorf("cannot scan users from db: %v", err)
	}
	return users, nil
}

func (c *UserStore) GetUser(user *models.User) error {
	err := c.DB.QueryRowx("select * from users where email=$1", user.Email).StructScan(user)
	if err != nil {
		logErr(err)
		return fmt.Errorf("cannot scan user from db: %v", err)
	}
	return nil
}

func (c *UserStore) Register(user *models.User) error {
	_, err := c.DB.Exec("insert into users (email,firstname, lastname, login ,password) values ($1,$2,$3, $4, $5)",
		user.Email, user.Firstname, user.Lastname, user.Login, user.Password)
	if err != nil {
		return fmt.Errorf("cannot affect rows in users in db: %v", err)
	}
	return nil
}

func (c *UserStore) ChangePassword(user *models.User, newPassword string) error {
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

// func Register(db *sqlx.DB, user models.User) int {
// 	err := db.QueryRow("insert into teammate (firstname,lastname,password) values ($1,$2,$3) RETURNING id_user;",
// 		user.Firstname, user.Lastname, user.Password).Scan(&user.ID)
// 	logErr(err)
// 	models.AddUser(&user)
// 	return user.ID
// }

func (c *UserStore) Login(userСhecking *models.User) error {
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

func (c *UserStore) GetPet(user *models.User) (models.Pet, error) {
	pet := models.Pet{}
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

type Pet struct {
	ID      int     `db:"user_id"`
	Name    string  `db:"name"`
	PetType string  `db:"animal_type"`
	Breed   string  `db:"breed"`
	Age     int     `db:"age"`
	Weight  float32 `db: "weight"`
	Gender  string  `db:"gender"`
}
