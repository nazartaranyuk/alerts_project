package domain

import "time"

type TelegramBotMessage struct {
	ChatTitle string    `json:"chat_title"`
	Message   string    `json:"message"`
	DateIso   time.Time `json:"date_iso"`
}
