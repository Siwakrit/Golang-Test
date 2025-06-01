package services

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"regexp"
	"time"

	"github.com/yourusername/api/proto"
	"github.com/yourusername/internal/auth"
	"github.com/yourusername/internal/db"
	"github.com/yourusername/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// UserService implements the gRPC UserService
type UserService struct {
	proto.UnimplementedUserServiceServer
	db         *db.MongoDB
	jwtManager *auth.JWTManager
}

// NewUserService creates a new user service instance
func NewUserService(db *db.MongoDB, jwtManager *auth.JWTManager) *UserService {
	return &UserService{
		db:         db,
		jwtManager: jwtManager,
	}
}

// Login authenticates a user and returns a JWT token
func (s *UserService) Login(ctx context.Context, req *proto.LoginRequest) (*proto.LoginResponse, error) {
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

// Register creates a new user account
func (s *UserService) Register(ctx context.Context, req *proto.RegisterRequest) (*proto.User, error) {
	// Validate inputs
	if req.Email == "" || req.Password == "" || req.Name == "" {
		return nil, status.Error(codes.InvalidArgument, "name, email and password are required")
	}

	// Validate email format
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(req.Email) {
		return nil, status.Error(codes.InvalidArgument, "invalid email format")
	}

	// Check password strength (at least 8 chars, with letters and numbers)
	if len(req.Password) < 8 {
		return nil, status.Error(codes.InvalidArgument, "password must be at least 8 characters")
	}
	hasLetter := regexp.MustCompile(`[a-zA-Z]`).MatchString(req.Password)
	hasNumber := regexp.MustCompile(`[0-9]`).MatchString(req.Password)
	if !hasLetter || !hasNumber {
		return nil, status.Error(codes.InvalidArgument, "password must contain both letters and numbers")
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

// Logout invalidates a JWT token
func (s *UserService) Logout(ctx context.Context, req *proto.LogoutRequest) (*proto.Empty, error) {
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

// GetProfile retrieves a user's profile
func (s *UserService) GetProfile(ctx context.Context, req *proto.GetProfileRequest) (*proto.User, error) {
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
func (s *UserService) UpdateProfile(ctx context.Context, req *proto.UpdateProfileRequest) (*proto.User, error) {
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
func (s *UserService) DeleteProfile(ctx context.Context, req *proto.DeleteProfileRequest) (*proto.Empty, error) {
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

// ListUsers retrieves a list of users with filtering and pagination
func (s *UserService) ListUsers(ctx context.Context, req *proto.ListUsersRequest) (*proto.ListUsersResponse, error) {
	// Set default pagination values if not provided
	page := int64(1)
	limit := int64(10)

	if req.Page > 0 {
		page = int64(req.Page)
	}

	if req.Limit > 0 && req.Limit <= 100 {
		limit = int64(req.Limit)
	}

	skip := (page - 1) * limit

	// Build filter
	filter := bson.M{"is_deleted": false}

	if req.NameFilter != "" {
		filter["name"] = bson.M{"$regex": req.NameFilter, "$options": "i"}
	}

	if req.EmailFilter != "" {
		filter["email"] = bson.M{"$regex": req.EmailFilter, "$options": "i"}
	}

	// Get total count
	total, err := s.db.UsersCol.CountDocuments(ctx, filter)
	if err != nil {
		return nil, status.Error(codes.Internal, "database error")
	}

	// Get users
	findOptions := options.Find().
		SetSkip(skip).
		SetLimit(limit).
		SetSort(bson.M{"created_at": -1})

	cursor, err := s.db.UsersCol.Find(ctx, filter, findOptions)
	if err != nil {
		return nil, status.Error(codes.Internal, "database error")
	}
	defer cursor.Close(ctx)

	// Decode users
	var users []*proto.User
	for cursor.Next(ctx) {
		var user models.User
		if err := cursor.Decode(&user); err != nil {
			return nil, status.Error(codes.Internal, "error decoding user")
		}

		users = append(users, &proto.User{
			Id:        user.ID.Hex(),
			Name:      user.Name,
			Email:     user.Email,
			CreatedAt: user.CreatedAt.Format(time.RFC3339),
			UpdatedAt: user.UpdatedAt.Format(time.RFC3339),
		})
	}

	if err := cursor.Err(); err != nil {
		return nil, status.Error(codes.Internal, "cursor error")
	}

	return &proto.ListUsersResponse{
		Users: users,
		Total: int32(total),
		Page:  int32(page),
		Limit: int32(limit),
	}, nil
}

// ResetPassword initiates a password reset process
func (s *UserService) ResetPassword(ctx context.Context, req *proto.ResetPasswordRequest) (*proto.Empty, error) {
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
func (s *UserService) ResetPasswordConfirm(ctx context.Context, req *proto.ResetPasswordConfirmRequest) (*proto.Empty, error) {
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

// getUserIDFromContext extracts the user ID from the authorization token
func (s *UserService) getUserIDFromContext(ctx context.Context) (string, error) {
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
