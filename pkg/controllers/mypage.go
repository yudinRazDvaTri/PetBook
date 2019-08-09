package controllers

import (
	"github.com/dpgolang/PetBook/pkg/view"
	"net/http"
)

func (c *Controller) MyPageGetHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		view.GenerateHTML(w, nil, "mypage")
	}
}
