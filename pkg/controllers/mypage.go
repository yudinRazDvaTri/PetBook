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
	Name string
	Age string
	PetType string
	Breed string
	Description string
	Weight string
	Gender string
	LogoPath string

}

type BlogData struct {
	BlogData []models.Blog
}

type MypageDataVet struct {
	Name string
	Category string
	Qualification string
	Surname       string
	Certificates string
	LogoPath string
}

func (c *Controller) MyPageGetHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := context.Get(r, "id").(int)
		role:=c.UserStore.GetUserRole(userID)
		if role=="pet"{
			c.MyPagePetGetHandler(userID,w,r)
		}else if role=="vet"{
			c.MyPageVetGetHandler(userID,w,r)
		}

	}
}

func (c *Controller) MyPagePetGetHandler(userID int,w http.ResponseWriter, r *http.Request) {

		pet, err := c.UserStore.GetPet(userID)
		path := c.MediaStore.GetLogo(userID)
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
		if err != nil {
			http.Redirect(w, r, "/petcabinet", http.StatusFound)
			return
		}

		photos := c.MediaStore.GetExistedGallery(userID)

		view.GenerateHTML(w, "My page", "navbar")
		view.GenerateHTML(w, myPageData, "mypage")
		view.GenerateHTML(w, photos, "gallery_main")
		view.GenerateTimeHTML(w, blog, "blog")
		view.GenerateHTML(w, nil, "footer")

}

func (c *Controller) MyPageVetGetHandler(userID int,w http.ResponseWriter, r *http.Request) {
		vet, err := c.UserStore.GetVet(userID)
		path := c.MediaStore.GetLogo(userID)
		var myPageData MypageDataVet
		myPageData.LogoPath = path
		myPageData.Name = vet.Name
		myPageData.Surname=vet.Surname
		myPageData.Category=vet.Category
		myPageData.Certificates=vet.Certificates
		myPageData.Qualification=vet.Qualification


		blog, err := c.BlogStore.GetVetBlog(userID)
		for i := 0; i < len(blog); i++ {
			blog[i].LogoPath = path
		}
		if err != nil {
			logger.Error(err)
			return
		}
		if err != nil {
			http.Redirect(w, r, "/vetcabinet", http.StatusFound)
			return
		}

		photos := c.MediaStore.GetExistedGallery(userID)

		view.GenerateHTML(w, "My page", "navbar")
		view.GenerateHTML(w, myPageData, "mypage_vet")
		view.GenerateHTML(w, photos, "gallery_main")
		view.GenerateTimeHTML(w, blog, "blog")
		view.GenerateHTML(w, nil, "footer")
}

func (c *Controller) MyPageOtherUsersHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		petsId := vars["id"]
		id,err := strconv.Atoi(petsId)
		role:=c.UserStore.GetUserRole(id)
		if role=="pet"{
			c.MyPagePetGetHandler(id,w,r)
		}else if role=="vet"{
			c.MyPageVetGetHandler(id,w,r)
		}
		if err != nil {
			http.Redirect(w, r, "/", http.StatusFound)
			return
		}

	}

}


//func (c *Controller) MyPageGetHandler() http.HandlerFunc {
//	return func(w http.ResponseWriter, r *http.Request) {
//		userID := context.Get(r, "id").(int)
//		pet, err := c.UserStore.GetPet(userID)
//		path:=c.MediaStore.GetLogo(userID)
//		var myPageData MypageData
//
//		myPageData.LogoPath=path
//		myPageData.Name=pet.Name
//		myPageData.Age=pet.Age
//		myPageData.PetType=pet.PetType
//		myPageData.Weight=pet.Weight
//		myPageData.Description=pet.Description
//		myPageData.Gender=pet.Gender
//		myPageData.Breed=pet.Breed
//
//		blog,err := c.BlogStore.GetPetBlog(userID)
//		for i := 0; i < len(blog); i++{
//			blog[i].LogoPath=path
//		}
//		if err!=nil{
//			logger.Error(err)
//			return
//		}
//		if err != nil {
//			http.Redirect(w, r, "/petcabinet", http.StatusFound)
//			return
//		}
//
//		photos:=c.MediaStore.GetExistedGallery(userID)
//
//		view.GenerateHTML(w, "MYPAGE", "navbarBlack")
//		view.GenerateHTML(w, myPageData, "mypage")
//		view.GenerateHTML(w, photos, "gallery_main")
//		view.GenerateTimeHTML(w, blog, "blog")
//		view.GenerateHTML(w, nil, "footer")
//	}
//
//}
//
//
//func (c *Controller) MyPageOtherUsersHandler() http.HandlerFunc {
//	return func(w http.ResponseWriter, r *http.Request) {
//		vars := mux.Vars(r)
//		petsId := vars["id"]
//		petID, err := strconv.Atoi(petsId)
//		pet, err := c.UserStore.GetPet(petID)
//		path:=c.MediaStore.GetLogo(petID)
//		var myPageData MypageData
//		myPageData.LogoPath=path
//		myPageData.Name=pet.Name
//		myPageData.Age=pet.Age
//		myPageData.PetType=pet.PetType
//		myPageData.Weight=pet.Weight
//		myPageData.Description=pet.Description
//		myPageData.Gender=pet.Gender
//		myPageData.Breed=pet.Breed
//
//		blog,err := c.BlogStore.GetPetBlog(petID)
//		for i := 0; i < len(blog); i++{
//			blog[i].LogoPath=path
//		}
//		if err!=nil{
//			logger.Error(err)
//			return
//		}
//		if err != nil {
//			http.Redirect(w, r, "/petcabinet", http.StatusFound)
//			return
//		}
//		photos:=c.MediaStore.GetExistedGallery(petID)
//		view.GenerateHTML(w, "MYPAGE", "navbarBlack")
//		view.GenerateHTML(w, myPageData, "mypage")
//		view.GenerateHTML(w, photos, "gallery_main")
//		view.GenerateTimeHTML(w, blog, "blog")
//		view.GenerateHTML(w, nil, "footer")
//	}
//
//}