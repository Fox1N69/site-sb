#!/bin/bash

cd sitesb/site-sb/server || exit

git pull origin main --force

go build -o bin/server ./cmd/app/main.go

pm2 restart server
