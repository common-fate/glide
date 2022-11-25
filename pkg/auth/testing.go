package auth

import (
	"context"

	"github.com/common-fate/common-fate/pkg/identity"
)

// TestingSetUserID allows the user ID to be set in the context for testing purposes.
func TestingSetUserID(ctx context.Context, userID string) context.Context {
	return context.WithValue(ctx, userIDContext, userID)
}

// TestingSetUserID allows the user ID to be set in the context for testing purposes.
func TestingSetUser(ctx context.Context, user identity.User) context.Context {
	return context.WithValue(ctx, userContext, &user)
}

// TestingSetIsAdmin allows the isAdmin to be set in the context for testing purposes.
func TestingSetIsAdmin(ctx context.Context, isAdmin bool) context.Context {
	return context.WithValue(ctx, adminContext, isAdmin)
}
