package user

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"github.com/yourusername/internal/models"
)

// getUserIDFromContext extracts the user ID from the authorization token
func (s *Service) getUserIDFromContext(ctx context.Context) (string, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", status.Error(codes.Unauthenticated, "metadata not provided")
	}

	values := md["authorization"]
	if len(values) == 0 {
		return "", status.Error(codes.Unauthenticated, "authorization token not provided")
	}

	tokenString := values[0]
	if len(tokenString) <= 7 || tokenString[:7] != "Bearer " {
		return "", status.Error(codes.Unauthenticated, "invalid token format")
	}

	tokenString = tokenString[7:] // Remove "Bearer " prefix

	// Check if token is blacklisted
	var blacklistedToken models.Token
	err := s.db.TokensCol.FindOne(ctx, bson.M{"token": tokenString}).Decode(&blacklistedToken)
	if err == nil {
		return "", status.Error(codes.Unauthenticated, "token has been revoked")
	} else if err != mongo.ErrNoDocuments {
		return "", status.Error(codes.Internal, "database error")
	}

	// Verify the token
	claims, err := s.jwtManager.Verify(tokenString)
	if err != nil {
		return "", status.Error(codes.Unauthenticated, "invalid token")
	}

	return claims.UserID, nil
}
