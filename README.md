# Banking App (Go)

A simple, secure, and extensible banking application built with Go, Gin, GORM, and PostgreSQL. The app supports user registration, authentication (JWT), account management, and money transfers between accounts. It is containerized for easy deployment using Docker Compose.

## Features
- User registration and login (JWT authentication)
- Secure password hashing (bcrypt)
- Account creation, balance management, and currency support
- Money transfers between accounts (atomic, transactional)
- Entry logging for all account operations
- RESTful API with OpenAPI/Swagger documentation
- Built-in database migrations
- Containerized with Docker and Docker Compose

## Prerequisites
- [Docker](https://www.docker.com/get-started)
- [Docker Compose](https://docs.docker.com/compose/)

## Getting Started

### 1. Clone the repository
```bash
git clone https://github.com/ahmedkhaeld/banking-app.git
cd banking-app
```

### 2. Configure Environment Variables
Create a `.env` file in the project root with the following variables:

```
PORT=8080
GIN_MODE=release
JWT_SECRET_KEY=your-32-char-random-secret-key-here
DB_SOURCE=postgres://banking_user:banking_pass@db:5432/banking?sslmode=disable
```
- **Note:** `JWT_SECRET_KEY` must be at least 32 characters for security.

### 3. Build and Run with Docker Compose
This will build the Go app, start the PostgreSQL database, run migrations, and launch the API server.

```bash
docker-compose up --build
```
- The API server will be available at [http://localhost:8080](http://localhost:8080)
- The database will be available internally as `db:5432` (see `docker-compose.yml`)

### 4. Access the API Documentation
Interactive Swagger UI is available at:

[http://localhost:8080/docs/index.html](http://localhost:8080/docs/index.html)

You can explore and test all endpoints directly from the browser.

See Swagger UI for full details, request/response schemas, and authentication requirements.

## Project Structure
- `main.go` — Application entrypoint
- `db/` — Database connection, migrations, and models
- `internal/` — Business logic, services, controllers, and routes
- `common/` — Shared utilities and types
- `docs/` — API documentation (Swagger/OpenAPI)
- `Dockerfile`, `docker-compose.yml` — Containerization and orchestration

## Development Notes
- The app uses GORM for ORM and migrations.
- All primary keys are UUIDs for scalability and uniqueness.
- Passwords are hashed with bcrypt and never stored in plaintext.
- JWT tokens are required for all protected endpoints (see Swagger docs for details).
- Database migrations are run automatically on startup.
