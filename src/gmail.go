package main

import (
	"context"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"golang.org/x/oauth2/google"
	_ "log"
	"os"

	"google.golang.org/api/gmail/v1"
	"google.golang.org/api/option"
)

func getGmailService(ctx context.Context) (*gmail.Service, error) {
	b, err := os.ReadFile(credentialsFile)
	if err != nil {
		return nil, fmt.Errorf("unable to read credentials file: %v", err)
	}

	config, err := google.ConfigFromJSON(b, "https://www.googleapis.com/auth/gmail.modify")
	if err != nil {
		return nil, fmt.Errorf("unable to parse credentials file: %v", err)
	}

	tok, err := tokenFromFile(tokenFile)
	if err != nil {
		tok = getTokenFromWeb(config)
		saveToken(tokenFile, tok)
	}

	client := config.Client(ctx, tok)
	srv, err := gmail.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		return nil, fmt.Errorf("unable to create Gmail service: %v", err)
	}
	return srv, nil
}

func markAsRead(srv *gmail.Service, messageID string) error {
	_, err := srv.Users.Messages.Modify(gmailUserID, messageID, &gmail.ModifyMessageRequest{
		RemoveLabelIds: []string{"UNREAD"},
	}).Do()
	if err != nil {
		return fmt.Errorf("unable to mark message as read: %v", err)
	}
	return nil
}

func checkMail(srv *gmail.Service, bot *tgbotapi.BotAPI) error {
	r, err := srv.Users.Messages.List(gmailUserID).MaxResults(1).LabelIds("INBOX").Q("is:unread").Do()
	if err != nil {
		return fmt.Errorf("unable to retrieve messages: %v", err)
	}

	for _, m := range r.Messages {
		msg, err := srv.Users.Messages.Get(gmailUserID, m.Id).Do()
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
		if err := sendMessage(bot, chatID, messageText); err != nil {
			return fmt.Errorf("unable to send message: %v", err)
		}

		if err := markAsRead(srv, m.Id); err != nil {
			return fmt.Errorf("unable to mark message as read: %v", err)
		}
	}

	return nil
}
