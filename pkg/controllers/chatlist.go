package controllers

import (
	"github.com/dpgolang/PetBook/pkg/logger"
	"github.com/dpgolang/PetBook/pkg/models"
	"github.com/dpgolang/PetBook/pkg/view"
	"github.com/gorilla/context"
	"log"
	"net/http"
)

func (c *Controller) ChatsGetHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := context.Get(r, "id").(int)
		chats, err := c.ChatStore.GetChats(userID)
		if err != nil {
			logger.Error(err)
		}
		var viewChats []models.ChatToView
		for _, chat := range chats {
			username, err := c.PetStore.DisplayName(chat.ToID)
			if err != nil {
				logger.Error(err)
			}
			viewChat := models.ChatToView{
				ToID:      chat.ToID,
				Username:  username,
				Message:   chat.Message,
				CreatedAt: chat.CreatedAt.Format("02-01-2006 15:04:05"),
			}
			viewChats = append(viewChats, viewChat)
		}
		log.Println(viewChats)
		view.GenerateTimeHTML(w, "Chats", "navbar")
		view.GenerateTimeHTML(w, viewChats, "chatlist")
	}
}
