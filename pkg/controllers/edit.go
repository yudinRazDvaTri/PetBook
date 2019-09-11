package controllers

import (
	"github.com/dpgolang/PetBook/pkg/logger"
	"github.com/dpgolang/PetBook/pkg/models"
	"github.com/dpgolang/PetBook/pkg/utilerr"
	"github.com/dpgolang/PetBook/pkg/view"
	"github.com/gorilla/context"
	"net/http"
)

type Editstr struct {
	Pet      models.Pet
	Vet      models.Vet
	LogoPath string
}
//Handler displays the user settings page, depending on the role
// in which he can upload the logo, media, change characteristics
func (c *Controller) EditPageHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := context.Get(r, "id").(int)
		user, err := c.UserStore.GetUser(id)
		if err != nil {
			logger.Error(err)
			http.Error(w, "can't get user", http.StatusInternalServerError)
			return
		}

		path, err := c.MediaStore.GetLogo(id)
		if err != nil {
			switch e := err.(type) {
			case *utilerr.LogoDoesNotExist:
				break
			default:
				logger.Error(e)
				http.Redirect(w, r, err.Error(), http.StatusInternalServerError)
				return
			}
		}

		var filename string
		var edit Editstr

		if user.Role == "pet" {
			pet, _ := c.UserStore.GetPet(id)
			edit.Pet = pet
			edit.LogoPath = path
			filename = "edit"
		} else if user.Role == "vet" {
			vet, _ := c.UserStore.GetVet(id)
			edit.Vet = vet
			edit.LogoPath = path
			filename = "editVet"
		}

		gallery, err := c.MediaStore.GetExistedGallery(id)
		if err != nil {
			logger.Error(err, "Error occurred while getting user gallery.\n")
		}
		view.GenerateHTML(w, "Settings", "navbar")
		view.GenerateHTML(w, edit, filename)
		view.GenerateHTML(w, gallery, "gallery_edit")
		view.GenerateHTML(w, nil, "footer")
	}
}
//Handler for updating\recording the characteristics of the animal\veterinarian depending on his role
func (c *Controller) ProfileUpdateHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := context.Get(r, "id").(int)
		err := r.ParseForm()
		if err != nil {
			logger.Error(err, "Error occurred while getting user gallery.\n")
		}
		user, err := c.UserStore.GetUser(id)
		if err != nil {
			logger.Error(err)
			http.Error(w, "can't get user", http.StatusInternalServerError)
			return
		}
		if user.Role == "pet" {
			pet := &models.Pet{}
			pet.ID = id
			pet.Name = r.FormValue("name")
			pet.PetType = r.FormValue("animal_type")
			pet.Breed = r.FormValue("breed")
			pet.Age = r.FormValue("age")
			pet.Weight = r.FormValue("weight")
			pet.Gender = r.FormValue("gender")
			pet.Description = r.FormValue("description")
			err := c.PetStore.UpdatePet(pet)
			if err != nil {
				logger.Error(err)
				http.Error(w, "can't update pet", http.StatusInternalServerError)
				return
			}
		} else if user.Role == "vet" {
			vet := &models.Vet{}
			vet.ID = id
			vet.Name = r.FormValue("name")
			vet.Surname = r.FormValue("surname")
			vet.Qualification = r.FormValue("qualification")
			vet.Category = r.FormValue("category")
			vet.Certificates = r.FormValue("certificates")
			err := c.VetStore.UpdateVet(vet)
			if err != nil {
				logger.Error(err)
				http.Error(w, "can't update Vet", http.StatusInternalServerError)
				return
			}
		}
		http.Redirect(w, r, "/", 301)
	}
}
