package middleware

import (
	"context"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"github.com/yourusername/internal/auth"
	"github.com/yourusername/internal/db"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// AuthInterceptor is a gRPC interceptor for authentication
type AuthInterceptor struct {
	jwtManager *auth.JWTManager
	db         *db.MongoDB
}

// NewAuthInterceptor creates a new auth interceptor
func NewAuthInterceptor(jwtManager *auth.JWTManager, db *db.MongoDB) *AuthInterceptor {
	return &AuthInterceptor{
		jwtManager: jwtManager,
		db:         db,
	}
}

// Unary returns a server interceptor function to authenticate and authorize unary RPC
func (interceptor *AuthInterceptor) Unary() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) { // Skip auth for login and register methods
		if info.FullMethod == "/proto.UserService/Login" ||
			info.FullMethod == "/proto.UserService/Register" ||
			info.FullMethod == "/proto.UserService/ResetPassword" ||
			info.FullMethod == "/proto.UserService/ResetPasswordConfirm" {
			return handler(ctx, req)
		}

		// Check authorization for all other methods
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, status.Error(codes.Unauthenticated, "metadata is not provided")
		}

		values := md["authorization"]
		if len(values) == 0 {
			return nil, status.Error(codes.Unauthenticated, "authorization token is not provided")
		}

		accessToken := values[0]
		if !strings.HasPrefix(accessToken, "Bearer ") {
			return nil, status.Error(codes.Unauthenticated, "invalid token format")
		}

		tokenStr := strings.TrimPrefix(accessToken, "Bearer ")

		// Check if token is blacklisted
		var blacklistedToken struct{}
		err := interceptor.db.TokensCol.FindOne(
			ctx,
			bson.M{"token": tokenStr},
		).Decode(&blacklistedToken)

		if err == nil {
			return nil, status.Error(codes.Unauthenticated, "token has been revoked")
		} else if err != mongo.ErrNoDocuments {
			return nil, status.Error(codes.Internal, "database error")
		}

		// Verify the token
		claims, err := interceptor.jwtManager.Verify(tokenStr)
		if err != nil {
			return nil, status.Errorf(codes.Unauthenticated, "invalid token: %v", err)
		}

		// Add user ID to the context
		ctxWithUserID := context.WithValue(ctx, "user_id", claims.UserID)
		return handler(ctxWithUserID, req)
	}
}
