package utils

import (
	"context"
	"encoding/json"
	"fmt"
	"golang.org/x/oauth2"
	"log"
	"os"
)

func TokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, fmt.Errorf("unable to open token file: %v", err)
	}
	defer f.Close()

	tok := &oauth2.Token{}
	if err := json.NewDecoder(f).Decode(tok); err != nil {
		return nil, fmt.Errorf("unable to decode token file: %v", err)
	}
	return tok, nil
}

func GetTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser:\n%v\n", authURL)

	var authCode string
	fmt.Print("Enter the authorization code: ")
	if _, err := fmt.Scan(&authCode); err != nil {
		log.Fatalf("Unable to read authorization code: %v", err)
	}

	tok, err := config.Exchange(context.Background(), authCode)
	if err != nil {
		log.Fatalf("Unable to retrieve token: %v", err)
	}
	return tok
}

func SaveToken(file string, token *oauth2.Token) {
	f, err := os.Create(file)
	if err != nil {
		log.Fatalf("Unable to create token file: %v", err)
	}
	defer f.Close()

	json.NewEncoder(f).Encode(token)
}
