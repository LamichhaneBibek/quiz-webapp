package service

import (
	"github.com/LamichhaneBibek/quiz-webapp/internal/collection"
	"github.com/LamichhaneBibek/quiz-webapp/internal/entity"
)

type QuizService struct {
	quizCollection *collection.QuizCollection
}

func NewQuizService(quizCollection *collection.QuizCollection) *QuizService {
	return &QuizService{
		quizCollection: quizCollection,
	}
}

func (qs *QuizService) GetAllQuizzes() ([]entity.Quiz, error) {
	return qs.quizCollection.GetAllQuizzes()
}