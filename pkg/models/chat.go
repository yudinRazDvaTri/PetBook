package models

import (
	//"github.com/dpgolang/PetBook/pkg/logger"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/jmoiron/sqlx"
	"time"
)

type MessageToView struct {
	ToID      int    `json:"toid", omitempty`
	FromID    int    `json:"fromid", omitempty`
	Username  string `json:"username"`
	Message   string `json:"message"`
	CreatedAt string `json:"created_at"`
}

type Message struct {
	ID        int       `db:"id"`
	ToID      int       `json:"toid" db:"to_id"`
	FromID    int       `json:"fromid" db:"from_id"`
	Text      string    `json:"text" db:"text"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

type Client struct {
	ID         int
	Connection *websocket.Conn
}

type Chat struct {
	Companion User
	Messages  []Message
}

type ChatStore struct {
	DB *sqlx.DB
}

type ChatStorer interface {
	GetMessages(toID, fromID int) ([]Message, error)
	SaveMessage(message *Message) error
}

func (c *ChatStore) GetMessages(toID, fromID int) ([]Message, error) {
	rows, err := c.DB.Query("select * from messages where (to_id=$1 and from_id=$2) or (from_id=$1 and to_id= $2) order by created_at", toID, fromID)
	if err != nil {
		return nil, fmt.Errorf("cannot make querry: %v", err)
	}
	defer rows.Close()
	messages := []Message{}
	err = sqlx.StructScan(rows, &messages)
	if err != nil {
		return nil, fmt.Errorf("cannot scan messages from db: %v", err)
	}
	return messages, nil
}

func (c *ChatStore) SaveMessage(message *Message) error {
	_, err := c.DB.Query(`insert into messages (to_id, from_id, text,created_at) select $1,$2,$3,$4
		WHERE NOT EXISTS (select 1 from messages  where to_id =$1 and from_id = $2 and text = $3 and created_at = $4)`,
		message.ToID, message.FromID, message.Text, message.CreatedAt)
	if err != nil {
		return fmt.Errorf("cannot insert message to messages in db: %v", err)
	}
	return nil
}

//TODO
// func (c *ChatStore) GetChats(user *User) ([]Chat, error) {
// 	rows, err := c.DB.Query(
// 		`select  to_id, text, created_at from (select to_id, text, created_at from messages where from_id = $1
// 			UNION
// 		select from_id, text,  created_at from messages where to_id = $2) as select_result order by created_at desc limit 1`, user.ID)
// 	if err != nil {
// 		return nil, fmt.Errorf("cannot make querry: %v", err)
// 	}
// 	defer rows.Close()
// 	messages := []Message{}
// 	err = sqlx.StructScan(rows, &messages)
// 	if err != nil {
// 		return nil, fmt.Errorf("cannot scan messages from db: %v", err)
// 	}
// 	return nil, nil
// }
