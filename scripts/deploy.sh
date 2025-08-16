#!/bin/bash

# Site-SB Auto Deployment Script
# This script automates the deployment process on Ubuntu server

set -e  # Exit on any error

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Configuration
PROJECT_NAME="site-sb"
PROJECT_DIR="/opt/${PROJECT_NAME}"
REPO_URL="https://github.com/your-username/site-sb.git"  # Replace with your repo
SERVICE_USER="siteuser"

# Functions
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

check_root() {
    if [[ $EUID -eq 0 ]]; then
        log_error "This script should not be run as root"
        exit 1
    fi
}

install_dependencies() {
    log_info "Installing system dependencies..."
    
    # Update system
    sudo apt update && sudo apt upgrade -y
    
    # Install basic tools
    sudo apt install -y curl wget git htop tree unzip nginx certbot python3-certbot-nginx
    
    # Install Docker
    if ! command -v docker &> /dev/null; then
        log_info "Installing Docker..."
        curl -fsSL https://get.docker.com -o get-docker.sh
        sudo sh get-docker.sh
        sudo usermod -aG docker $USER
        rm get-docker.sh
    else
        log_success "Docker already installed"
    fi
    
    log_success "Dependencies installed successfully"
}

setup_firewall() {
    log_info "Setting up firewall..."
    
    sudo ufw --force enable
    sudo ufw allow ssh
    sudo ufw allow 80
    sudo ufw allow 443
    
    log_success "Firewall configured"
}

setup_project() {
    log_info "Setting up project directory..."
    
    # Create project directory
    sudo mkdir -p $PROJECT_DIR
    sudo chown $USER:$USER $PROJECT_DIR
    
    # Clone or update repository
    if [ -d "$PROJECT_DIR/.git" ]; then
        log_info "Updating existing repository..."
        cd $PROJECT_DIR
        git pull origin main
    else
        log_info "Cloning repository..."
        git clone $REPO_URL $PROJECT_DIR
        cd $PROJECT_DIR
    fi
    
    # Create necessary directories
    mkdir -p backups
    mkdir -p logs
    
    # Set up environment file if it doesn't exist
    if [ ! -f ".env" ]; then
        log_warning ".env file not found. Creating from template..."
        cp .env.example .env
        log_warning "Please edit .env file with your configuration before proceeding"
        return 1
    fi
    
    log_success "Project setup completed"
}

setup_nginx() {
    log_info "Setting up Nginx configuration..."
    
    # Read domain from .env file
    if [ -f ".env" ]; then
        source .env
    fi
    
    if [ -z "$DOMAIN" ]; then
        log_error "DOMAIN not set in .env file"
        return 1
    fi
    
    # Create Nginx configuration
    sudo tee /etc/nginx/sites-available/$PROJECT_NAME > /dev/null <<EOF
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

    # Enable site
    sudo ln -sf /etc/nginx/sites-available/$PROJECT_NAME /etc/nginx/sites-enabled/
    sudo rm -f /etc/nginx/sites-enabled/default
    
    # Test Nginx configuration
    sudo nginx -t
    sudo systemctl restart nginx
    
    log_success "Nginx configured successfully"
}

setup_ssl() {
    log_info "Setting up SSL certificate..."
    
    if [ -f ".env" ]; then
        source .env
    fi
    
    if [ -z "$DOMAIN" ] || [ -z "$EMAIL" ]; then
        log_error "DOMAIN or EMAIL not set in .env file"
        return 1
    fi
    
    # Create certbot directory
    sudo mkdir -p /var/www/certbot
    
    # Get SSL certificate
    sudo certbot --nginx -d $DOMAIN -d www.$DOMAIN --email $EMAIL --agree-tos --non-interactive
    
    # Setup auto-renewal
    (sudo crontab -l 2>/dev/null; echo "0 12 * * * /usr/bin/certbot renew --quiet") | sudo crontab -
    
    log_success "SSL certificate configured"
}

deploy_application() {
    log_info "Deploying application..."
    
    cd $PROJECT_DIR
    
    # Create production config if it doesn't exist
    if [ ! -f "backend/config/config.prod.json" ]; then
        log_info "Creating production configuration..."
        cp backend/config/config.json backend/config/config.prod.json
        
        # Update production config
        if [ -f ".env" ]; then
            source .env
            sed -i "s/\"mode\": \"dev\"/\"mode\": \"release\"/g" backend/config/config.prod.json
            sed -i "s/\"host\": \"localhost\"/\"host\": \"postgres\"/g" backend/config/config.prod.json
            sed -i "s/\"pass\": \"8008\"/\"pass\": \"$DB_PASSWORD\"/g" backend/config/config.prod.json
        fi
    fi
    
    # Build and start services
    docker compose -f docker-compose.prod.yaml down || true
    docker compose -f docker-compose.prod.yaml up --build -d
    
    # Wait for services to start
    log_info "Waiting for services to start..."
    sleep 30
    
    # Check if services are running
    if docker compose -f docker-compose.prod.yaml ps | grep -q "Up"; then
        log_success "Application deployed successfully"
    else
        log_error "Application deployment failed"
        docker compose -f docker-compose.prod.yaml logs
        return 1
    fi
}

