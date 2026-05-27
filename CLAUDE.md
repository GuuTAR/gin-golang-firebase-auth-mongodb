# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Reusable Go project template using **Gin** as the HTTP framework

## Commands

```bash
# Run the server with live reload (rebuilds on every save)
air

# Run the server (no live reload)
go run ./cmd/server/main.go

# Build
go build -o bin/server ./cmd/server

# Format code
gofmt -w .

# Lint (requires golangci-lint)
golangci-lint run

# Install dependencies
go mod vendor

# Docker
docker build -t gin-golang-firebase-auth-mongodb .
docker run --rm -p 8080:8080 --env-file .env gin-golang-firebase-auth-mongodb
```

## CI

GitHub Actions runs on every push and pull request to `main` (see [`.github/workflows/ci.yml`](.github/workflows/ci.yml)):

| Job              | Steps                                                       |
| ---------------- | ----------------------------------------------------------- |
| **Build & Test** | `go mod download` → `go vet` → `go build` → `go test -race` |
| **Lint**         | `golangci-lint` (latest)                                    |

## Architecture

The project follows a layered architecture:

```
cmd/server/         # Entry point — initializes config, DB, Firebase, and Gin router
internal/
  config/           # App configuration loaded from env vars or config file
  middleware/       # Gin middleware (JWT verification, CORS, etc.)
  handlers/         # HTTP handlers — thin layer that delegates to services
  services/         # Business logic
  repository/       # Database data access layer
  models/           # Shared data models / DTOs
pkg/                # Reusable packages (Auth client, DB client)
```
