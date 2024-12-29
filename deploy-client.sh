#!/bin/bash

# Перейдите в директорию вашего приложения
cd sitesb/site-sb || exit

# Получите последние изменения
git pull origin main --force

# Соберите проект
yarn
yarn build

# Перезапустите PM2
pm2 restart client
