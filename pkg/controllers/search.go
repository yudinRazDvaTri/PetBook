package controllers

import (
	"github.com/dpgolang/PetBook/pkg/logger"
	"github.com/dpgolang/PetBook/pkg/models"
	"github.com/dpgolang/PetBook/pkg/models/forum"
	"github.com/dpgolang/PetBook/pkg/models/search"
	"github.com/dpgolang/PetBook/pkg/view"
	"github.com/gorilla/context"
	"net/http"
)
type FollowerPets struct {
	Name        string `json:"name" db:"name"'`
	Description string `json:"description" db:"description"'`
	UserID 		int    `json:"user_id" db:"user_id"`
}
type DataSearch struct{
	UserID int
	PetsFollowing []*models.FollowerPets
	Pets []*search.DispPet
}
func (c *Controller) RedirectSearchHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/search?section=user", http.StatusFound)
		return
	}
}
func (c *Controller) ViewSearchHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		view.GenerateHTML(w, "Search", "navbar")
		section := r.URL.Query().Get("section")
		if section == "user" {
			c.searchByUser(w, r)
		}
		if section == "animal" {
			c.searchByPet(w, r)
		}
		if section == "forum" {
			c.searchByForum(w, r)
		}
	}
}

func (c *Controller) searchByUser(w http.ResponseWriter, r *http.Request) {
	var (
		err error
		dataSearch DataSearch

	)
	userID:=context.Get(r, "id").(int)
	dataSearch.UserID = userID
	dataSearch.PetsFollowing,err = c.FollowersStore.GetFollowing(userID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		logger.Error(err)
		return
	}
	email := r.URL.Query().Get("email")
	if email != "" {
		pet:= c.SearchStore.GetByUser(userID,email)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			logger.Error(err)
			return
		}
		dataSearch.Pets=append(dataSearch.Pets, pet)
		view.GenerateHTML(w, dataSearch, "search_by_user")
		return
	}
	dataSearch.Pets, err = c.SearchStore.GetAllPets(userID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		logger.Error(err)
		return
	}
	view.GenerateHTML(w, dataSearch, "search_by_user")

}
func (c *Controller) searchByPet(w http.ResponseWriter, r *http.Request) {
	var (
		err error
		dataSearch DataSearch
	)
	dataSearch.UserID=context.Get(r, "id").(int)
	dataSearch.PetsFollowing,err = c.FollowersStore.GetFollowing(dataSearch.UserID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		logger.Error(err)
		return
	}
	m := make(map[string]interface{})
	queryStr := []string{"age", "animal_type", "breed", "weight", "gender", "name"}
	for _, str := range queryStr {
		val := r.URL.Query().Get(str)
		if val != "" {
			m[str] = val
		}
	}
	if len(m) == 0 {
		dataSearch.Pets, err = c.SearchStore.GetAllPets(dataSearch.UserID)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			logger.Error(err)
			return
		}
		view.GenerateHTML(w, dataSearch, "search_by_animals")
		return
	}
	dataSearch.Pets, err = c.SearchStore.GetFilterPets(m)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		logger.Error(err)
		return
	}
	view.GenerateHTML(w, dataSearch, "search_by_animals")

}
func (c *Controller) searchByForum(w http.ResponseWriter, r *http.Request) {
	search := r.URL.Query().Get("search")
	if search != "" {
		topics, err := c.SearchStore.GetTopicsBySearch(search)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			logger.Error(err)
			return
		}
		var viewTopics []forum.ViewTopic

		for _, topic := range topics {
			userName, err := c.PetStore.DisplayName(topic.UserID)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				logger.Error(err)
				return
			}
			viewTopic, err := c.ForumStore.NewViewTopic(userName, topic)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				logger.Error(err)
				return
			}
			viewTopics = append(viewTopics, viewTopic)
		}

		view.GenerateTimeHTML(w, viewTopics, "search_by_forum")
		return
	}
	topics, err := c.ForumStore.GetAllTopics()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		logger.Error(err)
		return
	}

	var viewTopics []forum.ViewTopic

	for _, topic := range topics {
		userName, err := c.PetStore.DisplayName(topic.UserID)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			logger.Error(err)
			return
		}
		viewTopic, err := c.ForumStore.NewViewTopic(userName, topic)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			logger.Error(err)
			return
		}
		viewTopics = append(viewTopics, viewTopic)
	}

	view.GenerateTimeHTML(w, viewTopics, "search_by_forum")

}
