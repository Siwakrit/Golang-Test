# User Management API - คู่มือภาษาไทย

API จัดการผู้ใช้ที่สมบูรณ์แบบที่สร้างด้วย gRPC, JWT authentication และ MongoDB เป็นฐานข้อมูล โปรเจคนี้แสดงให้เห็นถึงแนวทางการเขียนโปรแกรม Go สมัยใหม่ตามหลักการ Clean Architecture และนำเสนอความสามารถในการจัดการผู้ใช้อย่างครบถ้วน

## คุณลักษณะเด่น

- **ระบบการพิสูจน์ตัวตน**
  - ความปลอดภัยด้วย JWT authentication
  - การตรวจสอบและต่ออายุ token
  - ฟังก์ชัน logout พร้อมการยกเลิก token

- **การจัดการผู้ใช้**
  - การลงทะเบียนผู้ใช้พร้อมการตรวจสอบอีเมล
  - การจัดการโปรไฟล์ (สร้าง, อ่าน, อัปเดต, ลบ)
  - การเข้ารหัสรหัสผ่านด้วย bcrypt
  - ฟังก์ชันการรีเซ็ตรหัสผ่าน

- **คุณสมบัติด้านความปลอดภัย**
  - การจำกัดอัตราการเรียก API (ปรับแต่งได้, ค่าเริ่มต้น 5 ครั้ง/นาที)
  - การป้องกันการโจมตีแบบ brute force
  - การขึ้นบัญชีดำของ token
  - การตรวจสอบและทำความสะอาดข้อมูลนำเข้า

## เทคโนโลยีที่ใช้

- **Transport Layer**
  - gRPC - เฟรมเวิร์ค RPC ประสิทธิภาพสูง
  - Protocol Buffers - นิยามอินเตอร์เฟซที่เป็นกลางทางภาษา
  - gRPC Gateway - API แบบ RESTful JSON จากบริการ gRPC

- **Backend**
  - Go 1.23+ - ภาษาที่ทันสมัย, เร็ว, มีการตรวจสอบประเภทข้อมูลแบบ static
  - MongoDB - ฐานข้อมูล NoSQL แบบ document
  - JWT - JSON Web Tokens สำหรับการพิสูจน์ตัวตน

- **Development & Operations**
  - Docker - การสร้าง container
  - Docker Compose - การจัดการ multi-container

## การเริ่มต้นใช้งาน

### สิ่งที่ต้องติดตั้งก่อน

