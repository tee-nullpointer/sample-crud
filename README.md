# Sample CRUD

This is a simple example of a CRUD (Create, Read, Update, Delete) application using Gin Gonic and Postgres.

## Tech Stack
- Go 1.25+
- Gin Gonic
- GORM
- Zap Logger
- Postgres
- Redis
- Clean Architecture

## Quick Start

1. **Set up environment variables:**
   ```bash
   cp .env.example .env
   # Edit .env with your configuration
   ```
   
2. **Install dependencies:**
   ```bash
    go mod tidy
    ```

3. **Run the application:**
   ```bash
   go run cmd/main.go
   ```