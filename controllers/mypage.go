package controllers

import (
	"fmt"
	"github.com/gorilla/context"
	"net/http"
	"test/pkg/utils"
)

func (c *Controller) MyPageGetHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println(context.Get(r, "email"))
		// if user has pet with filled fields
		utils.GenerateHTML(w, nil, "mypage")
		// else generate "cabinetPet"
	}
}
