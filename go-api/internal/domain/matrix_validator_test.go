package domain_test

import (
	"errors"
	"math"
	"testing"

	"matrix-api-challenge/go-api/internal/domain"
)

func TestNewMatrixRejectsInvalidInput(t *testing.T) {
	testCases := []struct {
		name      string
		data      [][]float64
		wantError error
	}{
		{
			name:      "empty matrix",
			data:      [][]float64{},
			wantError: domain.ErrEmptyMatrix,
		},
		{
			name:      "ragged rows",
			data:      [][]float64{{1, 2}, {3}},
			wantError: domain.ErrInvalidMatrixShape,
		},
		{
			name:      "nan value",
			data:      [][]float64{{1, math.NaN()}},
			wantError: domain.ErrInvalidMatrixValue,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			_, err := domain.NewMatrix(testCase.data)
			if !errors.Is(err, testCase.wantError) {
				t.Fatalf("expected %v, got %v", testCase.wantError, err)
			}
		})
	}
}

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
