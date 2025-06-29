# User Management Microservice

Simple Go microservice for user management with PostgreSQL.

## Run the API

1. Copy environment file:
   ```bash
   cp .env.example .env
   ```

2. Run with Docker:
   ```bash
   make run
   ```

## Testing

Run integration tests:
```bash
make test-integration
```

## API Endpoints

- `POST /save` - Create user
- `GET /{id}` - Get user by ID