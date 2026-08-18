package application

import (
	"context"

	"matrix-api-challenge/go-api/internal/domain"
	"matrix-api-challenge/go-api/internal/domain/ports"
)

type FactorizeMatrixUseCase struct {
	qrFactorizer      ports.QRFactorizer
	statisticsGateway ports.StatisticsGateway
}

func NewFactorizeMatrixUseCase(
	qrFactorizer ports.QRFactorizer,
	statisticsGateway ports.StatisticsGateway,
) *FactorizeMatrixUseCase {
	return &FactorizeMatrixUseCase{
		qrFactorizer:      qrFactorizer,
		statisticsGateway: statisticsGateway,
	}
}

func (useCase *FactorizeMatrixUseCase) Execute(
	ctx context.Context,
	matrix domain.Matrix,
) (domain.FactorizeMatrixResult, error) {
	qrDecomposition, err := useCase.qrFactorizer.Factorize(matrix)
	if err != nil {
		return domain.FactorizeMatrixResult{}, err
	}

	statistics, err := useCase.statisticsGateway.ComputeStatistics(ctx, qrDecomposition)
	if err != nil {
		return domain.FactorizeMatrixResult{}, err
	}

	return domain.FactorizeMatrixResult{
		Input:      matrix,
		QR:         qrDecomposition,
		Statistics: statistics,
	}, nil
}
