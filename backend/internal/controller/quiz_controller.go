package controller

import (
	"github.com/LamichhaneBibek/quiz-webapp/internal/service"
	"github.com/gofiber/fiber/v2"
)


type QuizController struct {
	quizService *service.QuizService
}

func NewQuizController(quizService *service.QuizService) *QuizController {
	return &QuizController{
		quizService: quizService,
	}
}

func(qc *QuizController) GetQuizzes(c *fiber.Ctx) error {
	quizzes, err := qc.quizService.GetAllQuizzes()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}
	return c.JSON(fiber.Map{
		"quizzes": quizzes,
	})
}