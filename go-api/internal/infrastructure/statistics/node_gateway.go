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
)

type GatewayConfig struct {
	BaseURL string
	Timeout time.Duration
}

type NodeStatisticsGateway struct {
	baseURL    string
	httpClient *http.Client
}

func NewNodeStatisticsGateway(config GatewayConfig) *NodeStatisticsGateway {
	return &NodeStatisticsGateway{
		baseURL: strings.TrimRight(config.BaseURL, "/"),
		httpClient: &http.Client{
			Timeout: config.Timeout,
		},
	}
}

func (gateway *NodeStatisticsGateway) ComputeStatistics(
	decomposition domain.QRDecomposition,
) (domain.Statistics, error) {
	payload := statisticsRequest{
		Matrices: []namedMatrixPayload{
			{Name: "Q", Data: decomposition.Q.Data},
			{Name: "R", Data: decomposition.R.Data},
		},
	}

	requestBody, err := json.Marshal(payload)
	if err != nil {
		return domain.Statistics{}, fmt.Errorf("marshal statistics request: %w", err)
	}

	request, err := http.NewRequestWithContext(
		context.Background(),
		http.MethodPost,
		gateway.baseURL+"/api/v1/statistics",
		bytes.NewReader(requestBody),
	)
	if err != nil {
		return domain.Statistics{}, fmt.Errorf("create statistics request: %w", err)
	}

	request.Header.Set("Content-Type", "application/json")

	response, err := gateway.httpClient.Do(request)
	if err != nil {
		return domain.Statistics{}, domain.ErrStatisticsUnavailable
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return domain.Statistics{}, domain.ErrStatisticsUnavailable
	}

	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		return domain.Statistics{}, domain.ErrStatisticsUnavailable
	}

	var statisticsPayload statisticsResponse
	if err := json.Unmarshal(responseBody, &statisticsPayload); err != nil {
		return domain.Statistics{}, domain.ErrStatisticsUnavailable
	}

	return domain.Statistics{
		Max:               statisticsPayload.Max,
		Min:               statisticsPayload.Min,
		Average:           statisticsPayload.Average,
		Sum:               statisticsPayload.Sum,
		HasDiagonalMatrix: statisticsPayload.HasDiagonalMatrix,
		DiagonalMatrices:  statisticsPayload.DiagonalMatrices,
	}, nil
}
