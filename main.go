package main

import (
	"log"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	env "github.com/joho/godotenv"
)

func init() {
	err := env.Load(".env")
	if err != nil {
		log.Println("Error loading .env file")
	}
}

func main() {
	botToken := os.Getenv("TG_BOT_TOKEN")

	bot, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		log.Panic(err)
	}

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil {
			continue
		}

		if update.Message.IsCommand() {
			handleCommand(update.Message, bot)
		}
	}

}

func handleCommand(message *tgbotapi.Message, bot *tgbotapi.BotAPI) {
	switch message.Command() {
	case "start":
		startCommand(message, bot)
	}
}

func startCommand(message *tgbotapi.Message, bot *tgbotapi.BotAPI) {
	msg := tgbotapi.NewMessage(message.Chat.ID, "Welcome!")
	bot.Send(msg)
}
