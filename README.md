# go-uni

University REST API in Go for managing students, teachers, courses, and enrollments.

## Project Summary

This project demonstrates a layered backend architecture with:

- Clean handler/repository separation
- PostgreSQL integration
- Input validation and consistent JSON responses
- JWT-based authentication for write operations
- Swagger/OpenAPI documentation
- Request and error logging middleware

## Key Features (For Review)

1. CRUD APIs for core entities:
- Students
- Teachers
- Courses

2. Enrollment flow:
- Enroll a student in a course
- Unenroll a student from a course

3. Authentication and authorization:
- Token endpoint: POST /auth/token
- JWT middleware protects POST, PUT, DELETE routes
- Public read access for GET routes

4. Validation and error handling:
- Request body validation via go-playground/validator
- Unknown JSON fields are rejected
- Consistent error responses

5. API documentation:
- Swagger UI endpoint available
- Security scheme for Bearer token configured

6. Observability:
- Request logging (method, path, status, duration)
- Error logging for handler and database failure paths

7. Test coverage:
- Handler tests with sqlmock
- Core routes and error scenarios covered

## Tech Stack

- Go (net/http, ServeMux)
- PostgreSQL
- swaggo/swag + http-swagger
- go-playground/validator
- sqlmock + testify (tests)

## Project Structure

- cmd/server: application entrypoint
- internal/database: DB connection setup
- internal/env: environment config helpers
- internal/handlers: HTTP handlers and routing
- internal/repository: SQL data access layer
- internal/models: domain models and DTOs
- pkg/middleware: request logging and JWT auth middleware
- docs: generated OpenAPI docs
- migrations: database schema

## Environment Variables

Required/used by application:

- DB_ADDR: PostgreSQL DSN
- HTTP_ADDR: HTTP bind address (default :8080)
- AUTH_USERNAME: username for token issue endpoint
- AUTH_PASSWORD: password for token issue endpoint
- AUTH_JWT_SECRET: secret used to sign/verify JWT
- AUTH_TOKEN_TTL_MINUTES: token TTL in minutes (default 60)

Example .env:

```env
DB_ADDR=postgres://admin:postgres@localhost:5432/uni_db?sslmode=disable
HTTP_ADDR=:8080
AUTH_USERNAME=admin
AUTH_PASSWORD=admin
AUTH_JWT_SECRET=change-me
AUTH_TOKEN_TTL_MINUTES=60
```

Note: .env.example may contain older key names; runtime expects DB_ADDR.

## Run Locally

1. Start database:

```bash
docker compose up -d db
```

2. Run API:

```bash
go run ./cmd/server
```

3. Open Swagger UI:

```text
http://localhost:8080/swagger/index.html
```

## Authentication Flow

1. Request token:

POST /auth/token

```json
{
	"username": "admin",
	"password": "admin"
}
```

2. Use returned token in protected requests:

Authorization: Bearer <token>

## Endpoint Overview

Public (no token):

- GET /students
- GET /students/{id}
- GET /teachers
- GET /teachers/{id}
- GET /courses
- GET /courses/{id}
- POST /auth/token

Protected (JWT required):

- POST, PUT, DELETE /students
- POST, PUT, DELETE /teachers
- POST, PUT, DELETE /courses
- POST /students/{id}/courses/{course_id}
- DELETE /students/{id}/courses/{course_id}

## Testing

Run all tests:

```bash
go test ./...
```

## Swagger Regeneration

If annotations change:

```bash
swag init -g cmd/server/main.go -o docs
```

## Reviewer Notes

- API responses for success are wrapped as:

```json
{"data": ...}
```

- Errors are returned as:

```json
{"error": "..."}
```

- Authentication is intentionally simple (credential-based token issue) for project scope.
