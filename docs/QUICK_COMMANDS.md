# 🚀 Quick Commands for User Management API Project

This file contains all the necessary commands to build and work with the User Management API project, organized by workflow.

## 1. 🛠️ Installation of Required Tools

### 🔍 Verify Go Version
```powershell
go version
```

### 🍃 Verify MongoDB
```powershell
mongod --version
```

### 📦 Install Protocol Buffer Plugins for Go
```powershell
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@latest
```

## 2. 📂 Project Structure Creation

### 📁 Create Main and Sub Directories
```powershell
# Create main directory
mkdir -p Golang-Test
cd Golang-Test

# Create project structure
mkdir -p api/proto api/third_party/google/api api/third_party/google/protobuf
mkdir -p cmd/server
mkdir -p docs
mkdir -p internal/{auth,config,db,gateway,middleware,models,services/user,utils}
mkdir -p test/api
```

## 3. 📦 Go Module Setup

```powershell
# Initialize Go module
go mod init github.com/yourusername

# Install dependencies
go get -u github.com/golang-jwt/jwt/v5
go get -u go.mongodb.org/mongo-driver/mongo
go get -u golang.org/x/crypto/bcrypt
go get -u google.golang.org/grpc
go get -u google.golang.org/protobuf
go get -u github.com/joho/godotenv
go get -u github.com/grpc-ecosystem/grpc-gateway/v2
```

## 4. 📋 Create Proto File and Generate Code

```powershell
# Create Protocol Buffer definition
New-Item -Path "api/proto/user_service.proto" -ItemType File

# Generate Go code from Proto file
protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative api/proto/user_service.proto

# Generate REST gateway code (optional)
protoc -I . --grpc-gateway_out . --grpc-gateway_opt logtostderr=true --grpc-gateway_opt paths=source_relative api/proto/user_service.proto
```

## 5. 💻 Create Code Files

```powershell
# Create model files
New-Item -Path "internal/models/user.go" -ItemType File

# Create database management files
New-Item -Path "internal/db/mongodb.go" -ItemType File

# Create JWT management files
New-Item -Path "internal/auth/jwt.go" -ItemType File
New-Item -Path "internal/auth/context.go" -ItemType File

# Create middleware files
New-Item -Path "internal/middleware/auth_interceptor.go" -ItemType File
New-Item -Path "internal/middleware/rate_limiter.go" -ItemType File

# Create configuration files
New-Item -Path "internal/config/config.go" -ItemType File

# Create gateway files
New-Item -Path "internal/gateway/gateway.go" -ItemType File

# Create service logic files (organized by domain)
New-Item -Path "internal/services/user/service.go" -ItemType File
New-Item -Path "internal/services/user/auth.go" -ItemType File
New-Item -Path "internal/services/user/profile.go" -ItemType File
New-Item -Path "internal/services/user/registration.go" -ItemType File
New-Item -Path "internal/services/user/list.go" -ItemType File
New-Item -Path "internal/services/user/password.go" -ItemType File
New-Item -Path "internal/services/user/utils.go" -ItemType File

# Create utility files
New-Item -Path "internal/utils/constants.go" -ItemType File
New-Item -Path "internal/utils/token.go" -ItemType File
New-Item -Path "internal/utils/validation.go" -ItemType File

# Create server main file
New-Item -Path "cmd/server/main.go" -ItemType File

# Create Docker files
New-Item -Path "Dockerfile" -ItemType File
New-Item -Path "docker-compose.yml" -ItemType File
```

## 6. 🚀 Build and Run

### 🔨 Build the Application
```powershell
go build -o server.exe ./cmd/server
```

### ▶️ Run the Application
```powershell
.\server.exe
```

### 🐳 Build and Run with Docker
```powershell
docker compose up --build
```

## 7. 🧪 Testing the API

### 📋 List Available gRPC Services
```powershell
grpcurl -plaintext localhost:50052 list
```

### 🔍 List Methods of a Service
```powershell
grpcurl -plaintext localhost:50052 list proto.UserService
```

### 📡 Call a gRPC Method (Example: Register)
```powershell
grpcurl -plaintext -d '{
  "username": "testuser",
  "email": "test@example.com",
  "password": "Password123!"
}' localhost:50052 proto.UserService/Register
```

### 🔑 Login Example
```powershell
grpcurl -plaintext -d '{
  "email": "test@example.com",
  "password": "Password123!"
}' localhost:50052 proto.UserService/Login
```

## 8. 🐳 Docker Commands

### 🏗️ Build Docker Image
```powershell
docker build -t user-service .
```

### 🚢 Run Docker Container
```powershell
docker run -p 50052:50052 -p 8081:8081 user-service
```

### ⏹️ Stop All Containers
```powershell
docker compose down
```

### 🧹 Stop and Remove All Containers, Images, and Volumes
```powershell
docker compose down --rmi all -v
```

## 9. ⚙️ Environment Variables

Create a `.env` file in the project root with the following content:
```
PORT=:50052
GRPC_PORT=50052
HTTP_PORT=8081
MONGO_URI=mongodb://localhost:27017
DB_NAME=usermanagement
JWT_SECRET_KEY=your-secret-key-change-in-production
TOKEN_DURATION=24h
RATE_LIMIT=5
RATE_LIMIT_WINDOW=1m
```

## 10. 📝 Useful Git Commands

### 🌱 Initialize Git Repository
```powershell
git init
```

### ➕ Add All Files
```powershell
git add .
```

### 💾 Create a Commit
```powershell
git commit -m "Initial commit"
```

### 🔗 Add Remote Repository
```powershell
git remote add origin <your-repository-url>
```

### 🚀 Push Changes
```powershell
git push -u origin main
```
