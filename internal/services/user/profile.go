package user

import (
	"context"
	"regexp"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/yourusername/api/proto"
	"github.com/yourusername/internal/models"
)

// GetProfile retrieves a user's profile
func (s *Service) GetProfile(ctx context.Context, req *proto.GetProfileRequest) (*proto.User, error) {
	// Check authorization
	userID, err := s.getUserIDFromContext(ctx)
	if err != nil {
		return nil, err
	}

	// Convert the ID string to ObjectID
	objectID, err := primitive.ObjectIDFromHex(req.UserId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid user ID")
	}

	// Users can only access their own profile unless they have an admin role
	// In a real system, you'd check roles here
	if userID != req.UserId {
		return nil, status.Error(codes.PermissionDenied, "cannot access other user's profile")
	}

	// Find user
	var user models.User
	err = s.db.UsersCol.FindOne(ctx, bson.M{"_id": objectID, "is_deleted": false}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, status.Error(codes.NotFound, "user not found")
		}
		return nil, status.Error(codes.Internal, "database error")
	}
	// Return user profile
	return &proto.User{
		Id:        user.ID.Hex(),
		Name:      user.Name,
		Email:     user.Email,
		CreatedAt: user.CreatedAt.Format(time.RFC3339),
		UpdatedAt: user.UpdatedAt.Format(time.RFC3339),
	}, nil
}

// UpdateProfile updates a user's profile information
func (s *Service) UpdateProfile(ctx context.Context, req *proto.UpdateProfileRequest) (*proto.User, error) {
	// Check authorization
	userID, err := s.getUserIDFromContext(ctx)
	if err != nil {
		return nil, err
	}

	// Users can only update their own profile
	if userID != req.UserId {
		return nil, status.Error(codes.PermissionDenied, "cannot update other user's profile")
	}

	objectID, err := primitive.ObjectIDFromHex(req.UserId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid user ID")
	}

	// Prepare update fields
	update := bson.M{"updated_at": time.Now()}

	if req.Name != "" {
		update["name"] = req.Name
	}

	if req.Email != "" {
		// Validate email format
		emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
		if !emailRegex.MatchString(req.Email) {
			return nil, status.Error(codes.InvalidArgument, "invalid email format")
		}

		// Check if email is already used by another user
		count, err := s.db.UsersCol.CountDocuments(ctx, bson.M{
			"email": req.Email,
			"_id":   bson.M{"$ne": objectID},
		})
		if err != nil {
			return nil, status.Error(codes.Internal, "database error")
		}
		if count > 0 {
			return nil, status.Error(codes.AlreadyExists, "email already in use")
		}

		update["email"] = req.Email
	}

	// Update user
	result := s.db.UsersCol.FindOneAndUpdate(
		ctx,
		bson.M{"_id": objectID, "is_deleted": false},
		bson.M{"$set": update},
		options.FindOneAndUpdate().SetReturnDocument(options.After),
	)

	// Check for errors
	var user models.User
	if err := result.Decode(&user); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, status.Error(codes.NotFound, "user not found")
		}
		return nil, status.Error(codes.Internal, "failed to update user")
	}
	// Return updated user
	return &proto.User{
		Id:        user.ID.Hex(),
		Name:      user.Name,
		Email:     user.Email,
		CreatedAt: user.CreatedAt.Format(time.RFC3339),
		UpdatedAt: user.UpdatedAt.Format(time.RFC3339),
	}, nil
}

// DeleteProfile soft-deletes a user profile
func (s *Service) DeleteProfile(ctx context.Context, req *proto.DeleteProfileRequest) (*proto.Empty, error) {
	// Check authorization
	userID, err := s.getUserIDFromContext(ctx)
	if err != nil {
		return nil, err
	}

	// Users can only delete their own profile
	if userID != req.UserId {
		return nil, status.Error(codes.PermissionDenied, "cannot delete other user's profile")
	}

	objectID, err := primitive.ObjectIDFromHex(req.UserId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid user ID")
	}

	// Soft delete by setting is_deleted flag
	update := bson.M{
		"$set": bson.M{
			"is_deleted": true,
			"updated_at": time.Now(),
		},
	}

	result, err := s.db.UsersCol.UpdateOne(ctx, bson.M{"_id": objectID}, update)
	if err != nil {
		return nil, status.Error(codes.Internal, "database error")
	}

	if result.MatchedCount == 0 {
		return nil, status.Error(codes.NotFound, "user not found")
	}

	return &proto.Empty{}, nil
}
