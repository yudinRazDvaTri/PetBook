package handler

import (
	"github.com/dpgolang/PetBook/pkg/models"
	"github.com/dpgolang/PetBook/pkg/view"
	"log"
	"net/http"
)

func (c *Controller) CreatePetPostHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if err != nil {
			log.Println(err)
		}
		pet := &models.Pet{}
		pet.Name = r.FormValue("nickname")
		pet.PetType = r.FormValue("pet-type")
		pet.Breed = r.FormValue("breed")
		pet.Age = r.FormValue("age")
		pet.Weight = r.FormValue("weight")
		pet.Gender = r.FormValue("gender")
		pet.Description = r.FormValue("description")
		err = c.PetStore.RegisterPet(pet)
		if err != nil {
			log.Println(err)
		}
		http.Redirect(w, r, "/main", 301)
	}
}
func (c *Controller) CreatePetGetHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		view.GenerateHTML(w, nil, "cabinetPet")
	}
}
