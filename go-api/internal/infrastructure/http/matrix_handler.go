package http

import (
	"errors"
	"log"

	"github.com/gofiber/fiber/v2"

	"matrix-api-challenge/go-api/internal/application"
	"matrix-api-challenge/go-api/internal/domain"
	"matrix-api-challenge/go-api/internal/infrastructure/http/dto"
)

type MatrixHandler struct {
	factorizeMatrixUseCase *application.FactorizeMatrixUseCase
}

func NewMatrixHandler(factorizeMatrixUseCase *application.FactorizeMatrixUseCase) *MatrixHandler {
	return &MatrixHandler{
		factorizeMatrixUseCase: factorizeMatrixUseCase,
	}
}

func (handler *MatrixHandler) HealthCheck(context *fiber.Ctx) error {
	return context.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "ok",
		"service": "go-api",
	})
}

func (handler *MatrixHandler) FactorizeMatrix(context *fiber.Ctx) error {
	var request dto.MatrixQRRequest
	if err := context.BodyParser(&request); err != nil {
		return writeError(context, fiber.StatusBadRequest, "invalid JSON payload", "INVALID_JSON")
	}

	matrix, err := domain.NewMatrix(request.Matrix)
	if err != nil {
		return writeDomainError(context, err)
	}

	result, err := handler.factorizeMatrixUseCase.Execute(matrix)
	if err != nil {
		return writeDomainError(context, err)
	}

	return context.Status(fiber.StatusOK).JSON(dto.MatrixQRResponseFromDomain(result))
}

func writeDomainError(context *fiber.Ctx, err error) error {
	switch {
	case errors.Is(err, domain.ErrEmptyMatrix),
		errors.Is(err, domain.ErrInvalidMatrixShape),
		errors.Is(err, domain.ErrInvalidMatrixValue):
		return writeError(context, fiber.StatusBadRequest, err.Error(), errorCode(err))
	case errors.Is(err, domain.ErrMatrixNotFactorizable):
		return writeError(context, fiber.StatusUnprocessableEntity, err.Error(), errorCode(err))
	case errors.Is(err, domain.ErrStatisticsUnavailable):
		return writeError(context, fiber.StatusBadGateway, err.Error(), errorCode(err))
	default:
		log.Printf("unexpected error: %v", err)
		return writeError(context, fiber.StatusInternalServerError, "internal server error", "INTERNAL_ERROR")
	}
}

func errorCode(err error) string {
	switch {
	case errors.Is(err, domain.ErrEmptyMatrix):
		return "EMPTY_MATRIX"
	case errors.Is(err, domain.ErrInvalidMatrixShape):
		return "INVALID_MATRIX_SHAPE"
	case errors.Is(err, domain.ErrInvalidMatrixValue):
		return "INVALID_MATRIX_VALUE"
	case errors.Is(err, domain.ErrMatrixNotFactorizable):
		return "MATRIX_NOT_FACTORIZABLE"
	case errors.Is(err, domain.ErrStatisticsUnavailable):
		return "STATISTICS_UNAVAILABLE"
	default:
		return "INTERNAL_ERROR"
	}
}

func writeError(context *fiber.Ctx, statusCode int, message string, code string) error {
	return context.Status(statusCode).JSON(dto.ErrorResponse{
		Error: message,
		Code:  code,
	})
}
