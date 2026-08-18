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
		Statistics: mapStatisticsResponse(result.Statistics),
	}
}

func mapStatisticsResponse(statistics domain.Statistics) StatisticsResponse {
	return StatisticsResponse{
		Max:               roundToDecimalPlaces(statistics.Max),
		Min:               roundToDecimalPlaces(statistics.Min),
		Average:           roundToDecimalPlaces(statistics.Average),
		Sum:               roundToDecimalPlaces(statistics.Sum),
		HasDiagonalMatrix: statistics.HasDiagonalMatrix,
		DiagonalMatrices:  statistics.DiagonalMatrices,
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
	factor := math.Pow10(jsonDecimalPlaces)
	return math.Round(value*factor) / factor
}
