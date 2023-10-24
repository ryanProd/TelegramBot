package handler

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"

	_ "github.com/lib/pq"
)

// Constants and Struct for storing State
const stateZero string = "Please review our product."
const stateOne string = "Would you like to leave additional feedback?"
const close string = "Thank you for being a valued customer, if you have additional questions please visit our website https://nike.com. If you would like to reset this bot type /start"

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
const dbHost = "DB_HOST"
const dbName = "DB_NAME"
const dbUser = "DB_USER"
const dbPassword = "DB_PASSWORD"
const dbPort = "DB_PORT"

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

	db := connectDB()
	defer db.Close()

	_, _ = handleUserInput(update, db)
}

// in real implementation this function would call database to get User State and store responses
// we could also send requests to our ai endpoint to integrate ai functionality here
// store User state and send appropriate message back to them
func handleUserInput(update *Update, db *sql.DB) (string, error) {
	userId := update.Message.From.UserId
	userState, err := queryDB(db, userId)
	if err != nil {
		//actual error returned from queryDB function
		if err != sql.ErrNoRows {
			log.Printf(err.Error())
			return "", err
		} else {
			//queryDB returned sql.ErrNoRows which means user is not in database, insert into database
			log.Printf("USER NOT FOUND")
			insertRow(db, userId)
			sendMessage(update.Message.Chat.ChatId, stateZero)
			return "", nil
		}
	}
	//user already exists in database, handle cases depending on current state
	if update.Message.Text == "/start" {
		log.Printf("USER FOUND")
		updateRow(db, userId, stateZero, "", "")
		sendMessage(update.Message.Chat.ChatId, stateZero)
	} else {
		if userState.currState == stateZero {
			updateRow(db, userId, stateOne, update.Message.Text, "")
			sendMessage(update.Message.Chat.ChatId, stateOne)
		} else if userState.currState == stateOne {
			updateRow(db, userId, close, userState.stateZeroReply, update.Message.Text)
			sendMessage(update.Message.Chat.ChatId, close)
		} else {
			sendMessage(update.Message.Chat.ChatId, close)
		}
	}

	return "", err

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

// connect to DB
func connectDB() *sql.DB {
	connStr := fmt.Sprintf("user=%s password=%s dbname=%s host=%s", os.Getenv(dbUser), os.Getenv(dbPassword),
		os.Getenv(dbName), os.Getenv(dbHost))

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}

	var version string
	if err := db.QueryRow("select version()").Scan(&version); err != nil {
		panic(err)
	}

	return db
}

func queryDB(db *sql.DB, userId int) (*UserState, error) {
	var userState UserState
	selectQuery := fmt.Sprintf(`SELECT * FROM userstates WHERE user_id = %d`, userId)
	log.Printf(selectQuery)
	if err := db.QueryRow(selectQuery).Scan(&userState.UserId,
		&userState.currState, &userState.stateZeroReply, &userState.stateOneReply); err != nil {
		if err == sql.ErrNoRows {
			log.Printf(err.Error())
			return nil, err
		}
		return nil, err
	}
	return &userState, nil
}

func insertRow(db *sql.DB, userId int) {
	insertQuery := fmt.Sprintf(`INSERT INTO userstates (user_id, currState, stateZeroReply, stateOneReply)
	    VALUES ('%d', '%s', '%s', '%s')`, userId, stateZero, "", "")
	log.Printf(insertQuery)
	_, err := db.Exec(insertQuery)
	if err != nil {
		log.Printf(err.Error())
		return
	}
	log.Printf("User ID: %d Inserted", userId)
}

func updateRow(db *sql.DB, userId int, currState string, stateZeroReply string, stateOnereply string) {
	updateQuery := fmt.Sprintf(`UPDATE userstates SET currState = '%s', stateZeroReply = '%s', stateOneReply = '%s' 
	    WHERE user_id = '%d'`, currState, stateZeroReply, stateOnereply, userId)
	log.Printf(updateQuery)
	_, err := db.Exec(updateQuery)
	if err != nil {
		log.Printf(err.Error())
		return
	}
	log.Printf("User ID: %d Updated", userId)
}
