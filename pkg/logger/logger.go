package logger

import (
	"log"
	"os"
)

func Init() {
	// Открываем файл для записи логов
	file, err := os.OpenFile("bot.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Failed to open log file: %v", err)
	}

	// Устанавливаем вывод логов в файл
	log.SetOutput(file)
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}
