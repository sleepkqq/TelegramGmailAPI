package main

import (
	"context"
	"log"
	"time"

	"github.com/go-telegram-bot-api/telegram-bot-api"
)

func main() {
	ctx := context.Background()

	// Настройка Gmail API
	srv, err := getGmailService(ctx)
	if err != nil {
		log.Fatalf("Unable to create Gmail service: %v", err)
	}

	// Настройка Telegram бота
	bot, err := tgbotapi.NewBotAPI(telegramBotToken)
	if err != nil {
		log.Fatalf("Unable to create Telegram bot: %v", err)
	}
	bot.Debug = true

	// Проверка почты каждую минуту
	for {
		if err := checkMail(srv, bot); err != nil {
			log.Printf("Error checking mail: %v", err)
		}
		time.Sleep(1 * time.Minute)
	}
}
