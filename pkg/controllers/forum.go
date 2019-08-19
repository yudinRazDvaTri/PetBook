package controllers

import (
	"net/http"
	"strconv"

	"github.com/dpgolang/PetBook/pkg/models/forum"

	"github.com/dpgolang/PetBook/pkg/logger"
	"github.com/dpgolang/PetBook/pkg/view"
	"github.com/gorilla/context"
	"github.com/gorilla/mux"
)

func (c *Controller) TopicsHandler() http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
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

			view.GenerateTimeHTML(w, "Forum", "navbar")
			view.GenerateTimeHTML(w, viewTopics, "topics")
		}

		if r.Method == http.MethodPost {
			r.ParseForm()
			title := r.FormValue("title")
			description := r.FormValue("description")
			uid := context.Get(r, "id").(int)

			if err := c.ForumStore.CreateNewTopic(uid, title, description); err != nil {
				logger.Error(err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			http.Redirect(w, r, "/forum", http.StatusFound)
		}
	}
}

func (c *Controller) CommentsHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		vars := mux.Vars(r)
		idStr := vars["id"]
		topicID, err := strconv.Atoi(idStr)
		if err != nil {
			logger.Error(err)
		}

		if r.Method == http.MethodGet {
			comments, err := c.ForumStore.GetTopicComments(topicID)
			if err != nil {
				logger.Error(err)
			}
			topic, err := c.ForumStore.GetTopicByID(topicID)
			if err != nil {
				logger.Error(err)
			}
			var viewComments []forum.ViewComment

			for _, comment := range comments {
				userName, err := c.PetStore.DisplayName(comment.UserID)
				if err != nil {
					logger.Error(err)
				}
				viewComments = append(viewComments, forum.ViewComment{userName, comment})
			}

			ViewData := struct {
				ID           int
				Topic        forum.Topic
				ViewComments []forum.ViewComment
			}{
				topicID,
				topic,
				viewComments,
			}

			view.GenerateTimeHTML(w, "Topic", "navbar")
			view.GenerateTimeHTML(w, ViewData, "comments")
		}

		if r.Method == http.MethodPost {
			r.ParseForm()
			content := r.FormValue("content")
			uid := context.Get(r, "id").(int)

			if err := c.ForumStore.AddNewComment(topicID, uid, content); err != nil {
				logger.Error(err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			id := strconv.Itoa(topicID)
			http.Redirect(w, r, "/forum/topic/"+id+"/comments", http.StatusFound)
		}

	}
}
