package controllers

import (
	"github.com/dpgolang/PetBook/pkg/models"
	"github.com/dpgolang/PetBook/pkg/view"
	"github.com/gorilla/context"
	"net/http"
)

func (c *Controller) MyPageGetHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user := models.User{
			Email: context.Get(r,"email").(string),
		}
		_, err := c.UserStore.GetPet(&user)
		if err != nil {
			http.Redirect(w, r, "/petcabinet", http.StatusFound)
			return
		}
		view.GenerateHTML(w, nil, "mypage")
	}
}
