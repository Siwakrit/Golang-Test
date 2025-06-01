package user

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/yourusername/api/proto"
	"github.com/yourusername/internal/models"
	"github.com/yourusername/internal/utils"
)

// Register creates a new user account
func (s *Service) Register(ctx context.Context, req *proto.RegisterRequest) (*proto.User, error) {
	// Validate name
	if err := utils.ValidateName(req.Name); err != nil {
		return nil, err
	}

	// Validate email
	if err := utils.ValidateEmail(req.Email); err != nil {
		return nil, err
	}

	// Validate password
	if err := utils.ValidatePassword(req.Password); err != nil {
		return nil, err
	}

	// Check if email is already registered
	count, err := s.db.UsersCol.CountDocuments(ctx, bson.M{"email": req.Email})
	if err != nil {
		return nil, status.Error(codes.Internal, "database error")
	}
	if count > 0 {
		return nil, status.Error(codes.AlreadyExists, "email already registered")
	}

	// Hash password
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to hash password")
	}

	// Create new user
	now := time.Now()
	user := models.User{
		ID:           primitive.NewObjectID(),
		Name:         req.Name,
		Email:        req.Email,
		PasswordHash: string(passwordHash),
		CreatedAt:    now,
		UpdatedAt:    now,
		IsDeleted:    false,
	}

	_, err = s.db.UsersCol.InsertOne(ctx, user)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to create user")
	}
	// Return user without password
	return &proto.User{
		Id:        user.ID.Hex(),
		Name:      user.Name,
		Email:     user.Email,
		CreatedAt: user.CreatedAt.Format(time.RFC3339),
		UpdatedAt: user.UpdatedAt.Format(time.RFC3339),
	}, nil
}
