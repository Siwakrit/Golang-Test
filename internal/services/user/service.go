package user

import (
	"github.com/yourusername/api/proto"
	"github.com/yourusername/internal/auth"
	"github.com/yourusername/internal/db"
)

// Service implements the gRPC UserService
type Service struct {
	proto.UnimplementedUserServiceServer
	db         *db.MongoDB
	jwtManager *auth.JWTManager
}

// NewService creates a new user service instance
func NewService(db *db.MongoDB, jwtManager *auth.JWTManager) *Service {
	return &Service{
		db:         db,
		jwtManager: jwtManager,
	}
}

// This file contains the main UserService struct and constructor.
// The implementation of each method is in separate files organized by functionality:
// - auth.go: Login and Logout methods
// - registration.go: Register method
// - profile.go: GetProfile, UpdateProfile, and DeleteProfile methods
// - list.go: ListUsers method
// - password.go: ResetPassword and ResetPasswordConfirm methods
// - utils.go: getUserIDFromContext utility method
