#!/bin/bash

# Перейдите в директорию вашего приложения
cd /sitesb || exit

# Получите последние изменения
git pull origin main

# Соберите проект
npm install
npm run build

# Перезапустите PM2
pm2 restart your-app-name
