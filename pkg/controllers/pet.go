package controllers

import (
	"github.com/dpgolang/PetBook/pkg/logger"
	"github.com/dpgolang/PetBook/pkg/models"
	"github.com/dpgolang/PetBook/pkg/view"
	"github.com/gorilla/context"
	"net/http"
)

func (c *Controller) PetPostHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		id := context.Get(r, "id").(int)

		pet := &models.Pet{
			ID:          id,
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
			logger.Error(err, "Error occurred while trying to register pet.\n")
		}
		http.Redirect(w, r, "/mypage", http.StatusSeeOther)
	}
}
func (c *Controller) PetGetHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		view.GenerateHTML(w, nil, "cabinetPet")
	}
}
