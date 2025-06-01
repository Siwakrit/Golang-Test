# คู่มือการสร้างโปรเจค User Management API

คู่มือนี้จะอธิบายขั้นตอนทั้งหมดในการสร้างโปรเจค User Management API ตั้งแต่เริ่มต้น ทั้งการติดตั้ง dependencies, การสร้างโครงสร้างโปรเจค, การเขียนโค้ด, และการรันโปรเจค

## สารบัญ
1. [การติดตั้ง Prerequisites](#1-การติดตั้ง-prerequisites)
2. [สร้างโครงสร้างโปรเจค](#2-สร้างโครงสร้างโปรเจค)
3. [ตั้งค่า Go Module](#3-ตั้งค่า-go-module)
4. [สร้างไฟล์ Protocol Buffer](#4-สร้างไฟล์-protocol-buffer)
5. [สร้าง Models](#5-สร้าง-models)
6. [การจัดการฐานข้อมูล MongoDB](#6-การจัดการฐานข้อมูล-mongodb)
7. [สร้าง JWT Authentication](#7-สร้าง-jwt-authentication)
8. [สร้าง Middleware](#8-สร้าง-middleware)
9. [พัฒนา Service Logic](#9-พัฒนา-service-logic)
10. [สร้าง Server](#10-สร้าง-server)
11. [รัน และ ทดสอบโปรเจค](#11-รัน-และ-ทดสอบโปรเจค)
12. [การตั้งค่า Docker](#12-การตั้งค่า-docker)

## 1. การติดตั้ง Prerequisites

ก่อนเริ่มต้นเราต้องติดตั้งซอฟต์แวร์ที่จำเป็น:

### 1.1 ติดตั้ง Go

1. ดาวน์โหลด Go จาก [golang.org/dl](https://golang.org/dl/)
2. ติดตั้งตามคำแนะนำของระบบปฏิบัติการ
3. ตรวจสอบการติดตั้ง:
   ```powershell
   go version
   ```

### 1.2 ติดตั้ง MongoDB

1. ดาวน์โหลด MongoDB Community Edition จาก [mongodb.com](https://www.mongodb.com/try/download/community)
2. ติดตั้งตามคำแนะนำ (ให้เลือก "Install MongoDB as a Service" สำหรับ Windows)
3. ตรวจสอบการติดตั้ง:
   ```powershell
   mongod --version
   ```

### 1.3 ติดตั้ง Protocol Buffers Compiler

1. ดาวน์โหลด protoc จาก [github.com/protocolbuffers/protobuf/releases](https://github.com/protocolbuffers/protobuf/releases)
2. แตกไฟล์และเพิ่ม bin directory เข้าไปใน PATH
3. ติดตั้ง Go plugins สำหรับ protoc:
   ```powershell
   go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
   go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
   ```
4. ตรวจสอบการติดตั้ง:
   ```powershell
   protoc --version
   ```

### 1.4 ติดตั้ง gRPCurl (สำหรับการทดสอบ)

1. ดาวน์โหลดจาก [github.com/fullstorydev/grpcurl/releases](https://github.com/fullstorydev/grpcurl/releases)
2. แตกไฟล์และเพิ่ม executable เข้าไปใน PATH
3. ตรวจสอบการติดตั้ง:
   ```powershell
   grpcurl --version
   ```

## 2. สร้างโครงสร้างโปรเจค

สร้างโครงสร้างไดเรกทอรีสำหรับโปรเจค:

```powershell
# สร้างไดเรกทอรีหลักของโปรเจค
mkdir -p Golang-Test
cd Golang-Test

# สร้างโครงสร้างไดเรกทอรี
mkdir -p api/proto
mkdir -p cmd/server
mkdir -p docs
mkdir -p internal/{auth,config,db,middleware,models,services}
mkdir -p test
```

## 3. ตั้งค่า Go Module

ตั้งค่า Go module สำหรับการจัดการ dependencies:

```powershell
# ภายในไดเรกทอรี Golang-Test
go mod init github.com/yourusername
```

จากนั้นติดตั้ง dependencies ที่จำเป็น:

```powershell
go get -u github.com/golang-jwt/jwt/v5
go get -u go.mongodb.org/mongo-driver/mongo
go get -u golang.org/x/crypto/bcrypt
go get -u google.golang.org/grpc
go get -u google.golang.org/protobuf
go get -u github.com/joho/godotenv
```

## 4. สร้างไฟล์ Protocol Buffer

สร้างไฟล์ `api/proto/user_service.proto`:

```powershell
# สร้างไฟล์ proto
New-Item -Path "api/proto/user_service.proto" -ItemType File
```

เพิ่มเนื้อหาต่อไปนี้ในไฟล์:

```protobuf
syntax = "proto3";
package proto;

option go_package = "github.com/yourusername/api/proto";

// Common empty response
message Empty {}

// Auth related messages
message LoginRequest {
  string email = 1;
  string password = 2;
}

message LoginResponse {
  string token = 1;
  string user_id = 2;
}

message RegisterRequest {
  string email = 1;
  string password = 2;
  string name = 3;
}

message LogoutRequest {
  string token = 1;
}

message ResetPasswordRequest {
  string email = 1;
}

message ResetPasswordConfirmRequest {
  string token = 1;
  string new_password = 2;
}

// User related messages
message User {
  string id = 1;
  string name = 2;
  string email = 3;
  string created_at = 4;
  string updated_at = 5;
}

message GetProfileRequest {
  string user_id = 1;
}

message UpdateProfileRequest {
  string user_id = 1;
  string name = 2;
  string email = 3;
}

message DeleteProfileRequest {
  string user_id = 1;
}

message ListUsersRequest {
  string name_filter = 1;
  string email_filter = 2;
  int32 page = 3;
  int32 limit = 4;
}

message ListUsersResponse {
  repeated User users = 1;
  int32 total = 2;
  int32 page = 3;
  int32 limit = 4;
}

// The main service definition
service UserService {
  // Auth endpoints
  rpc Login(LoginRequest) returns (LoginResponse);
  rpc Register(RegisterRequest) returns (User);
  rpc Logout(LogoutRequest) returns (Empty);
  rpc ResetPassword(ResetPasswordRequest) returns (Empty);
  rpc ResetPasswordConfirm(ResetPasswordConfirmRequest) returns (Empty);
  
  // User management
  rpc GetProfile(GetProfileRequest) returns (User);
  rpc UpdateProfile(UpdateProfileRequest) returns (User);
  rpc DeleteProfile(DeleteProfileRequest) returns (Empty);
  rpc ListUsers(ListUsersRequest) returns (ListUsersResponse);
}
```

สร้าง Go code จากไฟล์ Protocol Buffer:

```powershell
protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative api/proto/user_service.proto
```

## 5. สร้าง Models

สร้างไฟล์ `internal/models/user.go`:

```powershell
New-Item -Path "internal/models/user.go" -ItemType File
```

เพิ่มเนื้อหาต่อไปนี้:

```go
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
```

## 6. การจัดการฐานข้อมูล MongoDB

สร้างไฟล์ `internal/db/mongodb.go`:

```powershell
New-Item -Path "internal/db/mongodb.go" -ItemType File
```

เพิ่มเนื้อหาต่อไปนี้:

```go
package db

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MongoDB represents the MongoDB client connection
type MongoDB struct {
	Client    *mongo.Client
	Database  *mongo.Database
	UsersCol  *mongo.Collection
	TokensCol *mongo.Collection
	ResetCol  *mongo.Collection
}

// NewMongoDB creates a new MongoDB connection
func NewMongoDB(uri, dbName string) (*MongoDB, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}

	// Check the connection
	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, err
	}

	database := client.Database(dbName)

	return &MongoDB{
		Client:    client,
		Database:  database,
		UsersCol:  database.Collection("users"),
		TokensCol: database.Collection("blacklisted_tokens"),
		ResetCol:  database.Collection("reset_tokens"),
	}, nil
}

// Close closes the MongoDB connection
func (m *MongoDB) Close() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := m.Client.Disconnect(ctx); err != nil {
		log.Printf("❌ Error disconnecting from MongoDB: %v", err)
	} else {
		log.Println("🔌 Disconnected from MongoDB")
	}
}
```

## 7. สร้าง JWT Authentication

สร้างไฟล์ `internal/auth/jwt.go`:

```powershell
New-Item -Path "internal/auth/jwt.go" -ItemType File
```

เพิ่มเนื้อหาต่อไปนี้:

```go
package auth

import (
    "errors"
    "time"

    "github.com/golang-jwt/jwt/v5"
)

// JWTManager handles JWT operations
type JWTManager struct {
    secretKey     []byte
    tokenDuration time.Duration
}

// UserClaims represents the claims in JWT token
type UserClaims struct {
    UserID string `json:"user_id"`
    Email  string `json:"email"`
    jwt.RegisteredClaims
}

// NewJWTManager creates a new JWT manager
func NewJWTManager(secretKey string, tokenDuration time.Duration) *JWTManager {
    return &JWTManager{
        secretKey:     []byte(secretKey),
        tokenDuration: tokenDuration,
    }
}

// Generate creates a new token for a specific user
func (manager *JWTManager) Generate(userID, email string) (string, error) {
    claims := UserClaims{
        UserID: userID,
        Email:  email,
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(time.Now().Add(manager.tokenDuration)),
            IssuedAt:  jwt.NewNumericDate(time.Now()),
        },
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString(manager.secretKey)
}

// Verify validates the token string and returns the claims
func (manager *JWTManager) Verify(tokenStr string) (*UserClaims, error) {
    token, err := jwt.ParseWithClaims(
        tokenStr,
        &UserClaims{},
        func(token *jwt.Token) (interface{}, error) {
            _, ok := token.Method.(*jwt.SigningMethodHMAC)
            if !ok {
                return nil, errors.New("unexpected token signing method")
            }
            return manager.secretKey, nil
        },
    )

    if err != nil {
        return nil, err
    }

    claims, ok := token.Claims.(*UserClaims)
    if !ok {
        return nil, errors.New("invalid token claims")
    }

    return claims, nil
}
```

## 8. สร้าง Middleware

### 8.1 สร้าง AuthInterceptor

สร้างไฟล์ `internal/middleware/auth_interceptor.go`:

```powershell
New-Item -Path "internal/middleware/auth_interceptor.go" -ItemType File
```

เพิ่มเนื้อหาต่อไปนี้:

```go
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
	) (interface{}, error) {
		// Skip auth for login and register methods
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
```

### 8.2 สร้าง RateLimiter

สร้างไฟล์ `internal/middleware/rate_limiter.go`:

```powershell
New-Item -Path "internal/middleware/rate_limiter.go" -ItemType File
```

เพิ่มเนื้อหาต่อไปนี้:

```go
package middleware

import (
	"context"
	"sync"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// RateLimiter implements a simple in-memory rate limiter
type RateLimiter struct {
	attempts map[string][]time.Time
	mu       sync.Mutex
	limit    int
	window   time.Duration
}

// NewRateLimiter creates a new rate limiter
func NewRateLimiter(limit int, window time.Duration) *RateLimiter {
	return &RateLimiter{
		attempts: make(map[string][]time.Time),
		limit:    limit,
		window:   window,
	}
}

// IsAllowed checks if a request is allowed based on the key
func (r *RateLimiter) IsAllowed(key string) bool {
	r.mu.Lock()
	defer r.mu.Unlock()

	now := time.Now()
	windowStart := now.Add(-r.window)

	// Clean up old attempts
	var validAttempts []time.Time
	for _, t := range r.attempts[key] {
		if t.After(windowStart) {
			validAttempts = append(validAttempts, t)
		}
	}

	r.attempts[key] = validAttempts

	// Check if we're over the limit
	if len(validAttempts) >= r.limit {
		return false
	}

	// Record this attempt
	r.attempts[key] = append(r.attempts[key], now)
	return true
}

// RateLimitInterceptor is a gRPC interceptor for rate limiting
func RateLimitInterceptor(limiter *RateLimiter) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		// Only rate limit the login method
		if info.FullMethod == "/proto.UserService/Login" {
			// Use client IP or some identifier from context as the key
			key := "login" // In production, use client IP or user identifier

			if !limiter.IsAllowed(key) {
				return nil, status.Errorf(codes.ResourceExhausted, "rate limit exceeded")
			}
		}

		return handler(ctx, req)
	}
}
```

## 9. พัฒนา Service Logic

### 9.1 สร้างไฟล์ Configuration

สร้างไฟล์ `internal/config/config.go`:

```powershell
New-Item -Path "internal/config/config.go" -ItemType File
```

เพิ่มเนื้อหาต่อไปนี้:

```go
package config

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

// Config holds application configuration
type Config struct {
	Port            string
	MongoURI        string
	DBName          string
	JWTSecretKey    string
	TokenDuration   time.Duration
	RateLimit       int
	RateLimitWindow time.Duration
}

// LoadConfig loads configuration from environment variables
func LoadConfig() *Config {
	// Load .env file if it exists
	if err := godotenv.Load(); err != nil {
		log.Println("⚠️ No .env file found or error loading it, using environment variables or defaults")
	} else {
		log.Println("📄 Successfully loaded .env file")
	}

	config := &Config{
		Port:            getEnv("PORT", ":50051"),
		MongoURI:        getEnv("MONGO_URI", "mongodb://localhost:27017"),
		DBName:          getEnv("DB_NAME", "usermanagement"),
		JWTSecretKey:    getEnv("JWT_SECRET_KEY", "your-secret-key"),
		TokenDuration:   getDurationEnv("TOKEN_DURATION", 24*time.Hour),
		RateLimit:       getIntEnv("RATE_LIMIT", 5),
		RateLimitWindow: getDurationEnv("RATE_LIMIT_WINDOW", time.Minute),
	}

	return config
}

// Helper functions
func getEnv(key, defaultValue string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		return defaultValue
	}
	return value
}

func getIntEnv(key string, defaultValue int) int {
	value, exists := os.LookupEnv(key)
	if !exists {
		return defaultValue
	}

	intValue, err := strconv.Atoi(value)
	if err != nil {
		return defaultValue
	}
	return intValue
}

func getDurationEnv(key string, defaultValue time.Duration) time.Duration {
	value, exists := os.LookupEnv(key)
	if !exists {
		return defaultValue
	}

	duration, err := time.ParseDuration(value)
	if err != nil {
		return defaultValue
	}
	return duration
}
```

### 9.2 สร้าง User Service

สร้างไฟล์ `internal/services/user_service.go`:

```powershell
New-Item -Path "internal/services/user_service.go" -ItemType File
```

เพิ่มเนื้อหาดังนี้ (นี่เป็นไฟล์ที่ยาวที่สุด เนื่องจากมีการจัดการ business logic ทั้งหมด):

```go
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

	// Users can only view their own profile unless they are admins (simplified version)
	if userID != req.UserId {
		return nil, status.Error(codes.PermissionDenied, "cannot view other user's profile")
	}

	objectID, err := primitive.ObjectIDFromHex(req.UserId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid user ID")
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

	// Return user without password
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
		var existingUser models.User
		err = s.db.UsersCol.FindOne(ctx, bson.M{
			"email":      req.Email,
			"_id":        bson.M{"$ne": objectID},
			"is_deleted": false,
		}).Decode(&existingUser)

		if err != nil && err != mongo.ErrNoDocuments {
			return nil, status.Error(codes.Internal, "database error")
		}

		if err != mongo.ErrNoDocuments {
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
	var updatedUser models.User
	if err := result.Decode(&updatedUser); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, status.Error(codes.NotFound, "user not found")
		}
		return nil, status.Error(codes.Internal, "failed to update user")
	}

	// Return updated user
	return &proto.User{
		Id:        updatedUser.ID.Hex(),
		Name:      updatedUser.Name,
		Email:     updatedUser.Email,
		CreatedAt: updatedUser.CreatedAt.Format(time.RFC3339),
		UpdatedAt: updatedUser.UpdatedAt.Format(time.RFC3339),
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

	// Soft delete by updating is_deleted flag
	result, err := s.db.UsersCol.UpdateOne(
		ctx,
		bson.M{"_id": objectID, "is_deleted": false},
		bson.M{"$set": bson.M{"is_deleted": true, "updated_at": time.Now()}},
	)

	if err != nil {
		return nil, status.Error(codes.Internal, "failed to delete user")
	}

	if result.ModifiedCount == 0 {
		return nil, status.Error(codes.NotFound, "user not found")
	}

	return &proto.Empty{}, nil
}

// ListUsers retrieves a list of users with filtering and pagination
func (s *UserService) ListUsers(ctx context.Context, req *proto.ListUsersRequest) (*proto.ListUsersResponse, error) {
	// Set default pagination values if not provided
	page := int64(1)
	if req.Page > 0 {
		page = int64(req.Page)
	}

	limit := int64(10)
	if req.Limit > 0 && req.Limit <= 100 {
		limit = int64(req.Limit)
	}

	// Create filter
	filter := bson.M{"is_deleted": false}

	// Add name filter if provided
	if req.NameFilter != "" {
		filter["name"] = bson.M{"$regex": req.NameFilter, "$options": "i"}
	}

	// Add email filter if provided
	if req.EmailFilter != "" {
		filter["email"] = bson.M{"$regex": req.EmailFilter, "$options": "i"}
	}

	// Count total matching users
	total, err := s.db.UsersCol.CountDocuments(ctx, filter)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to count users")
	}

	// Skip calculation for pagination
	skip := (page - 1) * limit

	// Query with pagination
	cursor, err := s.db.UsersCol.Find(
		ctx,
		filter,
		options.Find().
			SetSort(bson.M{"created_at": -1}).
			SetSkip(skip).
			SetLimit(limit),
	)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to query users")
	}
	defer cursor.Close(ctx)

	// Process results
	var users []*proto.User
	for cursor.Next(ctx) {
		var user models.User
		if err := cursor.Decode(&user); err != nil {
			return nil, status.Error(codes.Internal, "failed to decode user")
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
	if req.Email == "" {
		return nil, status.Error(codes.InvalidArgument, "email is required")
	}

	// Find user by email
	var user models.User
	err := s.db.UsersCol.FindOne(ctx, bson.M{"email": req.Email, "is_deleted": false}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			// Return success to prevent email enumeration attacks
			return &proto.Empty{}, nil
		}
		return nil, status.Error(codes.Internal, "database error")
	}

	// Generate random token
	token, err := generateRandomToken(32)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to generate token")
	}

	// Store token with expiration time
	resetToken := models.PasswordReset{
		ID:        primitive.NewObjectID(),
		UserID:    user.ID,
		Token:     token,
		CreatedAt: time.Now(),
		ExpiresAt: time.Now().Add(24 * time.Hour), // Token expires in 24 hours
		Used:      false,
	}

	// Insert reset token
	_, err = s.db.ResetCol.InsertOne(ctx, resetToken)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to create reset token")
	}

	// In a real application, send an email with the reset token
	// For this example, we'll just return success
	// fmt.Printf("Password reset token for %s: %s\n", user.Email, token)

	return &proto.Empty{}, nil
}

// ResetPasswordConfirm completes the password reset process
func (s *UserService) ResetPasswordConfirm(ctx context.Context, req *proto.ResetPasswordConfirmRequest) (*proto.Empty, error) {
	if req.Token == "" || req.NewPassword == "" {
		return nil, status.Error(codes.InvalidArgument, "token and new password are required")
	}

	// Validate password strength
	if len(req.NewPassword) < 8 {
		return nil, status.Error(codes.InvalidArgument, "password must be at least 8 characters")
	}
	hasLetter := regexp.MustCompile(`[a-zA-Z]`).MatchString(req.NewPassword)
	hasNumber := regexp.MustCompile(`[0-9]`).MatchString(req.NewPassword)
	if !hasLetter || !hasNumber {
		return nil, status.Error(codes.InvalidArgument, "password must contain both letters and numbers")
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

	// Hash the new password
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to hash password")
	}

	// Update the user's password
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
	userID := ctx.Value("user_id")
	if userID == nil {
		return "", status.Error(codes.Unauthenticated, "user ID not found in context")
	}

	userIDStr, ok := userID.(string)
	if !ok {
		return "", status.Error(codes.Internal, "invalid user ID format in context")
	}

	return userIDStr, nil
}

// Helper functions
func generateRandomToken(length int) (string, error) {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}
```

## 10. สร้าง Server

สร้างไฟล์ `cmd/server/main.go`:

```powershell
New-Item -Path "cmd/server/main.go" -ItemType File
```

เพิ่มเนื้อหาต่อไปนี้:

```go
package main

import (
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/yourusername/api/proto"
	"github.com/yourusername/internal/auth"
	"github.com/yourusername/internal/config"
	"github.com/yourusername/internal/db"
	"github.com/yourusername/internal/middleware"
	"github.com/yourusername/internal/services"
)

func main() {
	// Load configuration
	cfg := config.LoadConfig()
	// Initialize MongoDB
	mongoDB, err := db.NewMongoDB(cfg.MongoURI, cfg.DBName)
	if err != nil {
		log.Fatalf("❌ Failed to connect to MongoDB: %v", err)
	}
	defer mongoDB.Close()

	log.Println("🍃 Connected to MongoDB successfully")

	// Initialize JWT manager
	jwtManager := auth.NewJWTManager(cfg.JWTSecretKey, cfg.TokenDuration)

	// Create rate limiter
	rateLimiter := middleware.NewRateLimiter(cfg.RateLimit, cfg.RateLimitWindow)

	// Create auth interceptor
	authInterceptor := middleware.NewAuthInterceptor(jwtManager, mongoDB)

	// Create gRPC server with middlewares
	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			middleware.RateLimitInterceptor(rateLimiter),
			authInterceptor.Unary(),
		),
	) // Create and register user service
	userService := services.NewUserService(mongoDB, jwtManager)
	proto.RegisterUserServiceServer(grpcServer, userService)

	// Enable reflection for tools like grpcurl
	reflection.Register(grpcServer)

	// Start listening
	listener, err := net.Listen("tcp", cfg.Port)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	log.Printf("🚀 Server started on port %s", cfg.Port)
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("❌ Failed to serve: %v", err)
	}
}
```

## 11. รัน และ ทดสอบโปรเจค

### 11.1 สร้างไฟล์ .env

สร้างไฟล์ `.env` ในรูทไดเรกทอรี:

```powershell
New-Item -Path ".env" -ItemType File
```

เพิ่มเนื้อหาต่อไปนี้:

```
PORT=:50051
MONGO_URI=mongodb://localhost:27017
DB_NAME=usermanagement
JWT_SECRET_KEY=your-secure-secret-key-should-be-long-and-random
TOKEN_DURATION=24h
RATE_LIMIT=5
RATE_LIMIT_WINDOW=1m
```

### 11.2 Build และรันโปรเจค

Build โปรเจค:

```powershell
go build -o user-service.exe ./cmd/server
```

รันโปรเจค:

```powershell
./user-service.exe
```

### 11.3 ทดสอบ API ด้วย gRPCurl

#### ลงทะเบียนผู้ใช้ใหม่:

```powershell
grpcurl -d '{"name": "Test User", "email": "test@example.com", "password": "Password123"}' -plaintext localhost:50051 proto.UserService/Register
```

#### เข้าสู่ระบบ:

```powershell
grpcurl -d '{"email": "test@example.com", "password": "Password123"}' -plaintext localhost:50051 proto.UserService/Login
```

เก็บค่า token จากผลลัพธ์ไว้ใช้ในคำสั่งถัดไป

#### ดูข้อมูลโปรไฟล์:

```powershell
grpcurl -H "Authorization: Bearer YOUR_JWT_TOKEN" -d '{"user_id": "USER_ID"}' -plaintext localhost:50051 proto.UserService/GetProfile
```

#### แก้ไขโปรไฟล์:

```powershell
grpcurl -H "Authorization: Bearer YOUR_JWT_TOKEN" -d '{"user_id": "USER_ID", "name": "Updated Name"}' -plaintext localhost:50051 proto.UserService/UpdateProfile
```

#### ดูรายชื่อผู้ใช้:

```powershell
grpcurl -H "Authorization: Bearer YOUR_JWT_TOKEN" -d '{"page": 1, "limit": 10}' -plaintext localhost:50051 proto.UserService/ListUsers
```

#### ออกจากระบบ:

```powershell
grpcurl -H "Authorization: Bearer YOUR_JWT_TOKEN" -d '{"token": "YOUR_JWT_TOKEN"}' -plaintext localhost:50051 proto.UserService/Logout
```

## 12. การตั้งค่า Docker

### 12.1 สร้าง Dockerfile

สร้างไฟล์ `Dockerfile` ในรูทไดเรกทอรี:

```powershell
New-Item -Path "Dockerfile" -ItemType File
```

เพิ่มเนื้อหาต่อไปนี้:

```dockerfile
FROM golang:latest AS builder

WORKDIR /app

# Copy go.mod and go.sum
COPY go.mod ./
COPY go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build with CGO disabled for alpine compatibility
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o user-service ./cmd/server

# Use a smaller image for the final container
FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /app

# Copy the binary from the builder stage
COPY --from=builder /app/user-service .

# Expose the port
EXPOSE 50051

# Run the service
CMD ["./user-service"]
```

### 12.2 สร้าง docker-compose.yml

สร้างไฟล์ `docker-compose.yml` ในรูทไดเรกทอรี:

```powershell
New-Item -Path "docker-compose.yml" -ItemType File
```

เพิ่มเนื้อหาต่อไปนี้:

```yaml
version: '3.8'

services:
  mongo:
    image: mongo:latest
    container_name: gRPC-mongo
    ports:
      - "27017:27017"
    volumes:
      - mongo_data:/data/db
    environment:
      - MONGO_INITDB_DATABASE=usermanagement
    networks:
      - app-network
    restart: always
    healthcheck:
      test: ["CMD", "mongo", "--eval", "db.adminCommand('ping')"]
      interval: 10s
      timeout: 5s
      retries: 5

  user-service:
    container_name: gRPC-user-service
    build:
      context: .
      dockerfile: Dockerfile
    ports:
      - "50051:50051"
    environment:
      - PORT=:50051
      - MONGO_URI=mongodb://mongo:27017
      - DB_NAME=usermanagement
      - JWT_SECRET_KEY=your-secret-key
      - TOKEN_DURATION=24h
      - RATE_LIMIT=5
      - RATE_LIMIT_WINDOW=1m
    depends_on:
      - mongo
    networks:
      - app-network
    restart: always
    logging:
      driver: "json-file"
      options:
        max-size: "10m"
        max-file: "3"

networks:
  app-network:
    driver: bridge

volumes:
  mongo_data:
```

### 12.3 รันโปรเจคด้วย Docker Compose

เริ่มต้นทุกบริการ:

```powershell
docker-compose up -d
```

ดู logs:

```powershell
docker-compose logs -f
```

หยุดทุกบริการ:

```powershell
docker-compose down
```

## สรุป

ในคู่มือนี้ เราได้พัฒนา User Management API ที่สมบูรณ์ตั้งแต่เริ่มต้น โดยใช้ gRPC, MongoDB, และ JWT authentication ใน Golang โปรเจคนี้ประกอบด้วยคุณสมบัติต่างๆ ตามที่กำหนดไว้ในข้อกำหนด รวมถึง:

1. การยืนยันตัวตนและการจัดการผู้ใช้
2. การจัดการโปรไฟล์ผู้ใช้
3. การรองรับการค้นหาและแบ่งหน้า
4. การรองรับ rate limiting
5. การรีเซ็ตรหัสผ่าน
6. การจัดการ JWT tokens

โปรเจคนี้มีการแบ่งโครงสร้างที่ชัดเจนตามหลัก clean architecture ทำให้ง่ายต่อการบำรุงรักษาและขยายฟีเจอร์ในอนาคต
