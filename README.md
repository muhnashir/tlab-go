# 💰 E-Wallet API

RESTful API untuk aplikasi dompet digital (E-Wallet) yang dibangun dengan Golang. Aplikasi ini menyediakan fitur autentikasi, manajemen wallet, dan transfer antar pengguna.

## 📋 Daftar Isi

- [Teknologi yang Digunakan](#-teknologi-yang-digunakan)
- [Fitur Utama](#-fitur-utama)
- [Arsitektur](#-arsitektur)
- [Kebutuhan Sistem](#-kebutuhan-sistem)
- [Instalasi dan Menjalankan Proyek](#-instalasi-dan-menjalankan-proyek)
  - [Opsi 1: Menggunakan Docker (Recommended)](#opsi-1-menggunakan-docker-recommended)
  - [Opsi 2: Menjalankan Secara Lokal](#opsi-2-menjalankan-secara-lokal)
- [API Documentation](#-api-documentation)
- [Database Schema](#-database-schema)
- [Environment Variables](#-environment-variables)
- [Testing](#-testing)
- [Project Structure](#-project-structure)

---

## 🛠 Teknologi yang Digunakan

### Backend Framework & Libraries

- **[Go](https://golang.org/)** v1.24+ - Programming Language
- **[Fiber](https://gofiber.io/)** v2.52+ - Web Framework (Express-like untuk Go)
- **[MySQL](https://www.mysql.com/)** 8.0 - Relational Database
- **[GoQu](https://github.com/doug-martin/goqu)** v9.19+ - SQL Query Builder
- **[JWT](https://github.com/golang-jwt/jwt)** v5.3+ - JSON Web Token untuk Authentication
- **[Bcrypt](https://pkg.go.dev/golang.org/x/crypto/bcrypt)** - Password Hashing
- **[Godotenv](https://github.com/joho/godotenv)** - Environment Variable Management
- **[Migrate](https://github.com/golang-migrate/migrate)** v4.19+ - Database Migration Tool

### Development Tools

- **[Swagger](https://github.com/swaggo/swag)** v1.16+ - API Documentation
- **[Air](https://github.com/cosmtrek/air)** - Live Reload untuk Development
- **Docker** & **Docker Compose** - Containerization

### Architecture Pattern

- **Clean Architecture** - Separation of Concerns
- **Repository Pattern** - Data Access Layer
- **Dependency Injection** - Loose Coupling

---

## ✨ Fitur Utama

- ✅ **User Authentication**
  - Register user baru
  - Login dengan JWT token
  - Get user profile (protected)

- 💼 **Wallet Management**
  - Auto-create wallet saat register
  - Top-up saldo
  - Cek saldo wallet

- 💸 **Transaction**
  - Transfer antar wallet
  - Transaction history
  - Transaction status tracking

- 🔐 **Security**
  - Password hashing dengan bcrypt
  - JWT-based authentication
  - Protected routes dengan middleware

---

## 🏗 Arsitektur

Proyek ini menggunakan **Clean Architecture** dengan layer berikut:

```
┌─────────────────────────────────────────┐
│         Delivery Layer (HTTP)           │
│    (Handlers, Middleware, Router)       │
└──────────────┬──────────────────────────┘
               │
┌──────────────▼──────────────────────────┐
│         Service Layer                   │
│    (Business Logic)                     │
└──────────────┬──────────────────────────┘
               │
┌──────────────▼──────────────────────────┐
│         Repository Layer                │
│    (Database Access)                    │
└──────────────┬──────────────────────────┘
               │
┌──────────────▼──────────────────────────┐
│         Domain Layer                    │
│    (Entities, Interfaces)               │
└─────────────────────────────────────────┘
```

---

## 📦 Kebutuhan Sistem

### Untuk Menjalankan dengan Docker (Recommended):

- **Docker** versi 20.10 atau lebih baru
- **Docker Compose** versi 2.0 atau lebih baru
- **Make** (optional, untuk kemudahan perintah)

### Untuk Menjalankan Secara Lokal:

- **Go** versi 1.24 atau lebih baru
- **MySQL** versi 8.0 atau MariaDB
- **Make** (optional)
- **Air** (optional, untuk live reload)

---

## 🚀 Instalasi dan Menjalankan Proyek

### Opsi 1: Menggunakan Docker (Recommended)

Docker adalah cara tercepat dan termudah untuk menjalankan proyek ini karena semua dependensi sudah dikemas dalam container.

#### Quick Start (3 Langkah)

```bash
# 1. Clone repository
git clone <repository-url>
cd tlab-go

# 2. Jalankan aplikasi (otomatis membuat .env dari .env.example)
make up

# 3. Akses aplikasi
curl http://localhost:3000/health
```

#### Setup Detail

#### 1. Clone Repository

```bash
git clone <repository-url>
cd tlab-go
```

#### 2. Setup Environment Variables (Opsional)

**File `.env` akan otomatis dibuat dari `.env.example` saat menjalankan `make up`.**

Jika ingin mengubah konfigurasi default:

```bash
# Buat .env secara manual (opsional)
cp .env.example .env

# Edit file .env sesuai kebutuhan
nano .env
```

Atau gunakan perintah setup:

```bash
make setup  # Membuat .env dari .env.example
```

#### 3. Build dan Jalankan dengan Docker Compose

```bash
# Menggunakan Make (otomatis membuat .env jika belum ada)
make up

# Atau menggunakan docker-compose langsung
docker-compose up -d
```

> **💡 Tips**: Perintah `make up` akan otomatis membuat file `.env` dari `.env.example` jika file `.env` belum ada.

#### 4. Verifikasi Service Berjalan

```bash
# Cek status container
make ps
# atau
docker-compose ps

# Lihat logs
make logs
# atau
docker-compose logs -f

# Test API
curl http://localhost:3000/health
```

#### 5. Perintah Docker Lainnya

```bash
make help          # Lihat semua perintah yang tersedia
make setup         # Setup awal (buat .env dari .env.example)
make down          # Stop semua service
make restart       # Restart service
make logs-app      # Lihat logs aplikasi saja
make logs-db       # Lihat logs database saja
make shell-app     # Masuk ke container aplikasi
make shell-db      # Masuk ke MySQL shell
make clean         # Hapus semua container dan volume
make rebuild       # Clean rebuild
```

📖 **Dokumentasi lengkap Docker**: Lihat [README.Docker.md](./README.Docker.md)

---

### Opsi 2: Menjalankan Secara Lokal

#### 1. Install Go

```bash
# Cek versi Go
go version

# Jika belum terinstall, download dari https://golang.org/dl/
```

#### 2. Install MySQL

```bash
# Ubuntu/Debian
sudo apt update
sudo apt install mysql-server

# macOS (menggunakan Homebrew)
brew install mysql

# Start MySQL service
sudo systemctl start mysql    # Linux
brew services start mysql     # macOS
```

#### 3. Clone dan Setup Project

```bash
# Clone repository
git clone <repository-url>
cd tlab-go

# Install dependencies
go mod download
```

#### 4. Setup Database

```bash
# Login ke MySQL
mysql -u root -p

# Buat database
CREATE DATABASE wallet_api;
CREATE USER 'wallet_user'@'localhost' IDENTIFIED BY 'wallet_password';
GRANT ALL PRIVILEGES ON wallet_api.* TO 'wallet_user'@'localhost';
FLUSH PRIVILEGES;
EXIT;
```

#### 5. Setup Environment Variables

```bash
# Salin .env.example ke .env
cp .env.example .env

# Edit .env untuk konfigurasi lokal
nano .env
```

Ubah konfigurasi database di `.env`:

```env
DB_HOST=127.0.0.1
DB_PORT=3306
DB_USER=wallet_user
DB_PASSWORD=wallet_password
DB_NAME=wallet_api
JWT_SECRET=your_secret_key_here
APP_PORT=3000
```

#### 6. Jalankan Database Migration

```bash
# Install golang-migrate (jika belum ada)
# macOS
brew install golang-migrate

# Linux
curl -L https://github.com/golang-migrate/migrate/releases/download/v4.15.2/migrate.linux-amd64.tar.gz | tar xvz
sudo mv migrate /usr/local/bin/

# Jalankan migration
migrate -path migrations -database "mysql://wallet_user:wallet_password@tcp(127.0.0.1:3306)/wallet_api" up
```

#### 7. Jalankan Aplikasi

**Opsi A: Jalankan Langsung**

```bash
go run cmd/api/main.go
```

**Opsi B: Build dan Jalankan**

```bash
# Build binary
go build -o bin/api cmd/api/main.go

# Jalankan binary
./bin/api
```

**Opsi C: Menggunakan Air (Live Reload)**

```bash
# Install Air
go install github.com/cosmtrek/air@latest

# Jalankan dengan Air
air
```

#### 8. Verifikasi Aplikasi Berjalan

```bash
# Test endpoint
curl http://localhost:3000/api/health

# Atau buka di browser
# http://localhost:3000/swagger/index.html
```

---

## 📚 API Documentation

### Swagger UI

Setelah aplikasi berjalan, akses dokumentasi Swagger di:

- **Docker**: http://localhost:3000/swagger/index.html
- **Local**: http://localhost:3000/swagger/index.html

### Generate Swagger Docs (jika melakukan perubahan)

```bash
# Install swag
go install github.com/swaggo/swag/cmd/swag@latest

# Generate docs
swag init -g cmd/api/main.go
```

### API Endpoints

#### Authentication

| Method | Endpoint             | Description        | Auth Required |
| ------ | -------------------- | ------------------ | ------------- |
| POST   | `/api/auth/register` | Register user baru | ❌            |
| POST   | `/api/auth/login`    | Login user         | ❌            |
| GET    | `/api/users/profile` | Get profile user   | ✅            |

#### Wallet

| Method | Endpoint                | Description             | Auth Required |
| ------ | ----------------------- | ----------------------- | ------------- |
| GET    | `/api/wallets/balance`  | Get saldo wallet        | ✅            |
| POST   | `/api/wallets/topup`    | Top-up saldo            | ✅            |
| POST   | `/api/wallets/transfer` | Transfer ke wallet lain | ✅            |

#### Transaction

| Method | Endpoint            | Description             | Auth Required |
| ------ | ------------------- | ----------------------- | ------------- |
| GET    | `/api/transactions` | Get transaction history | ✅            |

---

## 🗄 Database Schema

### Users Table

```sql
CREATE TABLE users (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);
```

### Wallets Table

```sql
CREATE TABLE wallets (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    user_id BIGINT NOT NULL,
    balance DECIMAL(15, 2) NOT NULL DEFAULT 0.00,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    UNIQUE KEY (user_id)
);
```

### Transactions Table

```sql
CREATE TABLE transactions (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    sender_wallet_id BIGINT,
    receiver_wallet_id BIGINT,
    amount DECIMAL(15, 2) NOT NULL,
    status ENUM('pending', 'success', 'failed') DEFAULT 'pending',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (sender_wallet_id) REFERENCES wallets(id),
    FOREIGN KEY (receiver_wallet_id) REFERENCES wallets(id)
);
```

---

## 🔧 Environment Variables

Berikut adalah environment variables yang digunakan:

| Variable      | Description          | Default                             | Required |
| ------------- | -------------------- | ----------------------------------- | -------- |
| `DB_HOST`     | Database host        | `db` (Docker) / `127.0.0.1` (Local) | ✅       |
| `DB_PORT`     | Database port        | `3306`                              | ✅       |
| `DB_USER`     | Database username    | `wallet_user`                       | ✅       |
| `DB_PASSWORD` | Database password    | -                                   | ✅       |
| `DB_NAME`     | Database name        | `wallet_api`                        | ✅       |
| `JWT_SECRET`  | Secret key untuk JWT | -                                   | ✅       |
| `APP_PORT`    | Application port     | `3000`                              | ❌       |

---

## 🧪 Testing

```bash
# Jalankan semua tests
go test ./...

# Test dengan coverage
go test -cover ./...

# Test dengan verbose output
go test -v ./...

# Test specific package
go test ./internal/service/...
```

---

## 📁 Project Structure

```
.
├── cmd/
│   ├── api/              # Main application entry point
│   ├── migrate/          # Database migration runner
│   └── setup_db/         # Database setup utility
├── internal/
│   ├── delivery/
│   │   └── http/
│   │       ├── handler/  # HTTP handlers
│   │       ├── middleware/ # HTTP middleware
│   │       └── router.go # Route definitions
│   ├── domain/           # Domain entities & interfaces
│   ├── repository/       # Database repositories
│   ├── service/          # Business logic
│   └── pkg/
│       ├── database/     # Database connection
│       └── utils/        # Utility functions
├── migrations/           # SQL migration files
├── docs/                 # Swagger documentation
├── .env                  # Environment variables
├── .env.example          # Environment variables template
├── docker-compose.yaml   # Docker Compose configuration
├── Dockerfile            # Docker image definition
├── Makefile              # Build automation
├── go.mod                # Go module dependencies
└── README.md             # This file
```

---

## 🤝 Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

1. Fork the project
2. Create your feature branch (`git checkout -b feature/AmazingFeature`)
3. Commit your changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

---

## 📝 License

This project is licensed under the MIT License.

---

## 👨‍💻 Author

**Nashir**

---

## 📞 Support

Jika ada pertanyaan atau masalah, silakan buat issue di repository ini.

---

## 🎯 Roadmap

- [ ] Implement notification system
- [ ] Add transaction limits
- [ ] Implement QR code payment
- [ ] Add admin dashboard
- [ ] Implement webhook for external payments
- [ ] Add unit tests coverage > 80%
- [ ] Add integration tests
- [ ] Implement Redis caching
- [ ] Add rate limiting
- [ ] Implement logging with structured logs

---

**Happy Coding! 🚀**
