package statistics

type statisticsRequest struct {
	Matrices []namedMatrixPayload `json:"matrices"`
}

type namedMatrixPayload struct {
	Name string      `json:"name"`
	Data [][]float64 `json:"data"`
}

type statisticsResponse struct {
	Max               float64  `json:"max"`
	Min               float64  `json:"min"`
	Average           float64  `json:"average"`
	Sum               float64  `json:"sum"`
	HasDiagonalMatrix bool     `json:"hasDiagonalMatrix"`
	DiagonalMatrices  []string `json:"diagonalMatrices"`
}
