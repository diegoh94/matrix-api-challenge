package application_test

import (
	"context"
	"errors"
	"testing"

	"matrix-api-challenge/go-api/internal/application"
	"matrix-api-challenge/go-api/internal/domain"
)

type stubQRFactorizer struct {
	decomposition domain.QRDecomposition
	err           error
	called        bool
}

func (stub *stubQRFactorizer) Factorize(matrix domain.Matrix) (domain.QRDecomposition, error) {
	stub.called = true
	return stub.decomposition, stub.err
}

type stubStatisticsGateway struct {
	statistics domain.Statistics
	err        error
	lastCtx    context.Context
}

func (stub *stubStatisticsGateway) ComputeStatistics(
	ctx context.Context,
	_ []domain.NamedMatrix,
) (domain.Statistics, error) {
	stub.lastCtx = ctx
	return stub.statistics, stub.err
}

func TestFactorizeMatrixUseCaseReturnsCombinedResult(t *testing.T) {
	inputMatrix := domain.Matrix{
		Rows: 2,
		Cols: 2,
		Data: [][]float64{{1, 2}, {3, 4}},
	}

	expectedQR := domain.QRDecomposition{
		Q: domain.Matrix{Rows: 2, Cols: 2, Data: [][]float64{{1, 0}, {0, 1}}},
		R: domain.Matrix{Rows: 2, Cols: 2, Data: [][]float64{{1, 2}, {0, 4}}},
	}

	expectedStatistics := domain.Statistics{
		Max:               4,
		Min:               0,
		Average:           1.375,
		Sum:               11,
		HasDiagonalMatrix: true,
		DiagonalMatrices:  []string{"Q"},
	}

	qrFactorizer := &stubQRFactorizer{decomposition: expectedQR}
	statisticsGateway := &stubStatisticsGateway{statistics: expectedStatistics}

	useCase := application.NewFactorizeMatrixUseCase(qrFactorizer, statisticsGateway)

	result, err := useCase.Execute(context.Background(), inputMatrix)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !qrFactorizer.called {
		t.Fatal("expected QR factorizer to be called")
	}

	if result.Input.Rows != inputMatrix.Rows || result.QR.Q.Rows != 2 {
		t.Fatalf("unexpected result shape: %+v", result)
	}

	if result.Statistics.Sum != expectedStatistics.Sum {
		t.Fatalf("expected statistics to be merged into result, got %+v", result.Statistics)
	}
}

func TestFactorizeMatrixUseCaseStopsWhenStatisticsGatewayFails(t *testing.T) {
	qrFactorizer := &stubQRFactorizer{
		decomposition: domain.QRDecomposition{
			Q: domain.Matrix{Rows: 1, Cols: 1, Data: [][]float64{{1}}},
			R: domain.Matrix{Rows: 1, Cols: 1, Data: [][]float64{{1}}},
		},
	}

	statisticsGateway := &stubStatisticsGateway{err: domain.ErrStatisticsUnavailable}

	useCase := application.NewFactorizeMatrixUseCase(qrFactorizer, statisticsGateway)

	_, err := useCase.Execute(context.Background(), domain.Matrix{
		Rows: 1,
		Cols: 1,
		Data: [][]float64{{1}},
	})

	if !errors.Is(err, domain.ErrStatisticsUnavailable) {
		t.Fatalf("expected ErrStatisticsUnavailable, got %v", err)
	}
}

func TestFactorizeMatrixUseCaseDoesNotCallStatisticsWhenQRFails(t *testing.T) {
	qrFactorizer := &stubQRFactorizer{err: domain.ErrMatrixNotFactorizable}
	statisticsGateway := &stubStatisticsGateway{}

	useCase := application.NewFactorizeMatrixUseCase(qrFactorizer, statisticsGateway)

	_, err := useCase.Execute(context.Background(), domain.Matrix{
		Rows: 1,
		Cols: 1,
		Data: [][]float64{{1}},
	})

	if !errors.Is(err, domain.ErrMatrixNotFactorizable) {
		t.Fatalf("expected ErrMatrixNotFactorizable, got %v", err)
	}

	if statisticsGateway.lastCtx != nil {
		t.Fatal("statistics gateway must not be called when QR fails")
	}
}
