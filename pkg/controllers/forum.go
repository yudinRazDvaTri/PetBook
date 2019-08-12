package controllers

import (
	"github.com/dpgolang/PetBook/pkg/logger"
	"github.com/dpgolang/PetBook/pkg/models/forum"
	"github.com/dpgolang/PetBook/pkg/view"
	"github.com/gorilla/context"
	"net/http"
)

func (c *Controller) ViewTopicsHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		topics, err := c.ForumStore.GetAllTopics()
		if err != nil {
			logger.Error(err)
		}

		view.GenerateHTML(w, topics, "viewTopics")
	}
}

func (c *Controller) NewTopicHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			view.GenerateHTML(w, nil, "newTopic")
		}

		if r.Method == http.MethodPost {
			r.ParseForm()
			title := r.FormValue("title")
			description := r.FormValue("description")
			uid := context.Get(r, "id").(int)

			topic := &forum.Topic{
				UserID:      uid,
				Title:       title,
				Description: description,
			}
			if err := c.ForumStore.CreateNewTopic(topic); err != nil {
				logger.Error(err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			http.Redirect(w, r, "/forum/submit", http.StatusFound)
		}
	}
}