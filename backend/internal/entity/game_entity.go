package entity

import "github.com/google/uuid"

type Game struct {
	Id              uuid.UUID `json:"id"`
	Quiz            Quiz      `json:"quiz"`
	CurrentQuestion int       `json:"current_question"`
	Code            string    `json:"code"`
}
