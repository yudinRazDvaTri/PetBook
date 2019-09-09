package controllers

import (
	//"bytes"
	"encoding/json"
	"github.com/dpgolang/PetBook/pkg/logger"
	"github.com/dpgolang/PetBook/pkg/models"
	"github.com/dpgolang/PetBook/pkg/view"
	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	//"log"
	"net/http"
	"strconv"
	"time"
)

var clients = make(map[models.Client]bool) // connected clients

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (c *Controller) HandleChatConnectionGET() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var err error
		params := mux.Vars(r)
		toID, err := strconv.Atoi(params["id"])
		if err != nil {
			logger.Error(err)
			http.Redirect(w, r, "/chats", http.StatusNotFound)
			return
		}
		_, err = c.UserStore.GetPet(toID)
		if err != nil {
			logger.Error(err)
			http.Redirect(w, r, "/mypage", http.StatusNotFound)
			return
		}
		view.GenerateHTML(w, toID, "chat")
	}
}

func (c *Controller) HandleChatConnection() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ws, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			logger.Error(err)
			return
		}
		//defer ws.Close()

		fromID := context.Get(r, "id").(int)

		params := mux.Vars(r)
		toID, err := strconv.Atoi(params["id"])
		if err != nil {
			logger.Error(err)
			http.Redirect(w, r, "/chats", http.StatusNotFound)
			return
		}
		client := models.Client{
			ID:         fromID,
			Connection: ws,
		}
		// Register our new client
		clients[client] = true
		// Register our new client
		messages, err := c.ChatStore.GetMessages(toID, fromID)
		if err != nil {
			logger.Error("can't get messages: ", err)
		}
		for _, mes := range messages {
			msg := models.MessageToView{
				ToID:      mes.ToID,
				FromID:    mes.FromID,
				Message:   mes.Text,
				CreatedAt: mes.CreatedAt.Format("02-01-2006 15:04:05"),
			}
			msg.Username, err = c.PetStore.DisplayName(msg.FromID)
			if err != nil {
				logger.Error("cannot display name correctly: ", err)
			}
			err = ws.WriteJSON(msg)
			if err != nil {
				logger.Error("cannot write json from db: ", err)
			}
		}

		go c.HandleMessages(ws, client, toID)
	}
}

func (c *Controller) HandleMessages(ws *websocket.Conn, client models.Client, toID int) {
	go func() {
		for {
			var msg models.MessageToView
			var err error
			msg.Username, err = c.PetStore.DisplayName(client.ID)
			if err != nil {
				logger.Error("cannot display name correctly: ", err)
			}
			msg.FromID = client.ID
			msg.ToID = toID
			msg.CreatedAt = time.Now().Format("02-01-2006 15:04:05")
			err = ws.ReadJSON(&msg)
			if err != nil {
				logger.Error(err)
				delete(clients, client)
				ws.Close()
				return
			}

			messageCreatedAt, err := time.Parse("02-01-2006 15:04:05", msg.CreatedAt)
			if err != nil {
				logger.Error("something gone wrong while parsing message created_at:", err)
				continue
			}
			messageForSavingIntoDB := &models.Message{
				FromID:    msg.FromID,
				ToID:      msg.ToID,
				Text:      msg.Message,
				CreatedAt: messageCreatedAt,
			}
			err = c.ChatStore.SaveMessage(messageForSavingIntoDB)
			if err != nil {
				logger.Error(err)
			}
			for client := range clients {
				if client.ID == msg.FromID || client.ID == msg.ToID {
					err = client.Connection.WriteJSON(msg)
					if err != nil {
						logger.Error("send message error:", err)
						ws.Close()
						client.Connection.Close()
						delete(clients, client)
						return
					}
				}
			}
		}
	}()
}

func (c *Controller) HandleChatSearchConnection() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fromID := context.Get(r, "id").(int)
		params := mux.Vars(r)
		toID, err := strconv.Atoi(params["id"])
		if err != nil {
			logger.Error(err)
			http.Redirect(w, r, "/chats", http.StatusNotFound)
			return
		}

		strToID := strconv.FormatInt(int64(toID), 10)
		date, err := time.Parse("02-01-2006", params["date"])
		if err != nil {
			logger.Error(err)
			http.Redirect(w, r, "/chats/"+strToID, http.StatusNotFound)
			return
		}
		messages, err := c.ChatStore.GetMessagesByDate(toID, fromID, date)
		if err != nil {
			logger.Error("can't get messages: ", err)
		}
		messagesToView := []models.MessageToView{}
		for _, mes := range messages {
			msg := models.MessageToView{
				ToID:      mes.ToID,
				FromID:    mes.FromID,
				Message:   mes.Text,
				CreatedAt: mes.CreatedAt.Format("02-01-2006 15:04:05"),
			}
			msg.Username, err = c.PetStore.DisplayName(msg.FromID)
			if err != nil {
				logger.Error("cannot display name correctly: ", err)
			}

			messagesToView = append(messagesToView, msg)
		}
		marshalledMsgs, err := json.Marshal(messagesToView)
		if err != nil {
			logger.Error("can't marshall messages: ", err)
			http.Redirect(w, r, "/chats/"+strToID, http.StatusNotFound)
		}
		w.Write(marshalledMsgs)
	}
}
