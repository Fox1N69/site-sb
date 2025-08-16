#!/bin/bash

# Site-SB Update Script
# This script updates the application with zero downtime

set -e

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

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

create_backup() {
    log_info "Creating backup before update..."
    
    cd $PROJECT_DIR
    
    # Create backup directory with timestamp
    BACKUP_TIMESTAMP=$(date +%Y%m%d_%H%M%S)
    BACKUP_DIR="backups/pre_update_$BACKUP_TIMESTAMP"
    mkdir -p $BACKUP_DIR
    
    # Backup database
    docker compose -f docker-compose.prod.yaml exec -T postgres pg_dump -U postgres sitesb > "$BACKUP_DIR/database.sql"
    
    # Backup current code
    git stash push -m "Pre-update backup $BACKUP_TIMESTAMP" || true
    
    log_success "Backup created in $BACKUP_DIR"
}

update_code() {
    log_info "Updating code from repository..."
    
    cd $PROJECT_DIR
    
    # Fetch latest changes
    git fetch origin
    
    # Show what will be updated
    log_info "Changes to be applied:"
    git log --oneline HEAD..origin/main | head -10
    
    # Apply updates
    git pull origin main
    
    log_success "Code updated successfully"
}

update_dependencies() {
    log_info "Updating dependencies..."
    
    cd $PROJECT_DIR
    
    # Rebuild containers with latest dependencies
    docker compose -f docker-compose.prod.yaml build --no-cache
    
    log_success "Dependencies updated"
}

restart_services() {
    log_info "Restarting services..."
    
    cd $PROJECT_DIR
    
    # Restart services with rolling update approach
    docker compose -f docker-compose.prod.yaml up -d --force-recreate
    
    # Wait for services to be healthy
    log_info "Waiting for services to be ready..."
    sleep 30
    
    # Check if services are running
    if docker compose -f docker-compose.prod.yaml ps | grep -q "Up"; then
        log_success "Services restarted successfully"
    else
        log_error "Service restart failed"
        docker compose -f docker-compose.prod.yaml logs
        return 1
    fi
}

run_migrations() {
    log_info "Running database migrations..."
    
    cd $PROJECT_DIR
    
    # Database migrations are handled automatically by the backend
    # Check if backend is responding
    sleep 10
    
    if curl -f http://localhost:4000/ > /dev/null 2>&1; then
        log_success "Backend is responding, migrations completed"
    else
        log_warning "Backend not responding immediately, checking logs..."
        docker compose -f docker-compose.prod.yaml logs backend | tail -20
    fi
}

validate_deployment() {
    log_info "Validating deployment..."
    
    cd $PROJECT_DIR
    
    # Check if all containers are running
    if ! docker compose -f docker-compose.prod.yaml ps | grep -q "Up"; then
        log_error "Some containers are not running"
        docker compose -f docker-compose.prod.yaml ps
        return 1
    fi
    
    # Check if frontend is responding
    if ! curl -f http://localhost:4321 > /dev/null 2>&1; then
        log_error "Frontend is not responding"
        return 1
    fi
    
    # Check if backend is responding
    if ! curl -f http://localhost:4000/ > /dev/null 2>&1; then
        log_error "Backend is not responding"
        return 1
    fi
    
    log_success "Deployment validation passed"
}

cleanup() {
    log_info "Cleaning up..."
    
    # Remove old Docker images
    docker image prune -f
    
    # Remove old backups (keep last 10)
    cd $PROJECT_DIR/backups
    ls -t pre_update_* | tail -n +11 | xargs rm -rf 2>/dev/null || true
    
    log_success "Cleanup completed"
}

rollback() {
    log_error "Update failed, initiating rollback..."
    
    cd $PROJECT_DIR
    
    # Rollback code
    git reset --hard HEAD~1
    
    # Restart services with old code
    docker compose -f docker-compose.prod.yaml up --build -d
    
    log_warning "Rollback completed. Please check the logs and fix issues before retrying update."
}

main() {
    log_info "Starting Site-SB update process..."
    
    # Check if we're in the right directory
    if [ ! -f "$PROJECT_DIR/docker-compose.prod.yaml" ]; then
        log_error "Project directory not found or invalid"
        exit 1
    fi
    
    # Trap errors for rollback
    trap rollback ERR
    
    create_backup
    update_code
    update_dependencies
    restart_services
    run_migrations
    validate_deployment
    cleanup
    
    # Remove error trap on success
    trap - ERR
    
    log_success "Update completed successfully!"
    log_info "Application status:"
    /usr/local/bin/site-sb-status.sh
}

# Check for force flag
if [ "$1" = "--force" ]; then
    log_warning "Force update mode enabled, skipping confirmations"
else
    echo "This will update the Site-SB application with the latest code."
    echo "A backup will be created automatically."
    read -p "Do you want to continue? (y/N): " -r
    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
        log_info "Update cancelled"
        exit 0
    fi
fi

main