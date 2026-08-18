package statistics

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"time"

	"matrix-api-challenge/go-api/internal/domain"
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
	decomposition domain.QRDecomposition,
) (domain.Statistics, error) {
	requestBody, err := json.Marshal(statisticsRequestFromDecomposition(decomposition))
	if err != nil {
		return domain.Statistics{}, domain.ErrStatisticsUnavailable
	}

	responseBody, err := gateway.post(requestBody)
	if err != nil {
		return domain.Statistics{}, domain.ErrStatisticsUnavailable
	}

	var statisticsPayload statisticsResponse
	if err := json.Unmarshal(responseBody, &statisticsPayload); err != nil {
		return domain.Statistics{}, domain.ErrStatisticsUnavailable
	}

	return statisticsPayload.toDomainStatistics(), nil
}

func (gateway *NodeStatisticsGateway) post(requestBody []byte) ([]byte, error) {
	request, err := http.NewRequestWithContext(
		context.Background(),
		http.MethodPost,
		gateway.endpointURL,
		bytes.NewReader(requestBody),
	)
	if err != nil {
		return nil, err
	}

	request.Header.Set("Content-Type", "application/json")

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
