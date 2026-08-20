package statistics_test

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"matrix-api-challenge/go-api/internal/domain"
	"matrix-api-challenge/go-api/internal/infrastructure/auth"
	"matrix-api-challenge/go-api/internal/infrastructure/statistics"
)

func TestNodeStatisticsGatewayComputeStatistics(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(responseWriter http.ResponseWriter, request *http.Request) {
		if request.Method != http.MethodPost || request.URL.Path != "/api/v1/statistics" {
			t.Fatalf("unexpected request: %s %s", request.Method, request.URL.Path)
		}

		var payload struct {
			Matrices []struct {
				Name string          `json:"name"`
				Data [][]float64     `json:"data"`
			} `json:"matrices"`
		}

		if err := json.NewDecoder(request.Body).Decode(&payload); err != nil {
			t.Fatalf("decode request body: %v", err)
		}

		if len(payload.Matrices) != 2 || payload.Matrices[0].Name != "Q" || payload.Matrices[1].Name != "R" {
			t.Fatalf("unexpected matrices payload: %+v", payload.Matrices)
		}

		responseWriter.Header().Set("Content-Type", "application/json")
		json.NewEncoder(responseWriter).Encode(map[string]any{
			"max":               4.0,
			"min":               0.0,
			"average":           1.375,
			"sum":               11.0,
			"hasDiagonalMatrix": true,
			"diagonalMatrices":  []string{"Q"},
		})
	}))
	defer server.Close()

	gateway := statistics.NewNodeStatisticsGateway(statistics.GatewayConfig{
		BaseURL: server.URL,
		Timeout: 2 * time.Second,
	})

	result, err := gateway.ComputeStatistics(context.Background(), []domain.NamedMatrix{
		{Name: "Q", Matrix: domain.Matrix{
			Rows: 2,
			Cols: 2,
			Data: [][]float64{{1, 0}, {0, 1}},
		}},
		{Name: "R", Matrix: domain.Matrix{
			Rows: 2,
			Cols: 2,
			Data: [][]float64{{2, 3}, {0, 4}},
		}},
	})
	if err != nil {
		t.Fatalf("unexpected gateway error: %v", err)
	}

	if result.Max != 4 || result.Min != 0 || result.Sum != 11 {
		t.Fatalf("unexpected statistics values: %+v", result)
	}

	if !result.HasDiagonalMatrix || len(result.DiagonalMatrices) != 1 || result.DiagonalMatrices[0] != "Q" {
		t.Fatalf("unexpected diagonal result: %+v", result)
	}
}

func TestNodeStatisticsGatewayReturnsUnavailableOnServerError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(responseWriter http.ResponseWriter, _ *http.Request) {
		responseWriter.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	gateway := statistics.NewNodeStatisticsGateway(statistics.GatewayConfig{
		BaseURL: server.URL,
		Timeout: 2 * time.Second,
	})

	_, err := gateway.ComputeStatistics(context.Background(), nil)
	if !errors.Is(err, domain.ErrStatisticsUnavailable) {
		t.Fatalf("expected ErrStatisticsUnavailable, got %v", err)
	}
}

func TestNodeStatisticsGatewayForwardsBearerToken(t *testing.T) {
	const accessToken = "signed-token-from-go-api"

	server := httptest.NewServer(http.HandlerFunc(func(responseWriter http.ResponseWriter, request *http.Request) {
		if got := request.Header.Get("Authorization"); got != "Bearer "+accessToken {
			t.Fatalf("expected forwarded bearer token, got %q", got)
		}

		responseWriter.Header().Set("Content-Type", "application/json")
		json.NewEncoder(responseWriter).Encode(map[string]any{
			"max":               1.0,
			"min":               0.0,
			"average":           0.5,
			"sum":               1.0,
			"hasDiagonalMatrix": false,
			"diagonalMatrices":  []string{},
		})
	}))
	defer server.Close()

	gateway := statistics.NewNodeStatisticsGateway(statistics.GatewayConfig{
		BaseURL: server.URL,
		Timeout: 2 * time.Second,
	})

	ctx := auth.ContextWithToken(context.Background(), accessToken)

	_, err := gateway.ComputeStatistics(ctx, []domain.NamedMatrix{
		{Name: "Q", Matrix: domain.Matrix{Rows: 1, Cols: 1, Data: [][]float64{{1}}}},
		{Name: "R", Matrix: domain.Matrix{Rows: 1, Cols: 1, Data: [][]float64{{1}}}},
	})
	if err != nil {
		t.Fatalf("unexpected gateway error: %v", err)
	}
}
