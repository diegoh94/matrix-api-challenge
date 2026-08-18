package http

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	swagger "github.com/swaggo/fiber-swagger"

	_ "matrix-api-challenge/go-api/docs"

	"matrix-api-challenge/go-api/internal/infrastructure/auth"
	"matrix-api-challenge/go-api/internal/infrastructure/http/middleware"
)

type ServerConfig struct {
	AuthEnabled  bool
	TokenService *auth.TokenService
	APIKey       string
}

func NewServer(
	serverConfig ServerConfig,
	matrixHandler *MatrixHandler,
	authHandler *AuthHandler,
) *fiber.App {
	app := fiber.New(fiber.Config{
		AppName: "go-api",
	})

	app.Use(cors.New())
	app.Use(logRequest)

	app.Get("/docs/*", swagger.WrapHandler)

	app.Get("/health", matrixHandler.HealthCheck)

	if serverConfig.AuthEnabled {
		app.Post("/auth/token", authHandler.IssueToken)

		apiGroup := app.Group("/api/v1", middleware.JWTAuth(serverConfig.TokenService))
		apiGroup.Post("/matrix/qr", matrixHandler.FactorizeMatrix)

		log.Println("JWT authentication enabled")
	} else {
		app.Post("/api/v1/matrix/qr", matrixHandler.FactorizeMatrix)
		log.Println("JWT authentication disabled")
	}

	return app
}

func logRequest(ctx *fiber.Ctx) error {
	log.Printf("%s %s", ctx.Method(), ctx.Path())
	return ctx.Next()
}
