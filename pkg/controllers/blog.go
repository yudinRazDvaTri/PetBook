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
			return
		}

		tmpl, _ := template.ParseFiles("./web/templates/blog.html")
		tmpl.Execute(w, results)
	}

}

func (c *Controller) CreateBlogHandler() http.HandlerFunc{
	//if r.Method != http.MethodPost {
	//	http.Redirect(w, r, "/", http.StatusFound)
	//	return
	//}
	return func (w http.ResponseWriter, r *http.Request) {
		//var err error
		//if context.Get(r, "id") == nil {
		//	logger.Error(err)
		//	http.Redirect(w, r, "/login", http.StatusSeeOther)
		//	return
		//}
		id := context.Get(r, "id").(int)
		fn := r.FormValue("something")
		c.BlogStore.CreateBlog(fn, id)
		http.Redirect(w, r, "/", 301)
	}
}

func (c *Controller) DeleteBlogHandler() http.HandlerFunc{
	return func (w http.ResponseWriter, r *http.Request){
		//if r.Method != http.MethodGet {
		//	http.Redirect(w, r, "/mypage", 301)
		//	return
		//}
		blogid := r.FormValue("recordid")
		c.BlogStore.DeleteBlog(blogid)
		http.Redirect(w, r, "/", 301)
	}
}
