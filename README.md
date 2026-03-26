# URL Shortener Pro (API) 🚀

A robust and high-performance URL shortening API built with **Go**. This service features advanced JWT authentication, dynamic Role-Based Access Control (RBAC), real-time analytics, and Redis caching.

[![Go Version](https://img.shields.io)](https://golang.org)
[![Swagger](https://img.shields.io)](https://github.com)

## ✨ Key Features

- **Advanced Auth:** Secure authentication using JWT Bearer tokens.
- **Dynamic RBAC:** Flexible role and permission management system.
- **Role Toggling:** Ability to switch/update roles using Refresh Tokens.
- **Analytics:** Dedicated endpoint to track link performance and clicks.
- **Redirection:** High-speed redirecting engine.
- **User Dashboard:** Endpoints to retrieve user profile info and managed links.
- **API Docs:** Fully documented with Swagger (OpenAPI).

## 🛠 Tech Stack

- **Language:** Go 1.25.6
- **Framework:** [Echo v4](https://echo.labstack.com) (High performance)
- **Database:** PostgreSQL (using `pgx/v5`)
- **ORM:** [GORM](https://gorm.io)
- **Caching:** Redis (`go-redis/v9`)
- **Security:** JWT (`golang-jwt/v4`), Crypto (Bcrypt)
- **Validation:** Go-playground/validator
- **Configuration:** `caarlos0/env` & `godotenv`
- **Documentation:** Swag / Echo-Swagger

## 🚀 Getting Started

### 1. Prerequisites
- Go 1.25.6 or higher
- PostgreSQL & Redis instances

### 2. Installation
```bash
git clone https://github.com
cd URL-Shortener-Pro
```
### 3. Environment Setup
Create a .env file in the root directory:

DB_HOST=localhost
DB_PORT=5432
DB_USER=your_user
DB_PASSWORD=your_password
DB_NAME=url_shortener
REDIS_ADDR=localhost:6379
JWT_SECRET=your_super_secret_key

### 4. Running the App
``` bash
go mod tidy
go run main.go
```

API Documentation
Once the server is running, you can explore and test the API via Swagger UI:
👉 http://localhost:8080/swagger/index.html


Author
Sardor Innatov — GitHub
