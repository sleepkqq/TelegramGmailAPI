package api

import (
	"context"
	"encoding/base64"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/gmail/v1"
	"google.golang.org/api/option"
	"os"
	"strings"
	"telegram-gmail-api/config"
	"telegram-gmail-api/utils"
)

func GetGmailService(ctx context.Context) (*gmail.Service, error) {
	b, err := os.ReadFile(config.CredentialsFile)
	if err != nil {
		return nil, fmt.Errorf("unable to read credentials file: %v", err)
	}

	googleConfig, err := google.ConfigFromJSON(b, "https://www.googleapis.com/auth/gmail.modify")
	if err != nil {
		return nil, fmt.Errorf("unable to parse credentials file: %v", err)
	}

	tok, err := utils.TokenFromFile(config.TokenFile)
	if err != nil {
		tok = utils.GetTokenFromWeb(googleConfig)
		utils.SaveToken(config.TokenFile, tok)
	}

	client := googleConfig.Client(ctx, tok)
	srv, err := gmail.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		return nil, fmt.Errorf("unable to create Gmail service: %v", err)
	}
	return srv, nil
}

func sendMail(srv *gmail.Service, recipient, subject, body string) error {
	header := make(map[string]string)
	header["From"] = config.Sender
	header["To"] = recipient
	header["Subject"] = subject
	header["MIME-Version"] = "1.0"
	header["Content-Type"] = "text/plain; charset=\"UTF-8\""

	var msg strings.Builder
	for k, v := range header {
		msg.WriteString(k + ": " + v + "\r\n")
	}
	msg.WriteString("\r\n" + body)

	rawMessage := base64.URLEncoding.EncodeToString([]byte(msg.String()))

	message := &gmail.Message{
		Raw: rawMessage,
	}

	_, err := srv.Users.Messages.Send("me", message).Do()
	if err != nil {
		return fmt.Errorf("unable to send email: %v", err)
	}

	return nil
}

func markAsRead(srv *gmail.Service, messageID string) error {
	_, err := srv.Users.Messages.Modify(config.GmailUserID, messageID, &gmail.ModifyMessageRequest{
		RemoveLabelIds: []string{"UNREAD"},
	}).Do()
	if err != nil {
		return fmt.Errorf("unable to mark message as read: %v", err)
	}
	return nil
}

func checkMail(srv *gmail.Service, bot *tgbotapi.BotAPI) error {
	r, err := srv.Users.Messages.List(config.GmailUserID).MaxResults(1).LabelIds("INBOX").Q("is:unread").Do()
	if err != nil {
		return fmt.Errorf("unable to retrieve messages: %v", err)
	}

	for _, m := range r.Messages {
		msg, err := srv.Users.Messages.Get(config.GmailUserID, m.Id).Do()
		if err != nil {
			return fmt.Errorf("unable to retrieve message: %v", err)
		}

		subject := ""
		for _, header := range msg.Payload.Headers {
			if header.Name == "Subject" {
				subject = header.Value
				break
			}
		}

		messageText := fmt.Sprintf("New message: %s\n%s", subject, msg.Snippet)
		SendMessage(bot, config.ChatID, messageText)

		if err := markAsRead(srv, m.Id); err != nil {
			return fmt.Errorf("unable to mark message as read: %v", err)
		}
	}

	return nil
}
