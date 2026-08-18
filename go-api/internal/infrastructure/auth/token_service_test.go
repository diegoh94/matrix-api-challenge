package auth_test

import (
	"testing"
	"time"

	"matrix-api-challenge/go-api/internal/infrastructure/auth"
)

func TestTokenServiceCreateAndValidateAccessToken(t *testing.T) {
	service := auth.NewTokenService("shared-secret", time.Hour)

	token, expiresInSeconds, err := service.CreateAccessToken()
	if err != nil {
		t.Fatalf("create token: %v", err)
	}

	if token == "" || expiresInSeconds != int64(time.Hour.Seconds()) {
		t.Fatalf("unexpected token response: token=%q expiresIn=%d", token, expiresInSeconds)
	}

	if err := service.ValidateAccessToken(token); err != nil {
		t.Fatalf("expected created token to validate, got %v", err)
	}
}

func TestTokenServiceRejectsTokenSignedWithDifferentSecret(t *testing.T) {
	issuer := auth.NewTokenService("issuer-secret", time.Hour)
	validator := auth.NewTokenService("validator-secret", time.Hour)

	token, _, err := issuer.CreateAccessToken()
	if err != nil {
		t.Fatalf("create token: %v", err)
	}

	if err := validator.ValidateAccessToken(token); err == nil {
		t.Fatal("expected token signed by another service to be rejected")
	}
}
