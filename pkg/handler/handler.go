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

// Constants and Struct for storing State
const stateZero string = "Please review our product."
const stateOne string = "Would you like to leave additional feedback?"
const close string = "Thank you for being a valued customer, if you have additional questions please visit our website https://nike.com"

type UserState struct {
	UserId         int
	currState      string
	stateZeroReply string
	stateOneReply  string
}

// map that stores state with the UserId as keys
var userStates = make(map[int]*UserState, 20)

// Endpoint strings
const telegramApiBaseUrl string = "https://api.telegram.org/bot"
const telegramApiSendMessage string = "/sendMessage"
const telegramTokenEnv string = "BOT_TOKEN"

var telegramEndpoint string = telegramApiBaseUrl + os.Getenv(telegramTokenEnv) + telegramApiSendMessage

// Structs to decode Updates from telegram
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

// Entry point to Cloud Function
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

	_, _ = handleUserInput(update)
}

// in real implementation this function would call database to get User State and store responses
// we could also send requests to our ai endpoint to integrate ai functionality here
// store User state and send appropriate message back to them
func handleUserInput(update *Update) (string, error) {
	userId := update.Message.From.UserId
	val, keyExists := userStates[userId]
	if !keyExists || update.Message.Text == "/start" {
		var temp UserState
		temp.UserId = userId
		temp.currState = stateZero
		userStates[userId] = &temp
	} else {
		log.Printf(val.currState)
		if val.currState == stateZero {
			userStates[userId].currState = stateOne
			userStates[userId].stateZeroReply = update.Message.Text
		} else if val.currState == stateOne {
			userStates[userId].currState = close
			userStates[userId].stateOneReply = update.Message.Text
		}
	}
	responseBody, err := sendMessage(update.Message.Chat.ChatId, userStates[userId].currState)
	if err != nil {
		log.Printf(err.Error())
	} else {
		log.Printf("Message sent succesfully, Response Body: %s", responseBody)
	}
	return responseBody, err
}

// Destructure Update JSON into our structs
func decodeRequest(r *http.Request) (*Update, error) {
	var update Update
	if err := json.NewDecoder(r.Body).Decode(&update); err != nil {
		log.Printf("could not decode update %s", err.Error())
		return nil, err
	}
	return &update, nil
}

// Send Message back to User
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
