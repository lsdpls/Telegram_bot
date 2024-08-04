package handlers

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var (
	// Знаки препинания
	symbols = []string{".", ",", "!", "?", ";", ":", " - ", "—", `"`, "'", "[", "]", "(", ")"}

	// Союзы и др
	conjunctions = []string{" и ", " или ", " но ", "однако", " а ", " да ", " то ", " же ", "либо ", "чтобы ", "когда ", "как ", "только", "что ", "это ", " так ", " ну ", "где ", "если ", "короче"}

	// Частицы
	particles = map[string]struct{}{
		"на": {}, "с": {}, "о": {}, "в": {}, "к": {}, "не": {}, "у": {}, "при": {}, "ни": {}, "до": {}, "через": {},
		"над": {}, "под": {}, "перед": {}, "за": {}, "возле": {}, "мимо": {}, "после": {}, "от": {}, "для": {}, "про": {},
		"я": {},
	}
)

// MemeResponse представляет структуру ответа от API imgflip
type MemeResponse struct {
	Success bool `json:"success"`
	Data    struct {
		Memes []struct {
			ID       string `json:"id"`
			Name     string `json:"name"`
			URL      string `json:"url"`
			Width    int    `json:"width"`
			Height   int    `json:"height"`
			BoxCount int    `json:"box_count"`
		} `json:"memes"`
	} `json:"data"`
}

// CaptionResponse представляет структуру ответа на запрос caption_image от API imgflip
type CaptionResponse struct {
	Success bool `json:"success"`
	Data    struct {
		URL     string `json:"url"`
		PageURL string `json:"page_url"`
	} `json:"data"`
}

// HandleMeme обрабатывает команду /meme
func HandleMeme(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
	// Получаем список мемов
	resp, err := http.Get("https://api.imgflip.com/get_memes")
	if err != nil {
		sendMessage(bot, message, "Не удалось получить список мемов")
		return
	}
	defer resp.Body.Close()

	var memeResponse MemeResponse
	if err := json.NewDecoder(resp.Body).Decode(&memeResponse); err != nil {
		sendMessage(bot, message, "Ошибка при обработке списка мемов")
		return
	}

	if !memeResponse.Success || len(memeResponse.Data.Memes) == 0 {
		sendMessage(bot, message, "Не удалось обработать список мемов")
		return
	}

	// Выбираем случайный мем
	randomMeme := memeResponse.Data.Memes[rand.Intn(len(memeResponse.Data.Memes))]

	// Извлекаем текст из сообщения и удаляем упоминания бота и команды
	var text string
	text = extractTextFromMessage(bot, message, "/meme")
	text = strings.ToLower(text)

	// Разбиваем текст на осмысленные фразы
	phrases := splitIntoPhrases(text)

	// Выбираем случайные фразы для заполнения текстовых полей мема
	var selectedPhrases []string
	for i := 0; i < randomMeme.BoxCount && len(phrases) > 0; i++ {
		index := rand.Intn(len(phrases))
		selectedPhrases = append(selectedPhrases, phrases[index])
		phrases = append(phrases[:index], phrases[index+1:]...)
	}

	// Получаем логин и пароль из переменных окружения
	username := os.Getenv("IMGFLIP_USERNAME")
	password := os.Getenv("IMGFLIP_PASSWORD")

	// Формируем параметры для запроса создания мема
	params := url.Values{
		"template_id": {randomMeme.ID},
		"username":    {username},
		"password":    {password},
	}

	// Добавляем текстовые поля в формате boxes
	for i, phrase := range selectedPhrases {
		params.Add(fmt.Sprintf("boxes[%d][text]", i), phrase)
	}

	// Отправляем запрос на создание мема
	resp, err = http.PostForm("https://api.imgflip.com/caption_image", params)
	if err != nil {
		sendMessage(bot, message, "Не удалось создать мем")
		return
	}
	defer resp.Body.Close()

	var captionResponse CaptionResponse
	if err := json.NewDecoder(resp.Body).Decode(&captionResponse); err != nil || !captionResponse.Success {
		sendMessage(bot, message, "Ошибка при создании мема")
		return
	}

	// Отправляем созданный мем (картинку)
	sendMessagePhoto(bot, message, captionResponse.Data.URL)
}

// splitIntoPhrases разбивает текст на осмысленные фразы
func splitIntoPhrases(text string) []string {
	// Заменяем знаки препинания на разделители
	for _, symbol := range symbols {
		text = strings.ReplaceAll(text, symbol, " | ")
	}

	// Заменяем союзы и другие слова на разделители
	for _, conj := range conjunctions {
		text = strings.ReplaceAll(text, conj, " | ")
	}

	// Разделяем по разделителям
	phrases := strings.Split(text, "|")

	// Удаляем начальные и конечные пробелы в каждой фразе
	for i := range phrases {
		phrases[i] = strings.TrimSpace(phrases[i])
	}

	// Убираем пустые фразы
	var result []string
	for _, phrase := range phrases {
		if phrase != "" {
			result = append(result, phrase)
		}
	}

	// Разделяем слишком длинные фразы
	var finalResult []string
	maxLength := 40 // Максимальная длина фразы
	for _, phrase := range result {
		if len(phrase) > maxLength {
			words := strings.Fields(phrase)
			var currentPhrase string
			var tmpWord string
			for _, word := range words {
				if len(currentPhrase)+len(word)+1 > maxLength {
					if _, found := particles[tmpWord]; found {
						currentPhrase += " " + word
						finalResult = append(finalResult, currentPhrase)
						currentPhrase = ""
						continue
					}
					finalResult = append(finalResult, currentPhrase)
					currentPhrase = word
				} else {
					if currentPhrase != "" {
						currentPhrase += " "
					}
					currentPhrase += word
					tmpWord = word
				}
			}
			if currentPhrase != "" {
				finalResult = append(finalResult, currentPhrase)
			}
		} else {
			finalResult = append(finalResult, phrase)
		}
	}

	return finalResult
}
