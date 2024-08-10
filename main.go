package main

import (
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
	"telegram-gmail-api/api"
	"telegram-gmail-api/config"
)

func main() {
	ctx := context.Background()

	srv, err := api.GetGmailService(ctx)
	if err != nil {
		log.Fatalf("Unable to create Gmail service: %v", err)
	}

	bot, err := tgbotapi.NewBotAPI(config.TelegramBotToken)
	if err != nil {
		log.Fatalf("Unable to create Telegram bot: %v", err)
	}
	bot.Debug = true

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, _ := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil {
			chatID := update.Message.Chat.ID
			userMessage := update.Message.Text

			switch userMessage {
			case "/send":
				api.InitiateSendProcess(bot, chatID)
			case "/check":
				api.HandleCheckMail(srv, bot, chatID)
			default:
				api.HandleUserState(srv, bot, chatID, userMessage)
			}
		}
	}
}
