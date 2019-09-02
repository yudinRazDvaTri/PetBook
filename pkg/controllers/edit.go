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
}

type UserChar struct {
	Name string
	Email string
	Password string
}

func (c *Controller) EditHandler(w http.ResponseWriter, r *http.Request) {
	id := context.Get(r, "id").(int)
	user,_:=c.UserStore.GetUser(id)
	var usChar UserChar
	usChar.Name=user.Firstname
	usChar.Email=user.Email
	usChar.Password=user.Password
	path:=c.MediaStore.GetLogo(id)

	pet, _ := c.UserStore.GetPet(id)
	var edit Editstr
	edit.Name=user.Firstname
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

	view.GenerateHTML(w,"Settings","navbarBlack")
	view.GenerateHTML(w,edit,"edit")
	view.GenerateHTML(w,nil,"footer")
}
func (c *Controller) UpdateHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	id := context.Get(r, "id").(int)
	if err != nil {
		log.Println(err)
	}
	pet:= &models.Pet{}
	pet.ID=id
	pet.Name=r.FormValue("name")
	pet.PetType=     r.FormValue("animal_type")
	pet.Breed=       r.FormValue("breed")
	pet.Age=      r.FormValue("age")
	pet.Weight =   r.FormValue("weight")
	pet.Gender=      r.FormValue("gender")
	pet.Description= r.FormValue("description")
	c.PetStore.UpdatePet(pet)

	http.Redirect(w, r, "/mypage", 301)
}
