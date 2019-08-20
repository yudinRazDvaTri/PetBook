package controllers

import (
	"github.com/gorilla/context"
	"html/template"
	"net/http"
)

func (c *Controller) GetBlogHandler(w http.ResponseWriter, r *http.Request) {
	userID := context.Get(r, "id").(int)
	results := c.BlogStore.GetBlog(userID)
	tmpl, _ := template.ParseFiles("./web/templates/blog.html")
	tmpl.Execute(w, results)
}

func (c *Controller) CreateBlogHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
	if context.Get(r, "id") == nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	id := context.Get(r, "id").(int)
	fn := r.FormValue("something")
	c.BlogStore.CreateBlog(fn, id)
	http.Redirect(w, r, "/mypage", 301)
}

func (c *Controller) DeleteBlogHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Redirect(w, r, "/mypage", 501)
		return
	}
	blogid := r.FormValue("recordid")
	c.BlogStore.DeleteBlog(blogid)
	http.Redirect(w, r, "/mypage", 301)
}
