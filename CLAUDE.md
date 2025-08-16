# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Development Commands

### Frontend (Astro + React/Svelte)
- `cd frontend && yarn dev` - Start development server (port 4321)
- `cd frontend && yarn build` - Build for production (includes Astro check)
- `cd frontend && yarn preview` - Preview production build
- Frontend uses Yarn as package manager

### Backend (Go + Gin)
- `cd backend && go run cmd/app/main.go` - Start development server (port 4000)
- `cd backend && go build -o bin/app cmd/app/main.go` - Build binary
- `cd backend && go test -short ./...` - Run tests
- `cd backend && go mod tidy` - Clean up dependencies

### Full Stack Development
- `docker-compose up` - Start all services (frontend:4321, backend:4000, postgres:5432)
- `docker-compose up --build` - Rebuild and start all services

## Architecture Overview

### Backend Structure (Go)
- **Entry Point**: `cmd/app/main.go` - Initializes infrastructure and starts API server
- **API Layer**: `internal/api/` - Gin-based REST API with middleware
  - `api.go` - Main server setup and route registration
  - `v1/` - API version 1 handlers (auth_handler.go, post_handler.go)
- **Business Logic**: 
  - `internal/service/` - Business logic layer
  - `internal/repo/` - Data access layer with GORM
  - `internal/manager/` - Service and repository managers
- **Models**: `internal/models/` - Database entities (User, Post)
- **Infrastructure**: `infra/infra.go` - Database connections, config, logging
- **Common**: Shared utilities for HTTP middleware, validation, tokens, SMTP

### Frontend Structure (Astro)
- **Framework**: Astro with SSR enabled, using Node.js adapter
- **UI Libraries**: React + Svelte components, TailwindCSS, NextUI
- **Pages**: `/src/pages/` - File-based routing including admin panel
- **Components**: `/src/components/` - Reusable Astro/React/Svelte components
- **Layouts**: `/src/layouts/` - Page templates
- **Assets**: Fonts (Roboto, Wadik), SCSS styles, images

### Database
- PostgreSQL with GORM ORM
- Connection configured through environment variables
- Database migrations handled through GORM auto-migration

### Authentication & Authorization
- JWT-based authentication with middleware
- Session management using Redis store
- Role-based access control (admin routes protected)
- Admin panel requires authentication

### Key Configuration Files
- `backend/config/config.json` - Backend configuration
- `frontend/astro.config.mjs` - Astro configuration with React/Svelte integrations
- `docker-compose.yaml` - Multi-service deployment configuration

### API Endpoints
- `/v1/account/auth/register` - User registration
- `/v1/account/auth/login` - User login
- `/post/` - Get all posts
- `/post/:id` - Get specific post
- `/admin/post/create` - Create post (admin only)
- `/admin/post/update/:id` - Update post (admin only)
- `/admin/post/delete/:id` - Delete post (admin only)

### Development Notes
- Backend includes CORS middleware for frontend integration
- Frontend uses server-side rendering with standalone mode
- Docker setup includes multi-stage builds for optimization
- Both services are networked through Docker bridge network
- Database credentials: postgres/8008 (development only)