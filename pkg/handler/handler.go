package handler

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
)

const telegramApiBaseUrl string = "https://api.telegram.org/bot"
const telegramApiSendMessage string = "/sendMessage"
const telegramTokenEnv string = "BOT_TOKEN"

var telegramEndpoint string = telegramApiBaseUrl + os.Getenv(telegramTokenEnv) + telegramApiSendMessage

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

	if update.Message.Text == "/start" {
		response, err := sendMessage(update.Message.Chat.ChatId, "Hello Welcome to Ryan's Chat Bot!")
		if err != nil {
			log.Printf(err.Error())
		} else {
			log.Printf(response)
		}
	} else {
		response, err := sendMessage(update.Message.Chat.ChatId, update.Message.Text)
		if err != nil {
			log.Printf(err.Error())
		} else {
			log.Printf(response)
		}
	}
}

func decodeRequest(r *http.Request) (*Update, error) {
	var update Update
	if err := json.NewDecoder(r.Body).Decode(&update); err != nil {
		log.Printf("could not decode update %s", err.Error())
		return nil, err
	}
	return &update, nil
}

func sendMessage(chatId int, text string) (string, error) {
	log.Printf("Sending %s to ChatID: %d", text, chatId)

	response, err := http.PostForm(
		telegramEndpoint,
		url.Values{
			"chat_id": {strconv.Itoa(chatId)},
			"text":    {text},
		})
	if err != nil {
		log.Printf("Message to ChatID: %d failed to send. %s", chatId, err.Error())
	}
	defer response.Body.Close()

	body, responseErr := ioutil.ReadAll(response.Body)
	if responseErr != nil {
		log.Printf(responseErr.Error())
	}

	return string(body), responseErr

}
