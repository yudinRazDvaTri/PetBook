package controllers

import (
	"github.com/dpgolang/PetBook/pkg/logger"
	"github.com/dpgolang/PetBook/pkg/models"
	"github.com/dpgolang/PetBook/pkg/view"
	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

type MypageData struct {
	Name        string
	Age         string
	PetType     string
	Breed       string
	Description string
	Weight      string
	Gender      string
	LogoPath    string
}

type BlogData struct {
	BlogData []models.Blog
}

type MypageDataVet struct {
	Name          string
	Category      string
	Qualification string
	Surname       string
	Certificates  string
	LogoPath      string
}

func (c *Controller) MyPageGetHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := context.Get(r, "id").(int)
		role, err := c.UserStore.GetUserRole(userID)
		if err != nil {
			logger.Error(err, "Error occurred while getting user role in main page.\n")
			return
		}
		if role == "pet" {
			c.myPageDisplayPet(userID, w, r)
		} else if role == "vet" {
			c.myPageDisplayVet(userID, w, r)
		}
	}
}

func (c *Controller) myPageDisplayPet(userID int, w http.ResponseWriter, r *http.Request) {

	pet, err := c.UserStore.GetPet(userID)
	if err != nil {
		logger.Error(err)
		http.Redirect(w, r, "/petcabinet", http.StatusFound)
	}
	path, err := c.MediaStore.GetLogo(userID)
	if err != nil {
		logger.Error(err, "Error occurred while getting user gallery.\n")
	}
	var myPageData MypageData
	myPageData.LogoPath = path
	myPageData.Name = pet.Name
	myPageData.Age = pet.Age
	myPageData.PetType = pet.PetType
	myPageData.Weight = pet.Weight
	myPageData.Description = pet.Description
	myPageData.Gender = pet.Gender
	myPageData.Breed = pet.Breed

	blog, err := c.BlogStore.GetPetBlog(userID)
	for i := 0; i < len(blog); i++ {
		blog[i].LogoPath = path
	}
	if err != nil {
		logger.Error(err)
		return
	}
	photos, err := c.MediaStore.GetExistedGallery(userID)
	if err != nil {
		logger.Error(err, "Error occurred while getting user gallery.\n")
	}
	view.GenerateHTML(w, "My page", "navbar")
	view.GenerateHTML(w, myPageData, "mypage")
	view.GenerateHTML(w, photos, "gallery_main")
	view.GenerateTimeHTML(w, blog, "blog")
	view.GenerateHTML(w, nil, "footer")

}

func (c *Controller) myPageDisplayVet(userID int, w http.ResponseWriter, r *http.Request) {
	vet, err := c.UserStore.GetVet(userID)
	if err != nil {
		logger.Error(err)
		http.Redirect(w, r, "/vetcabinet", http.StatusFound)
		return
	}
	path, err := c.MediaStore.GetLogo(userID)
	if err != nil {
		logger.Error(err)
		http.Redirect(w, r, "/vetcabinet", http.StatusFound)
		return
	}
	var myPageData MypageDataVet
	myPageData.LogoPath = path
	myPageData.Name = vet.Name
	myPageData.Surname = vet.Surname
	myPageData.Category = vet.Category
	myPageData.Certificates = vet.Certificates
	myPageData.Qualification = vet.Qualification

	blog, err := c.BlogStore.GetVetBlog(userID)
	for i := 0; i < len(blog); i++ {
		blog[i].LogoPath = path
	}
	if err != nil {
		logger.Error(err)
		return
	}
	photos, err := c.MediaStore.GetExistedGallery(userID)
	if err != nil {
		logger.Error(err, "Error occurred while getting user gallery.\n")
	}
	view.GenerateHTML(w, "My page", "navbar")
	view.GenerateHTML(w, myPageData, "mypage_vet")
	view.GenerateHTML(w, photos, "gallery_main")
	view.GenerateTimeHTML(w, blog, "blog")
	view.GenerateHTML(w, nil, "footer")

}

func (c *Controller) DisplayOtherUsersHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		petsId := vars["id"]
		id, err := strconv.Atoi(petsId)
		if err != nil {
			logger.Error(err)
			http.Error(w, "inappropriate request", http.StatusBadRequest)
			return
		}
		role, err := c.UserStore.GetUserRole(id)
		if err != nil {
			logger.Error(err)
			http.Redirect(w, r, "/", http.StatusNotFound)
			return
		}
		if role == "pet" {
			c.myPageDisplayPet(id, w, r)
		} else if role == "vet" {
			c.myPageDisplayVet(id, w, r)
		}
	}

}
