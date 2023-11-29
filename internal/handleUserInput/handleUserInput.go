package handleUserInput

import (
	"database/sql"
	"log"

	"github.com/ryanProd/TelegramBot/internal/dbFunctions"
	"github.com/ryanProd/TelegramBot/internal/httpFunctions"
	"github.com/ryanProd/TelegramBot/internal/structs"
)

// Constants and Struct for storing State
const stateZero string = "Please review our product."
const stateOne string = "Would you like to leave additional feedback?"
const close string = "Thank you for being a valued customer, if you have additional questions please visit our website https://nike.com. If you would like to reset this bot type /start"

// entry point function to handle message from User
func HandleUserInput(update *structs.Update, db *sql.DB) (string, error) {
	userId := update.Message.From.UserId
	userState, err := dbFunctions.QueryDB(db, userId)
	if err != nil {
		//actual error returned from queryDB function
		if err != sql.ErrNoRows {
			log.Printf(err.Error())
			return "", err
		} else {
			//queryDB returned sql.ErrNoRows which means user is not in database, insert into database
			log.Printf("USER NOT FOUND")
			dbFunctions.InsertRow(db, userId)
			httpFunctions.PostMessage(update.Message.Chat.ChatId, stateZero)
			return "", nil
		}
	}
	//user already exists in database, handle cases depending on current state
	if update.Message.Text == "/start" {
		log.Printf("USER FOUND")
		dbFunctions.UpdateRow(db, userId, stateZero, "", "")
		httpFunctions.PostMessage(update.Message.Chat.ChatId, stateZero)
	} else {
		if userState.CurrState == stateZero {
			dbFunctions.UpdateRow(db, userId, stateOne, update.Message.Text, "")
			httpFunctions.PostMessage(update.Message.Chat.ChatId, stateOne)
		} else if userState.CurrState == stateOne {
			dbFunctions.UpdateRow(db, userId, close, userState.StateZeroReply, update.Message.Text)
			httpFunctions.PostMessage(update.Message.Chat.ChatId, close)
		} else {
			httpFunctions.PostMessage(update.Message.Chat.ChatId, close)
		}
	}

	return "", err

}
