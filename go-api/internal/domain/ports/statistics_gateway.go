package ports

import "matrix-api-challenge/go-api/internal/domain"

type StatisticsGateway interface {
	ComputeStatistics(decomposition domain.QRDecomposition) (domain.Statistics, error)
}
