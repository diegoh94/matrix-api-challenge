package application_test

import (
	"context"
	"errors"
	"testing"

	"matrix-api-challenge/go-api/internal/application"
	"matrix-api-challenge/go-api/internal/domain"
)

func TestRotateMatrixUseCaseReturnsCombinedResult(t *testing.T) {
	inputMatrix := domain.Matrix{
		Rows: 2,
		Cols: 3,
		Data: [][]float64{{1, 2, 3}, {4, 5, 6}},
	}

	expectedStatistics := domain.Statistics{
		Max:               6,
		Min:               1,
		Average:           3.5,
		Sum:               21,
		HasDiagonalMatrix: false,
		DiagonalMatrices:  nil,
	}

	statisticsGateway := &stubStatisticsGateway{statistics: expectedStatistics}

	useCase := application.NewRotateMatrixUseCase(statisticsGateway)

	result, err := useCase.Execute(context.Background(), inputMatrix, domain.Rotation90Degrees)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.Degrees != domain.Rotation90Degrees {
		t.Fatalf("expected 90 degrees, got %d", result.Degrees)
	}

	if result.Rotated.Rows != 3 || result.Rotated.Cols != 2 {
		t.Fatalf("unexpected rotated shape: rows=%d cols=%d", result.Rotated.Rows, result.Rotated.Cols)
	}

	if result.Statistics.Sum != expectedStatistics.Sum {
		t.Fatalf("expected statistics to be merged into result, got %+v", result.Statistics)
	}
}

func TestRotateMatrixUseCaseStopsWhenRotationAngleIsInvalid(t *testing.T) {
	statisticsGateway := &stubStatisticsGateway{}

	useCase := application.NewRotateMatrixUseCase(statisticsGateway)

	_, err := useCase.Execute(context.Background(), domain.Matrix{
		Rows: 1,
		Cols: 1,
		Data: [][]float64{{1}},
	}, 45)

	if !errors.Is(err, domain.ErrInvalidRotationAngle) {
		t.Fatalf("expected ErrInvalidRotationAngle, got %v", err)
	}

	if statisticsGateway.lastCtx != nil {
		t.Fatal("statistics gateway must not be called when rotation fails")
	}
}

func TestRotateMatrixUseCaseStopsWhenStatisticsGatewayFails(t *testing.T) {
	statisticsGateway := &stubStatisticsGateway{err: domain.ErrStatisticsUnavailable}

	useCase := application.NewRotateMatrixUseCase(statisticsGateway)

	_, err := useCase.Execute(context.Background(), domain.Matrix{
		Rows: 2,
		Cols: 2,
		Data: [][]float64{{1, 2}, {3, 4}},
	}, domain.Rotation180Degrees)

	if !errors.Is(err, domain.ErrStatisticsUnavailable) {
		t.Fatalf("expected ErrStatisticsUnavailable, got %v", err)
	}
}
