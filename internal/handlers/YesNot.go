package handlers

import (
	"encoding/json"
	"math/rand"
	"net/http"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// YesNoResponse представляет структуру ответа от API yesno.wtf
type YesNoResponse struct {
	Answer string `json:"answer"`
	Image  string `json:"image"`
}

// HandleYesNoT обрабатывает команду /yesnot
func HandleYesNoT(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
	var apiURL string
	if rand.Intn(10) == 0 {
		apiURL = "https://yesno.wtf/api?force=maybe"
	} else {
		apiURL = "https://yesno.wtf/api"
	}

	resp, err := http.Get(apiURL)
	if err != nil {
		sendMessage(bot, message, "Не удалось получить ответ от сервиса")
		return
	}
	defer resp.Body.Close()

	var yesNoResponse YesNoResponse
	if err := json.NewDecoder(resp.Body).Decode(&yesNoResponse); err != nil {
		sendMessage(bot, message, "Ошибка при обработке ответа сервиса")
		return
	}

	// Отправляем гифку
	gif := tgbotapi.NewAnimation(message.Chat.ID, tgbotapi.FileURL(yesNoResponse.Image))
	gif.ReplyToMessageID = message.MessageID
	bot.Send(gif)
}
