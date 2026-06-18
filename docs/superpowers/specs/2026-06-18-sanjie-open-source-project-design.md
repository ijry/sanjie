# Sanjie Open Source Project Design

## Goal

Turn `D:\Repos\sanjie` into an open-source project named `sanjie`, following the public-facing structure and documentation style of `D:\Repos\xyito\open\lyshop`, while keeping the scope aligned to this product: a joke-style 三界管理系统 mobile app built with `uni-app + uview-plus` and backed by `Go + SQLite`.

## Confirmed Defaults

- Project slug: `sanjie`
- Chinese name: `Sanjie 三界管理系统`
- Go module: `github.com/ijry/sanjie`
- Docs site base path: `/sanjie/`
- License: MIT
- Docs site: VitePress
- Main product targets: mobile app and backend API

## Repository Structure

The initial open-source scaffold should use this layout:

```text
sanjie/
├─ .github/
│  └─ workflows/
│     ├─ ci.yml
│     └─ docs.yml
├─ app/
├─ assets/
├─ docs/
│  └─ superpowers/
│     ├─ plans/
│     └─ specs/
├─ docs-site/
│  ├─ package.json
│  └─ docs/
│     ├─ index.md
│     ├─ README.md
│     ├─ guide/
│     │  ├─ features.md
│     │  └─ architecture.md
│     ├─ api/
│     │  └─ index.md
│     ├─ deploy/
│     │  └─ index.md
│     └─ dev/
│        └─ roadmap.md
├─ scripts/
├─ server/
├─ .gitignore
├─ AGENTS.md
├─ config.example.yaml
├─ LICENSE
└─ README.md
```

`app/` and `server/` can start as documented placeholders if the application implementation has not been generated yet. The root project should still be usable as an open-source repository immediately: a reader can understand the purpose, run docs locally, and see how app/server will be implemented.

## Root Documentation

`README.md` should mirror the practical shape of `lyshop`:

- Project title and one-sentence description.
- CI, Docs, and MIT license badges.
- Short product screenshot placeholder section pointing to `assets/`.
- Feature list focused on the 三界 domain.
- Tech stack table.
- Directory structure.
- Quick start for backend, mobile app, and docs site.
- Link list for docs, API, deployment, and development plan.
- MIT license reference.

The README should not claim production readiness. It should position the repository as an open-source demo/product prototype.

## Docs Site

Use VitePress under `docs-site/`.

Required docs pages:

- `docs-site/docs/index.md`: docs home page with project intro and entry cards.
- `docs-site/docs/README.md`: quick start page.
- `docs-site/docs/guide/features.md`: feature modules, including workbench, 勾魂, 投胎, 审批, 生死簿, 十八层地狱, 愿望工单, 孟婆汤, 审计日志.
- `docs-site/docs/guide/architecture.md`: high-level architecture for uni-app, Go API, SQLite, and docs-site.
- `docs-site/docs/api/index.md`: API groups and response envelope.
- `docs-site/docs/deploy/index.md`: local development, SQLite default, and future deployment notes.
- `docs-site/docs/dev/roadmap.md`: implementation phases, referencing the existing full implementation plan.

The VitePress config should:

- Set title to `Sanjie`.
- Set description to `玩笑性质的三界移动办公系统`.
- Set base to `/sanjie/`.
- Include nav entries for Guide, API, Deploy, and Development.
- Include sidebar groups matching the docs pages.
- Point social link to `https://github.com/ijry/sanjie`.

Generated VitePress output such as `.vitepress/dist` must not be committed.

## CI and Docs Workflows

Create `.github/workflows/ci.yml`:

- Trigger on push and pull request to `master`.
- Run Go tests when `server/go.mod` exists.
- Build docs-site with Node 20.
- Keep app build out of CI until the app scaffold is implemented.

Create `.github/workflows/docs.yml`:

- Trigger on push to `master` when `docs-site/**`, `.github/workflows/docs.yml`, or root docs files change.
- Build VitePress docs.
- Upload `docs-site/docs/.vitepress/dist` to GitHub Pages.
- Deploy only on push to `master`.

## Collaboration Rules

Create `AGENTS.md` adapted from `lyshop`:

- Git commit messages use Chinese.
- Functional changes must update `docs-site`.
- Prefer upgrading existing API semantics before adding new endpoints.
- Docs describe the current architecture directly.
- Mobile UI uses uview-plus components first.

## Configuration

Create `config.example.yaml`:

```yaml
server:
  addr: ":8080"

database:
  driver: "sqlite"
  dsn: "sanjie.db"

demo:
  seed: true
```

This is enough for the planned Go + SQLite backend.

## Git Ignore

`.gitignore` should cover:

- Node dependencies and build outputs.
- VitePress dist/cache.
- Go binaries and coverage files.
- SQLite runtime database files.
- local env/config files.
- editor and OS noise.

Do not ignore the existing `docs/superpowers` planning files.

## Open Questions Resolved

- The project will use `/sanjie/` as the docs base path.
- The initial docs-site will not include live demos yet.
- The root open-source scaffold can be created before app/server implementation.
- The existing full implementation plan remains under `docs/superpowers/plans/`.

## Acceptance Criteria

- Root `README.md`, `LICENSE`, `AGENTS.md`, `.gitignore`, and `config.example.yaml` exist.
- `docs-site` can install dependencies and build VitePress docs.
- `.github/workflows/ci.yml` and `docs.yml` exist and are suitable for GitHub.
- Docs pages explain the project, architecture, API shape, deployment, and roadmap.
- Existing implementation plan remains intact.
- No generated docs output is committed.

