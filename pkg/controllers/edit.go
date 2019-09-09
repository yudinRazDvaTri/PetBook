package controllers

import (
	"github.com/dpgolang/PetBook/pkg/models"
	"github.com/dpgolang/PetBook/pkg/view"
	"github.com/gorilla/context"
	"log"
	"net/http"
)

type Editstr struct {
	Name string
	Email string
	Password string
	PetName string
	Age string
	PetType string
	Breed string
	Description string
	Weight string
	Gender string
	LogoPath string
	VetName string
	Category string
	Qualification string
	Surname       string
	Certificates string
}

func (c *Controller) EditPageHandler(w http.ResponseWriter, r *http.Request) {
	id := context.Get(r, "id").(int)
	user,_:=c.UserStore.GetUser(id)
	path:=c.MediaStore.GetLogo(id)
	var filename string
	var edit Editstr

	if user.Role=="pet"{
		pet, _ := c.UserStore.GetPet(id)
		edit.Name=user.Login
		edit.Email=user.Email
		edit.Password=user.Password
		edit.PetName=pet.Name
		edit.Age=pet.Age
		edit.PetType=pet.PetType
		edit.Breed=pet.Breed
		edit.Description=pet.Description
		edit.Weight=pet.Weight
		edit.Gender=pet.Gender
		edit.LogoPath=path
		filename="edit"
	}else if user.Role=="vet" {
		vet,_:=c.UserStore.GetVet(id)
		edit.Name=user.Login
		edit.Email=user.Email
		edit.Password=user.Password
		edit.VetName=vet.Name
		edit.Surname=vet.Surname
		edit.Qualification=vet.Qualification
		edit.Category=vet.Category
		edit.Certificates=vet.Certificates
		edit.LogoPath=path
		filename="editVet"
	}

	gallery:=c.MediaStore.GetExistedGallery(id)

	view.GenerateHTML(w,"Settings","navbar")
	view.GenerateHTML(w,edit,filename)
	view.GenerateHTML(w,gallery,"gallery_edit")
	view.GenerateHTML(w,nil,"footer")
}

func (c *Controller) ProfileUpdateHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	id := context.Get(r, "id").(int)
	if err != nil {
		log.Println(err)
	}
	user,_:=c.UserStore.GetUser(id)
	if user.Role=="pet" {
		pet := &models.Pet{}
		pet.ID = id
		pet.Name = r.FormValue("name")
		pet.PetType = r.FormValue("animal_type")
		pet.Breed = r.FormValue("breed")
		pet.Age = r.FormValue("age")
		pet.Weight = r.FormValue("weight")
		pet.Gender = r.FormValue("gender")
		pet.Description = r.FormValue("description")
		c.PetStore.UpdatePet(pet)
	}else if user.Role=="vet"{
		vet := &models.Vet{}
		vet.ID = id
		vet.Name = r.FormValue("name")
		vet.Surname = r.FormValue("surname")
		vet.Qualification = r.FormValue("qualification")
		vet.Category = r.FormValue("category")
		vet.Certificates = r.FormValue("certificates")
		c.VetStore.UpdateVet(vet)
	}

	http.Redirect(w, r, "/", 301)
}

//func (c *Controller) ProfileUpdateHandler(w http.ResponseWriter, r *http.Request) {
//	err := r.ParseForm()
//	id := context.Get(r, "id").(int)
//	if err != nil {
//		log.Println(err)
//	}
//	pet:= &models.Pet{}
//	pet.ID=id
//	pet.Name=r.FormValue("name")
//	pet.PetType=     r.FormValue("animal_type")
//	pet.Breed=       r.FormValue("breed")
//	pet.Age=      r.FormValue("age")
//	pet.Weight =   r.FormValue("weight")
//	pet.Gender=      r.FormValue("gender")
//	pet.Description= r.FormValue("description")
//	c.PetStore.UpdatePet(pet)
//
//	http.Redirect(w, r, "/", 301)
//}
