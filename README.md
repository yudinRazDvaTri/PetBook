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

Context is passed to handler from middleware and there are two keys - `id` and `pet`. 
`id` is an int value and it represents an id of logged in user.
Get user_id in your handlers like this:
`context.Get(r, "id").(int)`

`pet` is a boolean value and it represents whether user has a pet or not.
Get user_id in your handlers like this:
`context.Get(r, "pet").(bool)`

### You have to check the value of `pet` in your `GET` and `POST` handlers! Users who didn't fill information about their pets
can't create new topics or leave comments but they still can look at existing topics, comments etc.
For example, you can pass value of `pet` from `GET` handler to template and check it there. According to value display needed content.
In `POST` handlers, which work with user's pet, you have to return from function if `pet` value is false.