package ports

import "matrix-api-challenge/go-api/internal/domain"

type QRFactorizer interface {
	Factorize(matrix domain.Matrix) (domain.QRDecomposition, error)
}
