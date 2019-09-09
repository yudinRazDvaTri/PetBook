package main

import (
	"github.com/dpgolang/PetBook/gomigrations"
	"github.com/dpgolang/PetBook/pkg/authentication"
	"github.com/dpgolang/PetBook/pkg/controllers"
	"github.com/dpgolang/PetBook/pkg/driver"
	"github.com/dpgolang/PetBook/pkg/logger"
	"github.com/dpgolang/PetBook/pkg/models"
	"github.com/dpgolang/PetBook/pkg/models/forum"
	"github.com/dpgolang/PetBook/pkg/models/search"
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

	basicRouter := mux.NewRouter()

	storeUser := models.UserStore{DB: db}
	storePet := models.PetStore{DB: db}
	storeRefreshToken := models.RefreshTokenStore{DB: db}
	storeForum := forum.ForumStore{DB: db}
	storeSearch := search.SearchStore{DB: db}
	storeBlog := models.BlogStore{DB: db}
	storeChat := models.ChatStore{DB: db}
	storeFollowers:=models.FollowersStore{DB: db}
	storeMedia := models.MediaStore{DB: db}
	storeVet := models.VetStore{DB: db}


	controller := controllers.Controller{
		PetStore:          &storePet,
		UserStore:         &storeUser,
		ForumStore:        &storeForum,
		SearchStore:       &storeSearch,
		RefreshTokenStore: &storeRefreshToken,
		BlogStore:         &storeBlog,
		ChatStore:         &storeChat,
		FollowersStore:    &storeFollowers,
		MediaStore:        &storeMedia,
		VetStore:         &storeVet,
	}

	basicRouter.HandleFunc("/register", controller.RegisterPostHandler()).Methods("POST")
	basicRouter.HandleFunc("/register", controller.RegisterGetHandler()).Methods("GET")

	basicRouter.HandleFunc("/login", controller.LoginPostHandler()).Methods("POST")
	basicRouter.HandleFunc("/login", controller.LoginGetHandler()).Methods("GET")
	basicRouter.HandleFunc("/loginGoogle", controller.LoginGoogleGetHandler()).Methods("GET")
	basicRouter.HandleFunc("/loginGoogleCallback", controller.GoogleCallback()).Methods("GET")
	basicRouter.HandleFunc("/logout", controller.LogoutGetHandler()).Methods("GET")

	authRouter := basicRouter.PathPrefix("/").Subrouter()
	authRouter.Use(authentication.AuthMiddleware(&storeRefreshToken))

	petRouter := authRouter.PathPrefix("/").Subrouter()
	petRouter.Use(authentication.PetMiddleware(&storeUser))

	//authRouter.HandleFunc("/mypage", controller.MyPageGetHandler()).Methods("GET")
	authRouter.HandleFunc("/petcabinet", controller.PetPostHandler()).Methods("POST")
	authRouter.HandleFunc("/petcabinet", controller.PetGetHandler()).Methods("GET")
  authRouter.HandleFunc("/vetcabinet", controller.VetGetHandler()).Methods("POST")
  authRouter.HandleFunc("/vetcabinet", controller.VetGetHandler()).Methods("GET")
	authRouter.HandleFunc("/search", controller.ViewSearchHandler()).Queries("section", "{section}").Methods("GET")
	authRouter.HandleFunc("/search", controller.RedirectSearchHandler()).Methods("GET")

	authRouter.HandleFunc("/topics", controller.TopicsGetHandler()).Methods("GET")
	authRouter.HandleFunc("/topics", controller.TopicsPostHandler()).Methods("POST")
	authRouter.HandleFunc("/topics/{topicID}", controller.CommentsGetHandler()).Methods("GET")
	authRouter.HandleFunc("/topics/{topicID}/comments", controller.CommentPostHandler()).Methods("POST")
	authRouter.HandleFunc("/topics/{topicID}/comments/{commentID}/ratings", controller.CommentsRatingHandler()).Methods("POST")

	authRouter.HandleFunc("/chats", controller.ChatsGetHandler()).Methods("GET")
	authRouter.HandleFunc("/chats/{id}/delete", controller.DeleteChatHandler()).Methods("POST")
	authRouter.HandleFunc("/chats/{id}", controller.HandleChatConnectionGET()).Methods("GET")
	authRouter.HandleFunc("/ws", controller.HandleChatConnection())
	go controller.HandleMessages()

	//authRouter.HandleFunc("/search", controller.ViewSearchHandler()).Methods("GET")
	authRouter.HandleFunc("/", controller.MyPageGetHandler()).Methods("GET")
  
	authRouter.HandleFunc("/process", controller.CreateBlogHandler()).Methods("POST")
	authRouter.HandleFunc("/delete", controller.DeleteBlogHandler()).Methods("GET")
	authRouter.HandleFunc("/upload", controllers.UploadLogo()).Methods("POST")
	authRouter.HandleFunc("/edit", controller.EditPageHandler).Methods("GET")
	authRouter.HandleFunc("/edit", controller.ProfileUpdateHandler).Methods("POST")
  authRouter.HandleFunc("/uploadmedia", controller.UploadMedia()).Methods("POST")
	authRouter.HandleFunc("/{id}", controller.MyPageOtherUsersHandler()).Methods("GET")
	authRouter.HandleFunc("/deleteimg", controller.DeleteImgHandler()).Methods("Post")
  
  authRouter.HandleFunc("/mypage/{follow:followers|following}", controller.GetFollowerHandler()).Methods("GET")
	authRouter.HandleFunc("/mypage/{follow:followers|following}", controller.PostFollowerHandler()).Methods("POST")

	basicRouter.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./web/static/"))))

	if err := http.ListenAndServe(":" + os.Getenv("APP_PORT"), basicRouter); err != nil {
		logger.FatalError(err, "Error occurred, while trying to listen and serve a server")
	}
}
