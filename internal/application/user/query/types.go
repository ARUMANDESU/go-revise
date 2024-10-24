package query

type User struct {
	ID       string   `json:"id"`
	ChatID   int64    `json:"chat_id"`
	Settings Settings `json:"settings"`
}

type Settings struct {
	Language     string       `json:"language"`
	ReminderTime ReminderTime `json:"reminder_time"`
}

type ReminderTime struct {
	Hour   uint8 `json:"hour"`
	Minute uint8 `json:"minute"`
}
