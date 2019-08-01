package main

import (
	//"database/sql"
	//"fmt"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	//"github.com/lib/pq"
	"PetBook/controllers"
	"PetBook/driver"
	"log"
	"net/http"
	"os"
	//"PetBook/models"
	"PetBook/store"
)

func logErr(err error) {
	if err != nil {
		log.Println(err)
	}
}

func main() {
	var db *sqlx.DB
	db = driver.ConnectDB()

	router := mux.NewRouter()

	storeUser := store.UserStore{DB: db}
	controller := controllers.Controller{US: &storeUser}

	router.HandleFunc("/register", controller.RegisterPostHandler()).Methods("POST")
	router.HandleFunc("/register", controller.RegisterGetHandler()).Methods("GET")

	router.HandleFunc("/login", controller.LoginPostHandler()).Methods("POST")
	router.HandleFunc("/login", controller.LoginGetHandler()).Methods("GET")

	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/",
		http.FileServer(http.Dir("./web/static/"))))

	loggedRouter := handlers.LoggingHandler(os.Stdout, router)
	log.Fatal(http.ListenAndServe(":8181", loggedRouter))

	//err := storeUser.Login(user)
	// if err != nil {
	// 	log.Println(err)
	// 	return
	// }
	// log.Println("user logged in")
}
