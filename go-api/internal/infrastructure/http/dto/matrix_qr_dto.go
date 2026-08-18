package dto

type MatrixQRRequest struct {
	Matrix [][]float64 `json:"matrix"`
}

type ErrorResponse struct {
	Error string `json:"error"`
	Code  string `json:"code"`
}

type AuthTokenRequest struct {
	APIKey string `json:"apiKey"`
}

type AuthTokenResponse struct {
	Token     string `json:"token"`
	ExpiresIn int64  `json:"expiresIn"`
}

type MatrixDimensionsResponse struct {
	Rows int `json:"rows"`
	Cols int `json:"cols"`
}

type MatrixValuesResponse struct {
	Q [][]float64 `json:"Q"`
	R [][]float64 `json:"R"`
}

type StatisticsResponse struct {
	Max               float64  `json:"max"`
	Min               float64  `json:"min"`
	Average           float64  `json:"average"`
	Sum               float64  `json:"sum"`
	HasDiagonalMatrix bool     `json:"hasDiagonalMatrix"`
	DiagonalMatrices  []string `json:"diagonalMatrices"`
}

type MatrixQRResponse struct {
	Input      MatrixDimensionsResponse `json:"input"`
	QR         MatrixValuesResponse     `json:"qr"`
	Statistics StatisticsResponse       `json:"statistics"`
}
