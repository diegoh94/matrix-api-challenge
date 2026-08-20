package dto

type MatrixRotateRequest struct {
	Matrix  [][]float64 `json:"matrix"`
	Degrees int         `json:"degrees"`
}

type MatrixRotateResponse struct {
	Input      MatrixDimensionsResponse `json:"input"`
	Degrees    int                      `json:"degrees"`
	Rotated    RotatedMatrixResponse    `json:"rotated"`
	Statistics StatisticsResponse       `json:"statistics"`
}

type RotatedMatrixResponse struct {
	Rows int           `json:"rows"`
	Cols int           `json:"cols"`
	Data [][]float64   `json:"data"`
}
