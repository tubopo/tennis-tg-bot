package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"math/rand"

	"github.com/tubopo/tennis-tg-bot/translations"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Training struct {
	Date        time.Time
	Place       string
	Level       string
	Participant string
}

type UserState struct {
	UserID          int64
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
		sendMessage(chatID, "welcome")
	case "new_training":
		userStates[chatID] = &UserState{
			UserID:          message.From.ID,
			State:           "awaiting_place",
			CurrentTraining: &Training{},
		}
		sendPlaceSelection(chatID)
	case "view_trainings":
		viewTrainings(chatID)
	default:
		sendMessage(chatID, "unknown_command")
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
		editMessage(chatID, messageID, "selected_opt_1", data)
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
		editMessage(chatID, messageID, "selected_opt_2", timeSlot, level)
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
			tgbotapi.NewInlineKeyboardButtonData(translations.Get("btn_place_1"), translations.Get("btn_place_1")),
			tgbotapi.NewInlineKeyboardButtonData(translations.Get("btn_place_2"), translations.Get("btn_place_2")),
		),
	)

	msg := tgbotapi.NewMessage(chatID, translations.Get("select_place"))
	msg.ReplyMarkup = keyboard
	bot.Send(msg)
}

func viewTrainings(chatID int64) {
	if len(trainings) == 0 {
		sendMessage(chatID, "no_trainings")
		return
	}

	groupedTrainings := make(map[string]map[string][]string)
	for _, training := range trainings {
		if groupedTrainings[training.Place] == nil {
			groupedTrainings[training.Place] = make(map[string][]string)
		}
		timeSlot := training.Date.Format("02.01.2006 15:04")
		groupedTrainings[training.Place][timeSlot] = append(groupedTrainings[training.Place][timeSlot], training.Participant)
	}

	var message strings.Builder
	message.WriteString(translations.Get("view_trainings_header") + "\n\n")

	for place, timeSlots := range groupedTrainings {
		message.WriteString(fmt.Sprintf("%s %s:\n", translations.Get("place"), place))
		for timeSlot, participants := range timeSlots {
			message.WriteString(fmt.Sprintf("  %s %s:\n", translations.Get("date_time"), timeSlot))
			for _, participant := range participants {
				message.WriteString(fmt.Sprintf("    + %s %s\n", translations.Get("participant"), participant))
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

	if userState.CurrentTraining.Place == translations.Get("btn_place_1") {
		if dayOfWeek == time.Tuesday || dayOfWeek == time.Wednesday || dayOfWeek == time.Thursday {
			keyboardRows = append(keyboardRows,
				tgbotapi.NewInlineKeyboardRow(
					tgbotapi.NewInlineKeyboardButtonData(
						fmt.Sprintf("%s: 20:40-22:10", translations.Get("btn_level_1")),
						fmt.Sprintf("%s|20:40-22:10", translations.Get("btn_level_1")),
					),
				),
			)
		}
		if dayOfWeek >= time.Monday && dayOfWeek <= time.Friday {
			keyboardRows = append(keyboardRows,
				tgbotapi.NewInlineKeyboardRow(
					tgbotapi.NewInlineKeyboardButtonData(
						fmt.Sprintf("%s: 19:00-20:30", translations.Get("btn_level_2")),
						fmt.Sprintf("%s|19:00-20:30", translations.Get("btn_level_2")),
					),
				),
			)
		}
	} else if userState.CurrentTraining.Place == translations.Get("btn_place_2") {
		if dayOfWeek == time.Tuesday || dayOfWeek == time.Thursday {
			keyboardRows = append(keyboardRows,
				tgbotapi.NewInlineKeyboardRow(
					tgbotapi.NewInlineKeyboardButtonData(
						fmt.Sprintf("%s: 19:00-20:30", translations.Get("btn_level_2")),
						fmt.Sprintf("%s|19:00-20:30", translations.Get("btn_level_2")),
					),
				),
			)
		}
	}

	if len(keyboardRows) == 0 {
		sendMessage(chatID, "no_time_slots")
		return
	}

	keyboard := tgbotapi.NewInlineKeyboardMarkup(keyboardRows...)

	msg := tgbotapi.NewMessage(chatID, fmt.Sprintf(translations.Get("select_time_slot"),
		currentDate.Format("02.01.2006"), userState.CurrentTraining.Place))
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

	chatMember, err := bot.GetChatMember(tgbotapi.GetChatMemberConfig{ChatConfigWithUser: tgbotapi.ChatConfigWithUser{ChatID: chatID, UserID: userState.UserID}})
	if err != nil {
		sendMessage(chatID, "error_getting_user")
		return
	}

	userState.CurrentTraining.Participant = getDisplayName(*chatMember.User)

	trainings[userState.UserID] = *userState.CurrentTraining

	message := fmt.Sprintf("%s\n\n%s %s\n%s %s\n%s %s\n%s %s",
		translations.Get("training_confirmed"),
		translations.Get("place"), userState.CurrentTraining.Place,
		translations.Get("level"), userState.CurrentTraining.Level,
		translations.Get("date_time"), userState.CurrentTraining.Date.Format("02.01.2006 15:04"),
		translations.Get("participant"), userState.CurrentTraining.Participant)

	sendMessage(chatID, message)

	delete(userStates, chatID)
}

func sendMessage(chatID int64, key string, args ...interface{}) {
	text := translations.Get(key)
	if len(args) > 0 {
		text = fmt.Sprintf(text, args...)
	}
	msg := tgbotapi.NewMessage(chatID, text)
	if _, err := bot.Send(msg); err != nil {
		log.Printf("Error sending message: %v", err)
	}
}

func editMessage(chatID int64, messageID int, key string, args ...interface{}) {
	text := translations.Get(key)
	if len(args) > 0 {
		text = fmt.Sprintf(text, args...)
	}
	editMsg := tgbotapi.NewEditMessageText(chatID, messageID, text)
	editMsg.ReplyMarkup = &tgbotapi.InlineKeyboardMarkup{
		InlineKeyboard: [][]tgbotapi.InlineKeyboardButton{},
	}
	if _, err := bot.Send(editMsg); err != nil {
		log.Printf("Error editing message: %v", err)
	}
}

var userNameCache = make(map[int64]string)

func getDisplayName(user tgbotapi.User) string {
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
