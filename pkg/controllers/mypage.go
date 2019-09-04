package controllers

import (
	"github.com/dpgolang/PetBook/pkg/logger"
	"github.com/dpgolang/PetBook/pkg/models"
	"github.com/dpgolang/PetBook/pkg/view"
	"github.com/gorilla/context"
	"net/http"
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

func (c *Controller) MyPageGetHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := context.Get(r, "id").(int)
		pet, err := c.UserStore.GetPet(userID)
		path:=c.MediaStore.GetLogo(userID)
		var myPageData MypageData

		myPageData.LogoPath=path
		myPageData.Name=pet.Name
		myPageData.Age=pet.Age
		myPageData.PetType=pet.PetType
		myPageData.Weight=pet.Weight
		myPageData.Description=pet.Description
		myPageData.Gender=pet.Gender
		myPageData.Breed=pet.Breed

		blog,err := c.BlogStore.GetBlog(userID)
		for i := 0; i < len(blog); i++{
			blog[i].LogoPath=path
		}
		if err!=nil{
			logger.Error(err)
			return
		}
		if err != nil {
			http.Redirect(w, r, "/petcabinet", http.StatusFound)
			return
		}
		view.GenerateHTML(w, "MYPAGE", "navbarBlack")
		view.GenerateHTML(w, myPageData, "mypage")
		view.GenerateHTML(w, nil, "gallery_main")
		view.GenerateTimeHTML(w, blog, "blog")
		view.GenerateHTML(w, nil, "footer")
	}

}

