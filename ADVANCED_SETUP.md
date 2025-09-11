# Advanced Music Player - Go Backend Setup

This document explains how to set up the advanced version of the music player with Go Fiber, PostgreSQL, and other advanced features.

## Prerequisites

- Go 1.19 or higher
- PostgreSQL 13 or higher
- Docker (optional, for containerization)
- Node.js 16+ (for frontend)

## Project Initialization

1. Create a new directory for the advanced project:
```bash
mkdir advanced-music-player
cd advanced-music-player
```

2. Initialize the Go module:
```bash
go mod init github.com/yourusername/advanced-music-player
```

3. Install required dependencies:
```bash
go get github.com/gofiber/fiber/v2
go get gorm.io/gorm
go get gorm.io/driver/postgres
go get github.com/golang-jwt/jwt/v5
go get github.com/go-redis/redis/v8
go get github.com/disintegration/imaging
```

## Directory Structure

```
advanced-music-player/
├── cmd/
│   └── server/
│       └── main.go
├── internal/
│   ├── api/
│   │   └── handlers.go
│   ├── auth/
│   │   └── auth.go
│   ├── database/
│   │   └── database.go
│   ├── models/
│   │   ├── user.go
│   │   ├── song.go
│   │   └── playlist.go
│   ├── storage/
│   │   └── storage.go
│   └── utils/
│       └── utils.go
├── configs/
│   └── config.go
├── migrations/
├── go.mod
└── go.sum
```

## Configuration

Create a `.env` file in the root directory:

```env
# Server Configuration
SERVER_PORT=8080
SERVER_HOST=localhost

# Database Configuration
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=yourpassword
DB_NAME=musicplayer
DB_SSL_MODE=disable

# JWT Configuration
JWT_SECRET=your_jwt_secret_key
JWT_EXPIRATION=24h

# Redis Configuration
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=
REDIS_DB=0

# Storage Configuration
STORAGE_PATH=./songs
STORAGE_MAX_SIZE=50MB
```

## Database Setup

1. Create the PostgreSQL database:
```sql
CREATE DATABASE musicplayer;
```

2. Run the initial migrations (example):
```sql
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(50) UNIQUE NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE songs (
    id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    artist VARCHAR(255) NOT NULL,
    album VARCHAR(255),
    genre VARCHAR(100),
    duration INTEGER,
    file_path VARCHAR(500) NOT NULL,
    cover_art_url VARCHAR(500),
    uploaded_by INTEGER REFERENCES users(id),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

## Running the Application

1. Start the PostgreSQL and Redis services:
```bash
# If using Docker
docker-compose up -d postgres redis
```

2. Run the Go application:
```bash
go run cmd/server/main.go
```

## API Endpoints

The advanced version will include these endpoints:

### Authentication
- `POST /api/v1/auth/register` - User registration
- `POST /api/v1/auth/login` - User login
- `POST /api/v1/auth/logout` - User logout
- `GET /api/v1/auth/me` - Get current user

### Songs
- `GET /api/v1/songs` - List songs
- `GET /api/v1/songs/:id` - Get song details
- `POST /api/v1/songs` - Upload song
- `PUT /api/v1/songs/:id` - Update song
- `DELETE /api/v1/songs/:id` - Delete song

### Playlists
- `GET /api/v1/playlists` - List playlists
- `GET /api/v1/playlists/:id` - Get playlist
- `POST /api/v1/playlists` - Create playlist
- `PUT /api/v1/playlists/:id` - Update playlist
- `DELETE /api/v1/playlists/:id` - Delete playlist

### Streaming
- `GET /api/v1/stream/:id` - Stream song

## Docker Setup

Create a `docker-compose.yml` file:

```yaml
version: '3.8'

services:
  postgres:
    image: postgres:13
    environment:
      POSTGRES_DB: musicplayer
      POSTGRES_USER: postgres
      POSTGRES_PASSWORD: yourpassword
    ports:
      - "5432:5432"
    volumes:
      - postgres_data:/var/lib/postgresql/data

  redis:
    image: redis:7-alpine
    ports:
      - "6379:6379"
    volumes:
      - redis_data:/data

  app:
    build: .
    ports:
      - "8080:8080"
    environment:
      - DB_HOST=postgres
      - REDIS_HOST=redis
    depends_on:
      - postgres
      - redis

volumes:
  postgres_data:
  redis_data:
```

## Frontend Integration

The frontend will be a separate Next.js application that connects to this backend API.

This setup provides a foundation for building a full-featured music streaming service with modern technologies.