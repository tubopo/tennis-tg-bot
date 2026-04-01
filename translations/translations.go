package translations

var Messages = map[string]string{
	"welcome": `🏓 Привет, любитель маленького мячика!

/new_training - Го постучать?
/view_trainings - Глянь, кто записался!

Ну что, готов размять ракетку? 😄`,

	"select_place":   "🏢 Пожалуйста, выберите место для тренировки:",
	"selected_opt_1": "Вы выбрали %s",
	"selected_opt_2": "Вы выбрали %s (%s)",

	"select_time_slot": "🕒 Выберите время для %s в %s:",
	"no_time_slots":    "😕 Не сегодня. Пожалуйста, попробуйте другой день.",
	"training_confirmed": `🎉 Отлично! Ваша тренировка подтверждена! 🏓

Вот детали вашей тренировки:`,
	"place":                 "🏢 Место:",
	"level":                 "🏅 Уровень:",
	"date_time":             "🕒 Дата и время:",
	"participant":           "👤 Участник:",
	"no_trainings":          "📅 Пока нет запланированных тренировок. Самое время забронировать одну! /new_training",
	"btn_place_1":           "Смолячкова, 9",
	"btn_place_2":           "Ленина, 27",
	"btn_level_1":           "Уровень 1",
	"btn_level_2":           "Уровень 2",
	"unknown_command":       "🤔 Извините, я не понимаю. Используйте /new_training для новой тренировки или /view_trainings для просмотра всех.",
	"error_getting_user":    "❌ Ошибка при получении информации о пользователе. Пожалуйста, попробуйте еще раз позже.",
	"view_trainings_header": "📋 Запланированные тренировки:",
	"no_current_training":   "❓ Нет текущей тренировки в процессе бронирования. Используйте /new_training, чтобы начать.",
}

func Get(key string) string {
	if msg, ok := Messages[key]; ok {
		return msg
	}
	return key
}
