package main

import (
	"log"

	"github.com/gofiber/fiber/v2"

	"matrix-api-challenge/go-api/internal/application"
	"matrix-api-challenge/go-api/internal/config"
	"matrix-api-challenge/go-api/internal/infrastructure/auth"
	httpadapter "matrix-api-challenge/go-api/internal/infrastructure/http"
	"matrix-api-challenge/go-api/internal/infrastructure/qr"
	"matrix-api-challenge/go-api/internal/infrastructure/statistics"
)

func main() {
	appConfig, err := config.Load()
	if err != nil {
		log.Fatalf("load config: %v", err)
	}

	app := buildApplication(appConfig)

	log.Printf("go-api listening on port %s", appConfig.Port)

	if err := app.Listen(":" + appConfig.Port); err != nil {
		log.Fatalf("start server: %v", err)
	}
}

func buildApplication(appConfig config.Config) *fiber.App {
	qrFactorizer := qr.NewGonumQRFactorizer()
	statisticsGateway := statistics.NewNodeStatisticsGateway(statistics.GatewayConfig{
		BaseURL: appConfig.NodeAPIURL,
		Timeout: appConfig.NodeAPITimeout,
	})

	factorizeMatrixUseCase := application.NewFactorizeMatrixUseCase(qrFactorizer, statisticsGateway)
	matrixHandler := httpadapter.NewMatrixHandler(factorizeMatrixUseCase)

	serverConfig := httpadapter.ServerConfig{
		AuthEnabled: appConfig.AuthEnabled(),
		APIKey:      appConfig.APIKey,
	}

	if serverConfig.AuthEnabled {
		serverConfig.TokenService = auth.NewTokenService(appConfig.JWTSecret, appConfig.JWTExpiration)
	}

	authHandler := (*httpadapter.AuthHandler)(nil)
	if serverConfig.AuthEnabled {
		authHandler = httpadapter.NewAuthHandler(serverConfig.TokenService, appConfig.APIKey)
	}

	return httpadapter.NewServer(serverConfig, matrixHandler, authHandler)
}
