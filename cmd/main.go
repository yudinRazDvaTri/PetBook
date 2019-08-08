package main

import (
	"github.com/dpgolang/PetBook/gomigrations"
	"github.com/dpgolang/PetBook/pkg/authentication"
	"github.com/dpgolang/PetBook/pkg/controllers"
	"github.com/dpgolang/PetBook/pkg/driver"
	"github.com/dpgolang/PetBook/pkg/handler"
	_ "github.com/dpgolang/PetBook/pkg/logger"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/urfave/negroni"

	"log"
	"net/http"
	"os"
)


func main() {

	db := driver.ConnectDB()
	err := gomigrations.Migrate(db)
	if err != nil {
		log.Fatal("Migration failed.")
	}

	router := mux.NewRouter()

	storeUser := controllers.UserStore{DB: db}
	controller := handler.Controller{UserStore: &storeUser}
	storePet:=controllers.PetStore{DB:db}
	controllerPet:=handler.Controller{PetStore:&storePet}

	router.HandleFunc("/register", controller.RegisterPostHandler()).Methods("POST")
	router.HandleFunc("/register", controller.RegisterGetHandler()).Methods("GET")

	router.HandleFunc("/login", controller.LoginPostHandler()).Methods("POST")
	router.HandleFunc("/login", controller.LoginGetHandler()).Methods("GET")

	router.HandleFunc("/petcabinet", controllerPet.CreatePetPostHandler()).Methods("POST")
	router.HandleFunc("/petcabinet", controllerPet.CreatePetGetHandler()).Methods("GET")


	router.Handle("/mypage", negroni.New(
		negroni.HandlerFunc(authentication.ValidateTokenMiddleware),
		negroni.Wrap(http.HandlerFunc(controller.MyPageGetHandler())),
	))


	router.PathPrefix("/assets/").Handler(http.StripPrefix("/assets/",
		http.FileServer(http.Dir("./web/assets/"))))
	router.Handle("/", negroni.New(
		negroni.HandlerFunc(authentication.ValidateTokenMiddleware),
		negroni.Wrap(http.HandlerFunc(controller.MyPageGetHandler())),
	))
	loggedRouter := handlers.LoggingHandler(os.Stdout, router)
	log.Fatal(http.ListenAndServe(":8080", loggedRouter))

}
