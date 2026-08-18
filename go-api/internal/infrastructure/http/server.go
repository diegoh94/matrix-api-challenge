package http

import (
	"log"

	"github.com/gofiber/fiber/v2"
)

func NewServer(matrixHandler *MatrixHandler) *fiber.App {
	app := fiber.New(fiber.Config{
		AppName: "go-api",
	})

	app.Use(requestLogger)

	app.Get("/health", matrixHandler.HealthCheck)
	app.Post("/api/v1/matrix/qr", matrixHandler.FactorizeMatrix)

	return app
}

func requestLogger(context *fiber.Ctx) error {
	log.Printf("%s %s", context.Method(), context.Path())
	return context.Next()
}
