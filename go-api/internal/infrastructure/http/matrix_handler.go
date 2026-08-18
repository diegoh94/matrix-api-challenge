package http

import (
	"errors"
	"log"

	"github.com/gofiber/fiber/v2"

	"matrix-api-challenge/go-api/internal/domain"
	"matrix-api-challenge/go-api/internal/infrastructure/http/dto"
)

type matrixFactorizer interface {
	Execute(matrix domain.Matrix) (domain.FactorizeMatrixResult, error)
}

type MatrixHandler struct {
	factorizeMatrixUseCase matrixFactorizer
}

func NewMatrixHandler(factorizeMatrixUseCase matrixFactorizer) *MatrixHandler {
	return &MatrixHandler{
		factorizeMatrixUseCase: factorizeMatrixUseCase,
	}
}

func (handler *MatrixHandler) HealthCheck(ctx *fiber.Ctx) error {
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "ok",
		"service": "go-api",
	})
}

func (handler *MatrixHandler) FactorizeMatrix(ctx *fiber.Ctx) error {
	var request dto.MatrixQRRequest
	if err := ctx.BodyParser(&request); err != nil {
		return respondWithError(ctx, fiber.StatusBadRequest, "invalid JSON payload", "INVALID_JSON")
	}

	matrix, err := domain.NewMatrix(request.Matrix)
	if err != nil {
		return respondWithDomainError(ctx, err)
	}

	result, err := handler.factorizeMatrixUseCase.Execute(matrix)
	if err != nil {
		return respondWithDomainError(ctx, err)
	}

	return ctx.Status(fiber.StatusOK).JSON(dto.MatrixQRResponseFromDomain(result))
}

func respondWithDomainError(ctx *fiber.Ctx, err error) error {
	if spec, ok := resolveDomainError(err); ok {
		return respondWithError(ctx, spec.statusCode, err.Error(), spec.code)
	}

	log.Printf("unexpected error: %v", err)
	return respondWithError(ctx, fiber.StatusInternalServerError, "internal server error", "INTERNAL_ERROR")
}

func respondWithError(ctx *fiber.Ctx, statusCode int, message string, code string) error {
	return ctx.Status(statusCode).JSON(dto.ErrorResponse{
		Error: message,
		Code:  code,
	})
}

type domainErrorSpec struct {
	statusCode int
	code       string
}

func resolveDomainError(err error) (domainErrorSpec, bool) {
	switch {
	case errors.Is(err, domain.ErrEmptyMatrix):
		return domainErrorSpec{fiber.StatusBadRequest, "EMPTY_MATRIX"}, true
	case errors.Is(err, domain.ErrInvalidMatrixShape):
		return domainErrorSpec{fiber.StatusBadRequest, "INVALID_MATRIX_SHAPE"}, true
	case errors.Is(err, domain.ErrInvalidMatrixValue):
		return domainErrorSpec{fiber.StatusBadRequest, "INVALID_MATRIX_VALUE"}, true
	case errors.Is(err, domain.ErrMatrixNotFactorizable):
		return domainErrorSpec{fiber.StatusUnprocessableEntity, "MATRIX_NOT_FACTORIZABLE"}, true
	case errors.Is(err, domain.ErrStatisticsUnavailable):
		return domainErrorSpec{fiber.StatusBadGateway, "STATISTICS_UNAVAILABLE"}, true
	default:
		return domainErrorSpec{}, false
	}
}
