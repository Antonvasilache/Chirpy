# Chirpy 🐦

Chirpy is a Twitter-like RESTful API service built with **Go** and **PostgreSQL**. It includes secure user authentication with JWT and refresh tokens, chirp (post) creation and retrieval with profanity filtering, webhooks for membership tier upgrades ("Chirpy Red"), admin metrics, and static asset serving.

---

## 🚀 Features

- **User Authentication & Authorization**:
  - Secure password hashing using `bcrypt`.
  - JWT-based authentication for access tokens (1-hour lifespan).
  - Hexadecimal refresh tokens (stored in PostgreSQL, 60-day lifespan) for rotating access tokens.
  - User profile management (email & password update).
- **Chirps (Posts) Management**:
  - Create, read, and delete chirps (up to 140 characters).
  - Profanity filter that censors specific banned words (`kerfuffle`, `sharbert`, `fornax`).
  - Retrieve all chirps or filter by author (`?author_id=<uuid>`).
  - Sorting support (`?sort=asc` or `?sort=desc`).
  - Author-only chirp deletion permissions.
- **Webhooks & Premium Upgrades**:
  - Integration with Polka webhook endpoint (`POST /api/polka/webhooks`).
  - Upgrades users to **Chirpy Red** (`is_chirpy_red = true`).
  - Authenticated via custom API Key header.
- **Admin & Monitoring**:
  - Static file server at `/app/` with hit counter middleware.
  - Admin dashboard displaying visit metrics at `/admin/metrics`.
  - Development reset endpoint (`POST /admin/reset`) to wipe user data and reset metrics (guarded by `PLATFORM=dev`).
  - Health check endpoint (`GET /api/healthz`).

---

## 🛠 Tech Stack

- **Backend**: Go (1.23+) using standard library `net/http` routing
- **Database**: PostgreSQL
- **SQL / Query Generation**: [sqlc](https://sqlc.dev/)
- **Authentication**: `golang-jwt/jwt/v5`, `golang.org/x/crypto/bcrypt`
- **Environment Management**: `godotenv`

---

## 📁 Project Structure

```text
Chirpy/
├── assets/                  # Static assets (images, logos)
├── internal/
│   ├── auth/                # JWT, password hashing, token validation, API key helpers
│   ├── database/            # Generated sqlc Go models and query functions
│   └── helpers/             # JSON response utilities, profanity filtering
├── sql/
│   ├── queries/             # SQL query files (users.sql, chirps.sql, refresh_tokens.sql)
│   └── schema/              # Database migration schemas (001 to 005)
├── handler_*.go             # HTTP handler implementations
├── index.html               # Frontend demo page
├── main.go                  # Server setup and route registration
├── sqlc.yaml                # sqlc configuration
├── types.go                 # Shared API types and request/response structs
├── go.mod                   # Go module definitions
└── go.sum                   # Checksums for dependencies
```

---

## ⚙️ Environment Variables

Create a `.env` file in the root of the project:

```env
PORT=8080
DB_URL=postgres://<username>:<password>@localhost:5432/<dbname>?sslmode=disable
PLATFORM=dev
JWT_SECRET=your_jwt_secret_key_here
POLKA_KEY=your_polka_api_key_here
```

### Description of Variables
| Variable | Description |
|---|---|
| `DB_URL` | PostgreSQL connection string. |
| `PLATFORM` | Environment mode (`dev` enables the `/admin/reset` endpoint). |
| `JWT_SECRET` | Secret key used to sign and verify JWT access tokens. |
| `POLKA_KEY` | API key required in the `Authorization` header for Polka webhooks. |

---

## 🚦 Getting Started

### Prerequisites
- [Go](https://golang.org/doc/install) (version 1.23 or higher)
- [PostgreSQL](https://www.postgresql.org/download/)
- [sqlc](https://docs.sqlc.dev/en/latest/overview/install.html) (optional, if modifying SQL queries)
- [goose](https://github.com/pressly/goose) (optional, for running schema migrations)

### 1. Clone the repository
```bash
git clone https://github.com/Antonvasilache/Chirpy.git
cd Chirpy/Chirpy
```

### 2. Configure Environment Variables
```bash
cp .env.example .env # or create .env as shown above
```

### 3. Database Migrations
Apply the schemas in `sql/schema/` to your PostgreSQL database:
```bash
cd sql/schema
goose postgres "postgres://<username>:<password>@localhost:5432/<dbname>?sslmode=disable" up
```

### 4. Build and Run
```bash
go run .
```
The server will start listening at `http://localhost:8080`.

---

## 📡 API Reference

### Health & Admin
| Method | Endpoint | Description | Auth Required |
|---|---|---|---|
| `GET` | `/api/healthz` | Readiness health check | No |
| `GET` | `/admin/metrics` | Returns HTML metrics page with file server hit count | No |
| `POST` | `/admin/reset` | Resets metrics counter and clears all users (only in `dev`) | No |

### Users & Authentication
| Method | Endpoint | Description | Auth Required |
|---|---|---|---|
| `POST` | `/api/users` | Register a new user | No |
| `POST` | `/api/login` | Login and obtain JWT & refresh token | No |
| `PUT` | `/api/users` | Update email and password | Bearer JWT |
| `POST` | `/api/refresh` | Exchange refresh token for a new access token | Bearer Refresh Token |
| `POST` | `/api/revoke` | Revoke a refresh token | Bearer Refresh Token |

### Chirps (Posts)
| Method | Endpoint | Description | Auth Required |
|---|---|---|---|
| `GET` | `/api/chirps` | List all chirps (optional: `?author_id=<id>&sort=asc\|desc`) | No |
| `GET` | `/api/chirps/{chirpID}` | Get a chirp by its UUID | No |
| `POST` | `/api/chirps` | Create a new chirp (max 140 characters) | Bearer JWT |
| `DELETE` | `/api/chirps/{chirpID}` | Delete a chirp (author only) | Bearer JWT |

### Webhooks
| Method | Endpoint | Description | Auth Required |
|---|---|---|---|
| `POST` | `/api/polka/webhooks` | Upgrade user to Chirpy Red (`user.upgraded` event) | `ApiKey <POLKA_KEY>` |

---

## 🧪 Testing

Run automated tests (such as unit tests in `internal/auth/`):
```bash
go test ./...
```
