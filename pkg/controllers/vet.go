package controllers

import (
	"github.com/dpgolang/PetBook/pkg/logger"
	"github.com/dpgolang/PetBook/pkg/models"
	"github.com/dpgolang/PetBook/pkg/view"
	"github.com/gorilla/context"
	"net/http"
	"regexp"
)

// TODO: check input values
func (c *Controller) VetPostHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		id := context.Get(r, "id").(int)
		if matched, err := regexp.Match(patternOnlyNum, []byte(r.FormValue("age"))); !matched || err != nil {
			if err != nil {
				logger.Error(err, "Error occurred while trying to match login.\n")
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			http.Redirect(w, r, "/vetcabinet", http.StatusSeeOther)
			return
		}
		vet := &models.Vet{
			ID:          id,
			Name:        r.FormValue("nickname"),
			Qualification:     r.FormValue("qualification"),
			Surname:       r.FormValue("surname"),
			Category:         r.FormValue("category"),
			Certificates:      r.FormValue("certificates"),
		}
		err = c.VetStore.RegisterVet(vet)
		if err != nil {
			logger.Error(err, "Error occurred while trying to register pet.\n")
		}
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

func (c *Controller) VetGetHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vetType,_ := c.VetStore.GetVetEnums()
		view.GenerateHTML(w, vetType, "cabinetVet")
	}
}