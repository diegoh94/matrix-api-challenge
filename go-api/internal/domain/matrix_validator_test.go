package domain_test

import (
	"testing"

	"matrix-api-challenge/go-api/internal/domain"
)

func TestNewMatrixClonesInputData(t *testing.T) {
	originalData := [][]float64{
		{1, 2},
		{3, 4},
	}

	matrix, err := domain.NewMatrix(originalData)
	if err != nil {
		t.Fatalf("unexpected validation error: %v", err)
	}

	originalData[0][0] = 99

	if matrix.Data[0][0] == 99 {
		t.Fatal("matrix data must be isolated from external mutations")
	}
}
