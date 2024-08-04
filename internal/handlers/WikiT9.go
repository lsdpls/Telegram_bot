package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// YandexSpellerResponse представляет структуру ответа от API Яндекс Спеллера
type YandexSpellerResponse struct {
	Code        int      `json:"code"`
	Pos         int      `json:"pos"`
	Row         int      `json:"row"`
	Col         int      `json:"col"`
	Len         int      `json:"len"`
	Word        string   `json:"word"`
	Suggestions []string `json:"s"`
}

// WikiSearchResponse представляет структуру ответа от API поиска Википедии
type WikiSearchResponse struct {
	Query struct {
		Search []struct {
			Title   string `json:"title"`
			Snippet string `json:"snippet"`
			PageID  int    `json:"pageid"`
		} `json:"search"`
	} `json:"query"`
}

// WikiSummaryResponse представляет структуру ответа от API Википедии
type WikiSummaryResponse struct {
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

// HandleWikit9 обрабатывает команду /wiki
func HandleWikit9(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
	// Извлекаем текст команды из сообщения
	text := extractTextFromMessage(bot, message, "/wiki")

	if text == "" {
		sendMessage(bot, message, "Предоставьте текст")
		return
	}

	// Сначала отправляем текст на проверку орфографии
	correctedText, err := checkSpelling(text)
	if err != nil {
		sendMessage(bot, message, fmt.Sprintf("Ошибка проверки запроса: %v", err))
		return
	}

	// Затем выполняем поиск в Википедии по исправленному тексту
	searchResults, err := searchWiki(correctedText)
	if err != nil || len(searchResults.Query.Search) == 0 {
		sendMessage(bot, message, fmt.Sprintf("Ошибка поиска в Википедии: %v", err))
		return
	}

	// Получаем первый результат поиска и запрашиваем его краткое содержание
	wikiResponse, err := fetchWikiSummary(searchResults.Query.Search[0].Title)
	if err != nil {
		sendMessage(bot, message, fmt.Sprintf("Ошибка запроса к Википедии: %v", err))
		return
	}

	// Обработка ответа типа "disambiguation"
	if wikiResponse.Type == "disambiguation" {
		msgText := fmt.Sprintf("Тема неоднозначна. Подробнее: <a href=\"%s\">Ссылка</a>", wikiResponse.ContentURLs.Desktop.Page)
		sendMessageParse(bot, message, msgText)
		return
	}

	// Формируем текст сообщения с извлеченным текстом и ссылкой на статью
	msgText := fmt.Sprintf("%s\n\n<a href=\"%s\">Читать далее</a>", wikiResponse.Extract, wikiResponse.ContentURLs.Desktop.Page)

	// Если есть изображение, отправляем его вместе с текстом
	if wikiResponse.Thumbnail.Source != "" {
		sendMessageParsePhoto(bot, message, wikiResponse.Thumbnail.Source, msgText)
	} else {
		sendMessageParse(bot, message, msgText)
	}
}

// checkSpelling проверяет орфографию текста с использованием API Яндекс Спеллера
func checkSpelling(text string) (string, error) {
	data := url.Values{}
	data.Set("text", text)
	data.Set("lang", "ru")

	resp, err := http.PostForm("https://speller.yandex.net/services/spellservice.json/checkText", data)
	if err != nil {
		return "", fmt.Errorf("не удалось отправить запрос: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("не удалось получить корректный ответ")
	}

	var spellerResponses []YandexSpellerResponse
	if err := json.NewDecoder(resp.Body).Decode(&spellerResponses); err != nil {
		return "", fmt.Errorf("ошибка при обработке ответа: %w", err)
	}

	correctedText := text
	for _, response := range spellerResponses {
		if len(response.Suggestions) > 0 {
			incorrectWord := response.Word
			correctWord := response.Suggestions[0]
			correctedText = strings.Replace(correctedText, incorrectWord, correctWord, -1)
		}
	}

	return correctedText, nil
}

// searchWiki выполняет поиск в Википедии по тексту запроса
func searchWiki(query string) (*WikiSearchResponse, error) {
	escapedQuery := strings.ReplaceAll(url.QueryEscape(query), "+", "%20")
	apiURL := fmt.Sprintf("https://ru.wikipedia.org/w/api.php?action=query&list=search&format=json&srsearch=%s", escapedQuery)

	resp, err := http.Get(apiURL)
	if err != nil {
		return nil, fmt.Errorf("не удалось отправить запрос к Википедии: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("не удалось выполнить поиск в Википедии")
	}

	var searchResponse WikiSearchResponse
	if err := json.NewDecoder(resp.Body).Decode(&searchResponse); err != nil {
		return nil, fmt.Errorf("ошибка при обработке ответа от Википедии: %w", err)
	}

	return &searchResponse, nil
}

// fetchWikiSummary запрашивает краткое содержание статьи из Википедии
func fetchWikiSummary(title string) (*WikiSummaryResponse, error) {
	escapedTitle := strings.ReplaceAll(url.QueryEscape(title), "+", "%20")
	apiURL := fmt.Sprintf("https://ru.wikipedia.org/api/rest_v1/page/summary/%s", escapedTitle)

	resp, err := http.Get(apiURL)
	if err != nil {
		return nil, fmt.Errorf("не удалось отправить запрос к Википедии: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("не удалось найти информацию по данной теме")
	}

	var wikiResponse WikiSummaryResponse
	if err := json.NewDecoder(resp.Body).Decode(&wikiResponse); err != nil {
		return nil, fmt.Errorf("ошибка при обработке ответа от Википедии: %w", err)
	}

	return &wikiResponse, nil
}
