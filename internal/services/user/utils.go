package user

import (
	"context"

	"github.com/yourusername/internal/auth"
)

// getUserIDFromContext extracts the user ID from the authorization token
// This is a wrapper around the centralized auth.GetUserIDFromContext function
func (s *Service) getUserIDFromContext(ctx context.Context) (string, error) {
	return auth.GetUserIDFromContext(ctx, s.db, s.jwtManager)
}
