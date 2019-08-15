package controllers

import (
	"github.com/dpgolang/PetBook/pkg/logger"
	"github.com/dpgolang/PetBook/pkg/view"
	"net/http"
)

func (c *Controller) ViewSearchHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		pets, err := c.SearchStore.GetAllPets()
		if err != nil {
			logger.Error(err)
		}

		view.GenerateHTML(w, pets, "searchAnimals")
	}
}
func (c *Controller) SearchHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if err != nil {
			logger.Error(err)
		}
		email := r.FormValue("search")
		pet, err := c.SearchStore.GetByUser(email)
		if err != nil {
			logger.Error(err)
		}
		view.GenerateHTML(w, pet, "viewAnimal")
	}
}
