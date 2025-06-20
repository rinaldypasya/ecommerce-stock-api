# 🛒 ecommerce-stock-api

A modular, scalable backend system for e-commerce and stock management built with Go, PostgreSQL, and Docker.

---

## ✅ Features

- 📦 Product listing with live stock availability
- 🛒 Order checkout with stock reservation
- 🔁 Background job to release expired orders
- 🏬 Shop and multi-warehouse support
- 🔄 Transfer stock between warehouses (concurrency-safe)
- 🔐 JWT-based user authentication
- 📊 Monitoring via Prometheus metrics
- 📁 Structured logging using Zerolog
- 🐳 Containerized with Docker & Docker Compose

---

## 🚀 Getting Started

### 🧱 Requirements
- Go 1.22+
- Docker + Docker Compose
- PostgreSQL

### 📦 Environment Variables
Create a `.env` file or set these manually:

```env
DB_DSN=postgres://postgres:postgres@localhost:5432/ecommerce?sslmode=disable
JWT_SECRET=your_super_secret_key
```

### 🐳 Run with Docker Compose

```bash
docker-compose up --build
```

App will be available at: `http://localhost:8080`

---

## 📚 API Endpoints

### 🔐 Auth
- `POST /register`
- `POST /login`

### 📦 Products
- `GET /products`

### 🛒 Orders
- `POST /order/checkout`

### 🏬 Shops
- `POST /shops`
- `GET /shops`

### 🏢 Warehouses
- `POST /warehouses`
- `POST /warehouses/status`
- `POST /warehouses/transfer`
- `GET /warehouses`

### 📈 Monitoring
- `GET /metrics` (Prometheus-compatible)

---

## 🧪 Testing

### Run unit + integration tests
```bash
go test ./internal/... -v
```

### Generate code coverage report
```bash
go test ./internal/... -coverprofile=coverage.out
go tool cover -html=coverage.out -o coverage.html
```

---

## ⚙️ Deployment

### With Docker Compose (dev/staging)
- Edit `docker-compose.yml`
- Use `.env` for secrets

---

## 🧠 Architecture

- Layered Clean Architecture: `handler → service → repository → db`
- GORM ORM
- Prometheus metrics / Zerolog logs
- Row-level locking with `SELECT FOR UPDATE` to prevent overselling
- Background job for order expiration using Go routine

---

## 👥 Contributors
- 💻 Rinaldy Pasya | rinaldypasya@gmail.com
