#!/bin/bash

# Site-SB Quick Install Script
# One-command installation for Ubuntu servers

set -e

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

# Configuration
REPO_URL="https://github.com/Fox1N69/site-sb.git" # Update with your repo
PROJECT_DIR="/opt/site-sb"

log_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

log_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

log_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Check if running as root
if [[ $EUID -eq 0 ]]; then
    log_error "This script should not be run as root"
    log_info "Please run as a regular user with sudo privileges"
    exit 1
fi

# Welcome message
echo -e "${GREEN}"
echo "=================================="
echo "    Site-SB Quick Installer"
echo "=================================="
echo -e "${NC}"
echo "This script will:"
echo "• Install Docker and dependencies"
echo "• Clone the Site-SB project"
echo "• Set up the environment"
echo "• Configure Nginx and SSL"
echo "• Deploy the application"
echo ""

# Get user input
read -p "Enter your domain name (e.g., example.com): " DOMAIN
read -p "Enter your email for SSL certificate: " EMAIL
read -s -p "Enter a strong database password: " DB_PASSWORD
echo ""
read -s -p "Enter a strong JWT secret (min 32 chars): " JWT_SECRET
echo ""
echo ""

log_info "Starting installation process..."

# Update system
log_info "Updating system packages..."
sudo apt update && sudo apt upgrade -y

# Install Docker
if ! command -v docker &>/dev/null; then
    log_info "Installing Docker..."
    curl -fsSL https://get.docker.com -o get-docker.sh
    sudo sh get-docker.sh
    sudo usermod -aG docker $USER
    rm get-docker.sh
    log_success "Docker installed"
else
    log_success "Docker already installed"
fi

# Install other dependencies
log_info "Installing additional dependencies..."
sudo apt install -y nginx certbot python3-certbot-nginx git curl

# Create project directory
log_info "Setting up project..."
sudo mkdir -p $PROJECT_DIR
sudo chown $USER:$USER $PROJECT_DIR

# Clone repository
log_info "Cloning repository..."
git clone $REPO_URL $PROJECT_DIR
cd $PROJECT_DIR

# Create environment file
log_info "Creating environment configuration..."
cat >.env <<EOF
DB_PASSWORD=$DB_PASSWORD
JWT_SECRET=$JWT_SECRET
NODE_ENV=production
DOMAIN=$DOMAIN
EMAIL=$EMAIL
EOF

# Create production config
log_info "Creating production configuration..."
cp backend/config/config.json backend/config/config.prod.json

# Update production config
sed -i "s/\"mode\": \"dev\"/\"mode\": \"release\"/g" backend/config/config.prod.json
sed -i "s/\"host\": \"localhost\"/\"host\": \"postgres\"/g" backend/config/config.prod.json
sed -i "s/\"pass\": \"8008\"/\"pass\": \"$DB_PASSWORD\"/g" backend/config/config.prod.json

# Setup firewall
log_info "Configuring firewall..."
sudo ufw --force enable
sudo ufw allow ssh
sudo ufw allow 80
sudo ufw allow 443

# Configure Nginx
log_info "Configuring Nginx..."
sudo tee /etc/nginx/sites-available/site-sb >/dev/null <<EOF
server {
    listen 80;
    server_name $DOMAIN www.$DOMAIN;
    server_tokens off;

    location /.well-known/acme-challenge/ {
        root /var/www/certbot;
        allow all;
    }

    location / {
        proxy_pass http://127.0.0.1:4321;
        proxy_set_header Host \$host;
        proxy_set_header X-Real-IP \$remote_addr;
        proxy_set_header X-Forwarded-For \$proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto \$scheme;
    }
}
EOF

# Enable site
sudo ln -sf /etc/nginx/sites-available/site-sb /etc/nginx/sites-enabled/
sudo rm -f /etc/nginx/sites-enabled/default
sudo nginx -t
sudo systemctl restart nginx

# Deploy application
log_info "Deploying application..."
mkdir -p backups logs

# Start services
docker compose -f docker-compose.prod.yaml up --build -d

# Wait for services
log_info "Waiting for services to start..."
sleep 45

# Check if services are running
if docker compose -f docker-compose.prod.yaml ps | grep -q "Up"; then
    log_success "Application deployed successfully"
else
    log_error "Application deployment failed"
    docker compose -f docker-compose.prod.yaml logs
    exit 1
fi

# Setup SSL
log_info "Setting up SSL certificate..."
sudo mkdir -p /var/www/certbot
sudo certbot --nginx -d $DOMAIN -d www.$DOMAIN --email $EMAIL --agree-tos --non-interactive

# Update Nginx config for HTTPS
sudo tee /etc/nginx/sites-available/site-sb >/dev/null <<EOF
server {
    listen 80;
    server_name $DOMAIN www.$DOMAIN;
    server_tokens off;

    location /.well-known/acme-challenge/ {
        root /var/www/certbot;
        allow all;
    }

    location / {
        return 301 https://\$host\$request_uri;
    }
}

