# Timelith Go Backend

Backend service for Timelith Telegram account manager.

## Features

- Telegram Session Management (gotd/td)
- Job Scheduler with cron support
- Message Dispatcher with retry logic
- REST API (Fiber)
- PostgreSQL database
- Redis queue
- **Web-Based Setup Wizard** - First-run configuration through browser interface

## Quick Start

### First Time Setup

При первом запуске автоматически запустится **веб-интерфейс мастера установки**:

```bash
# 1. Убедитесь, что PostgreSQL запущен
docker-compose up -d postgres

# 2. Запустите бэкенд
cd go-backend
go run cmd/server/main.go

# 3. В другом терминале запустите фронтенд
cd web-ui
npm install
npm run dev
```

После запуска откройте браузер и перейдите по адресу:

**http://localhost:3000**

Вы автоматически будете перенаправлены на веб-форму установки, где нужно указать:

**Шаг 1: Telegram API**
- App ID и App Hash (получите на https://my.telegram.org)

**Шаг 2: Сервер и База Данных**
- Порт сервера (по умолчанию: 8080)
- Окружение (production/development)
- Пароль PostgreSQL

**Шаг 3: Первый Администратор**
- Логин (минимум 3 символа)
- Пароль (минимум 6 символов)

Ключи безопасности (JWT_SECRET и ENCRYPTION_KEY) генерируются автоматически.

После завершения установки:
- Все настройки сохраняются в `.env`
- Создается первый администратор в БД
- **Перезапустите сервер** для применения изменений

### Последующие запуски

После первой настройки:

```bash
# Бэкенд
cd go-backend
go run cmd/server/main.go

# Фронтенд (в другом терминале)
cd web-ui
npm run dev
```

Откройте http://localhost:3000 и войдите с учетными данными администратора.

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
