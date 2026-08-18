package domain

import "math"

func NewMatrix(data [][]float64) (Matrix, error) {
	if err := validateMatrixData(data); err != nil {
		return Matrix{}, err
	}

	rowCount := len(data)
	columnCount := len(data[0])

	return Matrix{
		Rows: rowCount,
		Cols: columnCount,
		Data: data,
	}, nil
}

func validateMatrixData(data [][]float64) error {
	if len(data) == 0 {
		return ErrEmptyMatrix
	}

	columnCount := len(data[0])
	if columnCount == 0 {
		return ErrEmptyMatrix
	}

	for _, row := range data {
		if len(row) != columnCount {
			return ErrInvalidMatrixShape
		}

		for _, value := range row {
			if !isFiniteNumber(value) {
				return ErrInvalidMatrixValue
			}
		}
	}

	return nil
}

func isFiniteNumber(value float64) bool {
	return !math.IsNaN(value) && !math.IsInf(value, 0)
}
