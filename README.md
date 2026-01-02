# Wood Library API

Backend API cho hệ thống quản lý thư viện gỗ - xây dựng với Go + Gin + Firebase.

## Tech Stack

- Go 1.25
- Gin Framework
- Firebase (Firestore + Auth)
- Cloudinary (Image storage)

## Cài đặt

```bash
# Tải dependencies
go mod download

# Chạy server
go run main.go
```

## Cấu hình

1. Đặt file `serviceAccount.json` (Firebase credentials) vào thư mục root
2. Cấu hình Cloudinary trong `config/cloudinary.go`

## Cấu trúc thư mục

```
├── config/       # Cấu hình Firebase, Cloudinary
├── firestore/    # Firestore operations
├── handler/      # HTTP handlers
├── middleware/   # Auth middleware
├── models/       # Data models
├── router/       # Routes
└── service/      # Business logic
```

## API Endpoints

### Model API (`/model-api`)
| Method | Endpoint | Mô tả |
|--------|----------|-------|
| GET | `/version` | Lấy version model hiện tại |
| GET | `/list_versions` | Danh sách versions |
| POST | `/activate` | Kích hoạt version |
| POST | `/upload` | Upload model mới |

### Library API (`/library-api`) - Yêu cầu Auth
| Method | Endpoint | Mô tả |
|--------|----------|-------|
| POST | `/upload_image` | Upload hình ảnh |
| GET | `/database/list` | Danh sách bộ sưu tập |
| GET | `/database/get` | Chi tiết bộ sưu tập |
| POST | `/database/create` | Tạo bộ sưu tập |
| PUT | `/database/update/:id` | Cập nhật bộ sưu tập |
| DELETE | `/database/delete` | Xóa bộ sưu tập |
| GET | `/piece/list` | Danh sách mẫu gỗ |
| GET | `/piece/get` | Chi tiết mẫu gỗ |
| POST | `/piece/create` | Tạo mẫu gỗ |
| PUT | `/piece/update/:id` | Cập nhật mẫu gỗ |
| DELETE | `/piece/delete` | Xóa mẫu gỗ |

## Authentication

Sử dụng Firebase ID Token trong header:

```
Authorization: Bearer <firebase_id_token>
```

## Docker

```bash
docker build -t wood-library-api .
docker run -p 8080:8080 wood-library-api
```

## Environment

Server chạy mặc định trên port `8080`.