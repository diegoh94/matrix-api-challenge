package qr

import (
	"matrix-api-challenge/go-api/internal/domain"

	"gonum.org/v1/gonum/mat"
)

type GonumQRFactorizer struct{}

func NewGonumQRFactorizer() *GonumQRFactorizer {
	return &GonumQRFactorizer{}
}

func (factorizer *GonumQRFactorizer) Factorize(matrix domain.Matrix) (domain.QRDecomposition, error) {
	workMatrix := toDenseMatrix(matrix)

	var decomposition mat.QR
	decomposition.Factorize(workMatrix)

	var qDense mat.Dense
	decomposition.QTo(&qDense)

	var rDense mat.Dense
	decomposition.RTo(&rDense)

	return domain.QRDecomposition{
		Q: denseToMatrix(&qDense),
		R: denseToMatrix(&rDense),
	}, nil
}

func toDenseMatrix(matrix domain.Matrix) *mat.Dense {
	rows, cols := matrix.Rows, matrix.Cols
	buffer := make([]float64, rows*cols)

	for rowIndex, row := range matrix.Data {
		copy(buffer[rowIndex*cols:(rowIndex+1)*cols], row)
	}

	return mat.NewDense(rows, cols, buffer)
}

func denseToMatrix(denseMatrix *mat.Dense) domain.Matrix {
	rowCount, columnCount := denseMatrix.Dims()
	data := make([][]float64, rowCount)

	for rowIndex := 0; rowIndex < rowCount; rowIndex++ {
		row := make([]float64, columnCount)
		for columnIndex := 0; columnIndex < columnCount; columnIndex++ {
			row[columnIndex] = denseMatrix.At(rowIndex, columnIndex)
		}
		data[rowIndex] = row
	}

	return domain.Matrix{
		Rows: rowCount,
		Cols: columnCount,
		Data: data,
	}
}
