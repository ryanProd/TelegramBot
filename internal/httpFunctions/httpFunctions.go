package httpFunctions

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"

	"github.com/ryanProd/TelegramBot/internal/structs"
)

// Destructure Http Request(Message from User) that contains Telegram Update JSON into our structs
func DecodeRequest(r *http.Request) (*structs.Update, error) {
	var update structs.Update
	if err := json.NewDecoder(r.Body).Decode(&update); err != nil {
		log.Printf("could not decode update %s", err.Error())
		return nil, err
	}
	return &update, nil
}

const telegramApiBaseUrl string = "https://api.telegram.org/bot"
const telegramApiSendMessage string = "/sendMessage"
const telegramTokenEnv string = "BOT_TOKEN"

var telegramEndpoint string = telegramApiBaseUrl + os.Getenv(telegramTokenEnv) + telegramApiSendMessage

// Send Message with http Post back to User
func PostMessage(chatId int, text string) (string, error) {
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
