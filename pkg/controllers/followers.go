package controllers

import (
	"github.com/dpgolang/PetBook/pkg/logger"
	"github.com/dpgolang/PetBook/pkg/models"
	"github.com/dpgolang/PetBook/pkg/view"
	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
	"strings"
)
type DataFollow struct{
	UserID int
	PageAction string
	IsFollowing bool
	PetsFollowers []*models.FollowerPets
	PetsFollowing []*models.FollowerPets
	//CanFollowing func(userID int, followingUserID int, petsFollowing []*models.FollowerPets)bool
}
func (c *Controller) GetFollowerHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var Data DataFollow
		params := mux.Vars(r)
		follow:=params["follow"]
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
		Data = DataFollow {
			userID,
			strings.Title(follow),
			true,
			petsFollowers,
			petsFollowing,
			//CanFollow,
		}
		if follow=="followers"{
			Data.IsFollowing=false
			view.GenerateHTML(w, Data, "follower")
			return
		}
		view.GenerateHTML(w, Data, "follower")
		return
	}
}

func (c *Controller) PostFollowerHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var err error
		r.ParseForm()
		follow:=r.FormValue("follow")
		value := r.FormValue("followUserID")
		followUserID, err := strconv.Atoi(value)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			logger.Error(err)
			return
		}

		userID := context.Get(r, "id").(int)
		if follow=="Follow"{
			err = c.FollowersStore.Followed(userID,followUserID)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				logger.Error(err)
				return
			}

		}
		err = c.FollowersStore.UnFollowed(userID,followUserID)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			logger.Error(err)
			return
		}


	}
}
