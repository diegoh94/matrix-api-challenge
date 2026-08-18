package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"

	"matrix-api-challenge/go-api/internal/infrastructure/auth"
	"matrix-api-challenge/go-api/internal/infrastructure/http/dto"
)

const authTokenLocalKey = "authToken"

func JWTAuth(tokenService *auth.TokenService) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		authorizationHeader := ctx.Get("Authorization")
		if !strings.HasPrefix(authorizationHeader, "Bearer ") {
			return respondUnauthorized(ctx, "missing or invalid authorization header")
		}

		accessToken := strings.TrimSpace(strings.TrimPrefix(authorizationHeader, "Bearer "))
		if accessToken == "" {
			return respondUnauthorized(ctx, "missing or invalid authorization header")
		}

		if err := tokenService.ValidateAccessToken(accessToken); err != nil {
			return respondUnauthorized(ctx, "invalid or expired token")
		}

		ctx.Locals(authTokenLocalKey, accessToken)

		return ctx.Next()
	}
}

func AuthTokenFromRequest(ctx *fiber.Ctx) string {
	token, _ := ctx.Locals(authTokenLocalKey).(string)
	return token
}

func respondUnauthorized(ctx *fiber.Ctx, message string) error {
	return ctx.Status(fiber.StatusUnauthorized).JSON(dto.ErrorResponse{
		Error: message,
		Code:  "UNAUTHORIZED",
	})
}
