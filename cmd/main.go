package main

import (
	"github.com/dpgolang/PetBook/gomigrations"
	"github.com/dpgolang/PetBook/pkg/authentication"
	"github.com/dpgolang/PetBook/pkg/controllers"
	"github.com/dpgolang/PetBook/pkg/driver"
	"github.com/dpgolang/PetBook/pkg/logger"
	_ "github.com/dpgolang/PetBook/pkg/logger"
	"github.com/dpgolang/PetBook/pkg/models"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/urfave/negroni"

	"net/http"
	"os"
)

func main() {

	db := driver.ConnectDB()
	err := gomigrations.Migrate(db)
	if err != nil {
		logger.FatalError(err, "Migration failed.\n")
	}

	router := mux.NewRouter()

	storeUser := models.UserStore{DB: db}
	storePet := models.PetStore{DB: db}
	controller := controllers.Controller{
		PetStore:  &storePet,
		UserStore: &storeUser,
	}

	router.HandleFunc("/register", controller.RegisterPostHandler()).Methods("POST")
	router.HandleFunc("/register", controller.RegisterGetHandler()).Methods("GET")

	router.HandleFunc("/login", controller.LoginPostHandler()).Methods("POST")
	router.HandleFunc("/login", controller.LoginGetHandler()).Methods("GET")

	router.Handle("/mypage", negroni.New(
		negroni.HandlerFunc(authentication.ValidateTokenMiddleware),
		negroni.Wrap(http.HandlerFunc(controller.MyPageGetHandler())),
	)).Methods("GET")

	router.Handle("/petcabinet", negroni.New(
		negroni.HandlerFunc(authentication.ValidateTokenMiddleware),
		negroni.Wrap(http.HandlerFunc(controller.PetPutHandler())),
	)).Methods("PUT")

	router.Handle("/petcabinet", negroni.New(
		negroni.HandlerFunc(authentication.ValidateTokenMiddleware),
		negroni.Wrap(http.HandlerFunc(controller.PetPutHandler())),
	)).Methods("GET")


	router.Handle("/", negroni.New(
		negroni.HandlerFunc(authentication.ValidateTokenMiddleware),
		negroni.Wrap(http.HandlerFunc(controller.MyPageGetHandler())),
	))

	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/",
		http.FileServer(http.Dir("./web/static/"))))

	loggedRouter := handlers.LoggingHandler(os.Stdout, router)

	// Is it proper way to handle ListenAndServe() error?
	if err:= http.ListenAndServe(":8080", loggedRouter); err !=nil {
		logger.FatalError(err, "Error occurred, while trying to listen and serve a server")
	}
}