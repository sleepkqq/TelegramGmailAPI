# Telegram Bot for Gmail Notifications

This project is a Telegram bot that monitors a Gmail inbox and sends notifications to a Telegram chat whenever a new email arrives. The bot is built using Go and leverages the Gmail API for reading emails and the Telegram Bot API for sending messages.

## Features

- **Email Monitoring**: Continuously checks for new, unread emails in the Gmail inbox.
- **Real-time Notifications**: Sends Telegram messages with details of new emails as soon as they arrive.
- **Mark as Read**: Automatically marks emails as read in Gmail after sending the notification.

## Technologies Used

- **Go**: Programming language used for building the bot.
- **Telegram Bot API**: For sending notifications to a Telegram chat.
- **Gmail API**: For accessing and managing emails.
- **OAuth 2.0**: For secure authorization and access to Gmail.

## Setup

### Prerequisites

- Go installed on your system.
- A Google account and Telegram account.

### Getting Started

1. **Clone the Repository**:
   ```bash
   git clone <repository-url>
   cd <repository-folder>

2. **Set Up Credentials**:
   Place your `credentials.json` file**:
   - Ensure that your `credentials.json` file is located in the `src/credentials` directory. This file contains your OAuth 2.0 credentials from the Google Cloud Console.
   - The bot uses OAuth 2.0 for authentication. If you haven't created a `credentials.json` file, follow [this guide](https://developers.google.com/identity/protocols/oauth2) to generate one.

3. **Configure Your Environment**:
   - Place your `token.json` file, which contains the OAuth 2.0 token, in the `src/credentials` directory. If you don't have this file, it will be automatically created when you run the bot for the first time.
   - Ensure your `config.go` file contains the correct configuration for your bot and that it is placed in the root directory (e.g., `src/config.go`).

## Install Dependencies

Run the following command to install all necessary Go modules:

```bash
go mod tidy
```

## Run the Bot

To start the bot, run the following command:

```bash
go run src/main.go
```

## Usage

- The bot will start monitoring your Gmail inbox as soon as it is launched.
- It will send a Telegram message to the specified chat ID for each new email received in the "Primary" inbox.
- The emails will be marked as read after the notification is sent.

## License

This project is licensed under the MIT License.