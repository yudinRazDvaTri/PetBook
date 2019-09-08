package controllers

import (
	"github.com/dpgolang/PetBook/pkg/logger"
	"github.com/gorilla/context"
	"html/template"
	"net/http"
	"strconv"
)


func (c *Controller) GetCommentBlogHandler() http.HandlerFunc{
	return func (w http.ResponseWriter, r *http.Request){
		userID := context.Get(r, "id").(int)
		results,err := c.CommentBlogStore.GetPetCommentBlog(userID)
		if err !=nil{
			logger.Error(err)
			return
		}
		tmpl, _ := template.ParseFiles("./web/templates/blog.html")
		tmpl.Execute(w, results)
	}

}

func (c *Controller) CreateCommentBlogHandler() http.HandlerFunc{
	return func (w http.ResponseWriter, r *http.Request) {
		id := context.Get(r, "id").(int)
		blogid := r.FormValue("recordidcom")
		idblog,_ :=strconv.Atoi(blogid)
		fn := r.FormValue("createcommentblog")
		c.CommentBlogStore.CreateCommentBlog(fn, id,idblog)
		http.Redirect(w, r, "/", 301)
	}
}

func (c *Controller) DeleteCommentBlogHandler() http.HandlerFunc{
	return func (w http.ResponseWriter, r *http.Request){
		blogid := r.FormValue("recordid")
		c.CommentBlogStore.DeleteCommentBlog(blogid)
		http.Redirect(w, r, "/", 301)
	}
}

