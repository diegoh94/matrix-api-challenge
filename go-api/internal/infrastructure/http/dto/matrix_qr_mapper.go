package dto

import (
	"math"

	"matrix-api-challenge/go-api/internal/domain"
)

const jsonDecimalPlaces = 6

func MatrixQRResponseFromDomain(result domain.FactorizeMatrixResult) MatrixQRResponse {
	return MatrixQRResponse{
		Input: MatrixDimensionsResponse{
			Rows: result.Input.Rows,
			Cols: result.Input.Cols,
		},
		QR: MatrixValuesResponse{
			Q: roundMatrixValues(result.QR.Q.Data),
			R: roundMatrixValues(result.QR.R.Data),
		},
		Statistics: StatisticsResponse{
			Max:               roundToDecimalPlaces(result.Statistics.Max),
			Min:               roundToDecimalPlaces(result.Statistics.Min),
			Average:           roundToDecimalPlaces(result.Statistics.Average),
			Sum:               roundToDecimalPlaces(result.Statistics.Sum),
			HasDiagonalMatrix: result.Statistics.HasDiagonalMatrix,
			DiagonalMatrices:  result.Statistics.DiagonalMatrices,
		},
	}
}

func roundMatrixValues(matrix [][]float64) [][]float64 {
	roundedMatrix := make([][]float64, len(matrix))

	for rowIndex, row := range matrix {
		roundedMatrix[rowIndex] = make([]float64, len(row))

		for columnIndex, value := range row {
			roundedMatrix[rowIndex][columnIndex] = roundToDecimalPlaces(value)
		}
	}

	return roundedMatrix
}

func roundToDecimalPlaces(value float64) float64 {
	factor := math.Pow(10, jsonDecimalPlaces)
	return math.Round(value*factor) / factor
}
