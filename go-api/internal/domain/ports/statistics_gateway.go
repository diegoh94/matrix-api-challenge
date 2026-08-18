package ports

import (
	"context"

	"matrix-api-challenge/go-api/internal/domain"
)

type StatisticsGateway interface {
	ComputeStatistics(ctx context.Context, decomposition domain.QRDecomposition) (domain.Statistics, error)
}
