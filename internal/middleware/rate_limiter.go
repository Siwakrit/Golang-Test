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
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) { // Only rate limit the login method
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
