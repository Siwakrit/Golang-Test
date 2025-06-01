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
	"github.com/yourusername/internal/gateway"
	"github.com/yourusername/internal/middleware"
	"github.com/yourusername/internal/services/user"
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
	userService := user.NewService(mongoDB, jwtManager)
	proto.RegisterUserServiceServer(grpcServer, userService)

	// Enable reflection for tools like grpcurl
	reflection.Register(grpcServer)

	// Start gRPC server in a goroutine
	go func() {
		listener, err := net.Listen("tcp", cfg.Port)
		if err != nil {
			log.Fatalf("Failed to listen: %v", err)
		}
		log.Printf("🚀 gRPC server started on port %s", cfg.Port)
		if err := grpcServer.Serve(listener); err != nil {
			log.Fatalf("❌ Failed to serve: %v", err)
		}
	}()

	// Start HTTP REST gateway
	httpAddr := ":8080" // Default REST port
	if err := gateway.RunGatewayServer(cfg.Port, httpAddr); err != nil {
		log.Fatalf("❌ Failed to start gRPC-Gateway: %v", err)
	}
}
