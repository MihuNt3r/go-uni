# go-uni

University REST API in Go for managing students, teachers, courses, and enrollments.

## Project Summary

This project demonstrates a layered backend architecture with:

- Handler/repository separation
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
- Enrollments

2. Authentication and authorization:
- Token endpoint: POST /auth/token
- JWT middleware protects POST, PUT, DELETE routes
- Public read access for GET routes

3. Validation and error handling:
- Request body validation via go-playground/validator
- Unknown JSON fields are rejected
- Consistent error responses

4. API documentation:
- Swagger UI endpoint available
- Security scheme for Bearer token configured

5. Observability:
- Request logging (method, path, status, duration)
- Error logging for handler and database failure paths

6. Test coverage:
- Handler tests with sqlmock
- Core routes and error scenarios covered
