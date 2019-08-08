package handler

import (
	"fmt"
	"github.com/dpgolang/PetBook/pkg/view"
	"github.com/gorilla/context"
	"net/http"
)

func (c *Controller) MyPageGetHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println(context.Get(r, "email"))
		// if user has pet with filled fields
		view.GenerateHTML(w, nil, "mypage")
		// else generate "cabinetPet"
	}
}
