package controllers

import (
	"fmt"
	"github.com/dpgolang/PetBook/pkg/logger"
	"github.com/dpgolang/PetBook/pkg/models/forum"
	"github.com/dpgolang/PetBook/pkg/view"
	"github.com/gorilla/context"
	"net/http"
)

func (c *Controller) ForumHandler() http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			topics, err := c.ForumStore.GetAllTopics()
			if err != nil {
				logger.Error(err)
			}
			uid := context.Get(r, "id").(int)
			name, err := c.PetStore.DisplayName(uid)
			if err != nil {
				logger.Error(err)
			}

			fmt.Println(name)
			view.GenerateTimeHTML(w, "Forum", "navbar")
			view.GenerateTimeHTML(w, topics, "forum")
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
			http.Redirect(w, r, "/forum", http.StatusFound)
		}
	}
}
