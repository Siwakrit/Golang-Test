package user

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/yourusername/api/proto"
	"github.com/yourusername/internal/models"
)

// ResetPassword initiates a password reset process
func (s *Service) ResetPassword(ctx context.Context, req *proto.ResetPasswordRequest) (*proto.Empty, error) {
	// Find user by email
	var user models.User
	err := s.db.UsersCol.FindOne(ctx, bson.M{"email": req.Email, "is_deleted": false}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			// Don't reveal if email exists or not (security best practice)
			return &proto.Empty{}, nil
		}
		return nil, status.Error(codes.Internal, "database error")
	}

	// Generate a reset token
	tokenBytes := make([]byte, 32)
	if _, err := rand.Read(tokenBytes); err != nil {
		return nil, status.Error(codes.Internal, "failed to generate token")
	}
	token := hex.EncodeToString(tokenBytes)

	// Save token to database
	now := time.Now()
	expiresAt := now.Add(24 * time.Hour) // Token valid for 24 hours

	resetToken := models.PasswordReset{
		UserID:    user.ID,
		Token:     token,
		CreatedAt: now,
		ExpiresAt: expiresAt,
		Used:      false,
	}

	_, err = s.db.ResetCol.InsertOne(ctx, resetToken)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to save reset token")
	}

	// In a real implementation, send email with reset link here
	// For this exercise, we'll just log the token
	fmt.Printf("Password reset token for %s: %s\n", user.Email, token)

	return &proto.Empty{}, nil
}

// ResetPasswordConfirm completes the password reset process
func (s *Service) ResetPasswordConfirm(ctx context.Context, req *proto.ResetPasswordConfirmRequest) (*proto.Empty, error) {
	// Validate password
	if len(req.NewPassword) < 8 {
		return nil, status.Error(codes.InvalidArgument, "password must be at least 8 characters")
	}

	// Find the reset token
	var resetToken models.PasswordReset
	err := s.db.ResetCol.FindOne(ctx, bson.M{
		"token":      req.Token,
		"used":       false,
		"expires_at": bson.M{"$gt": time.Now()},
	}).Decode(&resetToken)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, status.Error(codes.InvalidArgument, "invalid or expired token")
		}
		return nil, status.Error(codes.Internal, "database error")
	}

	// Hash new password
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to hash password")
	}

	// Update user's password
	_, err = s.db.UsersCol.UpdateOne(
		ctx,
		bson.M{"_id": resetToken.UserID},
		bson.M{"$set": bson.M{"password_hash": string(passwordHash), "updated_at": time.Now()}},
	)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to update password")
	}

	// Mark token as used
	_, err = s.db.ResetCol.UpdateOne(
		ctx,
		bson.M{"_id": resetToken.ID},
		bson.M{"$set": bson.M{"used": true}},
	)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to update token status")
	}

	return &proto.Empty{}, nil
}
