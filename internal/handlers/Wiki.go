package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// WikiResponse представляет структуру ответа от API Википедии
type WikiResponse struct {
	Type        string `json:"type"`
	Extract     string `json:"extract"`
	ContentURLs struct {
		Desktop struct {
			Page string `json:"page"`
		} `json:"desktop"`
	} `json:"content_urls"`
	Thumbnail struct {
		Source string `json:"source"`
	} `json:"thumbnail"`
}

// Deprecated: is no longer used
// HandleWiki обрабатывает команду /wiki
func HandleWiki(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
	// Извлекаем текст команды из сообщения
	text := extractTextFromMessage(bot, message, "/wiki")

	if text == "" {
		sendMessage(bot, message, "Укажите тему для поиска в Википедии")
		return
	}

	// Явно заменяем пробелы на %20
	escapedText := strings.ReplaceAll(url.QueryEscape(text), "+", "%20")

	// Формируем URL для запроса к API Википедии
	apiURL := fmt.Sprintf("https://ru.wikipedia.org/api/rest_v1/page/summary/%s", escapedText)

	// Отправляем GET-запрос
	resp, err := http.Get(apiURL)
	if err != nil {
		sendMessage(bot, message, "Не удалось отправить запрос к Википедии")
		return
	}
	defer resp.Body.Close()

	// Обрабатываем ответ
	if resp.StatusCode != http.StatusOK {
		sendMessage(bot, message, "Не удалось найти информацию по данной теме")
		return
	}

	var wikiResponse WikiResponse
	if err := json.NewDecoder(resp.Body).Decode(&wikiResponse); err != nil {
		sendMessage(bot, message, "Ошибка при обработке ответа от Википедии")
		return
	}

	// Обработка ответа типа "disambiguation"
	if wikiResponse.Type == "disambiguation" {
		msgText := fmt.Sprintf("Тема неоднозначна. [Подробнее](%s)", wikiResponse.ContentURLs.Desktop.Page)
		msg := tgbotapi.NewMessage(message.Chat.ID, msgText)
		msg.ParseMode = "Markdown"
		msg.ReplyToMessageID = message.MessageID
		bot.Send(msg)
		return
	}

	// Формируем текст сообщения с извлеченным текстом и ссылкой на статью
	msgText := fmt.Sprintf("%s\n\n[Читать далее](%s)", wikiResponse.Extract, wikiResponse.ContentURLs.Desktop.Page)

	// Если есть изображение, отправляем его вместе с текстом
	if wikiResponse.Thumbnail.Source != "" {
		photo := tgbotapi.NewPhoto(message.Chat.ID, tgbotapi.FileURL(wikiResponse.Thumbnail.Source))
		photo.Caption = msgText
		photo.ParseMode = "Markdown"
		photo.ReplyToMessageID = message.MessageID
		bot.Send(photo)
	} else {
		msg := tgbotapi.NewMessage(message.Chat.ID, msgText)
		msg.ParseMode = "Markdown"
		msg.ReplyToMessageID = message.MessageID
		bot.Send(msg)
	}
}
