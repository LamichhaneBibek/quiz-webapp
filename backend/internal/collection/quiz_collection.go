package collection

import (
	"context"

	"github.com/LamichhaneBibek/quiz-webapp/internal/entity"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type QuizCollection struct {
	collection *mongo.Collection
}

func NewQuizCollection(collection *mongo.Collection) *QuizCollection {
	return &QuizCollection{
		collection: collection,
	}
}

func (qc *QuizCollection) InsertQuiz(quiz entity.Quiz) error {
	_, err := qc.collection.InsertOne(context.Background(), quiz)
	return err
}

func (qc *QuizCollection) GetQuizByID(id primitive.ObjectID) (*entity.Quiz, error) {
	var quiz entity.Quiz
	err := qc.collection.FindOne(context.Background(), bson.M{"_id": id}).Decode(&quiz)
	if err != nil {
		return nil, err
	}
	return &quiz, nil
}

func (qc *QuizCollection) GetAllQuizzes() ([]entity.Quiz, error) {
	cursor, err := qc.collection.Find(context.Background(), bson.M{})
	if err != nil {
		return nil, err
	}
	var quizzes []entity.Quiz
	err = cursor.All(context.Background(), &quizzes)
	if err != nil {
		return nil, err
	}
	return quizzes, nil
}