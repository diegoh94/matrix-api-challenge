package statistics

import "matrix-api-challenge/go-api/internal/domain"

func statisticsRequestFromDecomposition(decomposition domain.QRDecomposition) statisticsRequest {
	return statisticsRequest{
		Matrices: []namedMatrixPayload{
			{Name: "Q", Data: decomposition.Q.Data},
			{Name: "R", Data: decomposition.R.Data},
		},
	}
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
