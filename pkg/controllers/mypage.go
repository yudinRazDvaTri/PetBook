package controllers

import (
	"github.com/dpgolang/PetBook/pkg/logger"
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
	LogoPath []string

}

func (c *Controller) MyPageGetHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := context.Get(r, "id").(int)
		pet, err := c.UserStore.GetPet(userID)
		path:=GetImgLogo(userID)

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
		if err!=nil{
			logger.Error(err)
			return
		}
		//for i, v := range blog {
		//	v=blog[i]
		//	v.LogoPath=path[0]
		//}
		if err != nil {
			http.Redirect(w, r, "/petcabinet", http.StatusFound)
			return
		}
		view.GenerateHTML(w, "MYPAGE", "navbarBlack")
		view.GenerateHTML(w, myPageData, "mypage")
		view.GenerateTimeHTML(w, blog, "blog")
		view.GenerateHTML(w, nil, "footer")
	}

}

