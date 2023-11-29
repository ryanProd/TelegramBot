package handler

import (
	"log"
	"net/http"

	_ "github.com/lib/pq"
	"github.com/ryanProd/TelegramBot/internal/dbFunctions"
	"github.com/ryanProd/TelegramBot/internal/handleUserInput"
	"github.com/ryanProd/TelegramBot/internal/httpFunctions"
)

// Constants and Struct for storing State
const stateZero string = "Please review our product."
const stateOne string = "Would you like to leave additional feedback?"
const close string = "Thank you for being a valued customer, if you have additional questions please visit our website https://nike.com. If you would like to reset this bot type /start"

// Entry point to Cloud Function
func Handler(w http.ResponseWriter, r *http.Request) {

	var update, err = httpFunctions.DecodeRequest(r)
	if err != nil {
		log.Printf("could not decode update %s", err.Error())
		w.WriteHeader(500)
	}

	log.Printf("Update Received")
	log.Printf("Update ID is %d", update.UpdateId)
	log.Printf("Message ID is %d", update.Message.MessageId)
	log.Printf("Message Text is %s", update.Message.Text)
	log.Printf("User ID is %d", update.Message.From.UserId)
	log.Printf("Chat ID is %d", update.Message.Chat.ChatId)

	db := dbFunctions.ConnectDB()
	defer db.Close()

	if _, err = handleUserInput.HandleUserInput(update, db); err != nil {
		log.Printf("HandleUserInput returned error %s", err.Error())
		w.WriteHeader(500)
	} else {
		w.WriteHeader(200)
	}
}
