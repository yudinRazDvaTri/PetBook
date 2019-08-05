package main

import (
	_ "PetBook/init"
	"PetBook/pkg/utils"
	//"database/sql"
	//"fmt"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	//"github.com/jmoiron/sqlx"
	"PetBook/gomigrations"
	"github.com/urfave/negroni"
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

	//	var db *sqlx.DB
	db := driver.ConnectDB()
	err := gomigrations.Migrate(db)
	if err != nil {
		log.Fatal("Migration failed.")
	}

	router := mux.NewRouter()

	storeUser := store.UserStore{DB: db}
	controller := controllers.Controller{UserStore: &storeUser}

	router.HandleFunc("/register", controller.RegisterPostHandler()).Methods("POST")
	router.HandleFunc("/register", controller.RegisterGetHandler()).Methods("GET")

	router.HandleFunc("/login", controller.LoginPostHandler()).Methods("POST")
	router.HandleFunc("/login", controller.LoginGetHandler()).Methods("GET")

	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/",
		http.FileServer(http.Dir("./web/static/"))))

	router.Handle("/cabinetPet", negroni.New(
		negroni.HandlerFunc(utils.ValidateTokenMiddleware),
		negroni.Wrap(http.HandlerFunc(controller.CreatePetHandler())),
	))

	router.Handle("/mypage", negroni.New(
		negroni.HandlerFunc(utils.ValidateTokenMiddleware),
		negroni.Wrap(http.HandlerFunc(controller.MyPageGetHandler())),
	))

	router.Handle("/", negroni.New(
		negroni.HandlerFunc(utils.ValidateTokenMiddleware),
		negroni.Wrap(http.HandlerFunc(controller.MyPageGetHandler())),
	))

	loggedRouter := handlers.LoggingHandler(os.Stdout, router)
	log.Fatal(http.ListenAndServe(":8080", loggedRouter))

	//err := storeUser.Login(user)
	// if err != nil {
	// 	log.Println(err)
	// 	return
	// }
	// log.Println("user logged in")
}
