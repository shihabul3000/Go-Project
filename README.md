# SpotSync API

SpotSync is a clean-architecture Go API for parking and EV charging spot reservations. It is built for busy locations such as airports and malls where EV charging spaces are limited and capacity must stay correct even when multiple drivers reserve at the same time.

Live URL: add your Render, Railway, or Fly.io URL here after deployment.

## Features

- Public driver registration and login with bcrypt password hashing.
- JWT authentication with user id and role claims.
- Driver and admin roles with protected routes.
- Admin parking zone management, including capacity and hourly pricing.
- Public parking zone availability calculated from active reservations.
- Reservation creation with a GORM transaction and `FOR UPDATE` row lock on the parking zone.
- Drivers can view and cancel their own reservations.
- Admins can view all reservations with user and zone details.
- Centralized API response and error formatting.
- PostgreSQL connection pooling for production use.

## Tech Stack

- Go 1.22+
- Echo
- GORM with PostgreSQL
- go-playground/validator
- golang-jwt/jwt/v5
- bcrypt
- PostgreSQL on NeonDB, Supabase, Aiven, or local Postgres

## Architecture

The project follows clean architecture with manual dependency injection in `main.go`.

```text
HTTP request
  -> handler/      Bind and validate DTOs, read JWT context, return JSON
  -> service/      Business rules, password hashing, token creation, authorization checks
  -> repository/   GORM queries, transactions, row locks, preloads
  -> models/       Database table definitions
  -> dto/          Request and response shapes exposed by the API
```

Handlers never talk to the database directly, and repositories never know about HTTP.

## Local Setup

1. Install Go 1.22 or newer.
2. Create a PostgreSQL database.
3. Copy `.env.example` to `.env`.
4. Update `DATABASE_URL` and `JWT_SECRET`.
5. Download dependencies and run the API:

```bash
go mod tidy
go run .
```

The API runs on `http://localhost:8080` by default in local development. In production, use your deployed Vercel URL as the base URL.

Health check:

```bash
curl http://localhost:8080/health
curl https://your-vercel-project.vercel.app/health
```

Readiness check with database ping:

```bash
curl http://localhost:8080/ready
curl https://your-deployed-api.example.com/ready
```

## Environment Variables

| Variable | Required | Example |
| --- | --- | --- |
| `APP_ENV` | No | `development` |
| `PORT` | No | `8080` |
| `DATABASE_URL` | Yes | `postgres://postgres:postgres@localhost:5432/spotsync?sslmode=disable` |
| `JWT_SECRET` | Yes | `replace-with-a-long-random-secret` |
| `JWT_EXPIRES_IN` | No | `24h` |
| `BCRYPT_COST` | No | `12` |
| `ALLOWED_ORIGINS` | No | `*` or `https://your-frontend.com` |
| `DB_MAX_OPEN_CONNS` | No | `25` |
| `DB_MAX_IDLE_CONNS` | No | `10` |
| `DB_CONN_MAX_LIFETIME` | No | `30m` |

## API Endpoints

All responses follow:

```json
{
  "success": true,
  "message": "Operation message",
  "data": {}
}
```

Errors follow:

```json
{
  "success": false,
  "message": "Error message",
  "errors": {}
}
```

### Authentication

| Method | Path | Access | Purpose |
| --- | --- | --- | --- |
| `POST` | `/api/v1/auth/register` | Public | Register a driver or admin |
| `POST` | `/api/v1/auth/login` | Public | Login and receive a JWT |

### Health

| Method | Path | Access | Purpose |
| --- | --- | --- | --- |
| `GET` | `/` | Public | Show API landing metadata |
| `GET` | `/health` | Public | Confirm the HTTP server is alive |
| `GET` | `/ready` | Public | Confirm the API can reach the database |

Register request:

```json
{
  "name": "John Doe",
  "email": "john.doe@spotsync.com",
  "password": "securePassword123",
  "role": "driver"
}
```

