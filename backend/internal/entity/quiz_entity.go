package entity

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Quiz struct {
	ID        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Name      string             `json:"name"`
	Questions []QuizQuestion     `json:"questions"`
}

type QuizQuestion struct {
	ID      string       `json:"id"`
	Name    string       `json:"name"`
	Options []QuizOption `json:"options"`
}

type QuizOption struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	IsCorrect bool   `json:"isCorrect"`
}
