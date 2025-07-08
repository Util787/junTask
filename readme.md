# Test Task REST API

This is a RESTful API for managing users. It features user creation, retrieval, updating, and deletion, along with external API integration and caching.

## Main Features

- Retrieve all users with pagination and filtering
- Create users with automatic enrichment using external APIs:
  - https://api.agify.io/ (age)
  - https://api.genderize.io/ (gender)
  - https://api.nationalize.io/ (nationality)
- Partial user updates (only provided fields are changed)
- Redis caching
- Swagger UI for API documentation
- PostgreSQL database support

---

## Installation

### Prerequisites 📦
- [Go](https://golang.org/doc/install)
- [PostgreSQL](https://www.postgresql.org/download/)
- [Redis](https://redis.io/docs/getting-started/installation/)
- Optional: [Docker](https://docs.docker.com/get-docker/)

### 1. Clone Repository 📂
```bash
git clone https://github.com/Util787/user-manager-api
cd user-manager-api
```

### 2. Configure `.env` ⚙️
Create a `.env` file and configure according to your environment (`prod`, `dev`, or `local`):

```env
ENV=local
SERVER_PORT=8000
DB_HOST=localhost
DB_PORT=5432
DB_USERNAME=postgres
DB_PASSWORD=1111
DB_NAME=postgres
SSLMODE=disable
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=2222
REDIS_DB=0
```
## Optional: Docker Compose 🐳
Now you can run the entire project using Docker Compose.

```bash
docker-compose up --build
```

This will initialize PostgreSQL(migrations will be applied automatically), Redis, and the API server simultaneously in containers. (Api docs on [step 5](#5-api-documentation-
))

If you dont want to run with docker, skip this step and proceed to step 3.

### 3. Database Migration 🐘
Run migrations to set up your PostgreSQL database schema:

```bash
migrate -path sql/schema -database "postgres://postgres:1111@localhost:5432/postgres?sslmode=disable" up
```

### 4. Run the Application ▶️
Execute the following command from the project directory:

```bash
go run cmd/main.go
```

### 5. API Documentation 📘
Endpoints and usage are documented and available through Swagger at (example):

```
http://localhost:8000/swagger/index.html
```

---

## Additional Tech Stack 🛠️

- Gin (HTTP framework)
- Sqlx + Pgx (For postgres)
- Swagger (API docs)
- Mockery + Testify (testing)

