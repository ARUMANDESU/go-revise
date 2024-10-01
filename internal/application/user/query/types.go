package query

type User struct {
	ID       string   `json:"id"`
	ChatID   string   `json:"chat_id"`
	Settings Settings `json:"settings"`
}

type Settings struct {
	Language     string `json:"language"`
	ReminderTime struct {
		Hour   uint `json:"hour"`
		Minute uint `json:"minute"`
	} `json:"reminder_time"`
}
