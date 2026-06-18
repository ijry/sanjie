# Sanjie H5 Mock Demo Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Add a backend-free H5 mock demo mode for Sanjie so the mobile UI can be demonstrated from docs or local builds without starting the Go API.

**Architecture:** Keep existing API wrapper functions unchanged. Add a mock adapter behind `app/api/http.js` that intercepts `request()` when `VITE_SANJIE_MOCK=true`, serves in-memory seed data from `app/mock/`, and simulates common state-changing POST endpoints.

**Tech Stack:** uni-app, Vue 3, uview-plus, Vite env vars, VitePress docs.

---

## Files

- Create: `app/mock/data.js`
- Create: `app/mock/server.js`
- Modify: `app/api/http.js`
- Modify: `app/package.json`
- Modify: `package.json`
- Modify: `.github/workflows/ci.yml`
- Create: `docs-site/docs/demo/index.md`
- Modify: `docs-site/docs/.vitepress/config.mts`
- Modify: `README.md`
- Modify: `docs/runbook.md`
- Modify: `docs-site/docs/README.md`
- Modify: `docs-site/docs/deploy/index.md`

## Task 1: Mock Adapter

- [ ] Add mock seed data for users, metrics, alerts, capture tasks, life book, approvals, reincarnations, hell floors, wishes, soup, and audit logs.
- [ ] Add a request router that supports all mobile-read endpoints and common state-changing POST endpoints.
- [ ] Update `app/api/http.js` to call the mock router when `VITE_SANJIE_MOCK=true`.

## Task 2: Scripts and CI

- [ ] Add `dev:h5:mock` and `build:h5:mock` under `app/package.json`.
- [ ] Add root `dev:h5:mock` and `app:build:mock`.
- [ ] Add H5 mock build to CI.

## Task 3: Docs

- [ ] Add docs-site demo page describing mock mode and future online demo artifact.
- [ ] Link demo page from VitePress nav/sidebar and README.
- [ ] Update runbook and deploy docs with mock commands.

## Task 4: Verification

- [ ] Run `npm --prefix app run build:h5:mock`.
- [ ] Run `npm run verify`.

## Self-Review

- Spec coverage: mock request adapter, scripts, CI, and docs are covered.
- Placeholder scan: no placeholders are required.
- Type consistency: mock field names mirror the current mobile pages.
