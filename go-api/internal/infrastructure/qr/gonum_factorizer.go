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
	denseMatrix := mat.NewDense(matrix.Rows, matrix.Cols, nil)

	for rowIndex := 0; rowIndex < matrix.Rows; rowIndex++ {
		denseMatrix.SetRow(rowIndex, matrix.Data[rowIndex])
	}

	var qrFactorization mat.QR
	qrFactorization.Factorize(denseMatrix)

	var qMatrix mat.Dense
	qrFactorization.QTo(&qMatrix)

	var rMatrix mat.Dense
	qrFactorization.RTo(&rMatrix)

	return domain.QRDecomposition{
		Q: matrixFromDense(&qMatrix),
		R: matrixFromDense(&rMatrix),
	}, nil
}

func matrixFromDense(denseMatrix *mat.Dense) domain.Matrix {
	rowCount, columnCount := denseMatrix.Dims()
	data := make([][]float64, rowCount)

	for rowIndex := 0; rowIndex < rowCount; rowIndex++ {
		data[rowIndex] = make([]float64, columnCount)

		for columnIndex := 0; columnIndex < columnCount; columnIndex++ {
			data[rowIndex][columnIndex] = denseMatrix.At(rowIndex, columnIndex)
		}
	}

	return domain.Matrix{
		Rows: rowCount,
		Cols: columnCount,
		Data: data,
	}
}
