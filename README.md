Provides to pet-owners an ability to lead a social life of their pets just like people do in social networks.

# Have to add .env file with settings like those :

HOST_POSTGRES=localhost  
PORT_POSTGRES=5434  
POSTGRES_USER=postgres  
POSTGRES_PASSWORD=root  
POSTGRES_DB=pet-book  
APP_ADDRESS=http://localhost  
APP_PORT=8080  
GOOGLE_CALLBACK=loginGoogleCallback  
GOOGLE_CLIENT_ID=client_id  
GOOGLE_CLIENT_SECRET=client_secret  

# How to use middleware and context?
There are three routers: `basicRouter`, `authRouter` and `petRouter`.

`basicRouter` router is a really basic router, which is not wrapped into any middleware and handles user login, registration and logout.

`authRouter` router is wrapped into `AuthMiddleware` middleware. 
This router takes one parameter with type `*models.RefreshTokenStore`, which
is responsible for an access to database and operations on corresponding table.
This middleware checks whether user is authenticated or not.  
`authRouter := router.PathPrefix("/").Subrouter()`
`authRouter.Use(authentication.AuthMiddleware(&storeRefreshToken, &storeUser))`
 
 `petRouter` is a subrouter of `authRouter` and it is wrapped into `PetMiddleware`.
The point of this router is to restrict users from certain endpoints that are only allowed to users, who have registered their pets.
It takes one parameter with type `*models.UserStore`.

If you want to wrap your handler in middleware, you should simply handle functions on router like this:  
`authRouter.HandleFunc("/mypage", controller.MyPageGetHandler()).Methods("GET")`

Context is passed to handler from middleware and there is one key - `id`. 
`id` is an int value and it represents an id of logged in user.
Get user_id in your handlers like this:  
`context.Get(r, "id").(int)`


