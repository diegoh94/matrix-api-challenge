package http

import (
	"context"
	"errors"
	"log"

	"github.com/gofiber/fiber/v2"

	"matrix-api-challenge/go-api/internal/domain"
	"matrix-api-challenge/go-api/internal/infrastructure/auth"
	"matrix-api-challenge/go-api/internal/infrastructure/http/dto"
	"matrix-api-challenge/go-api/internal/infrastructure/http/middleware"
)

type matrixFactorizer interface {
	Execute(ctx context.Context, matrix domain.Matrix) (domain.FactorizeMatrixResult, error)
}

type MatrixHandler struct {
	factorizeMatrixUseCase matrixFactorizer
}

func NewMatrixHandler(factorizeMatrixUseCase matrixFactorizer) *MatrixHandler {
	return &MatrixHandler{
		factorizeMatrixUseCase: factorizeMatrixUseCase,
	}
}

// HealthCheck returns the service health status.
//
// @Summary Health check
// @Tags Health
// @Produce json
// @Success 200 {object} dto.HealthResponse
// @Router /health [get]
func (handler *MatrixHandler) HealthCheck(ctx *fiber.Ctx) error {
	return ctx.Status(fiber.StatusOK).JSON(dto.HealthResponse{
		Status:  "ok",
		Service: "go-api",
	})
}

// FactorizeMatrix performs QR decomposition and returns statistics.
//
// @Summary Factorize matrix using QR decomposition
// @Tags Matrix
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.MatrixQRRequest true "Matrix input"
// @Success 200 {object} dto.MatrixQRResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 422 {object} dto.ErrorResponse
// @Failure 502 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/v1/matrix/qr [post]
func (handler *MatrixHandler) FactorizeMatrix(ctx *fiber.Ctx) error {
	var request dto.MatrixQRRequest
	if err := ctx.BodyParser(&request); err != nil {
		return respondWithError(ctx, fiber.StatusBadRequest, "invalid JSON payload", "INVALID_JSON")
	}

	matrix, err := domain.NewMatrix(request.Matrix)
	if err != nil {
		return respondWithDomainError(ctx, err)
	}

	requestContext := auth.ContextWithToken(ctx.Context(), middleware.AuthTokenFromRequest(ctx))

	result, err := handler.factorizeMatrixUseCase.Execute(requestContext, matrix)
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
