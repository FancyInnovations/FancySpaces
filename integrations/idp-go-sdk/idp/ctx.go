package idp

import (
	"context"
	"log/slog"
)

type userCtxKey struct{}

// attachUserToCtx attaches the user to the context for later retrieval in handlers.
func attachUserToCtx(ctx context.Context, u *User) context.Context {
	return context.WithValue(ctx, userCtxKey{}, *u)
}

// UserFromCtx retrieves the user from the context. It returns nil if no user is found.
func UserFromCtx(ctx context.Context) *User {
	u, ok := ctx.Value(userCtxKey{}).(User)
	if !ok {
		slog.Warn("User not found in context")
		return nil
	}
	return &u
}
