package domain_test

import (
	"testing"

	"matrix-api-challenge/go-api/internal/domain"
)

func TestRotateMatrix90DegreesClockwise(t *testing.T) {
	matrix := domain.Matrix{
		Rows: 3,
		Cols: 2,
		Data: [][]float64{
			{1, 2},
			{3, 4},
			{5, 6},
		},
	}

	rotated, err := domain.RotateMatrix(matrix, domain.Rotation90Degrees)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if rotated.Rows != 2 || rotated.Cols != 3 {
		t.Fatalf("unexpected rotated shape: %dx%d", rotated.Rows, rotated.Cols)
	}

	expected := [][]float64{
		{5, 3, 1},
		{6, 4, 2},
	}

	for rowIndex, row := range expected {
		for columnIndex, value := range row {
			if rotated.Data[rowIndex][columnIndex] != value {
				t.Fatalf("unexpected value at [%d][%d]: got %v want %v", rowIndex, columnIndex, rotated.Data, expected)
			}
		}
	}
}

func TestRotateMatrix180Degrees(t *testing.T) {
	matrix := domain.Matrix{
		Rows: 2,
		Cols: 2,
		Data: [][]float64{
			{1, 2},
			{3, 4},
		},
	}

	rotated, err := domain.RotateMatrix(matrix, domain.Rotation180Degrees)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := [][]float64{
		{4, 3},
		{2, 1},
	}

	for rowIndex, row := range expected {
		for columnIndex, value := range row {
			if rotated.Data[rowIndex][columnIndex] != value {
				t.Fatalf("unexpected rotated matrix: %+v", rotated.Data)
			}
		}
	}
}

func TestRotateMatrixRejectsUnsupportedAngle(t *testing.T) {
	matrix := domain.Matrix{
		Rows: 1,
		Cols: 1,
		Data: [][]float64{{1}},
	}

	_, err := domain.RotateMatrix(matrix, 45)
	if err != domain.ErrInvalidRotationAngle {
		t.Fatalf("expected ErrInvalidRotationAngle, got %v", err)
	}
}
