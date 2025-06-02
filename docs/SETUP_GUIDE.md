# 🚀 User Management API Setup Guide (Golang + gRPC + MongoDB)

This guide will help you set up, build, and run the User Management API project using Go, gRPC, and MongoDB.

## 📋 Prerequisites

1. **Go** 🔧 (>= 1.23)
   - Download and install f### 🔧 Environment Considerations

Remember to set appropriate environment variables for production:
- 🔑 Use a strong JWT secret key
- ⚡ Configure appropriate rate limits
- 🔒 Set proper database credentials
- 🛡️ Consider implementing SSL/TLS for secure communication

## 🧰 Additional Commandss://golang.org/dl/](https://golang.org/dl/)
   - Verify installation:
     ```powershell
     go version
     ```

2. **Docker & Docker Compose** 🐳
   - Download and install from [https://www.docker.com/products/docker-desktop/](https://www.docker.com/products/docker-desktop/)
   - Verify installation:
     ```powershell
     docker --version
     docker compose version
     ```

3. **Protocol Buffers Compiler (protoc)** 📄
   - Download from [https://github.com/protocolbuffers/protobuf/releases](https://github.com/protocolbuffers/protobuf/releases)
   - Add the `bin` directory to your PATH
   - Install Go plugins:
     ```powershell
     go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
     go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
     go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@latest
     ```
   - Verify installation:
     ```powershell
     protoc --version
     ```

4. **gRPCurl (optional, for testing)** 🧪
   - Download from [https://github.com/fullstorydev/grpcurl/releases](https://github.com/fullstorydev/grpcurl/releases)
   - Add executable to your PATH
   - Verify installation:
     ```powershell
     grpcurl --version
     ```

## 📂 Project Structure

The project follows a clean architecture approach with domain-driven design principles:

```
Golang-Test/
 api/                        # API definitions
    package.go              # Package documentation
    proto/                  # Protocol Buffers definitions
       user_service.proto  # Service definitions
       *generated files*   # Auto-generated protobuf code
    third_party/            # Third-party proto files
 cmd/                        # Application entry points
    server/                 # Main server entrypoint
        main.go             # Server initialization
 internal/                   # Private application code
    auth/                   # Authentication components
    config/                 # Configuration management
    db/                     # Database connectivity
    gateway/                # REST API gateway
    middleware/             # gRPC interceptors
    models/                 # Data models
    services/               # Business logic by domain
       user/               # User service implementation
    utils/                  # Utility functions
 test/                       # Test files and collections
 docs/                       # Documentation
 Dockerfile                  # Container definition
 docker-compose.yml          # Multi-container setup
 go.mod / go.sum             # Go modules
 README.md                   # Project documentation
```

## 🚀 Quick Start with Docker (Recommended)

1. **Clone the repository** 📥 (if not already):
   ```powershell
   git clone <your-repo-url>
   cd Golang-Test
   ```

2. **Build and run with Docker Compose:** 🐳
   ```powershell
   docker compose up --build
   ```
   
   This will:
   - Start MongoDB container on port 27017
   - Build and start the gRPC service container
   - Expose ports:
     - 50051 for gRPC
     - 8080 for REST API (via gRPC Gateway)

3. **Verify the service is running:** ✅
   ```powershell
   grpcurl -plaintext localhost:50051 list
   ```
   
   You should see `proto.UserService` in the output.

## 🛠️ Manual Setup (For Development)

1. **Set up MongoDB:** 🍃
   - Install and start MongoDB locally, or
   - Use Docker: `docker run -d -p 27017:27017 --name mongodb mongo:latest`

2. **Configure environment:** ⚙️
   - Create `.env` file in project root:
     ```
     PORT=:50051
     MONGO_URI=mongodb://localhost:27017
     DB_NAME=usermanagement
     JWT_SECRET_KEY=your-secret-key
     TOKEN_DURATION=24h
     RATE_LIMIT=5
     RATE_LIMIT_WINDOW=1m
     ```

3. **Build and run:** 🚀
   ```powershell
   go build -o server.exe ./cmd/server
   .\server.exe
   ```

## 🔌 API Endpoints

### 🌐 gRPC Endpoints (Port 50051)

The service exposes these gRPC methods:

- **Authentication:** 🔐
  - `Register` - Create a new user account
  - `Login` - Authenticate and get JWT token
  - `Logout` - Invalidate JWT token

- **User Management:** 👤
  - `GetProfile` - Get user profile details
  - `UpdateProfile` - Update user profile information
  - `DeleteProfile` - Delete user account
  - `ListUsers` - List users with pagination (admin)

- **Password Management:** 🔑
  - `ResetPassword` - Request password reset
  - `ResetPasswordConfirm` - Complete password reset

### 🌍 REST Endpoints (Port 8080)

The gRPC-Gateway provides RESTful equivalents of all gRPC endpoints, for example:

- `POST /v1/auth/register` - Register user
- `POST /v1/auth/login` - Login user
- `GET /v1/users/profile` - Get user profile
- `PUT /v1/users/profile` - Update profile

> For complete API documentation, refer to the proto file: `api/proto/user_service.proto`

## 👨‍💻 Development Workflow

### 📝 Modifying the API

1. **Edit the proto file:**
   ```powershell
   # Edit api/proto/user_service.proto with your changes
   ```

2. **Generate updated code:**
   ```powershell
   # Generate Go code from proto file
   protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative api/proto/user_service.proto
   
   # Generate REST gateway code
   protoc -I . --grpc-gateway_out . --grpc-gateway_opt logtostderr=true --grpc-gateway_opt paths=source_relative api/proto/user_service.proto
   ```

3. **Implement the changes** in the corresponding service files.

### 🧪 Testing

- **Unit Tests:** ✅
  ```powershell
  go test ./internal/...
  ```

- **gRPC API Testing with gRPCurl:** 📡
  ```powershell
  # List all methods
  grpcurl -plaintext localhost:50051 list proto.UserService
  
  # Call specific method (example: Register)
  grpcurl -plaintext -d '{
    "username": "testuser",
    "email": "test@example.com",
    "password": "Password123!"
  }' localhost:50051 proto.UserService/Register
  ```

- **REST API Testing:** 🌐
  - Use Postman with the provided collection in `test/api/user-management-rest.postman_collection.json`
  - Or use curl:
    ```powershell
    Invoke-RestMethod -Uri "http://localhost:8080/v1/auth/register" -Method Post -ContentType "application/json" -Body '{"username":"testuser","email":"test@example.com","password":"Password123!"}'
    ```

## ⚠️ Common Issues and Solutions

### 🍃 MongoDB Connection Errors

- **Error:** Failed to connect to MongoDB
  - **Solution:** Ensure MongoDB is running, check the connection string, verify network access

### 🔐 JWT Authentication Issues

- **Error:** Invalid token
  - **Solution:** Check token expiration, verify the JWT_SECRET_KEY matches the one used to issue the token

### 🔌 gRPC Connection Failures

- **Error:** Cannot connect to gRPC server
  - **Solution:** Verify the server is running, check the port is correct and not in use, ensure no firewall blocking

### ⚡ Rate Limiting

- **Error:** Too many requests
  - **Solution:** The service has a default rate limit of 5 requests per minute. Adjust the `RATE_LIMIT` and `RATE_LIMIT_WINDOW` environment variables if needed.

## 🚢 Deployment

The service is Docker-ready and can be deployed to any container orchestration system:

### 🐳 Docker Standalone

```powershell
docker build -t user-service .
docker run -p 50051:50051 -p 8080:8080 user-service
```

### 🐙 Docker Compose (Production)

```powershell
docker compose up -d
```

### ⚙️ Environment Considerations

Remember to set appropriate environment variables for production:
- Use a strong JWT secret key
- Configure appropriate rate limits
- Set proper database credentials
- Consider implementing SSL/TLS for secure communication

## 🧰 Additional Commands

- **View server logs:** 📜
  ```powershell
  docker logs gRPC-user-service
  ```

- **Access MongoDB shell:** 🔍
  ```powershell
  docker exec -it gRPC-mongo mongosh
  ```

- **Stop all containers:** ⏹️
  ```powershell
  docker compose down
  ```

- **Clean restart (removes volumes):** 🔄
  ```powershell
  docker compose down -v
  docker compose up --build
  ```

---

For more detailed information, refer to the code comments and the main README.md file. 📚
