package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"telegram-gmail-api/config"
	"telegram-gmail-api/models"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"google.golang.org/api/gmail/v1"
	"gorm.io/gorm"
)

func HandleCheckMail(srv *gmail.Service, bot *tgbotapi.BotAPI, chatID int64) {
	if err := checkMail(srv, bot); err != nil {
		log.Printf("Error checking mail: %v", err)
		SendMessage(bot, chatID, "Failed to check mail.")
	} else {
		SendMessage(bot, chatID, "Mail checked successfully.")
	}
}

func InitiateSendProcess(bot *tgbotapi.BotAPI, chatID int64) {
	user := models.User{ChatID: chatID, State: config.StateAwaitingRecipient, Data: "{}"}
	config.DB.Save(&user)
	SendMessage(bot, chatID, "Please provide the recipient's email address.")
}

func HandleUserState(srv *gmail.Service, bot *tgbotapi.BotAPI, chatID int64, userMessage string) {
	var user models.User
	if err := config.DB.Where("chat_id = ?", chatID).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			InitiateSendProcess(bot, chatID)
			return
		}
		log.Printf("Error fetching user state: %v", err)
		return
	}

	var data map[string]string
	if err := json.Unmarshal([]byte(user.Data), &data); err != nil {
		log.Printf("Error unmarshaling user data: %v", err)
		return
	}

	switch user.State {
	case config.StateAwaitingRecipient:
		data[Recipient] = userMessage
		user.State = config.StateAwaitingTitle

	case config.StateAwaitingTitle:
		data[Subject] = userMessage
		user.State = config.StateAwaitingBody

	case config.StateAwaitingBody:
		data[Body] = userMessage
		if err := sendMail(srv, data[Recipient], data[Subject], data[Body]); err != nil {
			SendMessage(bot, chatID, fmt.Sprintf("Failed to send email: %v", err))
		} else {
			SendMessage(bot, chatID, "Email sent successfully!")
		}

		user.State = Completed
	}

	updatedData, _ := json.Marshal(data)
	user.Data = string(updatedData)
	if err := config.DB.Save(&user).Error; err != nil {
		log.Printf("Error updating user: %v", err)
	}

	SendMessage(bot, chatID, nextPromptMessage(user.State))
}

func SendMessage(bot *tgbotapi.BotAPI, chatID int64, text string) {
	msg := tgbotapi.NewMessage(chatID, text)
	_, err := bot.Send(msg)
	if err != nil {
		log.Printf("Error sending message: %v", err)
	}
}

func nextPromptMessage(state string) string {
	switch state {
	case config.StateAwaitingTitle:
		return "Please provide the email title."
	case config.StateAwaitingBody:
		return "Please provide the email body."
	default:
		return ""
	}
}
