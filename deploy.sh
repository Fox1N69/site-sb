#!/bin/bash

# Перейдите в директорию вашего приложения
cd sitesb || exit

# Получите последние изменения
git pull origin main

# Соберите проект
yarn
yarn build

# Перезапустите PM2
pm2 restart client
