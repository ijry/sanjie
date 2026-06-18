# Sanjie Embedded H5 Demo Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Provide a Lyshop-style embedded H5 demo in the VitePress docs site with complete frontend mock data.

**Architecture:** Build the uni-app H5 app in mock mode, copy the static output into `docs-site/docs/public/demo/`, and embed it from the docs demo page using an iframe phone frame. Keep mock data in frontend memory so static hosting needs no backend.

**Tech Stack:** uni-app, Vue 3, uview-plus, VitePress, GitHub Actions, Node.js scripts.

---

### Task 1: Demo Build Pipeline

**Files:**
- Create: `app/scripts/build-demo.js`
- Modify: `app/package.json`
- Modify: `docs-site/package.json`
- Modify: `package.json`

- [ ] Add a Node script that runs uni-app H5 build with `VITE_SANJIE_MOCK=true`, temporarily enforces hash router and relative public path, then copies `app/dist/build/h5` to `docs-site/docs/public/demo`.
- [ ] Add package scripts `build:demo`, `demo:build`, and root `docs:build-with-demo`.
- [ ] Run `npm --prefix app run build:demo` and confirm `docs-site/docs/public/demo/index.html` exists.

### Task 2: Complete Mock Dataset

**Files:**
- Modify: `app/mock/data.js`

- [ ] Expand mock data across users, alerts, souls, capture tasks, life book records, approvals, reincarnations, all 18 hell floors, hell sentences, wishes, soup inventory, soup records, and audit logs.
- [ ] Keep existing field names and state transitions compatible with current pages and mock server.
- [ ] Ensure IDs referenced across entities resolve.

### Task 3: Embedded Docs Demo

**Files:**
- Modify: `docs-site/docs/demo/index.md`
- Modify: `README.md`
- Modify: `docs-site/docs/deploy/index.md`

- [ ] Replace the plain demo instructions with a phone-frame iframe pointing to `/sanjie/demo/index.html#/pages/index/index`.
- [ ] Document direct demo URL, mock reset behavior, and how to rebuild embedded demo locally.
- [ ] Update deployment docs to explain that docs build includes embedded H5 demo output.

### Task 4: CI and Docs Publish

**Files:**
- Modify: `.github/workflows/docs.yml`
- Modify: `.github/workflows/ci.yml`

- [ ] Make Docs workflow run `npm run demo:build` before `docs:build`.
- [ ] Add CI coverage for `npm --prefix app run build:demo` so embedded demo generation stays verified.

### Task 5: Verify and Publish

**Files:**
- No source files expected beyond prior tasks.

- [ ] Run `npm run verify`.
- [ ] Run `npm run docs:build-with-demo`.
- [ ] Inspect `git status` for ignored build output behavior and intended tracked demo files.
- [ ] Commit with Chinese `head + body`.
- [ ] Push to `origin/master` and confirm GitHub Actions status.
