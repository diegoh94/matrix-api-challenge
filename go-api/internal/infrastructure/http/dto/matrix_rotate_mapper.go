package dto

import "matrix-api-challenge/go-api/internal/domain"

func MatrixRotateResponseFromDomain(result domain.RotateMatrixResult) MatrixRotateResponse {
	return MatrixRotateResponse{
		Input: MatrixDimensionsResponse{
			Rows: result.Input.Rows,
			Cols: result.Input.Cols,
		},
		Degrees: result.Degrees,
		Rotated: RotatedMatrixResponse{
			Rows: result.Rotated.Rows,
			Cols: result.Rotated.Cols,
			Data: roundMatrixValues(result.Rotated.Data),
		},
		Statistics: mapStatisticsResponse(result.Statistics),
	}
}
