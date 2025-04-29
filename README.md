# Product PostgreSQL Module

A module for managing PostgreSQL database connections.

## Features

- PostgreSQL connection management
- Environment-based configuration
- Connection pooling
- Graceful shutdown support

## Configuration

The module uses environment variables for database configuration:

```bash
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=product_db
DB_SSL_MODE=disable
```

## Usage

1. Import the module in your project:
```go
import "github.com/yourusername/product-postgres"
```

2. Initialize the database connection:
```go
cfg := config.NewDatabaseConfig()
db, err := database.NewDatabase(cfg)
if err != nil {
    log.Fatal(err)
}
defer db.Close()

// Use the database connection
// db.Exec(...)
// db.Raw(...)
// etc.
```

## Dependencies

- Go 1.21 or later
- GORM
- PostgreSQL driver for GORM

## License

MIT 