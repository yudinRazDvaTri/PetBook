package main

import (
	//"database/sql"
	"fmt"
	_ "github.com/gorilla/handlers"
	_ "github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	//"github.com/lib/pq"
	"log"
	//"net/http"
	//"os"
	"petbook/controllers"
	"petbook/driver"
	"petbook/models"
	//"petbook/store"
)

/*type ViewDataUsers struct {
	Title   string `json:"title"
	Content []models.User `json:"content"`
}
type ViewDataUser struct {
	Title   string `json:"title"`
	Content models.User   `json:"content"`
}*/

func logErr(err error) {
	if err != nil {
		log.Println(err)
	}
}

func main() {
	var db *sqlx.DB
	//controller := controllers.Controller{}
	db = driver.ConnectDB()

	user := &models.User{
		Email:     "newEMAIL@gmail.co213132m",
		Login:     "myLOGINnew",
		Password:  "333124124124141241241",
		Firstname: "name",
		Lastname:  "surname",
	}
	controllerUser := controllers.UserController{DB: db}
	//controllerUser.DB = db
	err := controllerUser.ChangePassword(user, "99999999999999")
	//_____________________________________________________
	// user := &models.User{}
	// err := GetUser(user, db, "user")
	// logErr(err)
	// err = ChangePassword(user, db, "1111111")
	logErr(err)
	fmt.Println(user)
}
