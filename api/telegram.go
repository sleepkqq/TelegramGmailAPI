package api

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"google.golang.org/api/gmail/v1"
	"log"
	"telegram-gmail-api/config"
)

func InitiateSendProcess(bot *tgbotapi.BotAPI, chatID int64) {
	config.UserStates[chatID] = config.StateAwaitingRecipient
	config.UserData[chatID] = make(map[string]string)
	sendMessage(bot, chatID, "Please provide the recipient's email address.")
}

func HandleCheckMail(srv *gmail.Service, bot *tgbotapi.BotAPI, chatID int64) {
	if err := checkMail(srv, bot); err != nil {
		log.Printf("Error checking mail: %v", err)
		sendMessage(bot, chatID, "Failed to check mail.")
	} else {
		sendMessage(bot, chatID, "Mail checked successfully.")
	}
}

func HandleUserState(srv *gmail.Service, bot *tgbotapi.BotAPI, chatID int64, userMessage string) {
	if state, exists := config.UserStates[chatID]; exists {
		switch state {
		case config.StateAwaitingRecipient:
			config.UserData[chatID]["recipient"] = userMessage
			config.UserStates[chatID] = config.StateAwaitingTitle
			sendMessage(bot, chatID, "Please provide the email title.")

		case config.StateAwaitingTitle:
			config.UserData[chatID]["subject"] = userMessage
			config.UserStates[chatID] = config.StateAwaitingBody
			sendMessage(bot, chatID, "Please provide the email body.")

		case config.StateAwaitingBody:
			config.UserData[chatID]["body"] = userMessage
			err := sendMail(srv, config.UserData[chatID]["recipient"], config.UserData[chatID]["subject"], config.UserData[chatID]["body"])
			if err != nil {
				sendMessage(bot, chatID, fmt.Sprintf("Failed to send email: %v", err))
			} else {
				sendMessage(bot, chatID, "Email sent successfully!")
			}

			delete(config.UserStates, chatID)
			delete(config.UserData, chatID)
		}
	}
}

func sendMessage(bot *tgbotapi.BotAPI, chatID int64, text string) {
	msg := tgbotapi.NewMessage(chatID, text)
	_, err := bot.Send(msg)
	if err != nil {
		return
	}
}
