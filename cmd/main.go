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
	storeChat := models.ChatStore{DB: db}
	storeMedia := models.MediaStore{DB: db}

	controller := controllers.Controller{
		PetStore:          &storePet,
		UserStore:         &storeUser,
		ForumStore:        &storeForum,
		SearchStore:       &storeSearch,
		RefreshTokenStore: &storeRefreshToken,
		BlogStore:         &storeBlog,
		ChatStore:         &storeChat,
		MediaStore:        &storeMedia,
	}

	router.HandleFunc("/register", controller.RegisterPostHandler()).Methods("POST")
	router.HandleFunc("/register", controller.RegisterGetHandler()).Methods("GET")

	router.HandleFunc("/login", controller.LoginPostHandler()).Methods("POST")
	router.HandleFunc("/login", controller.LoginGetHandler()).Methods("GET")
	router.HandleFunc("/loginGoogle", controller.LoginGoogleGetHandler()).Methods("GET")
	router.HandleFunc("/loginGoogleCallback", controller.GoogleCallback()).Methods("GET")
	router.HandleFunc("/logout", controller.LogoutGetHandler()).Methods("GET")

	subrouter := router.PathPrefix("/").Subrouter()
	subrouter.Use(authentication.AuthMiddleware(&storeRefreshToken, &storeUser))

	subrouter.HandleFunc("/mypage", controller.MyPageGetHandler()).Methods("GET")
	subrouter.HandleFunc("/petcabinet", controller.PetPostHandler()).Methods("POST")
	subrouter.HandleFunc("/petcabinet", controller.PetGetHandler()).Methods("GET")
	subrouter.HandleFunc("/search", controller.ViewSearchHandler()).Queries("section", "{section}").Methods("GET")
	subrouter.HandleFunc("/search", controller.RedirectSearchHandler()).Methods("GET")

	subrouter.HandleFunc("/topics", controller.TopicsGetHandler()).Methods("GET")
	subrouter.HandleFunc("/topics", controller.TopicsPostHandler()).Methods("POST")
	subrouter.HandleFunc("/topics/{topicID}", controller.CommentsGetHandler()).Methods("GET")
	subrouter.HandleFunc("/topics/{topicID}/comments", controller.CommentPostHandler()).Methods("POST")
	subrouter.HandleFunc("/topics/{topicID}/comments/{commentID}/ratings", controller.CommentsRatingHandler()).Methods("POST")

	subrouter.HandleFunc("/chats", controller.ChatsGetHandler()).Methods("GET")
	subrouter.HandleFunc("/chats/{id}/delete", controller.DeleteChatHandler()).Methods("POST")
	subrouter.HandleFunc("/chats/{id}", controller.HandleChatConnectionGET()).Methods("GET")
	subrouter.HandleFunc("/ws", controller.HandleChatConnection())
	go controller.HandleMessages()

	//subrouter.HandleFunc("/search", controller.ViewSearchHandler()).Methods("GET")
	subrouter.HandleFunc("/", controller.MyPageGetHandler()).Methods("GET")

	subrouter.HandleFunc("/process", controller.CreateBlogHandler()).Methods("POST")
	subrouter.HandleFunc("/delete", controller.DeleteBlogHandler()).Methods("GET")
	subrouter.HandleFunc("/upload", controller.UploadLogo()).Methods("POST")
	subrouter.HandleFunc("/edit", controller.EditPageHandler).Methods("GET")
	subrouter.HandleFunc("/edit", controller.ProfileUpdateHandler).Methods("POST")
	subrouter.HandleFunc("/uploadmedia", controller.UploadMedia()).Methods("POST")


	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/",
		http.FileServer(http.Dir("./web/static/"))))

	loggedRouter := handlers.LoggingHandler(os.Stdout, router)

	if err := http.ListenAndServe(":8080", loggedRouter); err != nil {
		logger.FatalError(err, "Error occurred, while trying to listen and serve a server")
	}
}
