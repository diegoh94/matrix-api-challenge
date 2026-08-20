package statistics

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"matrix-api-challenge/go-api/internal/domain"
	"matrix-api-challenge/go-api/internal/infrastructure/auth"
)

const statisticsPath = "/api/v1/statistics"

type GatewayConfig struct {
	BaseURL string
	Timeout time.Duration
}

type NodeStatisticsGateway struct {
	endpointURL string
	httpClient  *http.Client
}

func NewNodeStatisticsGateway(config GatewayConfig) *NodeStatisticsGateway {
	baseURL := strings.TrimRight(config.BaseURL, "/")

	return &NodeStatisticsGateway{
		endpointURL: baseURL + statisticsPath,
		httpClient: &http.Client{
			Timeout: config.Timeout,
		},
	}
}

func (gateway *NodeStatisticsGateway) ComputeStatistics(
	ctx context.Context,
	matrices []domain.NamedMatrix,
) (domain.Statistics, error) {
	requestBody, err := json.Marshal(statisticsRequestFromMatrices(matrices))
	if err != nil {
		return domain.Statistics{}, domain.ErrStatisticsUnavailable
	}

	responseBody, err := gateway.post(ctx, requestBody)
	if err != nil {
		return domain.Statistics{}, domain.ErrStatisticsUnavailable
	}

	var statisticsPayload statisticsResponse
	if err := json.Unmarshal(responseBody, &statisticsPayload); err != nil {
		return domain.Statistics{}, domain.ErrStatisticsUnavailable
	}

	return statisticsPayload.toDomainStatistics(), nil
}

func (gateway *NodeStatisticsGateway) post(ctx context.Context, requestBody []byte) ([]byte, error) {
	request, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		gateway.endpointURL,
		bytes.NewReader(requestBody),
	)
	if err != nil {
		return nil, err
	}

	request.Header.Set("Content-Type", "application/json")

	if accessToken, ok := auth.TokenFromContext(ctx); ok {
		request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))
	}

	response, err := gateway.httpClient.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, domain.ErrStatisticsUnavailable
	}

	return io.ReadAll(response.Body)
}
