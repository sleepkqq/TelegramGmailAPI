package main

import (
	"github.com/go-telegram-bot-api/telegram-bot-api"
)

func sendMessage(bot *tgbotapi.BotAPI, chatID int64, messageText string) error {
	msgToSend := tgbotapi.NewMessage(chatID, messageText)
	if _, err := bot.Send(msgToSend); err != nil {
		return err
	}
	return nil
}
