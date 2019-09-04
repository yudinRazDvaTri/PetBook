package controllers

import (
	"github.com/dpgolang/PetBook/pkg/logger"
	"github.com/dpgolang/PetBook/pkg/models/forum"
	"github.com/dpgolang/PetBook/pkg/view"
	"net/http"
)

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
	email := r.URL.Query().Get("email")
	if email != "" {
		pet, err := c.SearchStore.GetByUser(email)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			logger.Error(err)
			return
		}
		view.GenerateHTML(w, pet, "view_animal")
		return
	}
	pets, err := c.SearchStore.GetAllPets()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		logger.Error(err)
		return
	}
	view.GenerateHTML(w, pets, "search_by_user")

}
func (c *Controller) searchByPet(w http.ResponseWriter, r *http.Request) {
	m := make(map[string]interface{})
	queryStr := []string{"age", "animal_type", "breed", "weight", "gender", "name"}
	for _, str := range queryStr {
		val := r.URL.Query().Get(str)

		if val != "" {
			m[str] = val
		}
	}
	if len(m) == 0 {
		pets, err := c.SearchStore.GetAllPets()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			logger.Error(err)
			return
		}
		view.GenerateHTML(w, pets, "search_by_animals")
		return
	}
	filterPets, err := c.SearchStore.GetFilterPets(m)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		logger.Error(err)
		return
	}
	view.GenerateHTML(w, filterPets, "search_by_animals")

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
