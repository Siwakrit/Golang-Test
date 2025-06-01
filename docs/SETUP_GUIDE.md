# User Management API Setup Guide (Golang + gRPC + MongoDB)

This guide will help you set up, build, and run the User Management API project using Go, gRPC, and MongoDB.

## Prerequisites

1. **Go** (>= 1.23)
   - Download and install from [https://golang.org/dl/](https://golang.org/dl/)
   - Verify installation:
     ```powershell
     go version
     ```
2. **Docker & Docker Compose**
   - Download and install from [https://www.docker.com/products/docker-desktop/](https://www.docker.com/products/docker-desktop/)
   - Verify installation:
     ```powershell
     docker --version
     docker compose version
     ```
3. **Protocol Buffers Compiler (protoc)**
   - Download from [https://github.com/protocolbuffers/protobuf/releases](https://github.com/protocolbuffers/protobuf/releases)
   - Add the `bin` directory to your PATH
   - Install Go plugins:
     ```powershell
     go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
     go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
     ```
   - Verify installation:
     ```powershell
     protoc --version
     ```
4. **gRPCurl (optional, for testing)**
   - Download from [https://github.com/fullstorydev/grpcurl/releases](https://github.com/fullstorydev/grpcurl/releases)
   - Add executable to your PATH
   - Verify installation:
     ```powershell
     grpcurl --version
     ```

## Project Structure

```
Golang-Test/
├── api/proto/                  # Protobuf definitions & generated code
├── cmd/server/                 # Main server entrypoint
├── internal/                   # Internal packages (auth, db, models, services, etc.)
├── test/                       # Test collections (e.g. Postman)
├── Dockerfile                  # Docker build file
├── docker-compose.yml          # Docker Compose for local dev
├── go.mod / go.sum             # Go modules
└── docs/SETUP_GUIDE.md         # This guide
```

## Quick Start (Recommended)

1. **Clone the repository** (if not already):
   ```powershell
   git clone <your-repo-url>
   cd Golang-Test
   ```
2. **Build and run with Docker Compose:**
   ```powershell
   docker compose up --build
   ```
   - This will start both MongoDB and the gRPC user service.
   - The gRPC server will be available at `localhost:50051`.
   - MongoDB will be available at `localhost:27017`.

3. **(Optional) Run locally without Docker:**
   - Ensure MongoDB is running locally (default port 27017)
   - Copy `.env.example` to `.env` and adjust as needed (or set environment variables manually)
   - Build and run:
     ```powershell
     go build -o server.exe ./cmd/server
     .\server.exe
     ```

## Protobuf Code Generation

If you modify `api/proto/user_service.proto`, regenerate Go code:
```powershell
protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative api/proto/user_service.proto
```

## Environment Variables

See `docker-compose.yml` for all environment variables. Key variables:
- `PORT` (default: :50051)
- `MONGO_URI` (default: mongodb://mongo:27017)
- `DB_NAME` (default: usermanagement)
- `JWT_SECRET_KEY`, `TOKEN_DURATION`, etc.

## Testing the API

- Use gRPCurl or Postman (see `test/api/user-management-grpc.postman_collection.json`)
- Example (gRPCurl):
  ```powershell
  grpcurl -plaintext localhost:50051 list
  ```

## Useful Commands

- Build Docker image:
  ```powershell
  docker build -t user-service .
  ```
- Stop all containers:
  ```powershell
  docker compose down
  ```

---

For more details, see code comments and each package's README (if available).
