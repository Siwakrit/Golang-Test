# User Management API

A gRPC-based user management service with JWT authentication and MongoDB backend.

## Features

- User authentication with JWT
- User registration with password hashing
- User profile management (create, read, update, delete)
- Password reset functionality
- Rate limiting for login attempts
- Pagination and filtering for user listing

## Tech Stack

- **Protocol**: gRPC
- **Database**: MongoDB
- **Programming Language**: Golang
- **Authentication**: JWT

## Prerequisites

- Go 1.16+
- MongoDB
- Protocol Buffers compiler (`protoc`)

## Setup Instructions

1. **Clone the repository**

```bash
git clone https://github.com/yourusername/Golang-Test.git
cd Golang-Test
```

2. **Install dependencies**

```bash
go mod tidy
```

3. **Generate protobuf code**

```bash
protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative api/proto/user_service.proto
```

4. **Set up MongoDB**

Make sure MongoDB is running on your machine or update the connection string in the configuration.

5. **Environment variables**

You can configure the application using environment variables:

```powershell
$env:PORT=":50051"
$env:MONGO_URI="mongodb://localhost:27017"
$env:DB_NAME="usermanagement"
$env:JWT_SECRET_KEY="your-secure-secret-key"
$env:TOKEN_DURATION="24h"
$env:RATE_LIMIT="5"
$env:RATE_LIMIT_WINDOW="1m"
```

6. **Build and run**

```bash
go build -o user-service.exe ./cmd/server
./user-service.exe
```

## Testing with gRPCurl

You can use [gRPCurl](https://github.com/fullstorydev/grpcurl) to test the API:

1. **Register a new user**

```bash
grpcurl -d '{"name": "Test User", "email": "test@example.com", "password": "Password123"}' -plaintext localhost:50051 proto.UserService/Register
```

2. **Login**

```bash
grpcurl -d '{"email": "test@example.com", "password": "Password123"}' -plaintext localhost:50051 proto.UserService/Login
```

3. **Get profile**

```bash
grpcurl -H "Authorization: Bearer YOUR_JWT_TOKEN" -d '{"user_id": "USER_ID"}' -plaintext localhost:50051 proto.UserService/GetProfile
```

4. **Update profile**

```bash
grpcurl -H "Authorization: Bearer YOUR_JWT_TOKEN" -d '{"user_id": "USER_ID", "name": "Updated Name"}' -plaintext localhost:50051 proto.UserService/UpdateProfile
```

5. **List users**

```bash
grpcurl -H "Authorization: Bearer YOUR_JWT_TOKEN" -d '{"page": 1, "limit": 10}' -plaintext localhost:50051 proto.UserService/ListUsers
```

6. **Logout**

```bash
grpcurl -H "Authorization: Bearer YOUR_JWT_TOKEN" -d '{"token": "YOUR_JWT_TOKEN"}' -plaintext localhost:50051 proto.UserService/Logout
```

## API Documentation

### Authentication Services

- **Register**: Creates a new user account
- **Login**: Authenticates a user and returns a JWT token
- **Logout**: Invalidates a JWT token
- **ResetPassword**: Initiates a password reset process
- **ResetPasswordConfirm**: Completes the password reset process

### User Management Services

- **GetProfile**: Retrieves a user's profile
- **UpdateProfile**: Updates a user's profile information
- **DeleteProfile**: Soft-deletes a user profile
- **ListUsers**: Retrieves a list of users with filtering and pagination

## Architecture Overview

This service follows a clean architecture approach:

- **API Layer**: gRPC service definitions in Protocol Buffers
- **Service Layer**: Business logic implementation
- **Data Layer**: MongoDB database interaction
- **Auth**: JWT token generation, validation, and management
- **Middleware**: Interceptors for authentication and rate limiting

## Security Features

- Passwords are hashed using bcrypt
- JWT tokens for authentication
- Token blacklisting on logout
- Rate limiting for login attempts
- Input validation for all APIs

## Scaling Considerations

The service is designed to handle:
- 1,000 concurrent users
- ~100 requests per second
- Database capacity for ~100,000 user records
- Response times under 200ms for most operations