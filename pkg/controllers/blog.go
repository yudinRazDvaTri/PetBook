package controllers

import (
	"github.com/dpgolang/PetBook/pkg/logger"
	"github.com/gorilla/context"
	"html/template"
	"net/http"
)


func (c *Controller) GetBlogHandler() http.HandlerFunc{
	return func (w http.ResponseWriter, r *http.Request){
		userID := context.Get(r, "id").(int)
		results,err := c.BlogStore.GetPetBlog(userID)
		if err !=nil{
			logger.Error(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		tmpl, _ := template.ParseFiles("./web/templates/blog.html")
		tmpl.Execute(w, results)
	}

}

func (c *Controller) CreateBlogHandler() http.HandlerFunc{
	return func (w http.ResponseWriter, r *http.Request) {
		id := context.Get(r, "id").(int)
		fn := r.FormValue("something")
		err := c.BlogStore.CreateBlog(fn, id)
		if err!= nil {
			logger.Error(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/", 301)
	}
}

func (c *Controller) DeleteBlogHandler() http.HandlerFunc{
	return func (w http.ResponseWriter, r *http.Request){
		blogid := r.FormValue("recordid")
		err := c.BlogStore.DeleteBlog(blogid)
		if err!= nil {
			logger.Error(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/", 301)
	}
}
