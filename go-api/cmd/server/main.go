package main

import (
	"log"

	"matrix-api-challenge/go-api/internal/application"
	"matrix-api-challenge/go-api/internal/config"
	httpadapter "matrix-api-challenge/go-api/internal/infrastructure/http"
	"matrix-api-challenge/go-api/internal/infrastructure/qr"
	"matrix-api-challenge/go-api/internal/infrastructure/statistics"
)

func main() {
	appConfig, err := config.Load()
	if err != nil {
		log.Fatalf("load config: %v", err)
	}

	qrFactorizer := qr.NewGonumQRFactorizer()
	statisticsGateway := statistics.NewNodeStatisticsGateway(statistics.GatewayConfig{
		BaseURL: appConfig.NodeAPIURL,
		Timeout: appConfig.NodeAPITimeout,
	})

	factorizeMatrixUseCase := application.NewFactorizeMatrixUseCase(qrFactorizer, statisticsGateway)
	matrixHandler := httpadapter.NewMatrixHandler(factorizeMatrixUseCase)
	app := httpadapter.NewServer(matrixHandler)

	log.Printf("go-api listening on port %s", appConfig.Port)

	if err := app.Listen(":" + appConfig.Port); err != nil {
		log.Fatalf("start server: %v", err)
	}
}
