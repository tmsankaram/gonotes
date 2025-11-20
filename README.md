# GoNotes

A RESTful API server for managing notes and files, built with Go, Gin, and PostgreSQL.

## Features

- 📝 **Notes Management**: Create, read, update, and delete notes
- 👥 **User Authentication**: JWT-based authentication with TOTP support
- 📁 **File Storage**: In-memory file management system
- 🔒 **Security**: Password hashing, JWT tokens, and middleware protection
- 🗄️ **Database**: PostgreSQL with GORM ORM
- 🐳 **Docker**: Containerized PostgreSQL setup
- ⚡ **Graceful Shutdown**: Proper server shutdown handling
- 🎯 **Request Tracking**: Request ID middleware for debugging
- 📊 **Logging**: Structured logging middleware
- 🛡️ **Recovery**: Panic recovery middleware

## Tech Stack

- **Go** 1.25.4
- **Gin** - HTTP web framework
- **GORM** - ORM library
- **PostgreSQL** 16 - Database
- **JWT** - Authentication
- **Docker** - Containerization

## Project Structure

```bash
gonotes/
├── cmd/
│   └── server/
│       └── main.go           # Application entry point
├── internal/
│   ├── auth/                 # Authentication handlers and JWT
│   │   ├── handler.go
│   │   ├── jwt.go
│   │   └── middleware.go
│   ├── config/               # Configuration management
│   │   └── config.go
│   ├── db/                   # Database connection
│   │   └── db.go
│   ├── files/                # File management
│   │   ├── handler.go
│   │   ├── model.go
│   │   └── service.go
│   ├── middleware/           # HTTP middlewares
│   │   ├── logging.go
│   │   ├── recovery.go
│   │   └── requestid.go
│   ├── notes/                # Notes CRUD operations
│   │   ├── handler.go
│   │   ├── model.go
│   │   └── service.go
│   ├── pagination/           # Pagination utilities
│   │   └── pagination.go
│   ├── response/             # Response helpers
│   │   ├── error.go
│   │   └── success.go
│   ├── router/               # Route definitions
│   │   └── router.go
│   ├── users/                # User management
│   │   ├── model.go
│   │   ├── password.go
│   │   └── service.go
│   └── validator/            # Custom validators
│       └── custom.go
├── docker-compose.yml        # Docker Compose configuration
├── go.mod                    # Go module dependencies
└── README.md
```

## Prerequisites

- Go 1.25.4 or higher
- Docker and Docker Compose (for PostgreSQL)
- PostgreSQL 16 (if running locally without Docker)

## Installation

1. **Clone the repository**

```bash
git clone https://github.com/tmsankaram/gonotes.git
cd gonotes
```

2. **Install dependencies**

```bash
go mod download
```

3. **Set up environment variables**

Create a `.env` file in the root directory:

```env
PORT=8080
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASS=securepassword
DB_NAME=gonotes
JWT_SECRET=your-secret-key-here
```

4. **Start PostgreSQL with Docker**

```bash
docker-compose up -d
```

5. **Run the application**

```bash
go run cmd/server/main.go
```

The server will start on `http://localhost:8080`

## API Endpoints

### Health Check

```api
GET /health
```

### Authentication

```api
POST /auth/register    # Register a new user
POST /auth/login       # Login and get JWT token
GET  /auth/me          # Get current user (protected)
```

### Notes

```
GET    /notes          # Get all notes (with pagination)
GET    /notes/:id      # Get a specific note
POST   /notes          # Create a new note
PUT    /notes/:id      # Update a note
DELETE /notes/:id      # Delete a note
```

### Files

```api
GET    /files          # List all files
GET    /files/:id      # Get a specific file
POST   /files          # Upload a file
DELETE /files/:id      # Delete a file
```

## Usage Examples

### Register a User

```bash
curl -X POST http://localhost:8080/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "name": "John Doe",
    "email": "john@example.com",
    "password": "securepassword"
  }'
```

### Login

```bash
curl -X POST http://localhost:8080/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "john@example.com",
    "password": "securepassword"
  }'
```

### Create a Note

```bash
curl -X POST http://localhost:8080/notes \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "title": "My First Note",
    "content": "This is the content of my note"
  }'
```

### Get All Notes

```bash
curl -X GET "http://localhost:8080/notes?page=1&limit=10" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

## Development

### Running Tests

```bash
go test ./...
```

### Build for Production

```bash
go build -o gonotes cmd/server/main.go
```

### Run with Air (Hot Reload)

Install Air for hot reloading during development:

```bash
go install github.com/air-verse/air@latest
air
```

## Docker Support

### PostgreSQL Database

The project includes a `docker-compose.yml` file for PostgreSQL:

```bash
# Start PostgreSQL
docker-compose up -d

# Stop PostgreSQL
docker-compose down

# View logs
docker-compose logs -f
```

## Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `PORT` | 8080 | Server port |
| `DB_HOST` | localhost | Database host |
| `DB_PORT` | 5432 | Database port |
| `DB_USER` | postgres | Database user |
| `DB_PASS` | | Database password |
| `DB_NAME` | postgres | Database name |
| `JWT_SECRET` | | Secret key for JWT signing |

## Features Breakdown

### Authentication

- User registration with password hashing
- JWT-based authentication
- TOTP (Time-based One-Time Password) support
- Protected routes with middleware

### Notes Management

- Full CRUD operations
- Pagination support
- Custom validators (e.g., title validation)

### File Management

- In-memory file storage
- File upload and download
- File metadata tracking

### Middleware

- Request ID tracking
- Structured logging
- Panic recovery
- JWT authentication

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License.

## Author

**Mahadeva Sankaram Telidevara**

- GitHub: [@tmsankaram](https://github.com/tmsankaram)

## Acknowledgments

- Built with [Gin](https://github.com/gin-gonic/gin)
- ORM by [GORM](https://gorm.io/)
- JWT implementation using [golang-jwt](https://github.com/golang-jwt/jwt)
