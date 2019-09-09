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

// Returns a Page with List of Topics
func (c *Controller) TopicsGetHandler() http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		topics, err := c.ForumStore.GetAllTopics()
		if err != nil {
			logger.Error(err)
			http.Redirect(w, r, "/topics", http.StatusNotFound)
			return
		}

		var viewTopics []forum.ViewTopic

		for _, topic := range topics {
			userName, err := c.PetStore.DisplayName(topic.UserID)
			if err != nil {
				logger.Error(err)
				http.Redirect(w, r, "/topics", http.StatusInternalServerError)
				return
			}
			viewTopic, err := c.ForumStore.NewViewTopic(userName, topic)
			if err != nil {
				logger.Error(err)
				http.Redirect(w, r, "/topics", http.StatusInternalServerError)
				return
			}
			viewTopics = append(viewTopics, viewTopic)
		}

		view.GenerateTimeHTML(w, "Forum", "navbar")
		view.GenerateTimeHTML(w, viewTopics, "topics")
	}
}

// Process adding new Topic
func (c *Controller) TopicsPostHandler() http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		title := r.FormValue("title")
		description := r.FormValue("description")
		userID := context.Get(r, "id").(int)

		if err := c.ForumStore.CreateNewTopic(userID, title, description); err != nil {
			logger.Error(err)
			http.Error(w, "can't create new topic", http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/topics", http.StatusFound)
	}
}

// Returns a Page with Topic's Comments
func (c *Controller) CommentsGetHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := context.Get(r, "id").(int)

		vars := mux.Vars(r)
		topicIdStr := vars["topicID"]
		topicID, err := strconv.Atoi(topicIdStr)
		if err != nil {
			logger.Error(err)
			http.Error(w,"inappropriate request", http.StatusBadRequest)
			return
		}

		comments, err := c.ForumStore.GetTopicComments(topicID)
		if err != nil {
			logger.Error(err)
			http.Error(w, "can't get topic's comments", http.StatusInternalServerError)
			return
		}
		topic, err := c.ForumStore.GetTopicByID(topicID)
		if err != nil {
			logger.Error(err)
			http.Error(w, "No such topic", http.StatusNotFound)
			return
		}

		var viewComments []forum.ViewComment

		for _, comment := range comments {
			userName, err := c.PetStore.DisplayName(comment.UserID)
			if err != nil {
				logger.Error(err)
				http.Error(w, "can't get comment creator's name", http.StatusInternalServerError)
				return
			}
			viewComment, err := c.ForumStore.NewViewComment(userName, comment)
			if err != nil {
				logger.Error(err)
				http.Error(w, "can't get likes-field of comment", http.StatusInternalServerError)
				return
			}
			viewComments = append(viewComments, viewComment)
		}

		treeVComments, err := forum.TreeViewComments(viewComments)

		ViewData := struct {
			ContextUserID int
			Topic         forum.Topic
			TreeVComments []*forum.ViewComment
		}{
			userID,
			topic,
			treeVComments,
		}

		view.GenerateTimeHTML(w, "Topic", "navbar")
		view.GenerateTimeHTML(w, ViewData, "comments")
	}
}

// Process adding new Comment
func (c *Controller) CommentPostHandler() http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		var parentID int
		keys, ok := r.URL.Query()["parentID"]

		if !ok || len(keys[0]) < 1 {
			logger.Error("missing parentID for comment in URL")
			http.Error(w,"missing URL parameter", http.StatusInternalServerError)
			return
		}

		parentID, err := strconv.Atoi(keys[0])
		if err != nil {
			logger.Error("parentID from URL is not an integer")
			http.Error(w,"inappropriate URL parameter", http.StatusInternalServerError)
			return
		}

		vars := mux.Vars(r)
		topicIdStr := vars["topicID"]
		topicID, err := strconv.Atoi(topicIdStr)
		if err != nil {
			logger.Error(err)
			http.Error(w, "inappropriate url", http.StatusInternalServerError)
			return
		}

		r.ParseForm()
		content := r.FormValue("content")
		userID := context.Get(r, "id").(int)

		if err := c.ForumStore.AddNewComment(topicID, userID, parentID, content); err != nil {
			logger.Error(err)
			http.Error(w, "can't add comment", http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/topics/"+topicIdStr, http.StatusFound)
	}
}

// Process Like-action on Comment
func (c *Controller) CommentsRatingHandler() http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		topicIdStr := vars["topicID"]
		commentIdStr := vars["commentID"]

		commentID, err := strconv.Atoi(commentIdStr)
		if err != nil {
			logger.Error(err)
			http.Redirect(w, r, "/topics/"+topicIdStr, http.StatusInternalServerError)
			return
		}

		userID := context.Get(r, "id").(int)

		rateOk, err := c.ForumStore.RateComment(commentID, userID)
		if err != nil {
			logger.Error(err)
			http.Redirect(w, r, "/topics/"+topicIdStr, http.StatusInternalServerError)
			return
		}
		if !rateOk {
			http.Redirect(w, r, "/topics/"+topicIdStr, http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/topics/"+topicIdStr, http.StatusFound)
	}
}