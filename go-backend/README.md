# Timelith Go Backend

Backend service for Timelith Telegram account manager.

## Features

- Telegram Session Management (gotd/td)
- Job Scheduler with cron support
- Message Dispatcher with retry logic
- REST API (Fiber)
- PostgreSQL database
- Redis queue

## Development

```bash
# Install dependencies
go mod download

# Run migrations
go run cmd/server/main.go

# Build
make build

# Run
./bin/server
```

## Environment Variables

See `.env.example` in the root directory.

## API Endpoints

- `GET /api/health` - Health check
- `POST /api/auth/login` - Login
- `POST /api/auth/register` - Register
- `GET /api/accounts` - List accounts
- `POST /api/accounts` - Create account
- ... (see internal/api/router.go for full list)
