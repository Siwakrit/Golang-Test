package user

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/yourusername/api/proto"
	"github.com/yourusername/internal/models"
)

// ListUsers retrieves a list of users with filtering and pagination
func (s *Service) ListUsers(ctx context.Context, req *proto.ListUsersRequest) (*proto.ListUsersResponse, error) {
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
