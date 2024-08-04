package handlers

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// YandexSpellerResp представляет структуру ответа от API Яндекс Спеллера
type YandexSpellerResp struct {
	Code        int      `json:"code"`
	Pos         int      `json:"pos"`
	Row         int      `json:"row"`
	Col         int      `json:"col"`
	Len         int      `json:"len"`
	Word        string   `json:"word"`
	Suggestions []string `json:"s"`
}

// HandleT9 обрабатывает команду /t9
func HandleT9(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
	// Извлекаем текст команды из сообщения
	text := extractTextFromMessage(bot, message, "/t9")

	if text == "" {
		sendMessage(bot, message, "Предоставьте текст для проверки")
		return
	}

	// Формируем параметры запроса
	data := url.Values{}
	data.Set("text", text)
	data.Set("lang", "ru,en")
	data.Set("options", "6")

	// Отправляем POST-запрос к API Яндекс Спеллера
	resp, err := http.PostForm("https://speller.yandex.net/services/spellservice.json/checkText", data)
	if err != nil {
		sendMessage(bot, message, "Не удалось отправить запрос")
		return
	}
	defer resp.Body.Close()

	// Обрабатываем ответ
	if resp.StatusCode != http.StatusOK {
		sendMessage(bot, message, "Не удалось получить корректный ответ")
		return
	}

	var spellerResponses []YandexSpellerResp
	if err := json.NewDecoder(resp.Body).Decode(&spellerResponses); err != nil {
		sendMessage(bot, message, "Ошибка при обработке ответа")
		return
	}

	// Формируем исправленный текст
	correctedText := text
	for _, response := range spellerResponses {
		if len(response.Suggestions) > 0 {
			incorrectWord := response.Word
			correctWord := response.Suggestions[0]
			correctedText = strings.Replace(correctedText, incorrectWord, correctWord, -1)
		}
	}

	// Отправляем исправленный текст пользователю
	sendMessage(bot, message, correctedText)
}
