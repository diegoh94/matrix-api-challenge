package domain

type Matrix struct {
	Rows int
	Cols int
	Data [][]float64
}

type NamedMatrix struct {
	Name string
	Matrix
}

type QRDecomposition struct {
	Q Matrix
	R Matrix
}

type Statistics struct {
	Max                 float64
	Min                 float64
	Average             float64
	Sum                 float64
	HasDiagonalMatrix   bool
	DiagonalMatrices    []string
}

type FactorizeMatrixResult struct {
	Input      Matrix
	QR         QRDecomposition
	Statistics Statistics
}

type RotateMatrixResult struct {
	Input      Matrix
	Degrees    int
	Rotated    Matrix
	Statistics Statistics
}
