What is this about?
Our team is glad to represent you PetBook – a social network for our beloved pets: cats, dogs, hamsters and even iguanas! However those aren’t real pets – those are pet’s owners, extraordinary disguised by animal avatar. With PetBook you will  make plenty of new friendships, fulfil your feed with pet’s photos and funny captions, create, promote/keep up your own pet blog, discuss any topic on a forum according to your concern/problem. 
The main purpose of our application is to give an opportunity to all pet owners to communicate not only while going for a walk with their pets but in the worldwide network too.


Provides to pet-owners an ability to lead a social life of their pets just like people do in social networks.

# Have to add .env file with settings like those :

HOST_POSTGRES=db  
PORT_POSTGRES=5432  
POSTGRES_USER=postgres  
POSTGRES_PASSWORD=root  
POSTGRES_DB=pet-book  
APP_ADDRESS=http://localhost  
APP_PORT=8080  
GOOGLE_CALLBACK=loginGoogleCallback  
GOOGLE_CLIENT_ID=client_id  
GOOGLE_CLIENT_SECRET=client_secret  

# How to use middleware and context?
There are four routers: `basicRouter`, `authenticateRouter`, `authorizeRouter` and `petOrVetRouter`.

`basicRouter` is a basic router, which is not wrapped into any middleware and handles user login, registration and logout.

`authenticateRouter` is wrapped into `AuthenticateMiddleware` middleware. 
This router takes one parameter with type `*models.RefreshTokenStore`, which
is responsible for an access to database and operations on corresponding table.
This middleware checks whether user is authenticated or not and restricts users, who haven't logged in (haven't entered email and password), from endpoints
`authenticateRouter := basicRouter.PathPrefix("/").Subrouter()`
`authenticateRouter.Use(authentication.AuthenticateMiddleware(&storeRefreshToken))`

`authorizeRouter` is wrapped into `AuthorizeMiddleware` middleware.  
This router takes one parameter with type `*models.UserStore`.  
`AuthorizeMiddleware` allows user to request some endpoints if he chose a role.
In other way, user will be redirected to role page.
 
`petOrVetRouter` is a subrouter of `authorizeRouter` and it is wrapped into `petOrVetRouter`.  
The point of this router is to restrict users from certain endpoints that are only allowed to users, who have registered their entities (pet or vet).
This middleware is used with AuthenticateMiddleware and AuthorizeMiddleware.  
PetOrVetMiddleware is the maximum degree of authorization in this application.  
It takes one parameter with type `*models.UserStore`.

If you want to wrap your handler in middleware, you should simply handle functions on router like this:  
`authorizeRouter.HandleFunc("/mypage", controller.MyPageGetHandler()).Methods("GET")`

Context is passed to handler from middleware and there are two keys - `id`, `role`.   
`id` is an int value and it represents an id of logged in user.  
`role` is a string value and it representes th role of user.  
Get user_id in your handlers like this:  
`context.Get(r, "id").(int)`
