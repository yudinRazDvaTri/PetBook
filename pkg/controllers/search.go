package controllers

import (
	"github.com/dpgolang/PetBook/pkg/logger"
	"github.com/dpgolang/PetBook/pkg/models/forum"
	"github.com/dpgolang/PetBook/pkg/view"
	"net/http"
)
func (c *Controller) RedirectSearchHandler()http.HandlerFunc{
	return func(w http.ResponseWriter, r *http.Request){
		http.Redirect(w, r, "/search?section=user", http.StatusFound)
		return
	}
}
func (c *Controller) ViewSearchHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		view.GenerateHTML(w, "Search", "navbar")
		section := r.URL.Query().Get("section")
		if section == "user" {
			c.searchByUser(w,r)
		}
		if section == "animal" {
			c.searchByPet(w,r)
		}
		if section=="forum"{
			c.searchByForum(w,r)
		}
	}
}
func (c *Controller) searchByUser(w http.ResponseWriter, r *http.Request)  {
	email := r.URL.Query().Get("email")
	if email != "" {
		pet, err := c.SearchStore.GetByUser(email)
		if err != nil {
			logger.Error(err)
		}
		view.GenerateHTML(w, pet, "view_animal")
	} else {
		pets, err := c.SearchStore.GetAllPets()
		if err != nil {
			logger.Error(err)
		}
		view.GenerateHTML(w, pets, "search_by_user")
	}
}
func (c *Controller) searchByPet(w http.ResponseWriter, r *http.Request){
	m := make(map[string]string)
	age := r.URL.Query().Get("age")
	if age != "" {
		m["age"] = age
	}
	animalType := r.URL.Query().Get("type")
	if animalType != "" {
		m["animal_type"] = animalType
	}
	breed := r.URL.Query().Get("breed")
	if breed != "" {
		m["breed"] = breed
	}
	weight := r.URL.Query().Get("weight")
	if weight != "" {
		m["weight"] = weight
	}
	gender := r.URL.Query().Get("gender")
	if gender != "" {
		m["gender"] = gender
	}
	name := r.URL.Query().Get("name")
	if name != "" {
		m["name"] = name
	}
	if len(m) == 0 {
		pets, err := c.SearchStore.GetAllPets()
		if err != nil {
			logger.Error(err)
		}
		view.GenerateHTML(w, pets, "search_by_animals")
	} else {
		filterPets, err := c.SearchStore.GetFilterPets(m)
		if err != nil {
			logger.Error(err)
		}
		view.GenerateHTML(w, filterPets, "search_by_animals")
	}
}
func(c *Controller) searchByForum(w http.ResponseWriter, r * http.Request){
	search := r.URL.Query().Get("search")
	if search != "" {
		topics,err:=c.SearchStore.GetTopicsBySearch(search)
		if err != nil {
			logger.Error(err)
		}
		var viewTopics []forum.ViewTopic

		for _, topic := range topics {
			userName, err := c.PetStore.DisplayName(topic.UserID)
			if err != nil {
				logger.Error(err)
			}
			viewTopics = append(viewTopics, forum.ViewTopic{userName, topic})
		}

		view.GenerateTimeHTML(w,viewTopics,"search_by_forum")
	} else {
		topics, err := c.ForumStore.GetAllTopics()
		if err != nil {
			logger.Error(err)
		}

		var viewTopics []forum.ViewTopic

		for _, topic := range topics {
			userName, err := c.PetStore.DisplayName(topic.UserID)
			if err != nil {
				logger.Error(err)
			}
			viewTopics = append(viewTopics, forum.ViewTopic{userName, topic})
		}

		view.GenerateTimeHTML(w,viewTopics,"search_by_forum")
	}
}

