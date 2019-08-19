package controllers

import (
	"github.com/dpgolang/PetBook/pkg/logger"
	"github.com/dpgolang/PetBook/pkg/view"
	"net/http"
)

func (c *Controller) ViewSearchHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		section := r.URL.Query().Get("section")
		if section == "auto" {

		}
		if section == "user" {
			email := r.URL.Query().Get("email")
			if email != "" {
				pet, err := c.SearchStore.GetByUser(email)
				if err != nil {
					logger.Error(err)
				}
				view.GenerateHTML(w, pet, "view_animal")
			} else {
				pets, err := c.SearchStore.GetAllPets()
				if err != nil {
					logger.Error(err)
				}
				view.GenerateHTML(w, pets, "search_by_user")
			}

		}
		if section == "animal" {
			view.GenerateHTML(w, "Forum", "navbar")
			m := make(map[string]string)
			age := r.URL.Query().Get("age")
			if age != "" {
				m["age"] = age
			}
			animalType := r.URL.Query().Get("type")
			if animalType != "" {
				m["animal_type"] = animalType
			}
			breed := r.URL.Query().Get("breed")
			if breed != "" {
				m["breed"] = breed
			}
			weight := r.URL.Query().Get("weight")
			if weight != "" {
				m["weight"] = weight
			}
			gender := r.URL.Query().Get("gender")
			if gender != "" {
				m["gender"] = gender
			}
			name := r.URL.Query().Get("name")
			if name != "" {
				m["name"] = name
			}
			if len(m) == 0 {
				pets, err := c.SearchStore.GetAllPets()
				if err != nil {
					logger.Error(err)
				}
				view.GenerateHTML(w, pets, "search_by_animals")
			} else {
				filterPets, err := c.SearchStore.GetFilterPets(m)
				if err != nil {
					logger.Error(err)
				}
				view.GenerateHTML(w, filterPets, "search_by_animals")
			}
		}
		if section == "forum" {

		}
	}
}

