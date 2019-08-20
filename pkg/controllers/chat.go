package controllers

import (
	"github.com/dpgolang/PetBook/pkg/logger"
	"github.com/dpgolang/PetBook/pkg/models"
	"github.com/dpgolang/PetBook/pkg/view"
	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"net/http"
	"strconv"
)

var clients = make(map[models.Client]bool)      // connected clients
var broadcast = make(chan models.MessageToView) // broadcast channel

var toID int

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (c *Controller) HandleChatConnectionGET() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var err error
		view.GenerateHTML(w, nil, "chat")
		params := mux.Vars(r)
		toID, err = strconv.Atoi(params["id"])
		if err != nil {
			logger.Error(err)
			//http.Redirect(w, r, "/mypage", http.StatusNotFound)
		}
		_, err = c.UserStore.GetPet(toID)
		// if err != nil {
		// 	logger.Error(err)
		// 	http.Redirect(w, r, "/mypage", http.StatusNotFound)
		// }
	}
}

func (c *Controller) HandleChatConnection() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ws, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			logger.Error(err)
			return
		}
		defer ws.Close()

		fromID := context.Get(r, "id").(int)
		client := models.Client{
			ID:         fromID,
			Connection: ws,
		}

		// Register our new client
		clients[client] = true
		for {
			var msg models.MessageToView
			msg.Username, err = c.PetStore.DisplayName(fromID)
			msg.FromID = fromID
			msg.ToID = toID
			// Read in a new message as JSON and map it to a Message object
			err := ws.ReadJSON(&msg)
			if err != nil {
				logger.Error("error: %v", err)
				delete(clients, client)
				break
			}
			// Send the newly received message to the broadcast channel
			broadcast <- msg
		}
	}
}

func (c *Controller) HandleMessages() {
	var err error
	for {
		msg := <-broadcast
		for client := range clients {
			if client.ID == msg.FromID || client.ID == msg.ToID { // check that there are correct users to send and to get message
				err = client.Connection.WriteJSON(msg)
				if err != nil {
					logger.Error("send message error: %v", err)
					client.Connection.Close()
					delete(clients, client)
				}
			}
		}
	}
}
