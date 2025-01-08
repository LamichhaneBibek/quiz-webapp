package internal

import (
	"context"
	"log"
	"time"

	"github.com/LamichhaneBibek/quiz-webapp/internal/collection"
	"github.com/LamichhaneBibek/quiz-webapp/internal/controller"
	"github.com/LamichhaneBibek/quiz-webapp/internal/service"
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type App struct {
	httpServer  *fiber.App
	database    *mongo.Database
	quizService *service.QuizService
	netService  *service.NetService
}

func (a *App) Init() {
	a.setupDb()
	a.setupServices()
	a.setupHttp()

	log.Fatal(a.httpServer.Listen(":8000"))
}

func (a *App) setupHttp() {
	app := fiber.New()
	app.Use(cors.New())

	quizController := controller.NewQuizController(a.quizService)
	app.Get("/api/quizzes", quizController.GetQuizzes)

	wsController := controller.NewWebsocketController(a.netService)
	app.Get("/ws", websocket.New(wsController.HandleWS))

	a.httpServer = app
}

func (a *App) setupServices() {
	a.quizService = service.NewQuizService(collection.NewQuizCollection(a.database.Collection("quizzes")))
	a.netService = service.NewNetService(a.quizService)
}

func (a *App) setupDb() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatal(err)
	}
	a.database = client.Database("quiz")
}
