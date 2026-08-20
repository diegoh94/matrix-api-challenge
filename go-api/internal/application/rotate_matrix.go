package application

import (
	"context"

	"matrix-api-challenge/go-api/internal/domain"
	"matrix-api-challenge/go-api/internal/domain/ports"
)

type RotateMatrixUseCase struct {
	statisticsGateway ports.StatisticsGateway
}

func NewRotateMatrixUseCase(statisticsGateway ports.StatisticsGateway) *RotateMatrixUseCase {
	return &RotateMatrixUseCase{
		statisticsGateway: statisticsGateway,
	}
}

func (useCase *RotateMatrixUseCase) Execute(
	ctx context.Context,
	matrix domain.Matrix,
	degrees int,
) (domain.RotateMatrixResult, error) {
	rotated, err := domain.RotateMatrix(matrix, degrees)
	if err != nil {
		return domain.RotateMatrixResult{}, err
	}

	statistics, err := useCase.statisticsGateway.ComputeStatistics(ctx, []domain.NamedMatrix{
		{Name: "rotated", Matrix: rotated},
	})
	if err != nil {
		return domain.RotateMatrixResult{}, err
	}

	return domain.RotateMatrixResult{
		Input:      matrix,
		Degrees:    degrees,
		Rotated:    rotated,
		Statistics: statistics,
	}, nil
}
