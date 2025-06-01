# คำสั่งลัดสำหรับการสร้างโปรเจค User Management API

ไฟล์นี้รวบรวมคำสั่งทั้งหมดที่จำเป็นในการสร้างโปรเจค User Management API โดยเรียงตามลำดับการทำงาน

## 1. การติดตั้งเครื่องมือที่จำเป็น

### ตรวจสอบเวอร์ชัน Go
```powershell
go version
```

### ตรวจสอบ MongoDB
```powershell
mongod --version
```

### ติดตั้ง Protocol Buffer plugins สำหรับ Go
```powershell
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
```

## 2. สร้างโครงสร้างโปรเจค

### สร้างไดเรกทอรีหลักและย่อย
```powershell
# สร้างไดเรกทอรีหลัก
mkdir -p Golang-Test
cd Golang-Test

# สร้างโครงสร้างโปรเจค
mkdir -p api/proto
mkdir -p cmd/server
mkdir -p docs
mkdir -p internal/{auth,config,db,middleware,models,services}
mkdir -p test
```

## 3. ตั้งค่า Go Module

```powershell
# เริ่มต้น Go module
go mod init github.com/yourusername

# ติดตั้ง dependencies
go get -u github.com/golang-jwt/jwt/v5
go get -u go.mongodb.org/mongo-driver/mongo
go get -u golang.org/x/crypto/bcrypt
go get -u google.golang.org/grpc
go get -u google.golang.org/protobuf
go get -u github.com/joho/godotenv
```

## 4. สร้างไฟล์ Proto และ Generate code

```powershell
# สร้าง Protocol Buffer definition
New-Item -Path "api/proto/user_service.proto" -ItemType File

# Generate Go code จาก Proto file
protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative api/proto/user_service.proto
```

## 5. สร้างไฟล์โค้ด

```powershell
# สร้างไฟล์โมเดล
New-Item -Path "internal/models/user.go" -ItemType File

# สร้างไฟล์จัดการฐานข้อมูล
New-Item -Path "internal/db/mongodb.go" -ItemType File

# สร้างไฟล์การจัดการ JWT
New-Item -Path "internal/auth/jwt.go" -ItemType File

# สร้างไฟล์ Middleware
New-Item -Path "internal/middleware/auth_interceptor.go" -ItemType File
New-Item -Path "internal/middleware/rate_limiter.go" -ItemType File

# สร้างไฟล์ Configuration
New-Item -Path "internal/config/config.go" -ItemType File

# สร้างไฟล์ Service Logic
New-Item -Path "internal/services/user_service.go" -ItemType File

# สร้างไฟล์หลักของเซิร์ฟเวอร์
New-Item -Path "cmd/server/main.go" -ItemType File
```

## 6. สร้างไฟล์ Environment และ Docker

```powershell
# สร้างไฟล์ Environment
New-Item -Path ".env" -ItemType File

# สร้าง Dockerfile
New-Item -Path "Dockerfile" -ItemType File

# สร้างไฟล์ docker-compose
New-Item -Path "docker-compose.yml" -ItemType File
```

## 7. Build และรันโปรเจค

```powershell
# Build โปรเจค
go build -o user-service.exe ./cmd/server

# รันโปรเจค
./user-service.exe
```

## 8. ทดสอบ API ด้วย gRPCurl

```powershell
# ลงทะเบียนผู้ใช้ใหม่
grpcurl -d '{"name": "Test User", "email": "test@example.com", "password": "Password123"}' -plaintext localhost:50051 proto.UserService/Register

# เข้าสู่ระบบ
grpcurl -d '{"email": "test@example.com", "password": "Password123"}' -plaintext localhost:50051 proto.UserService/Login

# ดูโปรไฟล์ (ใส่ token และ user_id ที่ได้จากการล็อกอิน)
grpcurl -H "Authorization: Bearer YOUR_JWT_TOKEN" -d '{"user_id": "USER_ID"}' -plaintext localhost:50051 proto.UserService/GetProfile

# แก้ไขโปรไฟล์
grpcurl -H "Authorization: Bearer YOUR_JWT_TOKEN" -d '{"user_id": "USER_ID", "name": "Updated Name"}' -plaintext localhost:50051 proto.UserService/UpdateProfile

# ดูรายชื่อผู้ใช้
grpcurl -H "Authorization: Bearer YOUR_JWT_TOKEN" -d '{"page": 1, "limit": 10}' -plaintext localhost:50051 proto.UserService/ListUsers

# ออกจากระบบ
grpcurl -H "Authorization: Bearer YOUR_JWT_TOKEN" -d '{"token": "YOUR_JWT_TOKEN"}' -plaintext localhost:50051 proto.UserService/Logout
```

## 9. Docker Deployment

```powershell
# Build และรันด้วย Docker Compose
docker-compose up -d

# ดู logs
docker-compose logs -f

# หยุดการทำงานของ containers
docker-compose down
```
