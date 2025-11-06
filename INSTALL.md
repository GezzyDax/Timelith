# Installation Guide for Timelith

This guide will help you install and set up Timelith on your server.

## System Requirements

- Docker 20.10+
- Docker Compose 2.0+
- 2GB RAM minimum (4GB recommended)
- 10GB disk space
- Ubuntu 20.04+ / Debian 11+ / CentOS 8+ (or any Linux with Docker support)

## Step-by-Step Installation

### 1. Install Docker and Docker Compose

**Ubuntu/Debian:**
```bash
# Update package index
sudo apt-get update

# Install Docker
curl -fsSL https://get.docker.com -o get-docker.sh
sudo sh get-docker.sh

# Install Docker Compose
sudo apt-get install docker-compose-plugin

# Add your user to docker group (optional, to run without sudo)
sudo usermod -aG docker $USER
newgrp docker
```

**CentOS/RHEL:**
```bash
# Install Docker
sudo yum install -y yum-utils
sudo yum-config-manager --add-repo https://download.docker.com/linux/centos/docker-ce.repo
sudo yum install docker-ce docker-ce-cli containerd.io docker-compose-plugin

# Start Docker
sudo systemctl start docker
sudo systemctl enable docker
```

### 2. Get Telegram API Credentials

1. Go to https://my.telegram.org
2. Log in with your phone number
3. Click on "API Development Tools"
4. Create a new application:
   - App title: Timelith
   - Short name: timelith
   - Platform: Other
5. Save your `api_id` and `api_hash`

### 3. Clone and Configure

```bash
# Clone the repository
git clone https://github.com/yourusername/timelith.git
cd timelith

# Copy environment file
cp .env.example .env

# Generate secure keys
export SECRET_KEY_BASE=$(openssl rand -hex 64)
export GO_API_KEY=$(openssl rand -hex 32)
export SESSION_ENCRYPTION_KEY=$(openssl rand -hex 32)

# Edit .env file with your favorite editor
nano .env
```

**Required .env configuration:**
```env
# Telegram API (from step 2)
TELEGRAM_APP_ID=your_app_id_here
TELEGRAM_APP_HASH=your_app_hash_here

# Security keys (generated above)
SECRET_KEY_BASE=your_generated_secret_key_base
GO_API_KEY=your_generated_api_key
SESSION_ENCRYPTION_KEY=your_generated_encryption_key

# Database (change in production!)
POSTGRES_USER=timelith
POSTGRES_PASSWORD=your_secure_password_here
POSTGRES_DB=timelith_production

# Optional: Change ports if needed
APP_PORT=3000
GO_BACKEND_PORT=8080
```

### 4. Build and Start Services

```bash
# Build containers
docker-compose build

# Start services
docker-compose up -d

# Wait for services to start (about 30-60 seconds)
docker-compose logs -f
```

Press Ctrl+C when you see:
```
rails-app    | * Listening on http://0.0.0.0:3000
go-backend   | Server starting on port 8080
```

### 5. Initialize Database

```bash
# Run migrations
docker-compose exec rails-app bundle exec rails db:migrate

# Create seed data (admin user)
docker-compose exec rails-app bundle exec rails db:seed
```

### 6. Verify Installation

```bash
# Check all services are running
docker-compose ps

# Should show all services as "Up"
```

Access the application:
- Web UI: http://your-server-ip:3000
- Login: `admin@example.com` / `admin123`

**⚠️ Change the default password immediately after first login!**

## Production Setup (Optional but Recommended)

### 1. Setup Reverse Proxy with Nginx

Install Nginx:
```bash
sudo apt-get install nginx
```

Create Nginx config:
```bash
sudo nano /etc/nginx/sites-available/timelith
```

```nginx
server {
    listen 80;
    server_name your-domain.com;

    location / {
        proxy_pass http://localhost:3000;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}
```

Enable site:
```bash
sudo ln -s /etc/nginx/sites-available/timelith /etc/nginx/sites-enabled/
sudo nginx -t
sudo systemctl restart nginx
```

### 2. Setup SSL with Let's Encrypt

```bash
# Install Certbot
sudo apt-get install certbot python3-certbot-nginx

# Get certificate
sudo certbot --nginx -d your-domain.com

# Certificate will auto-renew
```

### 3. Setup Firewall

```bash
# Allow SSH
sudo ufw allow ssh

# Allow HTTP/HTTPS
sudo ufw allow 80/tcp
sudo ufw allow 443/tcp

# Enable firewall
sudo ufw enable
```

### 4. Setup Automatic Backups

Create backup script:
```bash
sudo nano /usr/local/bin/backup-timelith.sh
```

```bash
#!/bin/bash
BACKUP_DIR=/var/backups/timelith
DATE=$(date +%Y%m%d_%H%M%S)

mkdir -p $BACKUP_DIR

# Backup database
docker-compose -f /path/to/timelith/docker-compose.yml exec -T postgres \
  pg_dump -U timelith timelith_production > $BACKUP_DIR/db_$DATE.sql

# Compress
gzip $BACKUP_DIR/db_$DATE.sql

# Keep only last 7 days
find $BACKUP_DIR -name "*.sql.gz" -mtime +7 -delete
```

Make executable and add to cron:
```bash
sudo chmod +x /usr/local/bin/backup-timelith.sh
sudo crontab -e
```

Add daily backup at 2 AM:
```
0 2 * * * /usr/local/bin/backup-timelith.sh
```

## Troubleshooting

### Services won't start

Check logs:
```bash
docker-compose logs
```

### Port already in use

Change ports in `.env`:
```env
APP_PORT=3001
GO_BACKEND_PORT=8081
```

Then restart:
```bash
docker-compose down
docker-compose up -d
```

### Database connection errors

Reset database:
```bash
docker-compose down -v
docker-compose up -d
docker-compose exec rails-app bundle exec rails db:migrate db:seed
```

### Out of memory

Increase Docker memory limit or upgrade server RAM.

## Maintenance

### Update Timelith

```bash
cd timelith
git pull
docker-compose build
docker-compose up -d
docker-compose exec rails-app bundle exec rails db:migrate
```

### View logs

```bash
docker-compose logs -f
```

### Restart services

```bash
docker-compose restart
```

### Stop services

```bash
docker-compose down
```

## Support

For issues:
- Check existing GitHub issues
- Create new issue with:
  - System information
  - Error logs
  - Steps to reproduce

## Security Checklist

- [ ] Changed default admin password
- [ ] Set strong passwords in `.env`
- [ ] Enabled firewall
- [ ] Setup SSL certificate
- [ ] Configured automatic backups
- [ ] Restricted SSH access
- [ ] Kept system updated

Congratulations! Timelith is now installed and ready to use. 🎉
