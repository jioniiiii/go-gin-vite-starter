# Default target
.DEFAULT_GOAL := help

SERVER_DIR := src
WEB_DIR    := client
PORT      ?= 8080

.PHONY: help dev-api dev-web dev prod build clean

help:
	@echo "  make dev-api   - run API in dev mode"
	@echo "  make dev-web   - run Vite dev server"
	@echo "  make dev       - run web (bg) + API (fg)"
	@echo "  make build     - build SPA + Go binary"
	@echo "  make prod      - run API in production mode"
	@echo "  make clean     - remove build artifacts"

dev-api:
	cd $(SERVER_DIR) && ENV=development GIN_MODE=debug PORT=$(PORT) \
	ALLOW_ORIGINS=http://localhost:5173 \
	go run .

dev-web:
	cd $(WEB_DIR) && npm run dev

dev:
	@{ \
	  cd $(WEB_DIR) && npm run dev & \
	  WEB_PID=$$!; \
	  trap 'kill $$WEB_PID' EXIT INT TERM; \
	  cd "$(CURDIR)/$(SERVER_DIR)" && ENV=development GIN_MODE=debug PORT=$(PORT) \
	  ALLOW_ORIGINS=http://localhost:5173 go run .; \
	  wait $$WEB_PID; \
	}

build:
	cd $(WEB_DIR) && npm run build
	mkdir -p $(SERVER_DIR)/dist
	cp -R $(WEB_DIR)/dist/* $(SERVER_DIR)/dist/ 2>/dev/null || true
	cd $(SERVER_DIR) && GIN_MODE=release go build -o bin/app .

prod:
	cd $(SERVER_DIR) && ENV=production GIN_MODE=release PORT=$(PORT) go run .

clean:
	rm -rf $(WEB_DIR)/dist $(SERVER_DIR)/dist $(SERVER_DIR)/bin
