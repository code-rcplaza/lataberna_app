.PHONY: dev-backend dev-frontend test test-backend build build-backend build-frontend migrate

# ── Development ───────────────────────────────────────────────
dev-backend:
	cd backend && go run ./cmd/server

dev-frontend:
	cd frontend && pnpm dev

# ── Testing ───────────────────────────────────────────────────
test: test-backend

test-backend:
	cd backend && go test ./...

test-domain:
	cd backend && go test ./internal/domain/... ./internal/usecase/...

# ── Build ─────────────────────────────────────────────────────
build: build-backend build-frontend

build-backend:
	cd backend && go build -o ./bin/server ./cmd/server

build-frontend:
	cd frontend && pnpm build

# ── Database ──────────────────────────────────────────────────
migrate:
	cd backend && atlas schema apply --url "sqlite://forge.db" --to "file://migrations/001_initial.hcl" --auto-approve

migrate-diff:
	cd backend && atlas schema diff --from "sqlite://forge.db" --to "file://migrations/001_initial.hcl"

# ── Code generation ───────────────────────────────────────────
generate:
	cd backend && go generate ./...