Login request:

```json
{
  "email": "john.doe@spotsync.com",
  "password": "securePassword123"
}
```

### Parking Zones

| Method | Path | Access | Purpose |
| --- | --- | --- | --- |
| `GET` | `/api/v1/zones` | Public | List zones with availability |
| `GET` | `/api/v1/zones/:id` | Public | Get one zone with availability |
| `POST` | `/api/v1/zones` | Admin | Create a parking zone |
| `PATCH` | `/api/v1/zones/:id` | Admin | Update zone details or pricing |
| `DELETE` | `/api/v1/zones/:id` | Admin | Delete an inactive zone |

Create zone request:

```json
{
  "name": "Terminal 1 EV Charging",
  "type": "ev_charging",
  "total_capacity": 20,
  "price_per_hour": 5.5
}
```

### Reservations

| Method | Path | Access | Purpose |
| --- | --- | --- | --- |
| `POST` | `/api/v1/reservations` | Authenticated | Reserve a parking or EV spot |
| `GET` | `/api/v1/reservations/my-reservations` | Authenticated | View own reservations |
| `DELETE` | `/api/v1/reservations/:id` | Authenticated | Cancel own reservation; admins may cancel any |
| `GET` | `/api/v1/reservations` | Admin | View all reservations |

Create reservation request:

```json
{
  "zone_id": 5,
  "license_plate": "ABC-1234"
}
```

Use the JWT from login. Replace `BASE_URL` with `http://localhost:8080` locally or your Vercel URL in production:

```bash
curl -H "Authorization: Bearer YOUR_TOKEN" "$BASE_URL/api/v1/reservations/my-reservations"
```

## Concurrency Rule

The reservation repository creates reservations inside a database transaction. It locks the selected parking zone row with:

```go
tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&zone, zoneID)
```

After the row is locked, it counts active reservations for that zone, compares the count with `total_capacity`, and only then creates the reservation. This prevents two requests from taking the final EV charging spot at the same time.

## Testing

Run:

```bash
go test ./...
```

The included service tests cover password hashing, login rejection, license plate normalization, zone-full errors, and cancellation authorization.

## Deployment

### Vercel

This repository includes `vercel.json` with the Go framework preset and the Singapore function region (`sin1`) to stay close to the provided Neon database region.

1. Push the repository to GitHub.
2. Import the GitHub repository into Vercel.
3. In Vercel project settings, add these environment variables:
   - `DATABASE_URL`
   - `JWT_SECRET`
   - `JWT_EXPIRES_IN=24h`
   - `BCRYPT_COST=12`
   - `ALLOWED_ORIGINS=*`
   - `APP_ENV=production`
4. Deploy the project.
5. Test the health endpoint:

```bash
curl https://your-vercel-project.vercel.app/health
```

### Render

Render setup:

1. Push this repository to GitHub.
2. Create a PostgreSQL database on NeonDB, Supabase, or Aiven.
3. Create a Render Web Service from the GitHub repository.
4. Add these environment variables in Render:
   - `DATABASE_URL`
   - `JWT_SECRET`
   - `APP_ENV=production`
   - `ALLOWED_ORIGINS=https://your-frontend-or-domain.com`
5. Use `go build -o bin/spotsync-api .` as the build command.
6. Use `./bin/spotsync-api` as the start command.

`render.yaml` is included if you prefer Render blueprint deployment.

## Suggested Meaningful Commits

Use real commits as the project grows. A clean 10-commit flow could be:

1. `Initialize Go Echo project`
2. `Add configuration and database setup`
3. `Create GORM models for users zones and reservations`
4. `Add DTOs validation and response helpers`
5. `Implement user registration and login`
6. `Add JWT authentication and role middleware`
7. `Implement parking zone management`
8. `Add reservation transaction with row locking`
9. `Add service tests for auth and reservations`
10. `Document setup API and deployment steps`
