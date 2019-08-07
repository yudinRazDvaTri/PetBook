package handler

import (
	"fmt"
	"github.com/Khudienko/PetBook/pkg/tokens"
	"github.com/gorilla/context"
	"net/http"
)

func (c *Controller) MyPageGetHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println(context.Get(r, "email"))
		// if user has pet with filled fields
		tokens.GenerateHTML(w, nil, "mypage")
		// else generate "cabinetPet"
	}
}
