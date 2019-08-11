package controllers

import (
	"github.com/dpgolang/PetBook/pkg/logger"
	"github.com/dpgolang/PetBook/pkg/models"
	"github.com/dpgolang/PetBook/pkg/view"
	"github.com/gorilla/context"
	"net/http"
)

func (c *Controller) PetPutHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		user := &models.User{
			Email: context.Get(r, "email").(string),
		}

		err = c.UserStore.GetUser(user)
		if err != nil {
			logger.Error(err, "Error occurred while trying to register pet.\n")
			http.Redirect(w, r, "/mypage", http.StatusSeeOther)
			return
		}
		pet := &models.Pet{
			Name:        r.FormValue("nickname"),
			PetType:     r.FormValue("pet-type"),
			Breed:       r.FormValue("breed"),
			Age:         r.FormValue("age"),
			Weight:      r.FormValue("weight"),
			Gender:      r.FormValue("gender"),
			Description: r.FormValue("description"),
		}

		err = c.PetStore.RegisterPet(pet)
		if err != nil {
			logger.Error(err, "Error ocurrred while trying to register pet.\n")
		}
		http.Redirect(w, r, "/mypage", http.StatusSeeOther)
	}
}
func (c *Controller) PetGetHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		view.GenerateHTML(w, nil, "cabinetPet")
	}
}
