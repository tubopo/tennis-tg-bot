package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Training struct {
	Date  time.Time
	Place string
	Level string
}

var (
	bot          *tgbotapi.BotAPI
	trainings    []Training
	userState    map[int64]string
	userTraining map[int64]*Training
)

func main() {
	token := os.Getenv("TG_BOT_TOKEN")
	if token == "" {
		log.Fatal("TG_BOT_TOKEN environment variable is not set")
	}

	var err error
	bot, err = tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	userState = make(map[int64]string)
	userTraining = make(map[int64]*Training)

	for update := range updates {
		if update.Message != nil {
			handleMessage(update.Message)
		} else if update.CallbackQuery != nil {
			handleCallbackQuery(update.CallbackQuery)
		}
	}
}

func handleMessage(message *tgbotapi.Message) {
	chatID := message.Chat.ID

	switch message.Command() {
	case "start":
		sendMessage(chatID, "Welcome! Use /new_training to schedule a new training.")
	case "new_training":
		userState[chatID] = "awaiting_place"
		userTraining[chatID] = &Training{}
		sendPlaceSelection(chatID)
	default:
		sendMessage(chatID, "I'm not sure what you mean. Use /new_training to schedule a new training.")
	}
}

func handleCallbackQuery(query *tgbotapi.CallbackQuery) {
	chatID := query.Message.Chat.ID
	messageID := query.Message.MessageID
	data := query.Data

	switch userState[chatID] {
	case "awaiting_place":
		userTraining[chatID].Place = data
		userState[chatID] = "awaiting_time"
		editMessage(chatID, messageID, fmt.Sprintf("You selected: %s", data))
		showAvailableTimeSlots(chatID)

	case "awaiting_time":
		parts := strings.Split(data, "|")
		if len(parts) != 2 {
			log.Printf("Invalid time slot data: %s", data)
			return
		}
		level := parts[0]
		timeSlot := parts[1]

		training := userTraining[chatID]
		training.Level = level
		timeRange := strings.Split(timeSlot, "-")
		startTime, _ := time.Parse("15:04", timeRange[0])

		trainingDate := time.Now().Local()
		trainingDateTime := time.Date(
			trainingDate.Year(), trainingDate.Month(), trainingDate.Day(),
			startTime.Hour(), startTime.Minute(), 0, 0, trainingDate.Location(),
		)

		userTraining[chatID].Date = trainingDateTime
		editMessage(chatID, messageID, fmt.Sprintf("You selected: %s (%s)", timeSlot, level))
		confirmTraining(chatID, userTraining[chatID])
	}

	// Answer the callback query
	callback := tgbotapi.NewCallback(query.ID, "")
	if _, err := bot.Request(callback); err != nil {
		log.Printf("Error answering callback query: %v", err)
	}
}

func sendPlaceSelection(chatID int64) {
	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Смолячкова, 9", "Смолячкова, 9"),
			tgbotapi.NewInlineKeyboardButtonData("Ленина, 27", "Ленина, 27"),
		),
	)

	msg := tgbotapi.NewMessage(chatID, "Please select a place for the training:")
	msg.ReplyMarkup = keyboard
	bot.Send(msg)
}

func showAvailableTimeSlots(chatID int64) {
	training := userTraining[chatID]

	currentDate := time.Now().Local()
	dayOfWeek := currentDate.Weekday()

	var keyboardRows [][]tgbotapi.InlineKeyboardButton

	if training.Place == "Смолячкова, 9" {
		if dayOfWeek == time.Tuesday || dayOfWeek == time.Wednesday || dayOfWeek == time.Thursday {
			keyboardRows = append(keyboardRows,
				tgbotapi.NewInlineKeyboardRow(
					tgbotapi.NewInlineKeyboardButtonData("Уровень 1: 20:40-22:10", "Уровень 1|20:40-22:10"),
				),
			)
		}
		if dayOfWeek >= time.Monday && dayOfWeek <= time.Friday {
			keyboardRows = append(keyboardRows,
				tgbotapi.NewInlineKeyboardRow(
					tgbotapi.NewInlineKeyboardButtonData("Уровень 2: 19:00-20:30", "Уровень 2|19:00-20:30"),
				),
			)
		}
	} else if training.Place == "Ленина, 27" {
		if dayOfWeek == time.Tuesday || dayOfWeek == time.Thursday {
			keyboardRows = append(keyboardRows,
				tgbotapi.NewInlineKeyboardRow(
					tgbotapi.NewInlineKeyboardButtonData("Уровень 2: 19:00-20:30", "Уровень 2|19:00-20:30"),
				),
			)
		}
	}

	if len(keyboardRows) == 0 {
		sendMessage(chatID, "No time slots available for the selected place and day. Please try another day.")
		return
	}

	keyboard := tgbotapi.NewInlineKeyboardMarkup(keyboardRows...)

	msg := tgbotapi.NewMessage(chatID, fmt.Sprintf("Select a time slot for %s at %s:", currentDate.Format("02-01-2006"), training.Place))
	msg.ReplyMarkup = keyboard
	bot.Send(msg)

	userState[chatID] = "awaiting_time"
	userTraining[chatID].Date = currentDate
}

func confirmTraining(chatID int64, training *Training) {
	msg := fmt.Sprintf("Training confirmed:\nPlace: %s\nLevel: %s\nDate and Time: %s",
		training.Place, training.Level, training.Date.Format("02-01-2006 15:04"))
	sendMessage(chatID, msg)

	trainings = append(trainings, *training)

	delete(userState, chatID)
	delete(userTraining, chatID)
}

func sendMessage(chatID int64, text string) {
	msg := tgbotapi.NewMessage(chatID, text)
	if _, err := bot.Send(msg); err != nil {
		log.Printf("Error sending message: %v", err)
	}
}

func editMessage(chatID int64, messageID int, text string) {
	editMsg := tgbotapi.NewEditMessageText(chatID, messageID, text)
	editMsg.ReplyMarkup = &tgbotapi.InlineKeyboardMarkup{
		InlineKeyboard: [][]tgbotapi.InlineKeyboardButton{},
	}
	if _, err := bot.Send(editMsg); err != nil {
		log.Printf("Error editing message: %v", err)
	}
}
