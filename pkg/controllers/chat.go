package controllers

import (
	"github.com/dpgolang/PetBook/pkg/logger"
	"github.com/dpgolang/PetBook/pkg/models"
	"github.com/dpgolang/PetBook/pkg/view"
	"github.com/gorilla/context"
	//"github.com/gorilla/mux"
	"github.com/gorilla/websocket"

	"log"
	"net/http"
)

var clients = make(map[models.Client]bool)      // connected clients
var broadcast = make(chan models.MessageToView) // broadcast channel

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (c *Controller) HandleChatConnectionGET() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		view.GenerateHTML(w, nil, "chat")
	}
}

func (c *Controller) HandleChatConnection() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ws, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Fatal(err)
		}
		defer ws.Close()

		userID := context.Get(r, "id").(int)
		client := models.Client{
			UserID:     userID,
			Connection: ws,
		}
		// params := mux.Vars(r)
		// id, err := strconv.Atoi(params["id"])
		// Register our new client
		clients[client] = true
		for {
			var msg models.MessageToView
			// Read in a new message as JSON and map it to a Message object
			err := ws.ReadJSON(&msg)
			if err != nil {
				log.Printf("error: %v", err)
				delete(clients, client)
				break
			}
			// Send the newly received message to the broadcast channel
			broadcast <- msg
		}
	}
}

func (c *Controller) HandleMessages() {
	for {
		msg := <-broadcast
		for client := range clients {
			err := client.Connection.WriteJSON(msg)
			if err != nil {
				logger.Error("send message error: %v", err)
				client.Connection.Close()
				delete(clients, client)
			}
		}
	}
}
