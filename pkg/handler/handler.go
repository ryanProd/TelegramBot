package handler

import (
	"encoding/json"
	"log"
	"net/http"
)

type Update struct {
	UpdateId int     `json:"update_id"`
	Message  Message `json:"message"`
}

type Message struct {
	MessageId int    `json:"message_id"`
	Text      string `json:"text"`
	From      User   `json:"from"`
	Chat      Chat   `json:"chat"`
}

type Chat struct {
	ChatId int    `json:"id"`
	Type   string `json:"type"`
}

type User struct {
	UserId   int    `json:"id"`
	Username string `json:"username"`
}

func Handler(w http.ResponseWriter, r *http.Request) {
	var update, err = decodeRequest(r)
	if err != nil {
		log.Printf("could not decode update %s", err.Error())
	}
	log.Printf("Update Received")
	log.Printf("Update ID is %d", update.UpdateId)
	log.Printf("Message ID is %d", update.Message.MessageId)
	log.Printf("Message Text is %s", update.Message.Text)
	log.Printf("User ID is %d", update.Message.From.UserId)
	log.Printf("Chat ID is %d", update.Message.Chat.ChatId)

}

func decodeRequest(r *http.Request) (*Update, error) {
	var update Update
	if err := json.NewDecoder(r.Body).Decode(&update); err != nil {
		log.Printf("could not decode update %s", err.Error())
		return nil, err
	}
	return &update, nil
}
