package application

import (
	"matrix-api-challenge/go-api/internal/domain"
	"matrix-api-challenge/go-api/internal/domain/ports"
)

type FactorizeMatrixUseCase struct {
	qrFactorizer ports.QRFactorizer
}

func NewFactorizeMatrixUseCase(qrFactorizer ports.QRFactorizer) *FactorizeMatrixUseCase {
	return &FactorizeMatrixUseCase{
		qrFactorizer: qrFactorizer,
	}
}

func (useCase *FactorizeMatrixUseCase) Execute(matrix domain.Matrix) (domain.QRDecomposition, error) {
	return useCase.qrFactorizer.Factorize(matrix)
}
