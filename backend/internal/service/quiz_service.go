package service

import (
	"errors"

	"github.com/LamichhaneBibek/quiz-webapp/internal/collection"
	"github.com/LamichhaneBibek/quiz-webapp/internal/entity"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

func (qs QuizService) GetQuizById(id primitive.ObjectID) (*entity.Quiz, error) {
	return qs.quizCollection.GetQuizByID(id)
}

func (qs QuizService) UpdateQuiz(id primitive.ObjectID, name string, questions []entity.QuizQuestion) error {
	quiz, err := qs.quizCollection.GetQuizByID(id)
	if err != nil {
		return err
	}

	if quiz == nil {
		return errors.New("quiz not found")
	}

	quiz.Name = name
	quiz.Questions = questions
	return qs.quizCollection.UpdateQuiz(*quiz)
}
