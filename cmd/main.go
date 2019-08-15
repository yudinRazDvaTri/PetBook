package main

import (
	"github.com/dpgolang/PetBook/gomigrations"
	"github.com/dpgolang/PetBook/pkg/authentication"
	"github.com/dpgolang/PetBook/pkg/controllers"
	"github.com/dpgolang/PetBook/pkg/driver"
	"github.com/dpgolang/PetBook/pkg/logger"
	_ "github.com/dpgolang/PetBook/pkg/logger"
	"github.com/dpgolang/PetBook/pkg/models"
	"github.com/dpgolang/PetBook/pkg/models/forum"
	"github.com/dpgolang/PetBook/pkg/models/search"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
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
	storeRefreshToken := models.RefreshTokenStore{DB: db}
	storeForum := forum.ForumStore{DB: db}
	storeSearch := search.SearchStore{DB: db}
	storeBlog := models.BlogStore{DB: db}


	controller := controllers.Controller{
		PetStore:  &storePet,
		UserStore: &storeUser,
		ForumStore: &storeForum,
		SearchStore: &storeSearch,
		RefreshTokenStore: &storeRefreshToken,
		BlogStore: &storeBlog,
	}

	router.HandleFunc("/register", controller.RegisterPostHandler()).Methods("POST")
	router.HandleFunc("/register", controller.RegisterGetHandler()).Methods("GET")

	router.HandleFunc("/login", controller.LoginPostHandler()).Methods("POST")
	router.HandleFunc("/login", controller.LoginGetHandler()).Methods("GET")

	subrouter := router.PathPrefix("/").Subrouter()
	subrouter.Use(authentication.ValidateTokenMiddleware(&storeRefreshToken))

	subrouter.HandleFunc("/mypage", controller.MyPageGetHandler()).Methods("GET")
	subrouter.HandleFunc("/petcabinet", controller.PetPostHandler()).Methods("POST")
	subrouter.HandleFunc("/petcabinet", controller.PetGetHandler()).Methods("GET")

	subrouter.HandleFunc("/forum/topics", controller.ViewTopicsHandler()).Methods("GET")
	subrouter.HandleFunc("/forum/topics/new", controller.NewTopicHandler()).Methods("POST")
	subrouter.HandleFunc("/forum/topics/new", controller.NewTopicHandler()).Methods("GET")

	subrouter.HandleFunc("/search", controller.ViewSearchHandler()).Methods("GET")
	subrouter.HandleFunc("/search", controller.SearchHandler()).Methods("POST")
	subrouter.HandleFunc("/", controller.MyPageGetHandler())

	subrouter.HandleFunc("/",controller.GetBlogHandler)
	subrouter.HandleFunc("/process",controller.CreateBlogHandler)
	subrouter.HandleFunc("/delete",controller.DeleteBlogHandler)
	router.HandleFunc("/upload",controllers.UploadFile)



	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/",
		http.FileServer(http.Dir("./web/static/"))))

	loggedRouter := handlers.LoggingHandler(os.Stdout, router)

	// Is it proper way to handle ListenAndServe() error?
	if err:= http.ListenAndServe(":8080", loggedRouter); err !=nil {
		logger.FatalError(err, "Error occurred, while trying to listen and serve a server")
	}
}