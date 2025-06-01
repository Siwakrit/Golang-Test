package utils

import "time"

// รายการค่าคงที่ที่ใช้ทั่วทั้งแอปพลิเคชัน

const (
	// TokenValidDuration คือระยะเวลาที่ password reset token ยังใช้งานได้
	TokenValidDuration = 24 * time.Hour

	// MinPasswordLength คือความยาวขั้นต่ำของรหัสผ่านที่ยอมรับได้
	MinPasswordLength = 8

	// DefaultPageSize คือจำนวนรายการต่อหน้าเริ่มต้นสำหรับการแบ่งหน้า
	DefaultPageSize = 10

	// MaxPageSize คือจำนวนรายการต่อหน้าสูงสุดที่ยอมรับได้
	MaxPageSize = 100
)
