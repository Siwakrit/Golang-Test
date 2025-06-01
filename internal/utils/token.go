package utils

import (
	"crypto/rand"
	"encoding/hex"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// GenerateSecureToken สร้าง token แบบสุ่มที่ปลอดภัยสำหรับการใช้งานเช่น reset password
// bytes คือจำนวนไบต์ของ token ที่จะสร้าง (แนะนำให้ใช้อย่างน้อย 32 ไบต์)
func GenerateSecureToken(bytes int) (string, error) {
	tokenBytes := make([]byte, bytes)
	_, err := rand.Read(tokenBytes)
	if err != nil {
		return "", status.Error(codes.Internal, "failed to generate secure token")
	}

	return hex.EncodeToString(tokenBytes), nil
}
