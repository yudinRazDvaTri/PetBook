package controllers

import (
	"github.com/dpgolang/PetBook/pkg/logger"
	"github.com/dpgolang/PetBook/pkg/models"
	"github.com/dpgolang/PetBook/pkg/utilerr"
	"github.com/dpgolang/PetBook/pkg/view"
	"github.com/gorilla/context"
	"net/http"
	"regexp"
)


func (c *Controller) PetPostHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		role := context.Get(r, "role").(string)
		if role != "pet" {
			http.Redirect(w, r, "/mypage", http.StatusMethodNotAllowed)
			return
		}

		err := r.ParseForm()
		id := context.Get(r, "id").(int)
		if matched, err := regexp.Match(patternOnlyNum, []byte(r.FormValue("age"))); !matched || err != nil {
			if err != nil {
				logger.Error(err, "Error occurred while trying to match login.\n")
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			http.Redirect(w, r, "/petcabinet", http.StatusSeeOther)
			return
		}
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
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/mypage", http.StatusSeeOther)
	}
}

func (c *Controller) PetGetHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		role := context.Get(r, "role").(string)
		if role != "pet" {
			http.Redirect(w, r, "/mypage", http.StatusMethodNotAllowed)
			return
		}

		userId := context.Get(r, "id").(int)
		_, err := c.UserStore.GetPet(userId)
		if err == nil {
			http.Redirect(w, r, "/mypage", http.StatusSeeOther)
			return
		}

		switch e := err.(type) {
		case *utilerr.PetDoesNotExist:
			petType := c.PetStore.GetPetEnums()
			view.GenerateHTML(w, petType, "cabinetPet")
			return
		default:
			logger.Error(e)
			http.Error(w, e.Error(), http.StatusInternalServerError)
			return
		}
	}
}
