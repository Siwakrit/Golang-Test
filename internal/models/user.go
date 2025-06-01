package models

import (
    "time"

    "go.mongodb.org/mongo-driver/bson/primitive"
)

// User represents a user in the system
type User struct {
    ID           primitive.ObjectID `bson:"_id,omitempty"`
    Name         string             `bson:"name"`
    Email        string             `bson:"email"`
    PasswordHash string             `bson:"password_hash"`
    CreatedAt    time.Time          `bson:"created_at"`
    UpdatedAt    time.Time          `bson:"updated_at"`
    IsDeleted    bool               `bson:"is_deleted"`
}

// Token represents a blacklisted JWT token
type Token struct {
    ID        primitive.ObjectID `bson:"_id,omitempty"`
    TokenStr  string             `bson:"token"`
    ExpiresAt time.Time          `bson:"expires_at"`
}

// PasswordReset represents a password reset token
type PasswordReset struct {
    ID        primitive.ObjectID `bson:"_id,omitempty"`
    UserID    primitive.ObjectID `bson:"user_id"`
    Token     string             `bson:"token"`
    CreatedAt time.Time          `bson:"created_at"`
    ExpiresAt time.Time          `bson:"expires_at"`
    Used      bool               `bson:"used"`
}
