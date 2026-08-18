package domain

import "math"

func NewMatrix(data [][]float64) (Matrix, error) {
	if err := validateRectangularMatrix(data); err != nil {
		return Matrix{}, err
	}

	return Matrix{
		Rows: len(data),
		Cols: len(data[0]),
		Data: data,
	}, nil
}

func validateRectangularMatrix(data [][]float64) error {
	if len(data) == 0 {
		return ErrEmptyMatrix
	}

	expectedColumnCount := len(data[0])
	if expectedColumnCount == 0 {
		return ErrEmptyMatrix
	}

	for _, row := range data {
		if len(row) != expectedColumnCount {
			return ErrInvalidMatrixShape
		}

		for _, value := range row {
			if math.IsNaN(value) || math.IsInf(value, 0) {
				return ErrInvalidMatrixValue
			}
		}
	}

	return nil
}