setup_monitoring() {
    log_info "Setting up monitoring and maintenance scripts..."
    
    # Create status script
    sudo tee /usr/local/bin/site-sb-status.sh > /dev/null <<'EOF'
#!/bin/bash
echo "=== Site-SB Status Report ==="
echo "Date: $(date)"
echo ""
echo "=== System Resources ==="
echo "Memory usage:"
free -h
echo ""
echo "Disk usage:"
df -h /
echo ""
echo "=== Docker Containers ==="
cd /opt/site-sb
docker compose -f docker-compose.prod.yaml ps
echo ""
echo "=== Service Logs (last 10 lines) ==="
echo "Backend logs:"
docker compose -f docker-compose.prod.yaml logs --tail=10 backend
echo ""
echo "Frontend logs:"
docker compose -f docker-compose.prod.yaml logs --tail=10 frontend
EOF

    # Create backup script
    sudo tee /usr/local/bin/site-sb-backup.sh > /dev/null <<'EOF'
#!/bin/bash
BACKUP_DIR="/opt/site-sb/backups"
DATE=$(date +%Y%m%d_%H%M%S)
BACKUP_FILE="sitesb_backup_$DATE.sql"

echo "Starting backup at $(date)"
cd /opt/site-sb
docker compose -f docker-compose.prod.yaml exec -T postgres pg_dump -U postgres sitesb > "$BACKUP_DIR/$BACKUP_FILE"

if [ $? -eq 0 ]; then
    echo "Backup completed successfully: $BACKUP_FILE"
    gzip "$BACKUP_DIR/$BACKUP_FILE"
    echo "Backup compressed: $BACKUP_FILE.gz"
    find "$BACKUP_DIR" -name "*.gz" -mtime +30 -delete
    echo "Old backups cleaned up"
else
    echo "Backup failed!"
    exit 1
fi
EOF

    # Make scripts executable
    sudo chmod +x /usr/local/bin/site-sb-status.sh
    sudo chmod +x /usr/local/bin/site-sb-backup.sh
    
    # Setup cron jobs
    (crontab -l 2>/dev/null; echo "0 2 * * * /usr/local/bin/site-sb-backup.sh >> /var/log/site-sb-backup.log 2>&1") | crontab -
    (crontab -l 2>/dev/null; echo "0 9 * * 1 /usr/local/bin/site-sb-status.sh >> /var/log/site-sb-status.log 2>&1") | crontab -
    
    # Create log rotation
    sudo tee /etc/logrotate.d/site-sb > /dev/null <<EOF
/var/log/site-sb*.log {
    daily
    missingok
    rotate 14
    compress
    delaycompress
    notifempty
    copytruncate
}
EOF

    # Create useful aliases
    tee -a ~/.bashrc > /dev/null <<EOF

# Site-SB aliases
alias site-sb-start='cd /opt/site-sb && docker compose -f docker-compose.prod.yaml up -d'
alias site-sb-stop='cd /opt/site-sb && docker compose -f docker-compose.prod.yaml down'
alias site-sb-restart='cd /opt/site-sb && docker compose -f docker-compose.prod.yaml restart'
alias site-sb-logs='cd /opt/site-sb && docker compose -f docker-compose.prod.yaml logs -f'
alias site-sb-status='/usr/local/bin/site-sb-status.sh'
alias site-sb-backup='/usr/local/bin/site-sb-backup.sh'
EOF

    log_success "Monitoring and maintenance setup completed"
}

main() {
    log_info "Starting Site-SB deployment..."
    
    check_root
    
    # Check if .env file exists for domain configuration
    if [ ! -f ".env" ] && [ ! -f "$PROJECT_DIR/.env" ]; then
        log_warning "No .env file found. Please create one from .env.example first"
        log_info "Run: cp .env.example .env && nano .env"
        exit 1
    fi
    
    install_dependencies
    setup_firewall
    setup_project || exit 1
    setup_nginx
    setup_ssl
    deploy_application
    setup_monitoring
    
    log_success "Site-SB deployment completed successfully!"
    log_info "Your site should be available at: https://$DOMAIN"
    log_info "Use 'site-sb-status' to check application status"
    log_info "Use 'site-sb-logs' to view application logs"
    
    # Show final status
    /usr/local/bin/site-sb-status.sh
}

# Run main function
main "$@"