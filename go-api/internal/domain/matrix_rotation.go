package domain

const (
	Rotation90Degrees  = 90
	Rotation180Degrees = 180
	Rotation270Degrees = 270
)

func RotateMatrix(matrix Matrix, degrees int) (Matrix, error) {
	switch degrees {
	case Rotation90Degrees:
		return rotate90Clockwise(matrix), nil
	case Rotation180Degrees:
		return rotate180(matrix), nil
	case Rotation270Degrees:
		return rotate90CounterClockwise(matrix), nil
	default:
		return Matrix{}, ErrInvalidRotationAngle
	}
}

func rotate90Clockwise(matrix Matrix) Matrix {
	rows, cols := matrix.Rows, matrix.Cols
	rotatedData := make([][]float64, cols)

	for columnIndex := 0; columnIndex < cols; columnIndex++ {
		rotatedData[columnIndex] = make([]float64, rows)

		for rowIndex := 0; rowIndex < rows; rowIndex++ {
			rotatedData[columnIndex][rows-1-rowIndex] = matrix.Data[rowIndex][columnIndex]
		}
	}

	return Matrix{
		Rows: cols,
		Cols: rows,
		Data: rotatedData,
	}
}

func rotate90CounterClockwise(matrix Matrix) Matrix {
	return rotate90Clockwise(rotate90Clockwise(rotate90Clockwise(matrix)))
}

func rotate180(matrix Matrix) Matrix {
	rows, cols := matrix.Rows, matrix.Cols
	rotatedData := make([][]float64, rows)

	for rowIndex := 0; rowIndex < rows; rowIndex++ {
		rotatedData[rowIndex] = make([]float64, cols)

		for columnIndex := 0; columnIndex < cols; columnIndex++ {
			rotatedData[rowIndex][columnIndex] = matrix.Data[rows-1-rowIndex][cols-1-columnIndex]
		}
	}

	return Matrix{
		Rows: rows,
		Cols: cols,
		Data: rotatedData,
	}
}
