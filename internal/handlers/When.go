package handlers

import (
	"fmt"
	"math/rand"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// HandleWhen обрабатывает команду /when
func HandleWhen(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
	randomHours := rand.Intn(24)
	randomMinutes := rand.Intn(60)
	randomSeconds := rand.Intn(60)
	response := "через " +
		fmt.Sprintf("%d ч. ", randomHours) +
		fmt.Sprintf("%d мин. ", randomMinutes) +
		fmt.Sprintf("%d сек.", randomSeconds)

	sendMessage(bot, message, response)
}
