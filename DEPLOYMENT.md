# 🚀 Инструкция по развертыванию Site-SB на Ubuntu Server

Полная инструкция по развертыванию веб-приложения Site-SB на чистом Ubuntu Server с настройкой всех необходимых компонентов.

## 📋 Содержание

1. [Подготовка сервера](#подготовка-сервера)
2. [Установка Docker и Docker Compose](#установка-docker-и-docker-compose)
3. [Клонирование проекта](#клонирование-проекта)
4. [Настройка Nginx и SSL](#настройка-nginx-и-ssl)
5. [Запуск приложения](#запуск-приложения)
6. [Мониторинг и обслуживание](#мониторинг-и-обслуживание)
7. [Автоматизация развертывания](#автоматизация-развертывания)

## 📋 Требования

- Ubuntu Server 20.04+ (или Ubuntu 22.04 LTS рекомендуется)
- Минимум 2 GB RAM
- Минимум 20 GB свободного места на диске
- Доменное имя (опционально, для SSL)
- SSH доступ к серверу

---

## 1. Подготовка сервера

### 1.1 Подключение к серверу

```bash
ssh root@your-server-ip
# или
ssh ubuntu@your-server-ip
```

### 1.2 Обновление системы

```bash
# Обновление списка пакетов
sudo apt update

# Обновление установленных пакетов
sudo apt upgrade -y

# Установка необходимых утилит
sudo apt install -y curl wget git htop tree unzip
```

### 1.3 Создание пользователя (если используете root)

```bash
# Создание нового пользователя
sudo adduser siteuser

# Добавление в группу sudo
sudo usermod -aG sudo siteuser

# Переключение на нового пользователя
su - siteuser
```

### 1.4 Настройка SSH ключей (рекомендуется)

```bash
# На локальной машине
ssh-keygen -t rsa -b 4096 -C "your-email@example.com"

# Копирование ключа на сервер
ssh-copy-id siteuser@your-server-ip
```

### 1.5 Настройка firewall

```bash
# Включение UFW
sudo ufw enable

# Разрешение SSH
sudo ufw allow ssh

# Разрешение HTTP и HTTPS
sudo ufw allow 80
sudo ufw allow 443

# Разрешение портов приложения (для разработки)
sudo ufw allow 4000
sudo ufw allow 4321

# Проверка статуса
sudo ufw status
```

---

## 2. Установка Docker и Docker Compose

### 2.1 Удаление старых версий Docker

```bash
sudo apt-get remove docker docker-engine docker.io containerd runc
```

### 2.2 Установка Docker

```bash
# Обновление apt и установка зависимостей
sudo apt-get update
sudo apt-get install -y \
    ca-certificates \
    curl \
    gnupg \
    lsb-release

# Добавление официального GPG ключа Docker
sudo mkdir -p /etc/apt/keyrings
curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo gpg --dearmor -o /etc/apt/keyrings/docker.gpg

# Добавление репозитория Docker
echo \
  "deb [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/docker.gpg] https://download.docker.com/linux/ubuntu \
  $(lsb_release -cs) stable" | sudo tee /etc/apt/sources.list.d/docker.list > /dev/null

# Установка Docker Engine
sudo apt-get update
sudo apt-get install -y docker-ce docker-ce-cli containerd.io docker-compose-plugin

# Добавление пользователя в группу docker
sudo usermod -aG docker $USER

# Применение изменений групп (перелогиньтесь или выполните)
newgrp docker

# Проверка установки
docker --version
docker compose version
```

### 2.3 Настройка автозапуска Docker

```bash
# Включение автозапуска Docker
sudo systemctl enable docker
sudo systemctl enable containerd

# Запуск Docker
sudo systemctl start docker

# Проверка статуса
sudo systemctl status docker
```

### 2.4 Тестирование Docker

```bash
# Тестовый запуск
docker run hello-world

# Проверка работы docker compose
docker compose --version
```

---

## 3. Клонирование проекта

### 3.1 Создание директории для проекта

```bash
# Создание директории
sudo mkdir -p /opt/site-sb
sudo chown $USER:$USER /opt/site-sb

# Переход в директорию
cd /opt/site-sb
```

### 3.2 Клонирование репозитория

```bash
# Клонирование проекта (замените на ваш репозиторий)
git clone https://github.com/your-username/site-sb.git .

# Проверка содержимого
ls -la
```

### 3.3 Настройка конфигурации

```bash
# Создание копии конфигурации для продакшена
cp backend/config/config.json backend/config/config.prod.json

# Редактирование конфигурации
nano backend/config/config.prod.json
```

Измените в `config.prod.json`:
```json
{
    "environment": {
        "mode": "release"
    },
    "server": {
        "port": 4000
    },
    "database": {
        "host": "postgres",
        "port": "5432",
        "user": "postgres",
        "pass": "ваш_пароль_БД",
        "name": "sitesb"
    },
    "secret": {
        "key": "ваш_секретный_ключ_для_JWT"
    }
}
```

---

## 4. Настройка Nginx и SSL

### 4.1 Установка Nginx

```bash
# Установка Nginx
sudo apt install -y nginx

# Запуск и автозапуск Nginx
sudo systemctl start nginx
sudo systemctl enable nginx

# Проверка статуса
sudo systemctl status nginx
```

### 4.2 Настройка Nginx для проекта

```bash
# Создание конфигурации сайта
sudo nano /etc/nginx/sites-available/site-sb
```

Содержимое файла:
```nginx
server {
    listen 80;
    server_name your-domain.com www.your-domain.com;
    server_tokens off;

    # Временная конфигурация для получения SSL сертификата
    location /.well-known/acme-challenge/ {
        root /var/www/certbot;
        allow all;
    }

    location / {
        return 301 https://$host$request_uri;
    }
}

server {
    listen 443 ssl http2;
    server_name your-domain.com www.your-domain.com;
    server_tokens off;

    # SSL конфигурация (будет настроена после получения сертификатов)
    ssl_certificate /etc/letsencrypt/live/your-domain.com/fullchain.pem;
    ssl_certificate_key /etc/letsencrypt/live/your-domain.com/privkey.pem;
    include /etc/letsencrypt/options-ssl-nginx.conf;
    ssl_dhparam /etc/letsencrypt/ssl-dhparams.pem;

    # Проксирование на frontend
    location / {
        proxy_pass http://localhost:4321;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
        proxy_set_header X-Forwarded-Host $host;
        proxy_set_header X-Forwarded-Port $server_port;
    }

    # Проксирование API на backend
    location /api/ {
        proxy_pass http://localhost:4000/;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }

    # Статические файлы
    location ~* \.(js|css|png|jpg|jpeg|gif|ico|svg|woff|woff2|ttf|eot)$ {
        proxy_pass http://localhost:4321;
        expires 1y;
        add_header Cache-Control "public, immutable";
    }
}
```

### 4.3 Активация конфигурации

```bash
# Активация сайта
sudo ln -s /etc/nginx/sites-available/site-sb /etc/nginx/sites-enabled/

# Удаление дефолтного сайта
sudo rm /etc/nginx/sites-enabled/default

# Проверка конфигурации
sudo nginx -t

# Перезапуск Nginx
sudo systemctl restart nginx
```

### 4.4 Установка Certbot для SSL сертификатов

```bash
# Установка Certbot
sudo apt install -y certbot python3-certbot-nginx

# Создание директории для Certbot
sudo mkdir -p /var/www/certbot

# Получение SSL сертификата (замените на ваш домен)
sudo certbot --nginx -d your-domain.com -d www.your-domain.com

# Настройка автообновления сертификатов
sudo crontab -e
```

Добавьте в crontab:
```
0 12 * * * /usr/bin/certbot renew --quiet
```

---

## 5. Запуск приложения

### 5.1 Создание production docker-compose

```bash
# Переход в директорию проекта
cd /opt/site-sb

# Создание production конфигурации
cp docker-compose.yaml docker-compose.prod.yaml
nano docker-compose.prod.yaml
```

Обновите `docker-compose.prod.yaml`:
```yaml
services:
  postgres:
    image: postgres:15-alpine
    environment:
      POSTGRES_USER: postgres
      POSTGRES_DB: sitesb
      POSTGRES_PASSWORD: ${DB_PASSWORD}
    volumes:
      - postgres_data:/var/lib/postgresql/data
      - ./backups:/backups
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U postgres"]
      interval: 10s
      timeout: 5s
      retries: 5
    networks:
      - app-network
    restart: unless-stopped

  backend:
    build:
      context: ./backend
      dockerfile: Dockerfile
    environment:
      - DB_HOST=postgres
      - DB_USER=postgres
      - DB_NAME=sitesb
      - DB_PASS=${DB_PASSWORD}
      - DB_PORT=5432
      - GIN_MODE=release
    ports:
      - "127.0.0.1:4000:4000"
    depends_on:
      postgres:
        condition: service_healthy
    volumes:
      - ./backend/config/config.prod.json:/config/config.json:ro
    networks:
      - app-network
    restart: unless-stopped

  frontend:
    build:
      context: ./frontend
      dockerfile: Dockerfile
    environment:
      - HOST=0.0.0.0
      - PORT=4321
      - NODE_ENV=production
      - API_URL=http://backend:4000
    ports:
      - "127.0.0.1:4321:4321"
    depends_on:
      - backend
    networks:
      - app-network
    restart: unless-stopped

volumes:
  postgres_data:

networks:
  app-network:
    driver: bridge
```

### 5.2 Создание .env файла

```bash
# Создание .env файла
nano .env
```

Содержимое `.env`:
```bash
# Database
DB_PASSWORD=ваш_сложный_пароль_для_БД

# Security
JWT_SECRET=ваш_очень_сложный_JWT_секрет

# Environment
NODE_ENV=production
```

### 5.3 Создание директорий

```bash
# Создание директории для бэкапов
mkdir -p backups

# Установка прав
chmod 755 backups
```

### 5.4 Сборка и запуск приложения

```bash
# Сборка и запуск в production режиме
docker compose -f docker-compose.prod.yaml up --build -d

# Проверка статуса контейнеров
docker compose -f docker-compose.prod.yaml ps

# Просмотр логов
docker compose -f docker-compose.prod.yaml logs -f
```

### 5.5 Создание первого администратора

```bash
# Подключение к контейнеру backend
docker compose -f docker-compose.prod.yaml exec backend /bin/sh

# В контейнере можно выполнить команды для создания админа
# (если в приложении есть соответствующие команды)
```

---

## 6. Мониторинг и обслуживание

### 6.1 Мониторинг сервисов

```bash
# Создание скрипта мониторинга
sudo nano /usr/local/bin/site-sb-status.sh
```

Содержимое скрипта:
```bash
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
echo ""

echo "=== Nginx Status ==="
sudo systemctl status nginx --no-pager
echo ""

echo "=== SSL Certificate Status ==="
sudo certbot certificates
```

```bash
# Сделать скрипт исполняемым
sudo chmod +x /usr/local/bin/site-sb-status.sh

# Запуск проверки статуса
/usr/local/bin/site-sb-status.sh
```

### 6.2 Создание скрипта бэкапа базы данных

```bash
# Создание скрипта бэкапа
sudo nano /usr/local/bin/site-sb-backup.sh
```

Содержимое скрипта:
```bash
#!/bin/bash

BACKUP_DIR="/opt/site-sb/backups"
DATE=$(date +%Y%m%d_%H%M%S)
BACKUP_FILE="sitesb_backup_$DATE.sql"

echo "Starting backup at $(date)"

# Создание бэкапа базы данных
cd /opt/site-sb
docker compose -f docker-compose.prod.yaml exec -T postgres pg_dump -U postgres sitesb > "$BACKUP_DIR/$BACKUP_FILE"

if [ $? -eq 0 ]; then
    echo "Backup completed successfully: $BACKUP_FILE"
    
    # Сжатие бэкапа
    gzip "$BACKUP_DIR/$BACKUP_FILE"
    echo "Backup compressed: $BACKUP_FILE.gz"
    
    # Удаление старых бэкапов (старше 30 дней)
    find "$BACKUP_DIR" -name "*.gz" -mtime +30 -delete
    echo "Old backups cleaned up"
else
    echo "Backup failed!"
    exit 1
fi
```

```bash
# Сделать скрипт исполняемым
sudo chmod +x /usr/local/bin/site-sb-backup.sh

# Настройка автоматического бэкапа
sudo crontab -e
```

Добавьте в crontab:
```
# Ежедневный бэкап в 2:00
0 2 * * * /usr/local/bin/site-sb-backup.sh >> /var/log/site-sb-backup.log 2>&1

# Еженедельная проверка статуса в понедельник в 9:00
0 9 * * 1 /usr/local/bin/site-sb-status.sh >> /var/log/site-sb-status.log 2>&1
```

### 6.3 Настройка логирования

```bash
# Настройка ротации логов
sudo nano /etc/logrotate.d/site-sb
```

Содержимое файла:
```
/var/log/site-sb*.log {
    daily
    missingok
    rotate 14
    compress
    delaycompress
    notifempty
    copytruncate
}
```

### 6.4 Команды для управления приложением

```bash
# Создание алиасов для удобного управления
echo "alias site-sb-start='cd /opt/site-sb && docker compose -f docker-compose.prod.yaml up -d'" >> ~/.bashrc
echo "alias site-sb-stop='cd /opt/site-sb && docker compose -f docker-compose.prod.yaml down'" >> ~/.bashrc
echo "alias site-sb-restart='cd /opt/site-sb && docker compose -f docker-compose.prod.yaml restart'" >> ~/.bashrc
echo "alias site-sb-logs='cd /opt/site-sb && docker compose -f docker-compose.prod.yaml logs -f'" >> ~/.bashrc
echo "alias site-sb-status='/usr/local/bin/site-sb-status.sh'" >> ~/.bashrc
echo "alias site-sb-backup='/usr/local/bin/site-sb-backup.sh'" >> ~/.bashrc

# Применить изменения
source ~/.bashrc
```

---

## 7. Автоматизация развертывания

### 7.1 Создание скрипта автоматического развертывания

```bash
# Создание скрипта полного развертывания
sudo nano /usr/local/bin/site-sb-deploy.sh
```

### 7.2 Настройка CI/CD (GitHub Actions)

```bash
# Создание директории для CI/CD
mkdir -p .github/workflows
nano .github/workflows/deploy.yml
```

### 7.3 Создание скрипта обновления

```bash
# Создание скрипта обновления
sudo nano /usr/local/bin/site-sb-update.sh
```

---

## 🔧 Полезные команды

```bash
# Просмотр логов всех сервисов
docker compose -f docker-compose.prod.yaml logs -f

# Рестарт определенного сервиса
docker compose -f docker-compose.prod.yaml restart backend

# Просмотр использования ресурсов
docker stats

# Очистка неиспользуемых Docker ресурсов
docker system prune -a

# Проверка портов
sudo netstat -tulpn | grep :80
sudo netstat -tulpn | grep :443

# Проверка SSL сертификата
openssl x509 -in /etc/letsencrypt/live/your-domain.com/cert.pem -text -noout
```

---

## 🚨 Решение проблем

### Проблема: Контейнеры не запускаются
```bash
# Проверка логов
docker compose -f docker-compose.prod.yaml logs

# Проверка портов
sudo netstat -tulpn | grep :4000
sudo netstat -tulpn | grep :4321
```

### Проблема: База данных недоступна
```bash
# Проверка статуса PostgreSQL контейнера
docker compose -f docker-compose.prod.yaml ps postgres

# Подключение к базе данных
docker compose -f docker-compose.prod.yaml exec postgres psql -U postgres -d sitesb
```

### Проблема: SSL сертификат не работает
```bash
# Проверка статуса Certbot
sudo certbot certificates

# Тестовое обновление сертификата
sudo certbot renew --dry-run
```

---

## ✅ Чек-лист развертывания

- [ ] Сервер обновлен и настроен
- [ ] Docker и Docker Compose установлены
- [ ] Проект склонирован
- [ ] Конфигурация настроена
- [ ] Nginx установлен и настроен
- [ ] SSL сертификаты получены
- [ ] Приложение запущено
- [ ] Мониторинг настроен
- [ ] Бэкапы настроены
- [ ] Firewall настроен
- [ ] Домен указывает на сервер

---

## 📞 Поддержка

При возникновении проблем:
1. Проверьте логи сервисов
2. Убедитесь, что все порты открыты
3. Проверьте конфигурацию Nginx
4. Убедитесь, что DNS записи настроены правильно

Успешного развертывания! 🚀