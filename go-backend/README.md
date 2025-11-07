# Timelith Go Backend

Backend service for Timelith Telegram account manager.

## Features

- Telegram Session Management (gotd/td)
- Job Scheduler with cron support
- Message Dispatcher with retry logic
- REST API (Fiber)
- PostgreSQL database
- Redis queue
- **Interactive Setup Wizard** - First-run configuration assistant

## Quick Start

### First Time Setup

При первом запуске автоматически запустится интерактивный мастер установки:

```bash
# Убедитесь, что PostgreSQL запущен
docker-compose up -d postgres

# Запустите приложение
go run cmd/server/main.go
```

Мастер установки попросит вас указать:

1. **Telegram API credentials** (получите на https://my.telegram.org)
   - App ID
   - App Hash

2. **Настройки сервера**
   - Порт (по умолчанию: 8080)
   - Окружение (production/development)

3. **Настройки базы данных**
   - Пароль PostgreSQL

4. **Первый администратор**
   - Логин
   - Пароль

Ключи безопасности (JWT_SECRET и ENCRYPTION_KEY) генерируются автоматически.

После завершения установки все настройки сохраняются в файл `.env`, и создается первый администратор в базе данных.

### Последующие запуски

После первой настройки просто запустите:

```bash
go run cmd/server/main.go
```

## Development

```bash
# Install dependencies
go mod download

# Build
make build

# Run
./bin/server
```

## Environment Variables

После первого запуска конфигурация сохраняется в `.env`.
Пример см. в `.env.example`.

## API Endpoints

- `GET /api/health` - Health check
- `POST /api/auth/login` - Login
- `POST /api/auth/register` - Register
- `GET /api/accounts` - List accounts
- `POST /api/accounts` - Create account
- ... (see internal/api/router.go for full list)
