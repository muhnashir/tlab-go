# Docker Setup untuk E-Wallet API

## 📦 Prerequisite

- Docker (versi 20.10+)
- Docker Compose (versi 2.0+)

## 🚀 Cara Menjalankan

### 1. Setup Environment Variables

Salin file `.env.example` menjadi `.env` dan sesuaikan konfigurasinya:

```bash
cp .env.example .env
```

Edit file `.env` dan sesuaikan nilai-nilai berikut:

```env
DB_HOST=db
DB_PORT=3306
DB_USER=wallet_user
DB_PASSWORD=your_secure_password
DB_NAME=wallet_api
JWT_SECRET=your_jwt_secret_key
APP_PORT=3000
```

### 2. Build dan Jalankan Container

```bash
# Build dan jalankan semua service
docker-compose up -d

# Lihat logs
docker-compose logs -f

# Lihat logs untuk service tertentu
docker-compose logs -f app
docker-compose logs -f db
```

### 3. Verifikasi Service Berjalan

```bash
# Cek status container
docker-compose ps

# Test API endpoint
curl http://localhost:3000/api/health
```

### 4. Menjalankan Migration

Jika perlu menjalankan migration secara manual:

```bash
# Masuk ke container app
docker-compose exec app sh

# Jalankan migration (jika ada tool migrate)
# migrate -path /app/migrations -database "mysql://user:pass@tcp(db:3306)/dbname" up
```

## 🛠️ Perintah Berguna

### Menghentikan Container

```bash
docker-compose down
```

### Menghentikan dan Menghapus Volume (⚠️ Data akan hilang)

```bash
docker-compose down -v
```

### Rebuild Image

```bash
docker-compose build --no-cache
docker-compose up -d
```

### Melihat Resource Usage

```bash
docker-compose stats
```

### Mengakses MySQL Database

```bash
# Via docker-compose
docker-compose exec db mysql -u wallet_user -p wallet_api

# Via MySQL client dari host (port 3306)
mysql -h 127.0.0.1 -P 3306 -u wallet_user -p wallet_api
```

## 🏗️ Arsitektur Docker

### Multi-Stage Build

Dockerfile menggunakan 2 stage:

1. **Builder Stage** (`golang:1.24-alpine`): Compile aplikasi
2. **Runtime Stage** (`alpine:latest`): Jalankan binary yang sudah dikompilasi

Keuntungan:

- Image size lebih kecil (~20MB vs ~300MB)
- Lebih aman (tidak ada source code atau build tools di production image)
- Startup lebih cepat

### Docker Compose Services

#### 1. Database Service (`db`)

- Image: `mysql:8.0`
- Port: `3306`
- Health check: Menggunakan `mysqladmin ping`
- Volume: Data persisten di `db_data`

#### 2. Application Service (`app`)

- Build dari `Dockerfile`
- Port: `3000` (host) → `8080` (container)
- Depends on: `db` dengan kondisi `service_healthy`
- Environment: Membaca dari file `.env`

## 🔒 Security Best Practices

1. ✅ Non-root user di container
2. ✅ Minimal base image (Alpine Linux)
3. ✅ CA certificates installed
4. ✅ Environment variables untuk credentials
5. ⚠️ Pastikan `.env` tidak di-commit ke Git

## 📊 Monitoring

### Melihat Logs Real-time

```bash
docker-compose logs -f --tail=100
```

### Inspect Container

```bash
docker-compose exec app ps aux
docker-compose exec app netstat -tlnp
```

## 🐛 Troubleshooting

### Database Connection Failed

```bash
# Check apakah DB sudah healthy
docker-compose ps

# Check logs database
docker-compose logs db

# Test koneksi dari app container
docker-compose exec app ping db
```

### Port Already in Use

```bash
# Lihat process yang menggunakan port 3000
lsof -i :3000

# Atau gunakan port lain di docker-compose.yaml
ports:
  - "8000:8080"  # Ganti 3000 menjadi 8000
```

### Image Size Terlalu Besar

```bash
# Check ukuran image
docker images | grep wallet

# Cleanup unused images
docker image prune -a
```

## 📝 Notes

- Database data akan persist di Docker volume `db_data`
- Aplikasi akan restart otomatis jika crash (`restart: unless-stopped`)
- Healthcheck memastikan aplikasi hanya start setelah database siap
- Port mapping: `localhost:3000` → `container:8080`
