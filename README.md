# Timelith - Self-Hosted Telegram Message Scheduler

Timelith is a self-hosted solution for scheduling and automating Telegram message delivery across multiple channels and accounts.

## 🏗️ Architecture

The project consists of 4 main components running in Docker containers:

- **Rails Web UI** (Port 3000) - Admin panel for managing accounts, templates, schedules
- **Go Backend** (Port 8080) - Telegram client manager, scheduler, and message dispatcher
- **PostgreSQL** - Database for all entities
- **Redis** - Job queue and caching

## 🚀 Quick Start

### Prerequisites

- Docker and Docker Compose
- Telegram API credentials (get them from https://my.telegram.org)

### Installation

1. Clone the repository:
```bash
git clone https://github.com/yourusername/timelith.git
cd timelith
```

2. Copy the example environment file and configure it:
```bash
cp .env.example .env
nano .env
```

3. Set required environment variables in `.env`:
```env
# Get these from https://my.telegram.org
TELEGRAM_APP_ID=your_app_id
TELEGRAM_APP_HASH=your_app_hash

# Generate secure random keys
SECRET_KEY_BASE=$(openssl rand -hex 64)
GO_API_KEY=$(openssl rand -hex 32)
SESSION_ENCRYPTION_KEY=$(openssl rand -hex 32)

# Database (change password in production!)
POSTGRES_PASSWORD=your_secure_password
```

4. Build and start the services:
```bash
docker-compose up -d
```

5. Wait for services to initialize (check logs):
```bash
docker-compose logs -f
```

6. Access the web interface:
```
http://localhost:3000
```

### Initial Setup

1. Create the first admin user:
```bash
docker-compose exec rails-app bundle exec rails console
```

```ruby
User.create!(
  email: 'admin@example.com',
  password: 'your_secure_password',
  role: 'admin'
)
```

2. Login at http://localhost:3000/login

## 📖 Usage

### 1. Add Telegram Account

1. Go to **Telegram Accounts** → **New Account**
2. Enter your phone number (with country code, e.g., +1234567890)
3. Wait for verification code in Telegram
4. Enter the code to authorize

### 2. Create Message Template

1. Go to **Message Templates** → **New Template**
2. Enter template name and content
3. Optionally add media (photo, video, document)
4. Configure parse mode (Markdown/HTML) if needed

### 3. Add Channels

**Option A: Manual Entry**
1. Go to **Channels** → **New Channel**
2. Enter channel details manually

**Option B: Sync from Telegram**
1. Go to **Channels** → **Sync from Telegram**
2. System will fetch all accessible channels

### 4. Create Schedule

1. Go to **Schedules** → **New Schedule**
2. Select:
   - Telegram account to use
   - Message template
   - Target channels
   - Schedule type:
     - **Interval**: Every X minutes
     - **Cron**: Custom cron expression
     - **Once**: One-time at specific date/time
3. Activate the schedule

### 5. Monitor Logs

- View sending history in **Send Logs**
- Filter by status (Sent/Failed)
- See detailed error messages for failed sends

## 🔧 Configuration

### Environment Variables

| Variable | Description | Required |
|----------|-------------|----------|
| `TELEGRAM_APP_ID` | Telegram API App ID | Yes |
| `TELEGRAM_APP_HASH` | Telegram API Hash | Yes |
| `DATABASE_URL` | PostgreSQL connection string | Yes |
| `REDIS_URL` | Redis connection string | Yes |
| `GO_API_KEY` | API key for Rails ↔ Go communication | Yes |
| `SECRET_KEY_BASE` | Rails secret key | Yes |
| `SESSION_ENCRYPTION_KEY` | Key for encrypting Telegram sessions | Yes |

### Schedule Types

**Interval**
```
Run every X minutes
Example: 30 minutes
```

**Cron**
```
Unix cron expression
Example: "0 9 * * *" (daily at 9:00 AM)
```

**Once**
```
Single execution at specific time
Example: 2024-12-31 23:59:59
```

## 🛠️ Development

### Project Structure

```
timelith/
├── docker-compose.yml       # Main orchestration file
├── .env                     # Environment configuration
├── rails-app/              # Rails Web UI
│   ├── app/
│   │   ├── controllers/    # Request handlers
│   │   ├── models/         # Database models
│   │   ├── views/          # HTML templates
│   │   └── services/       # Business logic
│   ├── config/             # Rails configuration
│   └── db/migrate/         # Database migrations
├── go-backend/             # Go Backend
│   ├── cmd/server/         # Main entry point
│   ├── internal/
│   │   ├── api/            # HTTP API handlers
│   │   ├── telegram/       # Telegram client manager
│   │   ├── scheduler/      # Job scheduler
│   │   └── database/       # Database connection
│   └── pkg/                # Shared packages
└── docker/                 # Docker configurations
```

### Running Locally

**Rails App:**
```bash
cd rails-app
bundle install
bundle exec rails db:migrate
bundle exec rails server
```

**Go Backend:**
```bash
cd go-backend
go mod download
go run cmd/server/main.go
```

## 🔒 Security

### Important Security Notes

1. **Change default passwords** in `.env`
2. **Use HTTPS** in production (setup reverse proxy with nginx/Caddy)
3. **Restrict access** to the admin panel (firewall, VPN)
4. **Backup database** regularly
5. **Rotate API keys** periodically
6. **Session encryption** - never commit `SESSION_ENCRYPTION_KEY`

### Session Storage

Telegram sessions are stored encrypted in the database. The `SESSION_ENCRYPTION_KEY` is used to encrypt/decrypt session data. **Keep this key secure and never commit it to version control.**

## 📊 Monitoring

### Check Service Status

```bash
docker-compose ps
```

### View Logs

```bash
# All services
docker-compose logs -f

# Specific service
docker-compose logs -f rails-app
docker-compose logs -f go-backend
```

### Database Access

```bash
docker-compose exec postgres psql -U timelith -d timelith_production
```

### Redis CLI

```bash
docker-compose exec redis redis-cli
```

## 🚨 Troubleshooting

### Rails app won't start

Check database connection and migrations:
```bash
docker-compose exec rails-app bundle exec rails db:migrate:status
```

### Go backend connection errors

Verify environment variables are set:
```bash
docker-compose exec go-backend env | grep TELEGRAM
```

### Messages not sending

1. Check account status in Telegram Accounts
2. Verify schedule is activated
3. Check send logs for error messages
4. Review go-backend logs

### Database errors

Reset database (⚠️ destroys all data):
```bash
docker-compose down -v
docker-compose up -d
docker-compose exec rails-app bundle exec rails db:migrate
```

## 📝 API Documentation

### Go Backend API

Base URL: `http://localhost:8080/api/v1`

All requests require `X-API-Key` header.

**Send Auth Code**
```bash
POST /telegram/auth/send_code
Content-Type: application/json

{
  "phone_number": "+1234567890",
  "account_id": 1
}
```

**Verify Code**
```bash
POST /telegram/auth/verify_code
Content-Type: application/json

{
  "account_id": 1,
  "code": "12345"
}
```

**Get Schedules**
```bash
GET /schedules
```

## 🤝 Contributing

Contributions are welcome! Please:

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Submit a pull request

## 📜 License

This project is licensed under the MIT License.

## 🆘 Support

For issues and questions:
- Open an issue on GitHub
- Check existing issues for solutions

## 🔮 Roadmap

- [ ] Web interface improvements (Bootstrap UI)
- [ ] Support for media messages (photos, videos)
- [ ] Inline keyboard buttons
- [ ] Message formatting (Markdown, HTML)
- [ ] Multi-account load balancing
- [ ] Rate limiting per account
- [ ] Webhook support
- [ ] Message templates with variables
- [ ] Scheduled message editing/deletion
- [ ] Channel analytics
- [ ] 2FA authentication support
- [ ] Proxy support for Telegram connections
- [ ] Backup/restore functionality

## ⚠️ Disclaimer

This software is provided as-is. Use responsibly and in accordance with Telegram's Terms of Service. The authors are not responsible for any misuse or violations of Telegram's policies.

## 🙏 Acknowledgments

- Built with [Ruby on Rails](https://rubyonrails.org/)
- Powered by [Go](https://golang.org/)
- Telegram MTProto via [gotd](https://github.com/gotd/td)
