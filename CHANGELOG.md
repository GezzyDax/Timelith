# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [1.0.0] - 2024-01-01

### Added
- Initial release of Timelith self-hosted Telegram message scheduler
- Rails Web UI with admin panel
- Go backend with Telegram client manager
- Docker Compose orchestration
- PostgreSQL database integration
- Redis for job queue and caching
- Telegram account management (add, authorize, disconnect)
- Message template system
- Channel management (manual and sync from Telegram)
- Schedule system with three types:
  - Interval-based scheduling
  - Cron expression scheduling
  - One-time scheduled messages
- Send logs with detailed status tracking
- Dashboard with statistics and recent activity
- User authentication system
- API endpoints for Rails ↔ Go communication
- Automatic message dispatching via scheduler
- Rate limiting and retry logic
- Session encryption for Telegram accounts
- Comprehensive documentation (README, INSTALL guide)
- Makefile for common operations

### Features
- Multi-account support
- Multi-channel broadcasting
- Message templating
- Flexible scheduling options
- Real-time status monitoring
- Activity logs and error tracking
- Responsive Bootstrap UI
- Docker containerization
- Self-hosted solution (no external dependencies)
- Secure API key authentication
- Database migrations
- Seed data for initial setup

### Documentation
- README.md with usage instructions
- INSTALL.md with step-by-step installation guide
- .env.example with all configuration options
- Inline code documentation
- API documentation in README

### Infrastructure
- Docker Compose with 4 services (Rails, Go, PostgreSQL, Redis)
- Health checks for all services
- Volume persistence for data
- Network isolation
- Environment-based configuration
- Graceful shutdown handling
- Automatic service recovery
- Log aggregation

## [Unreleased]

### Planned Features
- Web UI improvements
- Media message support (photos, videos, documents)
- Inline keyboard buttons
- Advanced message formatting (Markdown, HTML)
- Multi-account load balancing
- Per-account rate limiting
- Webhook support
- Template variables and placeholders
- Scheduled message editing/deletion
- Channel analytics and statistics
- 2FA authentication support
- Proxy support for Telegram connections
- Backup and restore functionality
- Multi-language support
- Dark mode UI theme
- Mobile-responsive improvements
- Export/import schedules
- Bulk operations
- Search and filtering
- Advanced cron builder UI
