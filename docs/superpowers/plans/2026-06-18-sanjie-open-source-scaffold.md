# Sanjie Open Source Scaffold Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Convert `D:\Repos\sanjie` into an open-source repository scaffold named `sanjie`, with root project documentation, GitHub workflows, configuration defaults, and a VitePress docs site.

**Architecture:** Keep the repository lightweight and immediately readable. The root files describe the public project; `docs-site/` owns published documentation; `app/` and `server/` remain product implementation roots for the planned uni-app mobile client and Go + SQLite API.

**Tech Stack:** Markdown, VitePress, GitHub Actions, Node 20, Go 1.22, uni-app, uview-plus, SQLite.

---

## Files

- Create: `.gitignore`
- Create: `AGENTS.md`
- Create: `LICENSE`
- Create: `README.md`
- Create: `config.example.yaml`
- Create: `.github/workflows/ci.yml`
- Create: `.github/workflows/docs.yml`
- Create: `docs-site/package.json`
- Create: `docs-site/docs/.vitepress/config.mts`
- Create: `docs-site/docs/index.md`
- Create: `docs-site/docs/README.md`
- Create: `docs-site/docs/guide/features.md`
- Create: `docs-site/docs/guide/architecture.md`
- Create: `docs-site/docs/api/index.md`
- Create: `docs-site/docs/deploy/index.md`
- Create: `docs-site/docs/dev/roadmap.md`
- Create: `app/.gitkeep`
- Create: `server/.gitkeep`
- Create: `assets/.gitkeep`
- Create: `scripts/.gitkeep`
- Preserve: `docs/superpowers/plans/2026-06-18-sanjie-mobile-full-implementation.md`
- Preserve: `docs/superpowers/specs/2026-06-18-sanjie-open-source-project-design.md`

## Task 1: Root Open Source Files

- [ ] **Step 1: Create `.gitignore`**

Create ignore rules for Node, VitePress output, Go build output, SQLite runtime files, local config, and editor noise.

- [ ] **Step 2: Create `LICENSE`**

Use the MIT License with copyright holder `ijry`.

- [ ] **Step 3: Create `config.example.yaml`**

Use:

```yaml
server:
  addr: ":8080"

database:
  driver: "sqlite"
  dsn: "sanjie.db"

demo:
  seed: true
```

- [ ] **Step 4: Create `AGENTS.md`**

Use Chinese collaboration rules adapted from `lyshop`: Chinese commit messages, docs-site sync for functional changes, prefer existing API semantics, docs describe current architecture, mobile UI prefers uview-plus.

- [ ] **Step 5: Create `README.md`**

Include project title, badges, intro, features, tech stack, directory structure, quick start, docs links, and MIT license reference.

## Task 2: VitePress Docs Site

- [ ] **Step 1: Create `docs-site/package.json`**

Use scripts:

```json
{
  "docs:dev": "vitepress dev docs",
  "docs:build": "vitepress build docs",
  "docs:preview": "vitepress preview docs"
}
```

Depend on `vitepress`.

- [ ] **Step 2: Create VitePress config**

Create `docs-site/docs/.vitepress/config.mts` with base `/sanjie/`, title `Sanjie`, GitHub social link, nav, and sidebar groups.

- [ ] **Step 3: Create docs home**

Create `docs-site/docs/index.md` with project intro, feature cards, quick links, and current status.

- [ ] **Step 4: Create quick start**

Create `docs-site/docs/README.md` with backend, app, and docs-site commands.

- [ ] **Step 5: Create feature guide**

Create `docs-site/docs/guide/features.md` covering all planned business modules.

- [ ] **Step 6: Create architecture guide**

Create `docs-site/docs/guide/architecture.md` explaining app, server, SQLite, docs-site, API envelope, and module boundaries.

- [ ] **Step 7: Create API guide**

Create `docs-site/docs/api/index.md` with response envelope and endpoint groups.

- [ ] **Step 8: Create deploy guide**

Create `docs-site/docs/deploy/index.md` with local development and future deployment notes.

- [ ] **Step 9: Create roadmap**

Create `docs-site/docs/dev/roadmap.md` referencing the full implementation plan and phased delivery.

## Task 3: GitHub Workflows and Placeholder Dirs

- [ ] **Step 1: Create `ci.yml`**

Create a workflow that builds docs-site and runs Go tests only when `server/go.mod` exists.

- [ ] **Step 2: Create `docs.yml`**

Create a GitHub Pages workflow that builds VitePress and deploys on push to `master`.

- [ ] **Step 3: Create placeholder dirs**

Create `.gitkeep` files in `app/`, `server/`, `assets/`, and `scripts/`.

## Task 4: Verification

- [ ] **Step 1: Check placeholders**

Run:

```powershell
rg -n "TBD|TODO|implement later|fill in details" .
```

Expected: no matches except unrelated third-party content if dependencies are installed.

- [ ] **Step 2: Install docs dependencies**

Run:

```powershell
cd docs-site
npm install
```

Expected: `package-lock.json` is created.

- [ ] **Step 3: Build docs**

Run:

```powershell
cd docs-site
npm run docs:build
```

Expected: VitePress build succeeds and writes ignored output to `docs-site/docs/.vitepress/dist`.

- [ ] **Step 4: Check git status**

Run:

```powershell
git status --short
```

Expected: changed files are the scaffold files and existing docs. If the directory is not a git repository, report that commits cannot be made.

---

## Self-Review

- Spec coverage: root docs, docs-site, CI/docs workflows, config, license, ignore rules, placeholders, and verification are covered.
- Placeholder scan: no placeholder markers are required in generated files.
- Type consistency: project slug, docs base, Go module, and GitHub URL use `sanjie`.

