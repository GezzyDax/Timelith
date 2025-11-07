# Timelith - Self-Hosted Telegram Account Manager

A self-hosted service for managing Telegram accounts, message templates, and scheduling broadcasts through a convenient web dashboard.

## 🚀 Features

- **Multi-Account Management**: Manage multiple Telegram accounts
- **Message Templates**: Create reusable message templates with variables
- **Channel Management**: Organize target channels, groups, and users
- **Smart Scheduling**: Cron-based scheduling with timezone support
- **Job Queue**: Background job processing with retry logic
- **Real-time Logs**: Monitor all message deliveries and errors
- **Web Dashboard**: Modern Next.js UI with TailwindCSS
- **Self-Hosted**: Complete Docker setup, no external SaaS dependencies

## 🏗️ Architecture

```
┌─────────────────────────────────────────┐
│ Web UI (Next.js + TypeScript)          │
│ • Admin Dashboard                       │
│ • CRUD: Accounts, Templates, Schedules │
└──────────────┬──────────────────────────┘
               │ REST API
               ▼
┌─────────────────────────────────────────┐
│ Go Backend                              │
│ • Telegram Session Manager              │
│ • Job Scheduler (cron)                  │
│ • Message Dispatcher                    │
│ • Rate Limiter                          │
└──────────────┬──────────────────────────┘
               │
               ▼
┌─────────────────────────────────────────┐
│ Infrastructure                          │
│ • PostgreSQL - Database                 │
│ • Redis - Queue & Cache                 │
└─────────────────────────────────────────┘
```

## 📋 Prerequisites

- Docker & Docker Compose
- Telegram API credentials (get from https://my.telegram.org)

## 🚀 Quick Start

1. **Clone the repository**
   ```bash
   git clone https://github.com/GezzyDax/Timelith.git
   cd Timelith
   ```

2. **Configure environment**
   ```bash
   cp .env.example .env
   # Edit .env and fill in your Telegram API credentials
   ```

3. **Start the services**
   ```bash
   docker compose up -d
   ```

4. **Access the dashboard**
   - Web UI: http://localhost:3000
   - Backend API: http://localhost:8080

5. **Create admin user**
   First, you need to create an admin user. You can do this by sending a POST request to the backend:
   ```bash
   curl -X POST http://localhost:8080/api/auth/register \
     -H "Content-Type: application/json" \
     -d '{"username":"admin","password":"your_password"}'
   ```

## 📦 Services

| Service | Port | Description |
|---------|------|-------------|
| Web UI | 3000 | Next.js dashboard |
| Go Backend | 8080 | REST API & Scheduler |
| PostgreSQL | 5432 | Database |
| Redis | 6379 | Queue & Cache |

## 🔧 Configuration

### Environment Variables

| Variable | Required | Description |
|----------|----------|-------------|
| `POSTGRES_PASSWORD` | Yes | PostgreSQL password |
| `TELEGRAM_APP_ID` | Yes | Telegram API ID |
| `TELEGRAM_APP_HASH` | Yes | Telegram API Hash |
| `JWT_SECRET` | Yes | JWT signing secret |
| `ENCRYPTION_KEY` | Yes | 32-byte key for session encryption |
| `ENVIRONMENT` | No | `production` or `development` |
| `NEXT_PUBLIC_API_URL` | No | Backend API URL |

### Getting Telegram API Credentials

1. Go to https://my.telegram.org
2. Log in with your phone number
3. Go to "API development tools"
4. Create a new application
5. Copy your `api_id` and `api_hash`

## 📚 Usage

### Adding a Telegram Account

1. Go to **Accounts** page
2. Click "Add Account"
3. Enter phone number (with country code, e.g., +1234567890)
4. You'll receive a code via Telegram
5. Enter the code to complete authentication

### Creating Message Templates

1. Go to **Templates** page
2. Click "Add Template"
3. Enter template name and message content
4. Save the template

### Setting Up a Schedule

1. Go to **Schedules** page
2. Click "Add Schedule"
3. Select account, template, and target channel
4. Set cron expression (e.g., `0 9 * * *` for daily at 9 AM)
5. Choose timezone
6. Save the schedule

### Cron Expression Examples

| Expression | Description |
|------------|-------------|
| `0 9 * * *` | Daily at 9:00 AM |
| `0 */2 * * *` | Every 2 hours |
| `0 9 * * 1` | Every Monday at 9:00 AM |
| `0 9,18 * * *` | Daily at 9:00 AM and 6:00 PM |

## 🛠️ Development

### Local Development (without Docker)

**Backend:**
```bash
cd go-backend
cp .env.example .env
go mod download
go run cmd/server/main.go
```

**Frontend:**
```bash
cd web-ui
npm install
npm run dev
```

### Building from Source

```bash
# Build backend
cd go-backend
make build

# Build frontend
cd web-ui
npm run build
```

## 🔒 Security

- All Telegram sessions are encrypted with AES-256 before storage
- JWT authentication for web dashboard
- API key support for external integrations
- No plaintext credentials in database
- All services run in isolated Docker network

## 📊 Monitoring

- View job execution logs in **Logs** page
- Check system status in **Dashboard**
- Monitor active schedules and their next run times
- Track account statuses and last login times

## 🐛 Troubleshooting

### Backend won't start
- Check if PostgreSQL is running: `docker compose ps`
- Verify environment variables in `.env`
- Check logs: `docker compose logs go-backend`

### Telegram authentication fails
- Verify `TELEGRAM_APP_ID` and `TELEGRAM_APP_HASH`
- Ensure phone number format includes country code
- Check if you're receiving the code in Telegram

### Messages not sending
- Check account status in **Accounts** page
- Verify schedule is active
- Check **Logs** for error messages
- Ensure rate limits aren't exceeded

## 🤝 Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## 📝 License

This project is licensed under the MIT License.

## ⚠️ Disclaimer

This tool is for legitimate use only. Ensure you comply with Telegram's Terms of Service and local laws. The authors are not responsible for misuse of this software.

## 🙏 Credits

Built with:
- [Go](https://golang.org/) - Backend
- [Next.js](https://nextjs.org/) - Frontend
- [gotd/td](https://github.com/gotd/td) - Telegram client
- [Fiber](https://gofiber.io/) - Web framework
- [PostgreSQL](https://www.postgresql.org/) - Database
- [Redis](https://redis.io/) - Queue & Cache

## 📧 Support

For issues and questions, please use the GitHub Issues page.
