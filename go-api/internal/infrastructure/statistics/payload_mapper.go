package statistics

import "matrix-api-challenge/go-api/internal/domain"

func statisticsRequestFromMatrices(matrices []domain.NamedMatrix) statisticsRequest {
	payloadMatrices := make([]namedMatrixPayload, len(matrices))

	for index, matrix := range matrices {
		payloadMatrices[index] = namedMatrixPayload{
			Name: matrix.Name,
			Data: matrix.Matrix.Data,
		}
	}

	return statisticsRequest{
		Matrices: payloadMatrices,
	}
}

func statisticsRequestFromDecomposition(decomposition domain.QRDecomposition) statisticsRequest {
	return statisticsRequestFromMatrices([]domain.NamedMatrix{
		{Name: "Q", Matrix: decomposition.Q},
		{Name: "R", Matrix: decomposition.R},
	})
}

func (response statisticsResponse) toDomainStatistics() domain.Statistics {
	return domain.Statistics{
		Max:               response.Max,
		Min:               response.Min,
		Average:           response.Average,
		Sum:               response.Sum,
		HasDiagonalMatrix: response.HasDiagonalMatrix,
		DiagonalMatrices:  response.DiagonalMatrices,
	}
}
