# Sanjie Config and Docker Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Make Sanjie configurable and locally deployable as an open-source demo using Go + SQLite, uni-app H5, Docker Compose, and Nginx.

**Architecture:** Add a small backend config package that reads `config.yaml` when present and supports environment variable overrides. Package the Go API and H5 static build into separate containers, then route `/api/` and the H5 SPA through Nginx for same-origin local deployment.

**Tech Stack:** Go, SQLite, npm, uni-app H5, Docker Compose, Nginx, VitePress.

---

## Files

- Create: `server/internal/config/config.go`
- Create: `server/internal/config/config_test.go`
- Modify: `server/main.go`
- Modify: `app/api/http.js`
- Create: `server/Dockerfile`
- Create: `app/Dockerfile`
- Create: `nginx.conf`
- Create: `docker-compose.yml`
- Modify: `README.md`
- Modify: `docs/runbook.md`
- Modify: `docs-site/docs/deploy/index.md`
- Modify: `docs-site/docs/README.md`

## Task 1: Backend Configuration

- [ ] Add `server/internal/config` with defaults matching `config.example.yaml`.
- [ ] Support `SANJIE_CONFIG`, `SANJIE_SERVER_ADDR`, `SANJIE_DATABASE_DSN`, and `SANJIE_DEMO_SEED`.
- [ ] Update `server/main.go` to load config, open SQLite via configured DSN, and seed only when enabled.
- [ ] Add tests for default config, YAML parsing, and environment overrides.

## Task 2: H5 Runtime API Base

- [ ] Update `app/api/http.js` so H5 can use same-origin API by setting `VITE_SANJIE_API_BASE_URL`.
- [ ] Keep local development default as `http://localhost:8080`.
- [ ] Ensure request paths work when the API base is an empty string.

## Task 3: Docker Local Deployment

- [ ] Add `server/Dockerfile` using a multi-stage Go build and SQLite data volume.
- [ ] Add `app/Dockerfile` building H5 and serving it with Nginx.
- [ ] Add root `nginx.conf` routing `/api/` to the backend and all other routes to the H5 container.
- [ ] Add root `docker-compose.yml` wiring `nginx`, `server`, and `app_h5`.

## Task 4: Documentation

- [ ] Update README quick start with Docker Compose.
- [ ] Update local runbook with config and Docker commands.
- [ ] Update docs-site quick start and deploy docs with the same commands.

## Task 5: Verification

- [ ] Run `npm run verify`.
- [ ] Run `go test ./...` in `server` if config tests fail under the root script.
- [ ] Report that git commits are unavailable when the directory is not a git repository.

## Self-Review

- Spec coverage: config loading, H5 API base, Docker deployment, and docs are covered.
- Placeholder scan: no placeholder markers are required.
- Type consistency: config keys match `config.example.yaml`; environment variables use the `SANJIE_` prefix.
