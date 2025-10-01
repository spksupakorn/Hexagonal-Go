# D&D Characters and Quests API (Echo + GORM + Postgres)

A hexagonal-architecture API for managing D&D characters and quests.

## Features

- Public vs registered access:
  - GET /characters and GET /quests returns public items for unauthenticated visitors; returns all active items if authenticated.
- Registered users:
  - Create, edit, delete their own characters and quests.
- Admin:
  - Manage predefined options (Classes, Races, Quest Levels).
  - On deletion of any option, related characters/quests are archived (not hard-deleted).
- Validation:
  - Description max 5000 characters.
  - Up to 10 images per character/quest (stored as JSON array of URLs).
- JWT auth with roles (user/admin).
- Unit tests for business logic (use cases).

## Tech

- Echo framework
- GORM with PostgreSQL
- Hexagonal architecture

## Getting started

1. Set environment variables:
   DB_PASSWORD: The password for the database user.
   DB_NAME: The name of the database your application will use.
   DB_SSLMODE: The SSL mode for connecting to the database (e.g., disable, require, verify-full).
   DB_TIMEZONE: The timezone setting for your database connection (e.g., UTC).
   JWT_SECRET: The secret key used to sign and verify JWT tokens for authentication.
   FILE_STORAGE_PATH: The directory path where uploaded files will be stored.
   MAX_FILE_SIZE: The maximum allowed size (in bytes) for uploaded files.
   DOMAIN: The domain name where your application is hosted (used for generating URLs, cookies, etc.).

2. Run Postgres and create database.

3. Install deps and run:
   ```bash
   go mod tidy
   go run ./cmd/api
   ```

4. Run migration database & pre data insert:
   go run ./internal/infrastructure/db/migrations/migration.go 

## Swagger / OpenAPI

Use [swaggo/swag](https://github.com/swaggo/swag) and [swaggo/echo-swagger](https://github.com/swaggo/echo-swagger).

- Install the `swag` CLI if you don't have it:
  ```bash
  go install github.com/swaggo/swag/cmd/swag@latest
  ```
- Generate docs (from repo root):
  ```bash
  swag init -g cmd/api/main.go -o ./docs
  ```
- Start the server:
  ```bash
  go run ./cmd/api
  ```
- Open the Swagger UI:
  - http://localhost:8080/swagger/index.html

Note:
- Annotations live in the handler files.
- Security scheme is `BearerAuth` (Authorization: `Bearer <token>`).
- GET list endpoints are public; create/update/delete require auth; admin endpoints require admin role.

## API Overview

- Auth:
  - POST /auth/login {username, password} -> {token}

- Public/Registered:
  - GET /characters
  - GET /quests
  - GET /options/classes
  - GET /options/races
  - GET /options/quest-levels

- Registered (Authorization: Bearer <token>):
  - POST /characters
  - PUT /characters/:id
  - DELETE /characters/:id
  - POST /quests
  - PUT /quests/:id
  - DELETE /quests/:id
  - POST /characters/:id/images
  - POST /quests/:id/images

- Admin (Authorization: Bearer <admin token>):
  - POST /admin/options/classes
  - PUT /admin/options/classes/:id
  - DELETE /admin/options/classes/:id
  - POST /admin/options/races
  - PUT /admin/options/races/:id
  - DELETE /admin/options/races/:id
  - POST /admin/options/quest-levels
  - PUT /admin/options/quest-levels/:id
  - DELETE /admin/options/quest-levels/:id

## Notes

- Accessibility: public | private
  - Unauthenticated users see only public.
  - Any authenticated user can fetch all active items (as per the specification).
- Status: active | archived
  - archived items are not returned by list endpoints and cannot be edited.

## Testing

```bash
go test ./internal/usecases -v
```

## Docker Compose

This project supports running with Docker Compose for easy setup and deployment.

### Usage

To start all services in detached mode, run:

```bash
docker compose up -d
```
