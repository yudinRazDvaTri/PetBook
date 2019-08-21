package controllers

import (
	"net/http"
	"sort"
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
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/forum", http.StatusFound)
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
		}

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
			ratedUsers, err := c.ForumStore.GetCommentRatings(comment.CommentID)
			if err != nil {
				logger.Error(err)
			}
			viewComments = append(viewComments, forum.ViewComment{userName, ratedUsers, comment})
		}

		sort.Sort(sort.Reverse(forum.ByRating(viewComments)))

		ViewData := struct {
			ContextUserID int
			Topic         forum.Topic
			ViewComments  []forum.ViewComment
		}{
			userID,
			topic,
			viewComments,
		}

		view.GenerateTimeHTML(w, "Topic", "navbar")
		view.GenerateTimeHTML(w, ViewData, "comments")
	}
}

// Process adding new Comment
func (c *Controller) CommentsPostHandler() http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		topicIdStr := vars["topicID"]
		topicID, err := strconv.Atoi(topicIdStr)
		if err != nil {
			logger.Error(err)
		}

		r.ParseForm()
		content := r.FormValue("content")
		userID := context.Get(r, "id").(int)

		if err := c.ForumStore.AddNewComment(topicID, userID, content); err != nil {
			logger.Error(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/forum/topic/"+topicIdStr+"/comments", http.StatusFound)
	}
}

// Process Like-action on Comment
func (c *Controller) CommentsRatingsHandler() http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		topicIdStr := vars["topicID"]
		commentIdStr := vars["commentID"]

		commentID, err := strconv.Atoi(commentIdStr)
		if err != nil {
			logger.Error(err)
		}

		userID := context.Get(r, "id").(int)

		rateOk, err := c.ForumStore.RateComment(commentID, userID)
		if err != nil {
			logger.Error(err)
		}
		if !rateOk {
			http.Redirect(w, r, "/forum/topic/"+topicIdStr+"/comments", http.StatusFound)
		}

		http.Redirect(w, r, "/forum/topic/"+topicIdStr+"/comments", http.StatusFound)
	}
}