server {
    listen 443 ssl http2;
    server_name $DOMAIN www.$DOMAIN;
    server_tokens off;

    ssl_certificate /etc/letsencrypt/live/$DOMAIN/fullchain.pem;
    ssl_certificate_key /etc/letsencrypt/live/$DOMAIN/privkey.pem;
    include /etc/letsencrypt/options-ssl-nginx.conf;
    ssl_dhparam /etc/letsencrypt/ssl-dhparams.pem;

    location / {
        proxy_pass http://127.0.0.1:4321;
        proxy_set_header Host \$host;
        proxy_set_header X-Real-IP \$remote_addr;
        proxy_set_header X-Forwarded-For \$proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto \$scheme;
        proxy_set_header X-Forwarded-Host \$host;
        proxy_set_header X-Forwarded-Port \$server_port;
    }

    location /api/ {
        proxy_pass http://127.0.0.1:4000/;
        proxy_set_header Host \$host;
        proxy_set_header X-Real-IP \$remote_addr;
        proxy_set_header X-Forwarded-For \$proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto \$scheme;
    }

    location ~* \.(js|css|png|jpg|jpeg|gif|ico|svg|woff|woff2|ttf|eot)$ {
        proxy_pass http://127.0.0.1:4321;
        expires 1y;
        add_header Cache-Control "public, immutable";
    }
}
EOF

sudo nginx -t
sudo systemctl reload nginx

# Setup monitoring scripts
log_info "Setting up monitoring and maintenance..."

# Create status script
sudo tee /usr/local/bin/site-sb-status.sh >/dev/null <<'EOF'
#!/bin/bash
echo "=== Site-SB Status Report ==="
echo "Date: $(date)"
echo ""
echo "=== System Resources ==="
free -h
echo ""
df -h /
echo ""
echo "=== Docker Containers ==="
cd /opt/site-sb
docker compose -f docker-compose.prod.yaml ps
EOF

# Create backup script
sudo tee /usr/local/bin/site-sb-backup.sh >/dev/null <<'EOF'
#!/bin/bash
BACKUP_DIR="/opt/site-sb/backups"
DATE=$(date +%Y%m%d_%H%M%S)
BACKUP_FILE="sitesb_backup_$DATE.sql"

cd /opt/site-sb
docker compose -f docker-compose.prod.yaml exec -T postgres pg_dump -U postgres sitesb > "$BACKUP_DIR/$BACKUP_FILE"
gzip "$BACKUP_DIR/$BACKUP_FILE"
find "$BACKUP_DIR" -name "*.gz" -mtime +30 -delete
EOF

sudo chmod +x /usr/local/bin/site-sb-status.sh
sudo chmod +x /usr/local/bin/site-sb-backup.sh

# Setup cron jobs
(
    crontab -l 2>/dev/null
    echo "0 2 * * * /usr/local/bin/site-sb-backup.sh"
) | crontab -
(
    sudo crontab -l 2>/dev/null
    echo "0 12 * * * /usr/bin/certbot renew --quiet"
) | sudo crontab -

# Create aliases
cat >>~/.bashrc <<'EOF'

# Site-SB aliases
alias site-sb-start='cd /opt/site-sb && docker compose -f docker-compose.prod.yaml up -d'
alias site-sb-stop='cd /opt/site-sb && docker compose -f docker-compose.prod.yaml down'
alias site-sb-restart='cd /opt/site-sb && docker compose -f docker-compose.prod.yaml restart'
alias site-sb-logs='cd /opt/site-sb && docker compose -f docker-compose.prod.yaml logs -f'
alias site-sb-status='/usr/local/bin/site-sb-status.sh'
alias site-sb-backup='/usr/local/bin/site-sb-backup.sh'
EOF

# Final verification
log_info "Performing final verification..."
sleep 10

if curl -f https://$DOMAIN >/dev/null 2>&1; then
    log_success "HTTPS site is accessible!"
else
    log_warning "HTTPS not accessible yet, checking HTTP..."
    if curl -f http://$DOMAIN >/dev/null 2>&1; then
        log_warning "HTTP accessible, HTTPS setup may need more time"
    else
        log_error "Site not accessible"
    fi
fi

# Installation complete
echo ""
echo -e "${GREEN}=================================="
echo "   Installation Complete! 🚀"
echo -e "==================================${NC}"
echo ""
echo "Your Site-SB application is now running!"
echo ""
echo "📱 Website: https://$DOMAIN"
echo "🔧 Admin: https://$DOMAIN/admin"
echo ""
echo "Useful commands:"
echo "• site-sb-status  - Check application status"
echo "• site-sb-logs    - View application logs"
echo "• site-sb-restart - Restart services"
echo "• site-sb-backup  - Create database backup"
echo ""
echo "⚠️  Important notes:"
echo "• Please logout and login again for Docker group to take effect"
echo "• Change default passwords in admin panel"
echo "• Review firewall settings for your specific needs"
echo "• Monitor application logs regularly"
echo ""
log_success "Installation completed successfully!"

# Show current status
/usr/local/bin/site-sb-status.sh

