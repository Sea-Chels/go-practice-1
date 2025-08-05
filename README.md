# Go API with JWT Authentication & PostgreSQL

A secure REST API built with Go featuring JWT authentication, PostgreSQL database, and Docker support.

## Features

- JWT-based authentication
- PostgreSQL database with migrations
- Docker Compose for easy local development
- Secure password hashing with bcrypt
- Input validation and sanitization
- CORS support
- Graceful shutdown
- Database connection pooling

## Prerequisites

- Go 1.23 or higher
- Docker and Docker Compose
- PostgreSQL (if running without Docker)

## Quick Start with Docker

1. Clone the repository:
```bash
git clone <repository-url>
cd go-practice-1
```

2. Start the application with Docker Compose:
```bash
# Production mode (optimized image)
make docker-up

# Development mode (with hot reload)
make docker-dev

# View logs
make docker-logs
```

3. The API will be available at `http://localhost:8080`

4. Check the health endpoint:
```bash
curl http://localhost:8080/health
# Response: {"status":"ok","database":"connected"}
```

## Local Development without Docker

1. Install PostgreSQL and create a database:
```sql
CREATE DATABASE school_db;
CREATE USER devuser WITH PASSWORD 'devpass123';
GRANT ALL PRIVILEGES ON DATABASE school_db TO devuser;
```

2. Copy the environment file:
```bash
cp .env.example .env
```

3. Update `.env` with your database credentials

4. Install dependencies:
```bash
go mod download
```

5. Run the application:
```bash
make run
# or
go run cmd/api/main.go
```

## API Endpoints

### Public Endpoints

#### Health Check
```bash
GET /health
```

#### Login
```bash
POST /auth/login
Content-Type: application/json

{
  "email": "admin@example.com",
  "password": "Admin123!"
}

# Response:
{
  "token": "eyJhbGc...",
  "expires_at": "2024-01-02T15:04:05Z"
}
```

### Protected Endpoints (Require JWT)

#### Get All Students
```bash
GET /students
Authorization: Bearer <jwt-token>

# Response:
{
  "students": [
    {
      "id": 1,
      "name": "Alice Johnson",
      "grade": 10,
      "created_at": "2024-01-01T10:00:00Z",
      "updated_at": "2024-01-01T10:00:00Z"
    }
  ],
  "count": 5
}
```

#### Create Student
```bash
POST /students
Authorization: Bearer <jwt-token>
Content-Type: application/json

{
  "name": "John Doe",
  "grade": 10
}

# Response:
{
  "id": 6,
  "name": "John Doe", 
  "grade": 10,
  "created_at": "2024-01-01T10:00:00Z",
  "updated_at": "2024-01-01T10:00:00Z"
}
```

## Default Credentials

The database is seeded with an admin user:
- Email: `admin@example.com`
- Password: `Admin123!`

## Project Structure

```
.
├── cmd/api/              # Application entry points
├── internal/             # Private application code
│   ├── auth/            # JWT authentication
│   ├── database/        # Database connection and operations
│   ├── handlers/        # HTTP request handlers
│   ├── models/          # Data models
│   └── utils/           # Utility functions
├── migrations/          # SQL migration files
├── docker/              # Docker-related files
├── .env.example         # Example environment variables
├── docker-compose.yml   # Docker Compose configuration
├── Dockerfile           # Docker image definition
├── Makefile            # Build and run commands
└── README.md           # This file
```

## Adding New Features

### Adding a New Endpoint

1. Create a new handler in `internal/handlers/`:
```go
// internal/handlers/teachers.go
func GetTeachersHandler(w http.ResponseWriter, r *http.Request) {
    // Implementation
}
```

2. Add the route in `cmd/api/main.go`:
```go
// Protected route
router.HandleFunc("/teachers", auth.JWTMiddleware(handlers.GetTeachersHandler)).Methods("GET")

// Public route
router.HandleFunc("/teachers", handlers.GetTeachersHandler).Methods("GET")
```

### Adding a New Migration

1. Create a new SQL file in `migrations/` with the next sequence number:
```bash
touch migrations/003_create_teachers_table.sql
```

2. Add your SQL:
```sql
CREATE TABLE IF NOT EXISTS teachers (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    subject VARCHAR(255) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE
);
```

3. Restart the application to run migrations automatically

### Adding New Models

1. Create a new model file in `internal/models/`:
```go
// internal/models/teacher.go
package models

import "time"

type Teacher struct {
    ID        int        `json:"id"`
    Name      string     `json:"name"`
    Subject   string     `json:"subject"`
    CreatedAt time.Time  `json:"created_at"`
    UpdatedAt time.Time  `json:"updated_at"`
    DeletedAt *time.Time `json:"deleted_at,omitempty"`
}
```

## Security Considerations

- Passwords are hashed using bcrypt with cost factor 14
- JWT tokens expire after 24 hours
- All database queries use prepared statements to prevent SQL injection
- Input validation on all user inputs
- CORS is configured (update `ALLOWED_ORIGINS` in production)
- Environment variables for sensitive configuration

## Database Access & Management

### Accessing PostgreSQL Database

When running with Docker, you can access the database in several ways:

1. **Using psql command line:**
```bash
# Connect to PostgreSQL container
docker exec -it go-practice-1-postgres-1 psql -U devuser -d school_db

# Common psql commands:
\dt              # List all tables
\d users         # Describe users table
\d students      # Describe students table
SELECT * FROM users;
SELECT * FROM students;
\q               # Quit psql
```

2. **Using Docker Compose run:**
```bash
docker-compose exec postgres psql -U devuser -d school_db
```

3. **Connection details for GUI tools (e.g., TablePlus, DBeaver, pgAdmin):**
- Host: `localhost`
- Port: `5432`
- Database: `school_db`
- Username: `devuser`
- Password: `devpass123`

### Database Migrations

Migrations are automatically run when the application starts. They are located in the `migrations/` directory and run in alphabetical order.

To manually run migrations or check migration status:

```bash
# View migration files
ls -la migrations/

# Connect to database and check if tables exist
docker exec -it go-practice-1-postgres-1 psql -U devuser -d school_db -c "\dt"

# View the structure of created tables
docker exec -it go-practice-1-postgres-1 psql -U devuser -d school_db -c "\d+ users"
docker exec -it go-practice-1-postgres-1 psql -U devuser -d school_db -c "\d+ students"
```

### Resetting the Database

If you need to reset the database and re-run migrations:

```bash
# Stop containers and remove volumes
docker-compose down -v

# Start fresh (migrations will run automatically)
docker-compose up -d
```

## Available Make Commands

```bash
make help         # Show all available commands
make docker-up    # Start Docker containers
make docker-down  # Stop Docker containers
make docker-dev   # Start in development mode
make docker-logs  # View container logs
make run          # Run locally (requires local PostgreSQL)
make test         # Run tests
make lint         # Run linter
make build        # Build the application
```

## Environment Variables

See `.env.example` for all available configuration options:

- `DB_*` - Database configuration
- `JWT_SECRET` - Secret key for JWT signing (change in production!)
- `PORT` - Server port (default: 8080)
- `ALLOWED_ORIGINS` - CORS allowed origins

## Troubleshooting

### Database Connection Issues
- Ensure PostgreSQL is running
- Check database credentials in `.env`
- Verify network connectivity between containers

### JWT Token Issues
- Ensure the token is included in the Authorization header
- Check token expiration
- Verify JWT_SECRET matches between requests

### Docker Issues
- Run `docker-compose logs` to see container logs
- Ensure ports 5432 and 8080 are not in use
- Try `docker-compose down -v` to reset volumes

## License

This project is for educational purposes.