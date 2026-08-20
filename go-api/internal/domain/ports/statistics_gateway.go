package ports

import (
	"context"

	"matrix-api-challenge/go-api/internal/domain"
)

type StatisticsGateway interface {
	ComputeStatistics(ctx context.Context, matrices []domain.NamedMatrix) (domain.Statistics, error)
}
