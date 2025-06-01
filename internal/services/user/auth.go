package user

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/yourusername/api/proto"
	"github.com/yourusername/internal/models"
)

// Login authenticates a user and returns a JWT token
func (s *Service) Login(ctx context.Context, req *proto.LoginRequest) (*proto.LoginResponse, error) {
	// Input validation
	if req.Email == "" || req.Password == "" {
		return nil, status.Error(codes.InvalidArgument, "email and password are required")
	}

	// Find user by email
	var user models.User
	err := s.db.UsersCol.FindOne(ctx, bson.M{"email": req.Email, "is_deleted": false}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, status.Error(codes.NotFound, "user not found")
		}
		return nil, status.Error(codes.Internal, "database error")
	}

	// Check password
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password))
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "invalid credentials")
	}

	// Generate token
	token, err := s.jwtManager.Generate(user.ID.Hex(), user.Email)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to generate token")
	}

	return &proto.LoginResponse{
		Token:  token,
		UserId: user.ID.Hex(),
	}, nil
}

// Logout invalidates a JWT token
func (s *Service) Logout(ctx context.Context, req *proto.LogoutRequest) (*proto.Empty, error) {
	if req.Token == "" {
		return nil, status.Error(codes.InvalidArgument, "token is required")
	}

	// Verify the token to get expiration time
	claims, err := s.jwtManager.Verify(req.Token)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid token")
	}

	// Store token in blacklist
	expiresAt := claims.ExpiresAt.Time
	_, err = s.db.TokensCol.InsertOne(ctx, models.Token{
		TokenStr:  req.Token,
		ExpiresAt: expiresAt,
	})
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to blacklist token")
	}

	return &proto.Empty{}, nil
}
