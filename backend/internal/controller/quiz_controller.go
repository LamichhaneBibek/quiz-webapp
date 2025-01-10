package controller

import (
	"github.com/LamichhaneBibek/quiz-webapp/internal/entity"
	"github.com/LamichhaneBibek/quiz-webapp/internal/service"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type QuizController struct {
	quizService *service.QuizService
}

func NewQuizController(quizService *service.QuizService) *QuizController {
	return &QuizController{
		quizService: quizService,
	}
}

func (qc *QuizController) GetQuizzes(c *fiber.Ctx) error {
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

func (c QuizController) GetQuizById(ctx *fiber.Ctx) error {
	quizIdStr := ctx.Params("quizId")
	quizId, err := primitive.ObjectIDFromHex(quizIdStr)
	if err != nil {
		return ctx.SendStatus(fiber.StatusBadRequest)
	}

	quiz, err := c.quizService.GetQuizById(quizId)
	if err != nil {
		return err
	}

	if quiz == nil {
		return ctx.SendStatus(fiber.StatusNotFound)
	}

	return ctx.JSON(quiz)
}

type UpdateQuizRequest struct {
	Name      string                `json:"name"`
	Questions []entity.QuizQuestion `json:"questions"`
}

func (c QuizController) UpdateQuizById(ctx *fiber.Ctx) error {

	quizIdStr := ctx.Params("quizId")
	quizId, err := primitive.ObjectIDFromHex(quizIdStr)
	if err != nil {
		return ctx.SendStatus(fiber.StatusBadRequest)
	}

	var req UpdateQuizRequest
	if err := ctx.BodyParser(&req); err != nil {
		return err
	}

	if err := c.quizService.UpdateQuiz(quizId, req.Name, req.Questions); err != nil {
		return err
	}

	return ctx.SendStatus(fiber.StatusOK)
}
