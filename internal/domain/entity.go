package domain

import "time"

type Error struct {
	Message string    `json:"message"`
	Time    time.Time `json:"time,omitempty"`
}

func NewError(message string, time time.Time) Error {
	return Error{message, time}
}
