package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"math/rand"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Training struct {
	Date        time.Time
	Place       string
	Level       string
	Participant string
}

type UserState struct {
	State           string
	CurrentTraining *Training
}

var (
	bot *tgbotapi.BotAPI
	// storage for the confirmed trainings
	trainings map[int64]Training
	// registration state for the current user, might not be completed
	userStates map[int64]*UserState
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

	trainings = make(map[int64]Training)
	userStates = make(map[int64]*UserState)

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
		userStates[chatID] = &UserState{
			State:           "awaiting_place",
			CurrentTraining: &Training{}, // init empty training
		}
		sendPlaceSelection(chatID)
	case "view_trainings":
		viewTrainings(chatID)
	default:
		sendMessage(chatID, "I'm not sure what you mean. Use /new_training to schedule a new training.")
	}
}

func handleCallbackQuery(query *tgbotapi.CallbackQuery) {
	chatID := query.Message.Chat.ID
	messageID := query.Message.MessageID
	data := query.Data

	userState, exists := userStates[chatID]
	if !exists {
		log.Printf("No user state for user %d", chatID)
		return
	}

	switch userState.State {
	case "awaiting_place":
		userState.CurrentTraining.Place = data
		userState.State = "awaiting_time"
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

		userState.CurrentTraining.Level = level
		timeRange := strings.Split(timeSlot, "-")
		startTime, _ := time.Parse("15:04", timeRange[0])

		trainingDate := time.Now().Local()
		trainingDateTime := time.Date(
			trainingDate.Year(), trainingDate.Month(), trainingDate.Day(),
			startTime.Hour(), startTime.Minute(), 0, 0, trainingDate.Location(),
		)

		userState.CurrentTraining.Date = trainingDateTime
		editMessage(chatID, messageID, fmt.Sprintf("You selected: %s (%s)", timeSlot, level))
		confirmTraining(chatID)
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

func viewTrainings(chatID int64) {
	if len(trainings) == 0 {
		sendMessage(chatID, "No trainings scheduled yet.")
		return
	}

	groupedTrainings := make(map[string]map[string][]string)
	for _, training := range trainings {
		if groupedTrainings[training.Place] == nil {
			groupedTrainings[training.Place] = make(map[string][]string)
		}
		timeSlot := training.Date.Format("02-01-2006 15:04")
		groupedTrainings[training.Place][timeSlot] = append(groupedTrainings[training.Place][timeSlot], training.Participant)
	}

	var message strings.Builder
	for place, timeSlots := range groupedTrainings {
		message.WriteString(fmt.Sprintf("📍 %s:\n", place))
		for timeSlot, participants := range timeSlots {
			message.WriteString(fmt.Sprintf("  🕒 %s:\n", timeSlot))
			for _, participant := range participants {
				message.WriteString(fmt.Sprintf("    + %s\n", participant))
			}
		}
		message.WriteString("\n")
	}

	sendMessage(chatID, message.String())
}

func showAvailableTimeSlots(chatID int64) {
	userState, exists := userStates[chatID]
	if !exists || userState.CurrentTraining == nil {
		log.Printf("No current training for user %d", chatID)
		return
	}

	currentDate := time.Now().Local()
	dayOfWeek := currentDate.Weekday()

	var keyboardRows [][]tgbotapi.InlineKeyboardButton

	if userState.CurrentTraining.Place == "Смолячкова, 9" {
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
	} else if userState.CurrentTraining.Place == "Ленина, 27" {
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

	msg := tgbotapi.NewMessage(chatID, fmt.Sprintf("Select a time slot for %s at %s:", currentDate.Format("02-01-2006"), userState.CurrentTraining.Place))
	msg.ReplyMarkup = keyboard
	bot.Send(msg)

	userState.State = "awaiting_time"
}

func confirmTraining(chatID int64) {
	userState, exists := userStates[chatID]
	if !exists || userState.CurrentTraining == nil {
		log.Printf("No current training for user %d", chatID)
		return
	}

	user, err := bot.GetChat(tgbotapi.ChatInfoConfig{ChatConfig: tgbotapi.ChatConfig{ChatID: chatID}})
	if err != nil {
		log.Printf("Error getting user info: %v", err)
		return
	}

	userState.CurrentTraining.Participant = getDisplayName(user)

	trainings[chatID] = *userState.CurrentTraining // Add or update the training

	msg := fmt.Sprintf("Training registration confirmed:\nPlace: %s\nLevel: %s\nDate and Time: %s\nParticipant: %s",
		userState.CurrentTraining.Place, userState.CurrentTraining.Level,
		userState.CurrentTraining.Date.Format("02-01-2006 15:04"), userState.CurrentTraining.Participant)

	sendMessage(chatID, msg)

	delete(userStates, chatID)
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

var userNameCache = make(map[int64]string)

func getDisplayName(user tgbotapi.Chat) string {
	if name, exists := userNameCache[user.ID]; exists {
		return name
	}

	var displayName string
	if user.FirstName != "" && user.LastName != "" {
		displayName = fmt.Sprintf("%s %s", user.FirstName, user.LastName)
	} else if user.FirstName != "" {
		displayName = user.FirstName
	} else if user.UserName != "" {
		displayName = user.UserName
	} else {
		displayName = generateRandomName()
	}

	userNameCache[user.ID] = displayName
	return displayName
}

var colors = []string{
	"Red", "Blue", "Green", "Yellow", "Purple", "Orange", "Pink", "Brown", "Gray", "Cyan",
}

var animals = []string{
	"Elephant", "Tiger", "Lion", "Giraffe", "Zebra", "Kangaroo", "Penguin", "Dolphin", "Koala", "Panda",
}

func generateRandomName() string {
	color := colors[rand.Intn(len(colors))]
	animal := animals[rand.Intn(len(animals))]
	return fmt.Sprintf("Unknown%s%s", color, animal)
}
