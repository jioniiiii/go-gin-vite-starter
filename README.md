# go-gin-vite-starter

Minimal starter for a **Go (Gin)** API with a **Vite/React/TypeScript** frontend.
Dev uses a Vite proxy (no CORS headaches). Prod is env-driven with CORS, timeouts, health checks, and graceful shutdown.

## Features

* Go (Gin) API with `release/debug` modes
* Env-driven config: `PORT`, `ALLOW_ORIGINS`, timeouts, trusted proxies
* Health endpoints: `/health`
* Example API: `/api/message → { "message": "success" }`
* Vite + React + TS with dev proxy to `:8080`
* Tiny Makefile: `dev`, `dev-api`, `dev-web`, `build`, `prod`

## Project Structure

```
.
├─ src/           # Go (package main)
│  ├─ main.go
│  └─ config.go
├─ client/        # Vite + React + TS
│  └─ src/App.tsx
├─ db/            # (optional) migrations
└─ Makefile
```

## Quick Start (Development)

```bash
# Backend
cd src
ENV=development GIN_MODE=debug ALLOW_ORIGINS=http://localhost:5173 go run .

# Frontend (new terminal)
cd client
npm install
npm run dev
```

* API: [http://localhost:8080/health](http://localhost:8080/health)
* Frontend: [http://localhost:5173](http://localhost:5173)

> `client/vite.config.ts` proxies `/api` → `http://localhost:8080`.

## Environment Variables

```
ENV=development            # development | production
GIN_MODE=debug             # debug | release | test
PORT=8080
ALLOW_ORIGINS=http://localhost:5173   # prod: https://your-frontend.com
TRUSTED_PROXIES=0.0.0.0/0
READ_TIMEOUT=10s
WRITE_TIMEOUT=15s
IDLE_TIMEOUT=60s
# Optional
DATABASE_URL=
JWT_SECRET=
```

## Makefile Commands

```bash
make dev-api   # run API (dev)
make dev-web   # run Vite (dev)
make dev       # run both; Ctrl+C stops both
make build     # build SPA + Go binary at src/bin/app
make prod      # run API in production mode
```

## Production Notes

* Frontend host: set `VITE_API_URL=https://api.your-domain.com`.
* API host: set `ALLOW_ORIGINS=https://your-frontend.com`, `GIN_MODE=release`, `ENV=production`.
* Optional single-service: serve built SPA from Go:

  ```go
  r.Static("/assets", "./dist/assets")
  r.NoRoute(func(c *gin.Context) { c.File("./dist/index.html") })
  ```

## API Routes

* `GET /health` → `{"status":"ok"}`

---

