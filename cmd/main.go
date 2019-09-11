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

	// Connecting to database and applying migrations.
	db := driver.ConnectDB()
	err := gomigrations.Migrate(db)
	if err != nil {
		logger.FatalError(err, "Migration failed.\n")
	}

	// Basic router, which is not wrapped into any middleware. Used for register, login and logout endpoints.
	basicRouter := mux.NewRouter()

	// Initializing models.
	storeUser := models.UserStore{DB: db}
	storePet := models.PetStore{DB: db}
	storeRefreshToken := models.RefreshTokenStore{DB: db}
	storeForum := forum.ForumStore{DB: db}
	storeSearch := search.SearchStore{DB: db}
	storeBlog := models.BlogStore{DB: db}
	storeChat := models.ChatStore{DB: db}
	storeFollowers := models.FollowersStore{DB: db}
	storeMedia := models.MediaStore{DB: db}
	storeVet := models.VetStore{DB: db}

	// Initializing controller.
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
		VetStore:          &storeVet,
	}

	// Endpoints.
	basicRouter.HandleFunc("/register", controller.RegisterPostHandler()).Methods("POST")
	basicRouter.HandleFunc("/register", controller.RegisterGetHandler()).Methods("GET")

	basicRouter.HandleFunc("/login", controller.LoginPostHandler()).Methods("POST")
	basicRouter.HandleFunc("/login", controller.LoginGetHandler()).Methods("GET")
	basicRouter.HandleFunc("/loginGoogle", controller.LoginGoogleGetHandler()).Methods("GET")
	basicRouter.HandleFunc("/loginGoogleCallback", controller.GoogleCallback()).Methods("GET")
	basicRouter.HandleFunc("/logout", controller.LogoutGetHandler()).Methods("GET")

	// Subrouter of basicRouter, which is wrapped into AuthenticateMiddleware.
	// You can access to endpoints of this router if you are logged in and selected a role.
	authenticateRouter := basicRouter.PathPrefix("/").Subrouter()
	authenticateRouter.Use(authentication.AuthenticateMiddleware(&storeRefreshToken))

	authenticateRouter.HandleFunc("/role", controller.RoleGetHandler()).Methods("GET")
	authenticateRouter.HandleFunc("/role", controller.RolePostHandler()).Methods("POST")

	// Subrouter of authenticateRouter, which is wrapped into AuthorizeMiddleware.
	// You can access to endpoints of this router if you are logged in, selected a role and filled information about you (pet or vet).
	authorizeRouter := authenticateRouter.PathPrefix("/").Subrouter()
	authorizeRouter.Use(authentication.AuthorizeMiddleware(&storeUser))

	authorizeRouter.HandleFunc("/petcabinet", controller.PetPostHandler()).Methods("POST")
	authorizeRouter.HandleFunc("/petcabinet", controller.PetGetHandler()).Methods("GET")
	authorizeRouter.HandleFunc("/vetcabinet", controller.VetPostHandler()).Methods("POST")
	authorizeRouter.HandleFunc("/vetcabinet", controller.VetGetHandler()).Methods("GET")
	authorizeRouter.HandleFunc("/search", controller.ViewSearchHandler()).Queries("section", "{section}").Methods("GET")
	authorizeRouter.HandleFunc("/search", controller.RedirectSearchHandler()).Methods("GET")

	authorizeRouter.HandleFunc("/topics", controller.TopicsGetHandler()).Methods("GET")
	authorizeRouter.HandleFunc("/topics/{topicID}", controller.CommentsGetHandler()).Methods("GET")

	authorizeRouter.HandleFunc("/", controller.MyPageGetHandler()).Methods("GET")
	authorizeRouter.HandleFunc("/mypage", controller.MyPageGetHandler()).Methods("GET")

	authorizeRouter.HandleFunc("/media/logo", controller.UploadLogo()).Methods("POST")
	authorizeRouter.HandleFunc("/mypage/edit", controller.EditPageHandler()).Methods("GET")
	authorizeRouter.HandleFunc("/mypage/edit", controller.ProfileUpdateHandler()).Methods("POST")
	authorizeRouter.HandleFunc("/media/gallery", controller.UploadMedia()).Methods("POST")
	authorizeRouter.HandleFunc("/users/{id}/", controller.DisplayOtherUsersHandler()).Methods("GET")
	authorizeRouter.HandleFunc("/media/{id}/delete", controller.DeleteImgHandler()).Methods("Post")

	petOrVetRouter := authorizeRouter.PathPrefix("/").Subrouter()
	petOrVetRouter.Use(authentication.PetOrVetMiddleware(&storeUser))

	petOrVetRouter.HandleFunc("/blogs", controller.CreateBlogHandler()).Methods("POST")
	petOrVetRouter.HandleFunc("/blogs/{id}", controller.DeleteBlogHandler()).Methods("GET")

	petOrVetRouter.HandleFunc("/chats", controller.ChatsGetHandler()).Methods("GET")
	petOrVetRouter.HandleFunc("/chats/{id}", controller.DeleteChatHandler()).Methods("POST") //does not work with method DELETE with overriding with js too
	petOrVetRouter.HandleFunc("/chats/{id}", controller.HandleChatConnectionGET()).Methods("GET")
	petOrVetRouter.HandleFunc("/chats/{id}/search/{date}", controller.HandleChatSearchConnection()).Methods("GET")
	petOrVetRouter.HandleFunc("/ws/{id}", controller.HandleChatConnection())

	petOrVetRouter.HandleFunc("/mypage/{follow:followers|following}", controller.PostFollowerHandler()).Methods("POST")

	petOrVetRouter.HandleFunc("/topics", controller.TopicsPostHandler()).Methods("POST")
	petOrVetRouter.HandleFunc("/topics/{topicID}/comments", controller.CommentPostHandler()).Methods("POST")
	petOrVetRouter.HandleFunc("/topics/{topicID}/comments/{commentID}/ratings", controller.CommentsRatingHandler()).Methods("POST")

	petOrVetRouter.HandleFunc("/mypage/{follow:followers|following}", controller.GetFollowerHandler()).Methods("GET")
	petOrVetRouter.HandleFunc("/mypage/{follow:followers|following}", controller.PostFollowerHandler()).Methods("POST")

	// Setting path for static files (.css, .js, pictures etc)
	basicRouter.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./web/static/"))))

	if err := http.ListenAndServe(":"+os.Getenv("APP_PORT"), basicRouter); err != nil {
		logger.FatalError(err, "Error occurred, while trying to listen and serve a server")
	}
}
