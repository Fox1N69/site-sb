#!/bin/bash

cd server || exit

git pull origin main

go build -o bin/server ./cmd/app/main.go

pm2 restart server
