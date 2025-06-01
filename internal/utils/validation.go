package utils

import (
	"regexp"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// ฟังก์ชันตรวจสอบความถูกต้องที่ใช้ทั่วทั้งแอปพลิเคชัน

// ValidateEmail ตรวจสอบว่าอีเมลมีรูปแบบที่ถูกต้องหรือไม่
func ValidateEmail(email string) error {
	if email == "" {
		return status.Error(codes.InvalidArgument, "email is required")
	}

	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(email) {
		return status.Error(codes.InvalidArgument, "invalid email format")
	}

	return nil
}

// ValidatePassword ตรวจสอบว่ารหัสผ่านตรงตามข้อกำหนดความปลอดภัยหรือไม่
func ValidatePassword(password string) error {
	if password == "" {
		return status.Error(codes.InvalidArgument, "password is required")
	}

	if len(password) < MinPasswordLength {
		return status.Error(codes.InvalidArgument, "password must be at least 8 characters")
	}

	hasLetter := regexp.MustCompile(`[a-zA-Z]`).MatchString(password)
	hasNumber := regexp.MustCompile(`[0-9]`).MatchString(password)

	if !hasLetter || !hasNumber {
		return status.Error(codes.InvalidArgument, "password must contain both letters and numbers")
	}

	return nil
}

// ValidateName ตรวจสอบว่าชื่อมีความยาวที่เหมาะสมและไม่เป็นค่าว่าง
func ValidateName(name string) error {
	if name == "" {
		return status.Error(codes.InvalidArgument, "name is required")
	}

	return nil
}
