package http

import (
	"github.com/gofiber/fiber/v2"

	"matrix-api-challenge/go-api/internal/infrastructure/auth"
	"matrix-api-challenge/go-api/internal/infrastructure/http/dto"
)

type AuthHandler struct {
	tokenService *auth.TokenService
	apiKey       string
}

func NewAuthHandler(tokenService *auth.TokenService, apiKey string) *AuthHandler {
	return &AuthHandler{
		tokenService: tokenService,
		apiKey:       apiKey,
	}
}

// IssueToken generates a JWT access token from a valid API key.
//
// @Summary Issue JWT access token
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body dto.AuthTokenRequest true "API key"
// @Success 200 {object} dto.AuthTokenResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /auth/token [post]
func (handler *AuthHandler) IssueToken(ctx *fiber.Ctx) error {
	var request dto.AuthTokenRequest
	if err := ctx.BodyParser(&request); err != nil {
		return respondWithError(ctx, fiber.StatusBadRequest, "invalid JSON payload", "INVALID_JSON")
	}

	if request.APIKey != handler.apiKey {
		return respondWithError(ctx, fiber.StatusUnauthorized, "invalid api key", "UNAUTHORIZED")
	}

	accessToken, expiresInSeconds, err := handler.tokenService.CreateAccessToken()
	if err != nil {
		return respondWithError(ctx, fiber.StatusInternalServerError, "internal server error", "INTERNAL_ERROR")
	}

	return ctx.Status(fiber.StatusOK).JSON(dto.AuthTokenResponse{
		Token:     accessToken,
		ExpiresIn: expiresInSeconds,
	})
}
