package main

import (
	//"database/sql"
	//"fmt"
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

func logErr(err error) {
	if err != nil {
		log.Println(err)
	}
}

func main() {
	var db *sqlx.DB
	db = driver.ConnectDB()

	user := &models.User{
		Email:     "asdsad@gmail.com",
		Login:     "mylogin",
		Password:  "123123",
		Firstname: "name",
		Lastname:  "surname",
	}
	controllerUser := controllers.UserController{DB: db}
	err := controllerUser.Login(user)
	if err != nil {
		log.Println(err)
		return
	}
	log.Println("user logged in")

}
