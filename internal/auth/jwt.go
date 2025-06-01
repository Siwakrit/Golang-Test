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
