package controllers

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/dpgolang/PetBook/pkg/logger"
	"github.com/dpgolang/PetBook/pkg/models"
	"github.com/dpgolang/PetBook/pkg/view"
	"github.com/gorilla/context"
	"github.com/gorilla/mux"
)

type dataFollow struct {
	UserID        int
	PageAction    string
	IsFollowing   bool
	PetsFollowers []*models.FollowerPets
	PetsFollowing []*models.FollowerPets
}

//Get all subscribers and animals that this user is subscribed to
func (c *Controller) GetFollowerHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var Data dataFollow
		params := mux.Vars(r)
		follow := params["follow"]
		userID := context.Get(r, "id").(int)
		petsFollowers, err := c.FollowersStore.GetFollowers(userID)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			logger.Error(err)
			return
		}
		petsFollowing, err := c.FollowersStore.GetFollowing(userID)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			logger.Error(err)
			return
		}
		Data = dataFollow{
			userID,
			strings.Title(follow),
			true,
			petsFollowers,
			petsFollowing,
		}
		if follow == "followers" {
			Data.IsFollowing = false
			view.GenerateHTML(w, "My community", "navbar")
			view.GenerateHTML(w, Data, "follower")
			view.GenerateHTML(w, nil, "footer")
			return
		}
		view.GenerateHTML(w, "My community", "navbar")
		view.GenerateHTML(w, Data, "follower")
		view.GenerateHTML(w, nil, "footer")
		return
	}
}

//Unsubscribe or subscribe to a user
func (c *Controller) PostFollowerHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var err error
		r.ParseForm()
		follow := r.FormValue("follow")
		value := r.FormValue("followUserID")
		followUserID, err := strconv.Atoi(value)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			logger.Error(err)
			return
		}
		userID := context.Get(r, "id").(int)
		if follow == "Follow" {
			err = c.FollowersStore.Followed(userID, followUserID)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				logger.Error(err)
				return
			}
			http.Redirect(w, r, r.Header.Get("Referer"), 302)
			return
		}
		err = c.FollowersStore.UnFollowed(userID, followUserID)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			logger.Error(err)
			return
		}
		http.Redirect(w, r, r.Header.Get("Referer"), 302)

	}
}
