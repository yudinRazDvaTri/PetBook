Provides to pet-owners an ability to lead a social life of their pets just like people do in social networks.

Have to att .env file with settings like those :

HOST_POSTGRES=localhost
PORT_POSTGRES=5434
POSTGRES_USER=postgres
POSTGRES_PASSWORD=root
POSTGRES_DB=pet-book

# How to use middleware and context?
Here is `subrouter`, which is already wrapped in authentcation middleware.
`subrouter := router.PathPrefix("/").Subrouter()
 subrouter.Use(mux.MiddlewareFunc(authentication.Content))`

If you want to wrap your handler in middleware, you should simply handle functions on `subrouter` like this:
`subrouter.HandleFunc("/mypage", controller.MyPageGetHandler()).Methods("GET")`

Context passes to handler from middleware and there is only one key - `id`, which is id of logged in user.
Get user_id in your handlers like this:
`context.Get(r, "id").(int)`