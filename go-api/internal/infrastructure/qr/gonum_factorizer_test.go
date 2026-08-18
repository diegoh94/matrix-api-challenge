package qr_test

import (
	"math"
	"testing"

	"matrix-api-challenge/go-api/internal/domain"
	"matrix-api-challenge/go-api/internal/infrastructure/qr"
)

func TestGonumQRFactorizerFactorize(t *testing.T) {
	inputMatrix, err := domain.NewMatrix([][]float64{
		{1, 2},
		{3, 4},
		{5, 6},
	})
	if err != nil {
		t.Fatalf("unexpected validation error: %v", err)
	}

	factorizer := qr.NewGonumQRFactorizer()
	decomposition, err := factorizer.Factorize(inputMatrix)
	if err != nil {
		t.Fatalf("unexpected factorization error: %v", err)
	}

	if decomposition.Q.Rows != 3 || decomposition.Q.Cols != 3 {
		t.Fatalf("expected Q to be 3x3, got %dx%d", decomposition.Q.Rows, decomposition.Q.Cols)
	}

	if decomposition.R.Rows != 3 || decomposition.R.Cols != 2 {
		t.Fatalf("expected R to be 3x2, got %dx%d", decomposition.R.Rows, decomposition.R.Cols)
	}

	assertMatrixProductMatchesOriginal(t, inputMatrix, decomposition)
}

func assertMatrixProductMatchesOriginal(
	t *testing.T,
	original domain.Matrix,
	decomposition domain.QRDecomposition,
) {
	t.Helper()

	const tolerance = 1e-9

	for rowIndex := 0; rowIndex < original.Rows; rowIndex++ {
		for columnIndex := 0; columnIndex < original.Cols; columnIndex++ {
			reconstructedValue := 0.0

			for innerIndex := 0; innerIndex < decomposition.Q.Cols; innerIndex++ {
				reconstructedValue += decomposition.Q.Data[rowIndex][innerIndex] * decomposition.R.Data[innerIndex][columnIndex]
			}

			originalValue := original.Data[rowIndex][columnIndex]
			if math.Abs(reconstructedValue-originalValue) > tolerance {
				t.Fatalf(
					"Q*R mismatch at (%d,%d): got %.8f, want %.8f",
					rowIndex,
					columnIndex,
					reconstructedValue,
					originalValue,
				)
			}
		}
	}
}