- **Go 1.23+** - [ดาวน์โหลดและติดตั้ง Go](https://golang.org/dl/)
- **MongoDB** - [ดาวน์โหลดและติดตั้ง MongoDB](https://www.mongodb.com/try/download/community)
  - ตรวจสอบว่าบริการ MongoDB กำลังทำงานอยู่ (บน Windows, ตรวจสอบใน Services app)
  - ตรวจสอบการติดตั้งด้วย `mongod --version`

- **Protocol Buffers Compiler**:
  - ดาวน์โหลดจาก [github.com/protocolbuffers/protobuf/releases](https://github.com/protocolbuffers/protobuf/releases)
  - เพิ่มลงใน PATH
  - ติดตั้ง Go plugins ที่จำเป็น:
    ```bash
    go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
    go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
    ```
  - ตรวจสอบการติดตั้งด้วย `protoc --version`

- **gRPCurl**: สำหรับทดสอบ gRPC endpoints
  - ดาวน์โหลดจาก [github.com/fullstorydev/grpcurl/releases](https://github.com/fullstorydev/grpcurl/releases)
  - เพิ่มลงใน PATH
  - ตรวจสอบการติดตั้งด้วย `grpcurl --version`

## โครงสร้างโปรเจค

```
├── api/               # นิยาม API
│   └── proto/         # นิยาม Protocol Buffers
├── cmd/               # จุดเริ่มต้นของแอปพลิเคชัน
│   └── server/        # การใช้งาน gRPC server
├── docs/              # ไฟล์เอกสาร
├── internal/          # รหัสแอปพลิเคชันส่วนตัว
│   ├── auth/          # ตรรกะการพิสูจน์ตัวตน
│   │   └── context.go # ฟังก์ชัน context การพิสูจน์ตัวตน
│   ├── config/        # การจัดการการตั้งค่า
│   ├── db/            # การเชื่อมต่อฐานข้อมูล
│   ├── middleware/    # gRPC interceptors
│   ├── models/        # โมเดลข้อมูล
│   ├── services/      # ตรรกะธุรกิจ
│   │   └── user/      # การใช้งานบริการที่เกี่ยวข้องกับผู้ใช้
│   │       ├── service.go       # การเริ่มต้นบริการผู้ใช้
│   │       ├── auth.go          # ฟังก์ชันการพิสูจน์ตัวตน
│   │       ├── profile.go       # การจัดการโปรไฟล์
│   │       ├── registration.go  # การลงทะเบียนผู้ใช้
│   │       ├── list.go          # ฟังก์ชันการแสดงรายการผู้ใช้
│   │       ├── password.go      # ฟังก์ชันรีเซ็ตรหัสผ่าน
│   │       └── utils.go         # ยูทิลิตี้เฉพาะบริการผู้ใช้
│   └── utils/         # ฟังก์ชันและตัวช่วยยูทิลิตี้
│       ├── constants.go  # ค่าคงที่แอปพลิเคชัน
│       ├── token.go      # ยูทิลิตี้การจัดการโทเค็น
│       └── validation.go # ยูทิลิตี้ตรวจสอบข้อมูลนำเข้า
├── test/              # ไฟล์ทดสอบ
├── .env               # ตัวแปรสภาพแวดล้อม (สร้างไฟล์นี้)
├── docker-compose.yml # การตั้งค่า Docker Compose
├── Dockerfile         # นิยามภาพ Docker
├── go.mod             # นิยามโมดูล Go
├── go.sum             # checksums โมดูล Go
└── README.md          # เอกสารโปรเจค
```

## คำสั่งลัดสำหรับการใช้งาน

### การตั้งค่าและเริ่มใช้งาน

#### 1. การติดตั้ง Dependencies

```bash
go mod tidy
```

#### 2. การสร้างไฟล์สภาพแวดล้อม

สร้างไฟล์ `.env` ที่โฟลเดอร์หลักด้วยข้อมูลต่อไปนี้:

```
PORT=:50051
MONGO_URI=mongodb://localhost:27017
DB_NAME=usermanagement
JWT_SECRET_KEY=your-secure-secret-key
TOKEN_DURATION=24h
RATE_LIMIT=5
RATE_LIMIT_WINDOW=1m
```

#### 3. การสร้าง Protocol Buffer Code

```bash
protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative api/proto/user_service.proto
```

#### 4. การสร้างและเริ่มการทำงานของบริการ

```bash
go build -o server.exe ./cmd/server
./server.exe
```

#### 5. การใช้งาน Docker

```bash
docker compose up --build -d
```

### การทดสอบ API

#### การลงทะเบียนผู้ใช้ใหม่

```bash
grpcurl -d '{"name": "ทดสอบ ผู้ใช้", "email": "test@example.com", "password": "Password123"}' -plaintext localhost:50051 proto.UserService/Register
```

#### การเข้าสู่ระบบ

```bash
grpcurl -d '{"email": "test@example.com", "password": "Password123"}' -plaintext localhost:50051 proto.UserService/Login
```

#### การดูข้อมูลผู้ใช้

```bash
grpcurl -H "Authorization: Bearer YOUR_JWT_TOKEN" -d '{"user_id": "USER_ID"}' -plaintext localhost:50051 proto.UserService/GetProfile
```

#### การแสดงรายการบริการที่มี

```bash
grpcurl -plaintext localhost:50051 list
```

## การปรับปรุงล่าสุด

โปรเจ็กต์นี้เพิ่งได้รับการปรับโครงสร้างครั้งใหญ่เพื่อพัฒนาการจัดระเบียบโค้ดและการบำรุงรักษา:

1. **การแยกเซอร์วิสแบบ Monolithic**: แยกไฟล์ `user_service.go` เดิมออกเป็นหลายไฟล์ตามโดเมน
   - แต่ละไฟล์มุ่งเน้นที่ฟังก์ชันการทำงานเฉพาะ (auth, profile, registration, ฯลฯ)
   - การแบ่งแยกหน้าที่และความรับผิดชอบที่ดีขึ้นและปรับปรุงความสามารถในการอ่าน

2. **โครงสร้างตามหลัก Domain-Driven**: จัดระเบียบโค้ดใหม่ตามหลักการออกแบบตามโดเมน
   - ย้ายจากโครงสร้างแบบแบนราบไปเป็นโครงสร้างตามโดเมน
   - อัปเดต imports จาก `services` เป็น `services/user`

3. **การแยกฟังก์ชันที่ใช้ร่วมกัน**: สร้างแพ็คเกจยูทิลิตี้ทั่วไปสำหรับฟังก์ชันที่ใช้ร่วมกัน
   - `internal/auth/context.go`: ฟังก์ชัน context การพิสูจน์ตัวตนส่วนกลาง
   - `internal/utils/validation.go`: ตัวช่วยตรวจสอบข้อมูลนำเข้า
   - `internal/utils/token.go`: การสร้างและจัดการโทเค็น
   - `internal/utils/constants.go`: ค่าคงที่ทั่วทั้งระบบ

4. **การเริ่มต้นเซอร์วิส**: อัปเดตการเริ่มต้นเซอร์วิสใน `main.go`
   - เปลี่ยนจาก `services.NewUserService()` เป็น `user.NewService()`

การปรับโครงสร้างนี้ช่วยปรับปรุงการบำรุงรักษาโค้ด ลดหนี้ทางเทคนิค และทำให้โค้ดเบสสามารถปรับขนาดได้ดีขึ้นสำหรับการพัฒนาในอนาคต

## เอกสารเพิ่มเติม

สำหรับข้อมูลเพิ่มเติม โปรดดูที่:

- [SETUP_GUIDE.md](SETUP_GUIDE.md) - คู่มือการตั้งค่าโดยละเอียด
- [QUICK_COMMANDS.md](QUICK_COMMANDS.md) - คำสั่งด่วนสำหรับอ้างอิง
- [README.md](../README.md) - เอกสารหลักภาษาอังกฤษ

สำหรับคำถามหรือการมีส่วนร่วม โปรดเปิด issue หรือ pull request บนที่เก็บ GitHub
