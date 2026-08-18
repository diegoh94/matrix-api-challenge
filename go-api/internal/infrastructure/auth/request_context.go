package auth

import (
	"context"
)

type contextKey struct {
}

var authTokenKey = contextKey{}

func ContextWithToken(parent context.Context, token string) context.Context {
	return context.WithValue(parent, authTokenKey, token)
}

func TokenFromContext(ctx context.Context) (string, bool) {
	token, ok := ctx.Value(authTokenKey).(string)
	return token, ok && token != ""
}
