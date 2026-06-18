# Sanjie Mobile Full Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Build a complete joke-style "三界管理系统" mobile app using uni-app + uview-plus, backed by a Go + SQLite API that implements all planned modules: workbench, soul capture, reincarnation, approvals, life book, hell floors, wish tickets, alerts, audit logs, users, and role-aware workflows.

**Architecture:** Use a monorepo with `app/` for the uni-app client and `server/` for the Go API. The backend owns all business state in SQLite and exposes JSON REST endpoints; the mobile client uses typed API wrappers, uview-plus components, and role-aware page actions. Start with stable backend foundations, then implement one vertical business module at a time from database to API to mobile UI.

**Tech Stack:** uni-app, Vue 3, uview-plus, Pinia or uni-app global store, Go 1.22+, SQLite, `database/sql`, `modernc.org/sqlite` or `github.com/mattn/go-sqlite3`, `net/http` or Chi router, Go tests, HBuilderX or CLI build flow.

---

## Scope

This plan implements the full mobile product, not only an MVP.

Included modules:

- Workbench dashboard and alerts.
- User login mock, current user, role-aware actions.
- Soul records and life book search.
- Soul capture task list, detail, start, complete, exception reporting.
- Reincarnation queue, detail, approve, reject, review, quota assignment, soup status update.
- Approval center for lifespan changes, reincarnation jumps, capture exceptions, hell reviews, deity travel filings.
- Eighteen hell floors, sentence records, floor dispatch, sentence review.
- Wish tickets, routing to deity owners, approve/reject/resolve.
- Mengpo soup inventory and issue records.
- Audit logs for sensitive operations.
- Seed data with humorous but stable demo records.

Explicitly out of scope for the first complete version:

- Real authentication, password hashing, JWT refresh, and third-party login.
- Real map SDK integration for capture locations.
- Push notifications.
- Cloud deployment.
- Multi-tenant support.

---

## Repository Layout

Create this structure:

```text
sanjie/
├─ app/
│  ├─ App.vue
│  ├─ main.js
│  ├─ manifest.json
│  ├─ pages.json
│  ├─ uni.scss
│  ├─ api/
│  │  ├─ http.js
│  │  ├─ dashboard.js
│  │  ├─ users.js
│  │  ├─ souls.js
│  │  ├─ lifeBook.js
│  │  ├─ capture.js
│  │  ├─ reincarnation.js
│  │  ├─ approvals.js
│  │  ├─ hell.js
│  │  ├─ wishes.js
│  │  ├─ soup.js
│  │  └─ audit.js
│  ├─ stores/
│  │  └─ user.js
│  ├─ utils/
│  │  ├─ status.js
│  │  ├─ format.js
│  │  └─ auth.js
│  ├─ components/
│  │  ├─ StatusTag.vue
│  │  ├─ MetricCard.vue
│  │  ├─ EmptyState.vue
│  │  ├─ FixedActionBar.vue
│  │  └─ SectionTitle.vue
│  └─ pages/
│     ├─ index/index.vue
│     ├─ login/index.vue
│     ├─ capture/index.vue
│     ├─ capture/detail.vue
│     ├─ reincarnation/index.vue
│     ├─ reincarnation/detail.vue
│     ├─ approval/index.vue
│     ├─ approval/detail.vue
│     ├─ life-book/search.vue
│     ├─ life-book/detail.vue
│     ├─ hell/index.vue
│     ├─ hell/sentence-detail.vue
│     ├─ wish/index.vue
│     ├─ wish/detail.vue
│     ├─ soup/index.vue
│     ├─ audit/index.vue
│     └─ mine/index.vue
├─ server/
│  ├─ go.mod
│  ├─ main.go
│  ├─ sanjie.db
│  ├─ internal/
│  │  ├─ app/
│  │  │  └─ app.go
│  │  ├─ db/
│  │  │  ├─ db.go
│  │  │  ├─ migrate.go
│  │  │  └─ seed.go
│  │  ├─ httpx/
│  │  │  ├─ response.go
│  │  │  └─ middleware.go
│  │  ├─ model/
│  │  │  ├─ model.go
│  │  │  └─ constants.go
│  │  ├─ repository/
│  │  │  ├─ users.go
│  │  │  ├─ dashboard.go
│  │  │  ├─ souls.go
│  │  │  ├─ life_book.go
│  │  │  ├─ capture.go
│  │  │  ├─ reincarnation.go
│  │  │  ├─ approvals.go
│  │  │  ├─ hell.go
│  │  │  ├─ wishes.go
│  │  │  ├─ soup.go
│  │  │  └─ audit.go
│  │  ├─ service/
│  │  │  ├─ dashboard.go
│  │  │  ├─ capture.go
│  │  │  ├─ reincarnation.go
│  │  │  ├─ approvals.go
│  │  │  ├─ life_book.go
│  │  │  ├─ hell.go
│  │  │  ├─ wishes.go
│  │  │  └─ soup.go
│  │  └─ handler/
│  │     ├─ routes.go
│  │     ├─ users.go
│  │     ├─ dashboard.go
│  │     ├─ souls.go
│  │     ├─ life_book.go
│  │     ├─ capture.go
│  │     ├─ reincarnation.go
│  │     ├─ approvals.go
│  │     ├─ hell.go
│  │     ├─ wishes.go
│  │     ├─ soup.go
│  │     └─ audit.go
│  └─ internal_test/
│     ├─ testdb.go
│     ├─ capture_test.go
│     ├─ reincarnation_test.go
│     ├─ approvals_test.go
│     ├─ life_book_test.go
│     ├─ hell_test.go
│     ├─ wishes_test.go
│     └─ dashboard_test.go
└─ docs/
   ├─ api.md
   └─ runbook.md
```

---

## Backend API Contract

Use this response envelope for every JSON endpoint:

```json
{
  "code": 0,
  "message": "ok",
  "data": {}
}
```

Use this error shape:

```json
{
  "code": 40001,
  "message": "invalid status transition",
  "data": null
}
```

Use mock auth headers:

```text
X-User-ID: 1
```

If the header is absent, default to user `1` for local demo. The current user determines role-sensitive action availability and audit log actor fields.

Endpoint list:

```text
GET    /api/health
GET    /api/users/me
POST   /api/users/switch

GET    /api/dashboard
GET    /api/alerts
POST   /api/alerts/{id}/handle

GET    /api/souls
GET    /api/souls/{id}

GET    /api/life-book
GET    /api/life-book/{id}
POST   /api/life-book/{id}/change-lifespan
POST   /api/life-book/{id}/freeze
GET    /api/life-book/{id}/audit-logs

GET    /api/capture-tasks
GET    /api/capture-tasks/{id}
POST   /api/capture-tasks/{id}/start
POST   /api/capture-tasks/{id}/complete
POST   /api/capture-tasks/{id}/exception
POST   /api/capture-tasks/{id}/cancel

GET    /api/reincarnations
GET    /api/reincarnations/{id}
POST   /api/reincarnations/{id}/approve
POST   /api/reincarnations/{id}/reject
POST   /api/reincarnations/{id}/review
POST   /api/reincarnations/{id}/assign-quota
POST   /api/reincarnations/{id}/soup
POST   /api/reincarnations/{id}/complete

GET    /api/approvals
GET    /api/approvals/{id}
POST   /api/approvals/{id}/approve
POST   /api/approvals/{id}/reject
POST   /api/approvals/{id}/transfer

GET    /api/hell/floors
GET    /api/hell/sentences
GET    /api/hell/sentences/{id}
POST   /api/hell/sentences/{id}/review
POST   /api/hell/floors/{id}/dispatch

GET    /api/wishes
GET    /api/wishes/{id}
POST   /api/wishes/{id}/route
POST   /api/wishes/{id}/resolve
POST   /api/wishes/{id}/reject

GET    /api/soup/inventory
POST   /api/soup/inventory/{id}/adjust
GET    /api/soup/records

GET    /api/audit-logs
```

---

## Database Schema

Implement migrations in Go, not external SQL tooling. The migration must create these tables exactly.

```sql
CREATE TABLE IF NOT EXISTS users (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  username TEXT NOT NULL UNIQUE,
  nickname TEXT NOT NULL,
  role TEXT NOT NULL,
  avatar TEXT NOT NULL DEFAULT '',
  created_at TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS souls (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  name TEXT NOT NULL,
  birth_info TEXT NOT NULL,
  death_info TEXT NOT NULL DEFAULT '',
  status TEXT NOT NULL,
  merit_score INTEGER NOT NULL DEFAULT 0,
  karma_score INTEGER NOT NULL DEFAULT 0,
  memory_residue INTEGER NOT NULL DEFAULT 0,
  relationship_risk TEXT NOT NULL DEFAULT 'low',
  created_at TEXT NOT NULL,
  updated_at TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS life_book_records (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  soul_id INTEGER NOT NULL,
  expected_lifespan INTEGER NOT NULL,
  actual_lifespan INTEGER NOT NULL,
  death_reason TEXT NOT NULL,
  fate_level TEXT NOT NULL,
  calamity_count INTEGER NOT NULL DEFAULT 0,
  locked INTEGER NOT NULL DEFAULT 0,
  risk_flag TEXT NOT NULL DEFAULT 'normal',
  updated_at TEXT NOT NULL,
  FOREIGN KEY (soul_id) REFERENCES souls(id)
);

CREATE TABLE IF NOT EXISTS capture_tasks (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  soul_id INTEGER NOT NULL,
  assignee_id INTEGER,
  status TEXT NOT NULL,
  location TEXT NOT NULL,
  scheduled_time TEXT NOT NULL,
  actual_time TEXT NOT NULL DEFAULT '',
  exception_type TEXT NOT NULL DEFAULT '',
  exception_note TEXT NOT NULL DEFAULT '',
  created_at TEXT NOT NULL,
  updated_at TEXT NOT NULL,
  FOREIGN KEY (soul_id) REFERENCES souls(id),
  FOREIGN KEY (assignee_id) REFERENCES users(id)
);

CREATE TABLE IF NOT EXISTS reincarnation_orders (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  soul_id INTEGER NOT NULL,
  target_realm TEXT NOT NULL,
  priority INTEGER NOT NULL DEFAULT 0,
  soup_status TEXT NOT NULL,
  quota_type TEXT NOT NULL DEFAULT '',
  status TEXT NOT NULL,
  review_note TEXT NOT NULL DEFAULT '',
  created_at TEXT NOT NULL,
  updated_at TEXT NOT NULL,
  FOREIGN KEY (soul_id) REFERENCES souls(id)
);

CREATE TABLE IF NOT EXISTS approvals (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  type TEXT NOT NULL,
  target_id INTEGER NOT NULL,
  title TEXT NOT NULL,
  status TEXT NOT NULL,
  applicant_id INTEGER,
  approver_id INTEGER,
  reason TEXT NOT NULL DEFAULT '',
  result_note TEXT NOT NULL DEFAULT '',
  created_at TEXT NOT NULL,
  updated_at TEXT NOT NULL,
  FOREIGN KEY (applicant_id) REFERENCES users(id),
  FOREIGN KEY (approver_id) REFERENCES users(id)
);

CREATE TABLE IF NOT EXISTS hell_floors (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  floor_no INTEGER NOT NULL UNIQUE,
  name TEXT NOT NULL,
  capacity INTEGER NOT NULL,
  occupancy INTEGER NOT NULL,
  equipment_status TEXT NOT NULL,
  load_level TEXT NOT NULL,
  updated_at TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS hell_sentences (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  soul_id INTEGER NOT NULL,
  floor_id INTEGER NOT NULL,
  crime_type TEXT NOT NULL,
  sentence_days INTEGER NOT NULL,
  pain_level INTEGER NOT NULL,
  review_status TEXT NOT NULL,
  equipment_id TEXT NOT NULL DEFAULT '',
  created_at TEXT NOT NULL,
  updated_at TEXT NOT NULL,
  FOREIGN KEY (soul_id) REFERENCES souls(id),
  FOREIGN KEY (floor_id) REFERENCES hell_floors(id)
);

CREATE TABLE IF NOT EXISTS wish_tickets (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  title TEXT NOT NULL,
  wish_type TEXT NOT NULL,
  requester_name TEXT NOT NULL,
  incense_amount INTEGER NOT NULL DEFAULT 0,
  merit_score INTEGER NOT NULL DEFAULT 0,
  status TEXT NOT NULL,
  assigned_deity TEXT NOT NULL DEFAULT '',
  result_note TEXT NOT NULL DEFAULT '',
  created_at TEXT NOT NULL,
  updated_at TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS soup_inventory (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  name TEXT NOT NULL,
  stock INTEGER NOT NULL,
  warning_stock INTEGER NOT NULL,
  formula_note TEXT NOT NULL DEFAULT '',
  updated_at TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS soup_records (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  soul_id INTEGER NOT NULL,
  order_id INTEGER,
  inventory_id INTEGER NOT NULL,
  dose INTEGER NOT NULL,
  memory_after INTEGER NOT NULL,
  operator_id INTEGER,
  created_at TEXT NOT NULL,
  FOREIGN KEY (soul_id) REFERENCES souls(id),
  FOREIGN KEY (order_id) REFERENCES reincarnation_orders(id),
  FOREIGN KEY (inventory_id) REFERENCES soup_inventory(id),
  FOREIGN KEY (operator_id) REFERENCES users(id)
);

CREATE TABLE IF NOT EXISTS alerts (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  level TEXT NOT NULL,
  title TEXT NOT NULL,
  content TEXT NOT NULL,
  handled INTEGER NOT NULL DEFAULT 0,
  created_at TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS audit_logs (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  actor_id INTEGER,
  action TEXT NOT NULL,
  target_type TEXT NOT NULL,
  target_id INTEGER NOT NULL,
  before_json TEXT NOT NULL DEFAULT '',
  after_json TEXT NOT NULL DEFAULT '',
  note TEXT NOT NULL DEFAULT '',
  created_at TEXT NOT NULL,
  FOREIGN KEY (actor_id) REFERENCES users(id)
);
```

---

## Status Constants

Define these constants once in `server/internal/model/constants.go` and mirror them in `app/utils/status.js`.

```text
roles: jade_emperor, yanwang, judge, impermanence, mengpo, heavenly_general, city_god, ghost_clerk, auditor

soul_status: alive, pending_capture, judging, hell, reincarnation_waiting, reincarnated

capture_status: pending, running, completed, exception, canceled

capture_exception_type: emergency_rescue, wrong_soul, escaped, not_found, ascended

reincarnation_status: queued, pending_soup, approved, rejected, reviewing, completed, exception

soup_status: pending, issued, failed, not_required

approval_type: lifespan_change, reincarnation_jump, capture_exception, hell_review, deity_travel

approval_status: pending, approved, rejected, transferred

hell_review_status: none, reviewing, changed, confirmed

wish_status: pending, routed, resolved, rejected

alert_level: info, warning, critical
```

---

## Implementation Tasks

### Task 1: Backend Module Scaffold

**Files:**

- Create: `server/go.mod`
- Create: `server/main.go`
- Create: `server/internal/app/app.go`
- Create: `server/internal/httpx/response.go`
- Create: `server/internal/httpx/middleware.go`
- Create: `server/internal/handler/routes.go`
- Test: `server/internal_test/testdb.go`

- [ ] **Step 1: Initialize Go module**

Run:

```powershell
cd server
go mod init sanjie/server
go get github.com/go-chi/chi/v5 modernc.org/sqlite
```

Expected: `go.mod` and `go.sum` exist.

- [ ] **Step 2: Create JSON response helpers**

Create `server/internal/httpx/response.go`:

```go
package httpx

import (
	"encoding/json"
	"net/http"
)

type Envelope struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func OK(w http.ResponseWriter, data interface{}) {
	Write(w, http.StatusOK, Envelope{Code: 0, Message: "ok", Data: data})
}

func Fail(w http.ResponseWriter, status int, code int, message string) {
	Write(w, status, Envelope{Code: code, Message: message, Data: nil})
}

func Write(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}
```

- [ ] **Step 3: Create middleware**

Create `server/internal/httpx/middleware.go`:

```go
package httpx

import (
	"context"
	"net/http"
	"strconv"
)

type contextKey string

const userIDKey contextKey = "userID"

func DemoAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID := int64(1)
		if raw := r.Header.Get("X-User-ID"); raw != "" {
			if parsed, err := strconv.ParseInt(raw, 10, 64); err == nil && parsed > 0 {
				userID = parsed
			}
		}
		ctx := context.WithValue(r.Context(), userIDKey, userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func CurrentUserID(r *http.Request) int64 {
	if id, ok := r.Context().Value(userIDKey).(int64); ok && id > 0 {
		return id
	}
	return 1
}
```

- [ ] **Step 4: Create app wiring**

Create `server/internal/app/app.go`:

```go
package app

import (
	"database/sql"
	"net/http"

	"sanjie/server/internal/handler"
	"sanjie/server/internal/httpx"
)

type App struct {
	DB     *sql.DB
	Router http.Handler
}

func New(db *sql.DB) *App {
	router := handler.Routes(db)
	return &App{
		DB:     db,
		Router: httpx.DemoAuth(router),
	}
}
```

- [ ] **Step 5: Create routes with health endpoint**

Create `server/internal/handler/routes.go`:

```go
package handler

import (
	"database/sql"
	"net/http"

	"github.com/go-chi/chi/v5"

	"sanjie/server/internal/httpx"
)

func Routes(db *sql.DB) http.Handler {
	r := chi.NewRouter()

	r.Get("/api/health", func(w http.ResponseWriter, r *http.Request) {
		httpx.OK(w, map[string]string{"status": "ok"})
	})

	return r
}
```

- [ ] **Step 6: Create main**

Create `server/main.go`:

```go
package main

import (
	"log"
	"net/http"

	"sanjie/server/internal/app"
	"sanjie/server/internal/db"
)

func main() {
	database, err := db.Open("sanjie.db")
	if err != nil {
		log.Fatal(err)
	}
	defer database.Close()

	if err := db.Migrate(database); err != nil {
		log.Fatal(err)
	}
	if err := db.Seed(database); err != nil {
		log.Fatal(err)
	}

	application := app.New(database)
	log.Println("sanjie server listening on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", application.Router))
}
```

- [ ] **Step 7: Create temporary DB helpers**

Create `server/internal_test/testdb.go`:

```go
package internal_test

import (
	"database/sql"
	"testing"

	_ "modernc.org/sqlite"

	"sanjie/server/internal/db"
)

func NewTestDB(t *testing.T) *sql.DB {
	t.Helper()

	database, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatalf("open test db: %v", err)
	}
	t.Cleanup(func() { _ = database.Close() })

	if err := db.Migrate(database); err != nil {
		t.Fatalf("migrate test db: %v", err)
	}
	if err := db.Seed(database); err != nil {
		t.Fatalf("seed test db: %v", err)
	}
	return database
}
```

- [ ] **Step 8: Run compile and expect db package missing**

Run:

```powershell
cd server
go test ./...
```

Expected: FAIL because `server/internal/db` is not implemented yet.

- [ ] **Step 9: Commit**

Run:

```powershell
git add server
git commit -m "chore: scaffold sanjie go server"
```

Expected: commit succeeds if this directory is a git repository.

### Task 2: Backend Database Migration and Seed Data

**Files:**

- Create: `server/internal/db/db.go`
- Create: `server/internal/db/migrate.go`
- Create: `server/internal/db/seed.go`
- Create: `server/internal/model/model.go`
- Create: `server/internal/model/constants.go`
- Modify: `server/internal_test/testdb.go`

- [ ] **Step 1: Define constants**

Create `server/internal/model/constants.go`:

```go
package model

const (
	RoleJadeEmperor     = "jade_emperor"
	RoleYanwang         = "yanwang"
	RoleJudge           = "judge"
	RoleImpermanence    = "impermanence"
	RoleMengpo          = "mengpo"
	RoleHeavenlyGeneral = "heavenly_general"
	RoleCityGod         = "city_god"
	RoleGhostClerk      = "ghost_clerk"
	RoleAuditor         = "auditor"

	SoulAlive                = "alive"
	SoulPendingCapture       = "pending_capture"
	SoulJudging              = "judging"
	SoulHell                 = "hell"
	SoulReincarnationWaiting = "reincarnation_waiting"
	SoulReincarnated         = "reincarnated"

	CapturePending   = "pending"
	CaptureRunning   = "running"
	CaptureCompleted = "completed"
	CaptureException = "exception"
	CaptureCanceled  = "canceled"

	ReincarnationQueued      = "queued"
	ReincarnationPendingSoup = "pending_soup"
	ReincarnationApproved    = "approved"
	ReincarnationRejected    = "rejected"
	ReincarnationReviewing   = "reviewing"
	ReincarnationCompleted   = "completed"
	ReincarnationException   = "exception"

	SoupPending     = "pending"
	SoupIssued      = "issued"
	SoupFailed      = "failed"
	SoupNotRequired = "not_required"

	ApprovalPending     = "pending"
	ApprovalApproved    = "approved"
	ApprovalRejected    = "rejected"
	ApprovalTransferred = "transferred"

	WishPending  = "pending"
	WishRouted   = "routed"
	WishResolved = "resolved"
	WishRejected = "rejected"
)
```

- [ ] **Step 2: Define shared models**

Create `server/internal/model/model.go`:

```go
package model

type User struct {
	ID        int64  `json:"id"`
	Username  string `json:"username"`
	Nickname  string `json:"nickname"`
	Role      string `json:"role"`
	Avatar    string `json:"avatar"`
	CreatedAt string `json:"createdAt"`
}

type Soul struct {
	ID               int64  `json:"id"`
	Name             string `json:"name"`
	BirthInfo        string `json:"birthInfo"`
	DeathInfo        string `json:"deathInfo"`
	Status           string `json:"status"`
	MeritScore       int    `json:"meritScore"`
	KarmaScore       int    `json:"karmaScore"`
	MemoryResidue    int    `json:"memoryResidue"`
	RelationshipRisk string `json:"relationshipRisk"`
	CreatedAt        string `json:"createdAt"`
	UpdatedAt        string `json:"updatedAt"`
}

type LifeBookRecord struct {
	ID               int64  `json:"id"`
	SoulID           int64  `json:"soulId"`
	SoulName         string `json:"soulName,omitempty"`
	ExpectedLifespan int    `json:"expectedLifespan"`
	ActualLifespan   int    `json:"actualLifespan"`
	DeathReason      string `json:"deathReason"`
	FateLevel         string `json:"fateLevel"`
	CalamityCount    int    `json:"calamityCount"`
	Locked           bool   `json:"locked"`
	RiskFlag         string `json:"riskFlag"`
	UpdatedAt        string `json:"updatedAt"`
}

type CaptureTask struct {
	ID            int64  `json:"id"`
	SoulID        int64  `json:"soulId"`
	SoulName      string `json:"soulName,omitempty"`
	AssigneeID    int64  `json:"assigneeId"`
	AssigneeName  string `json:"assigneeName,omitempty"`
	Status        string `json:"status"`
	Location      string `json:"location"`
	ScheduledTime string `json:"scheduledTime"`
	ActualTime    string `json:"actualTime"`
	ExceptionType string `json:"exceptionType"`
	ExceptionNote string `json:"exceptionNote"`
	CreatedAt     string `json:"createdAt"`
	UpdatedAt     string `json:"updatedAt"`
}

type ReincarnationOrder struct {
	ID          int64  `json:"id"`
	SoulID      int64  `json:"soulId"`
	SoulName    string `json:"soulName,omitempty"`
	TargetRealm string `json:"targetRealm"`
	Priority    int    `json:"priority"`
	SoupStatus  string `json:"soupStatus"`
	QuotaType   string `json:"quotaType"`
	Status      string `json:"status"`
	ReviewNote  string `json:"reviewNote"`
	CreatedAt   string `json:"createdAt"`
	UpdatedAt   string `json:"updatedAt"`
}
```

- [ ] **Step 3: Implement DB open**

Create `server/internal/db/db.go`:

```go
package db

import (
	"database/sql"

	_ "modernc.org/sqlite"
)

func Open(path string) (*sql.DB, error) {
	database, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, err
	}
	database.SetMaxOpenConns(1)
	if _, err := database.Exec(`PRAGMA foreign_keys = ON;`); err != nil {
		_ = database.Close()
		return nil, err
	}
	return database, nil
}
```

- [ ] **Step 4: Implement migrations**

Create `server/internal/db/migrate.go` using the schema from the Database Schema section. The file must execute each `CREATE TABLE IF NOT EXISTS` statement in order:

```go
package db

import "database/sql"

func Migrate(database *sql.DB) error {
	statements := []string{
		`CREATE TABLE IF NOT EXISTS users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			username TEXT NOT NULL UNIQUE,
			nickname TEXT NOT NULL,
			role TEXT NOT NULL,
			avatar TEXT NOT NULL DEFAULT '',
			created_at TEXT NOT NULL
		);`,
		`CREATE TABLE IF NOT EXISTS souls (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL,
			birth_info TEXT NOT NULL,
			death_info TEXT NOT NULL DEFAULT '',
			status TEXT NOT NULL,
			merit_score INTEGER NOT NULL DEFAULT 0,
			karma_score INTEGER NOT NULL DEFAULT 0,
			memory_residue INTEGER NOT NULL DEFAULT 0,
			relationship_risk TEXT NOT NULL DEFAULT 'low',
			created_at TEXT NOT NULL,
			updated_at TEXT NOT NULL
		);`,
		`CREATE TABLE IF NOT EXISTS life_book_records (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			soul_id INTEGER NOT NULL,
			expected_lifespan INTEGER NOT NULL,
			actual_lifespan INTEGER NOT NULL,
			death_reason TEXT NOT NULL,
			fate_level TEXT NOT NULL,
			calamity_count INTEGER NOT NULL DEFAULT 0,
			locked INTEGER NOT NULL DEFAULT 0,
			risk_flag TEXT NOT NULL DEFAULT 'normal',
			updated_at TEXT NOT NULL,
			FOREIGN KEY (soul_id) REFERENCES souls(id)
		);`,
		`CREATE TABLE IF NOT EXISTS capture_tasks (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			soul_id INTEGER NOT NULL,
			assignee_id INTEGER,
			status TEXT NOT NULL,
			location TEXT NOT NULL,
			scheduled_time TEXT NOT NULL,
			actual_time TEXT NOT NULL DEFAULT '',
			exception_type TEXT NOT NULL DEFAULT '',
			exception_note TEXT NOT NULL DEFAULT '',
			created_at TEXT NOT NULL,
			updated_at TEXT NOT NULL,
			FOREIGN KEY (soul_id) REFERENCES souls(id),
			FOREIGN KEY (assignee_id) REFERENCES users(id)
		);`,
		`CREATE TABLE IF NOT EXISTS reincarnation_orders (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			soul_id INTEGER NOT NULL,
			target_realm TEXT NOT NULL,
			priority INTEGER NOT NULL DEFAULT 0,
			soup_status TEXT NOT NULL,
			quota_type TEXT NOT NULL DEFAULT '',
			status TEXT NOT NULL,
			review_note TEXT NOT NULL DEFAULT '',
			created_at TEXT NOT NULL,
			updated_at TEXT NOT NULL,
			FOREIGN KEY (soul_id) REFERENCES souls(id)
		);`,
		`CREATE TABLE IF NOT EXISTS approvals (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			type TEXT NOT NULL,
			target_id INTEGER NOT NULL,
			title TEXT NOT NULL,
			status TEXT NOT NULL,
			applicant_id INTEGER,
			approver_id INTEGER,
			reason TEXT NOT NULL DEFAULT '',
			result_note TEXT NOT NULL DEFAULT '',
			created_at TEXT NOT NULL,
			updated_at TEXT NOT NULL,
			FOREIGN KEY (applicant_id) REFERENCES users(id),
			FOREIGN KEY (approver_id) REFERENCES users(id)
		);`,
		`CREATE TABLE IF NOT EXISTS hell_floors (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			floor_no INTEGER NOT NULL UNIQUE,
			name TEXT NOT NULL,
			capacity INTEGER NOT NULL,
			occupancy INTEGER NOT NULL,
			equipment_status TEXT NOT NULL,
			load_level TEXT NOT NULL,
			updated_at TEXT NOT NULL
		);`,
		`CREATE TABLE IF NOT EXISTS hell_sentences (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			soul_id INTEGER NOT NULL,
			floor_id INTEGER NOT NULL,
			crime_type TEXT NOT NULL,
			sentence_days INTEGER NOT NULL,
			pain_level INTEGER NOT NULL,
			review_status TEXT NOT NULL,
			equipment_id TEXT NOT NULL DEFAULT '',
			created_at TEXT NOT NULL,
			updated_at TEXT NOT NULL,
			FOREIGN KEY (soul_id) REFERENCES souls(id),
			FOREIGN KEY (floor_id) REFERENCES hell_floors(id)
		);`,
		`CREATE TABLE IF NOT EXISTS wish_tickets (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			title TEXT NOT NULL,
			wish_type TEXT NOT NULL,
			requester_name TEXT NOT NULL,
			incense_amount INTEGER NOT NULL DEFAULT 0,
			merit_score INTEGER NOT NULL DEFAULT 0,
			status TEXT NOT NULL,
			assigned_deity TEXT NOT NULL DEFAULT '',
			result_note TEXT NOT NULL DEFAULT '',
			created_at TEXT NOT NULL,
			updated_at TEXT NOT NULL
		);`,
		`CREATE TABLE IF NOT EXISTS soup_inventory (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL,
			stock INTEGER NOT NULL,
			warning_stock INTEGER NOT NULL,
			formula_note TEXT NOT NULL DEFAULT '',
			updated_at TEXT NOT NULL
		);`,
		`CREATE TABLE IF NOT EXISTS soup_records (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			soul_id INTEGER NOT NULL,
			order_id INTEGER,
			inventory_id INTEGER NOT NULL,
			dose INTEGER NOT NULL,
			memory_after INTEGER NOT NULL,
			operator_id INTEGER,
			created_at TEXT NOT NULL,
			FOREIGN KEY (soul_id) REFERENCES souls(id),
			FOREIGN KEY (order_id) REFERENCES reincarnation_orders(id),
			FOREIGN KEY (inventory_id) REFERENCES soup_inventory(id),
			FOREIGN KEY (operator_id) REFERENCES users(id)
		);`,
		`CREATE TABLE IF NOT EXISTS alerts (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			level TEXT NOT NULL,
			title TEXT NOT NULL,
			content TEXT NOT NULL,
			handled INTEGER NOT NULL DEFAULT 0,
			created_at TEXT NOT NULL
		);`,
		`CREATE TABLE IF NOT EXISTS audit_logs (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			actor_id INTEGER,
			action TEXT NOT NULL,
			target_type TEXT NOT NULL,
			target_id INTEGER NOT NULL,
			before_json TEXT NOT NULL DEFAULT '',
			after_json TEXT NOT NULL DEFAULT '',
			note TEXT NOT NULL DEFAULT '',
			created_at TEXT NOT NULL,
			FOREIGN KEY (actor_id) REFERENCES users(id)
		);`,
	}

	for _, statement := range statements {
		if _, err := database.Exec(statement); err != nil {
			return err
		}
	}
	return nil
}
```

- [ ] **Step 5: Implement seed data**

Create `server/internal/db/seed.go`. It must insert only when `users` is empty:

```go
package db

import (
	"database/sql"
	"time"
)

func Seed(database *sql.DB) error {
	var count int
	if err := database.QueryRow(`SELECT COUNT(*) FROM users`).Scan(&count); err != nil {
		return err
	}
	if count > 0 {
		return nil
	}

	now := time.Now().Format(time.RFC3339)
	tx, err := database.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	inserts := []string{
		`INSERT INTO users (username, nickname, role, avatar, created_at) VALUES
		('yudi', '玉帝', 'jade_emperor', '', '` + now + `'),
		('yanwang', '阎王', 'yanwang', '', '` + now + `'),
		('panguan', '判官', 'judge', '', '` + now + `'),
		('heiwu', '黑无常', 'impermanence', '', '` + now + `'),
		('mengpo', '孟婆', 'mengpo', '', '` + now + `');`,
		`INSERT INTO souls (name, birth_info, death_info, status, merit_score, karma_score, memory_residue, relationship_risk, created_at, updated_at) VALUES
		('张三', '甲子年三月初八 午时', '阳寿将尽', 'pending_capture', 12, 37, 0, 'low', '` + now + `', '` + now + `'),
		('李四', '乙丑年六月十二 子时', '待投胎', 'reincarnation_waiting', 88, 4, 12, 'medium', '` + now + `', '` + now + `'),
		('王五', '丙寅年九月廿一 酉时', '审判完成', 'hell', 3, 91, 0, 'low', '` + now + `', '` + now + `'),
		('赵六', '丁卯年正月十五 辰时', '阳寿异常', 'alive', 66, 6, 0, 'high', '` + now + `', '` + now + `');`,
		`INSERT INTO life_book_records (soul_id, expected_lifespan, actual_lifespan, death_reason, fate_level, calamity_count, locked, risk_flag, updated_at) VALUES
		(1, 73, 73, '正常到期', '普通命', 3, 0, 'normal', '` + now + `'),
		(2, 81, 81, '功德圆满', '富贵命', 1, 0, 'normal', '` + now + `'),
		(3, 58, 58, '恶业反噬', '坎坷命', 7, 0, 'warning', '` + now + `'),
		(4, 69, 96, '疑似被加寿', '关系命', 2, 1, 'critical', '` + now + `');`,
		`INSERT INTO capture_tasks (soul_id, assignee_id, status, location, scheduled_time, created_at, updated_at) VALUES
		(1, 4, 'pending', '人间东市第三巷', '` + now + `', '` + now + `', '` + now + `');`,
		`INSERT INTO reincarnation_orders (soul_id, target_realm, priority, soup_status, quota_type, status, review_note, created_at, updated_at) VALUES
		(2, 'human', 88, 'pending', '普通胎', 'queued', '功德较高但记忆残留偏高', '` + now + `', '` + now + `');`,
		`INSERT INTO hell_floors (floor_no, name, capacity, occupancy, equipment_status, load_level, updated_at) VALUES
		(1, '拔舌地狱', 1000, 940, 'warning', 'critical', '` + now + `'),
		(2, '剪刀地狱', 900, 421, 'normal', 'normal', '` + now + `'),
		(3, '铁树地狱', 800, 688, 'normal', 'warning', '` + now + `');`,
		`INSERT INTO hell_sentences (soul_id, floor_id, crime_type, sentence_days, pain_level, review_status, equipment_id, created_at, updated_at) VALUES
		(3, 1, '口舌是非', 360, 7, 'none', 'tongue-001', '` + now + `', '` + now + `');`,
		`INSERT INTO wish_tickets (title, wish_type, requester_name, incense_amount, merit_score, status, assigned_deity, created_at, updated_at) VALUES
		('求暴富但不想上班', 'wealth', '刘某', 18, 1, 'pending', '', '` + now + `', '` + now + `'),
		('求考试刚好及格', 'exam', '陈某', 6, 12, 'routed', '文昌帝君', '` + now + `', '` + now + `');`,
		`INSERT INTO soup_inventory (name, stock, warning_stock, formula_note, updated_at) VALUES
		('标准孟婆汤', 320, 80, '适合普通投胎批次', '` + now + `'),
		('加浓孟婆汤', 48, 60, '适合记忆残留较高魂魄', '` + now + `');`,
		`INSERT INTO alerts (level, title, content, handled, created_at) VALUES
		('critical', '拔舌地狱容量过高', '当前容量达到 94%，建议开启跨层分流。', 0, '` + now + `'),
		('warning', '加浓孟婆汤库存不足', '库存低于预警线，下一批投胎可能受影响。', 0, '` + now + `'),
		('critical', '生死簿高危变更', '赵六阳寿记录已冻结，疑似关系户操作。', 0, '` + now + `');`,
	}

	for _, statement := range inserts {
		if _, err := tx.Exec(statement); err != nil {
			return err
		}
	}
	return tx.Commit()
}
```

- [ ] **Step 6: Run backend tests**

Run:

```powershell
cd server
go test ./...
```

Expected: PASS.

- [ ] **Step 7: Commit**

Run:

```powershell
git add server/internal/db server/internal/model server/internal_test server/go.mod server/go.sum
git commit -m "feat: add sqlite schema and seed data"
```

### Task 3: Dashboard, Alerts, Users, and Audit Backend

**Files:**

- Create: `server/internal/repository/users.go`
- Create: `server/internal/repository/dashboard.go`
- Create: `server/internal/repository/audit.go`
- Create: `server/internal/service/dashboard.go`
- Create: `server/internal/handler/users.go`
- Create: `server/internal/handler/dashboard.go`
- Create: `server/internal/handler/audit.go`
- Modify: `server/internal/handler/routes.go`
- Test: `server/internal_test/dashboard_test.go`

- [ ] **Step 1: Write dashboard tests**

Create `server/internal_test/dashboard_test.go`:

```go
package internal_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"sanjie/server/internal/app"
)

func TestDashboardReturnsSeedMetrics(t *testing.T) {
	database := NewTestDB(t)
	application := app.New(database)

	req := httptest.NewRequest(http.MethodGet, "/api/dashboard", nil)
	res := httptest.NewRecorder()
	application.Router.ServeHTTP(res, req)

	if res.Code != http.StatusOK {
		t.Fatalf("status = %d", res.Code)
	}

	var body struct {
		Code int `json:"code"`
		Data struct {
			PendingCaptureCount int `json:"pendingCaptureCount"`
			PendingApprovalCount int `json:"pendingApprovalCount"`
			UnhandledAlertCount int `json:"unhandledAlertCount"`
		} `json:"data"`
	}
	if err := json.NewDecoder(res.Body).Decode(&body); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if body.Code != 0 {
		t.Fatalf("code = %d", body.Code)
	}
	if body.Data.PendingCaptureCount != 1 {
		t.Fatalf("pending capture count = %d", body.Data.PendingCaptureCount)
	}
	if body.Data.UnhandledAlertCount != 3 {
		t.Fatalf("unhandled alert count = %d", body.Data.UnhandledAlertCount)
	}
}
```

- [ ] **Step 2: Run failing test**

Run:

```powershell
cd server
go test ./internal_test -run TestDashboardReturnsSeedMetrics -v
```

Expected: FAIL with 404 for `/api/dashboard`.

- [ ] **Step 3: Implement repositories and handlers**

Create `server/internal/repository/dashboard.go`:

```go
package repository

import "database/sql"

type DashboardRepository struct {
	DB *sql.DB
}

type DashboardMetrics struct {
	PendingCaptureCount int `json:"pendingCaptureCount"`
	PendingApprovalCount int `json:"pendingApprovalCount"`
	ReincarnationQueueCount int `json:"reincarnationQueueCount"`
	CriticalHellFloorCount int `json:"criticalHellFloorCount"`
	LowSoupInventoryCount int `json:"lowSoupInventoryCount"`
	UnhandledAlertCount int `json:"unhandledAlertCount"`
	HighRiskLifeBookCount int `json:"highRiskLifeBookCount"`
}

func (r DashboardRepository) Metrics() (DashboardMetrics, error) {
	var m DashboardMetrics
	queries := []struct {
		dest *int
		sql  string
	}{
		{&m.PendingCaptureCount, `SELECT COUNT(*) FROM capture_tasks WHERE status = 'pending'`},
		{&m.PendingApprovalCount, `SELECT COUNT(*) FROM approvals WHERE status = 'pending'`},
		{&m.ReincarnationQueueCount, `SELECT COUNT(*) FROM reincarnation_orders WHERE status IN ('queued', 'reviewing', 'pending_soup')`},
		{&m.CriticalHellFloorCount, `SELECT COUNT(*) FROM hell_floors WHERE load_level = 'critical'`},
		{&m.LowSoupInventoryCount, `SELECT COUNT(*) FROM soup_inventory WHERE stock <= warning_stock`},
		{&m.UnhandledAlertCount, `SELECT COUNT(*) FROM alerts WHERE handled = 0`},
		{&m.HighRiskLifeBookCount, `SELECT COUNT(*) FROM life_book_records WHERE risk_flag = 'critical'`},
	}
	for _, q := range queries {
		if err := r.DB.QueryRow(q.sql).Scan(q.dest); err != nil {
			return m, err
		}
	}
	return m, nil
}
```

Create `server/internal/handler/dashboard.go`:

```go
package handler

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"sanjie/server/internal/httpx"
	"sanjie/server/internal/repository"
)

type DashboardHandler struct {
	DB *sql.DB
}

func (h DashboardHandler) Dashboard(w http.ResponseWriter, r *http.Request) {
	metrics, err := repository.DashboardRepository{DB: h.DB}.Metrics()
	if err != nil {
		httpx.Fail(w, http.StatusInternalServerError, 50001, err.Error())
		return
	}
	httpx.OK(w, metrics)
}

func (h DashboardHandler) Alerts(w http.ResponseWriter, r *http.Request) {
	rows, err := h.DB.Query(`SELECT id, level, title, content, handled, created_at FROM alerts ORDER BY id DESC`)
	if err != nil {
		httpx.Fail(w, http.StatusInternalServerError, 50002, err.Error())
		return
	}
	defer rows.Close()
	items := make([]map[string]interface{}, 0)
	for rows.Next() {
		var id int64
		var level, title, content, createdAt string
		var handled int
		if err := rows.Scan(&id, &level, &title, &content, &handled, &createdAt); err != nil {
			httpx.Fail(w, http.StatusInternalServerError, 50003, err.Error())
			return
		}
		items = append(items, map[string]interface{}{
			"id": id, "level": level, "title": title, "content": content,
			"handled": handled == 1, "createdAt": createdAt,
		})
	}
	httpx.OK(w, items)
}

func (h DashboardHandler) HandleAlert(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		httpx.Fail(w, http.StatusBadRequest, 40001, "invalid alert id")
		return
	}
	if _, err := h.DB.Exec(`UPDATE alerts SET handled = 1 WHERE id = ?`, id); err != nil {
		httpx.Fail(w, http.StatusInternalServerError, 50004, err.Error())
		return
	}
	httpx.OK(w, map[string]interface{}{"id": id})
}
```

Create `server/internal/handler/users.go`:

```go
package handler

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"sanjie/server/internal/httpx"
)

type UsersHandler struct {
	DB *sql.DB
}

func (h UsersHandler) Me(w http.ResponseWriter, r *http.Request) {
	userID := httpx.CurrentUserID(r)
	var id int64
	var username, nickname, role, avatar, createdAt string
	err := h.DB.QueryRow(`SELECT id, username, nickname, role, avatar, created_at FROM users WHERE id = ?`, userID).
		Scan(&id, &username, &nickname, &role, &avatar, &createdAt)
	if err != nil {
		httpx.Fail(w, http.StatusNotFound, 40401, "user not found")
		return
	}
	httpx.OK(w, map[string]interface{}{
		"id": id, "username": username, "nickname": nickname,
		"role": role, "avatar": avatar, "createdAt": createdAt,
	})
}

func (h UsersHandler) Switch(w http.ResponseWriter, r *http.Request) {
	var input struct {
		UserID int64 `json:"userId"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil || input.UserID <= 0 {
		httpx.Fail(w, http.StatusBadRequest, 40002, "invalid user id")
		return
	}
	httpx.OK(w, map[string]interface{}{"userId": input.UserID})
}
```

Create `server/internal/handler/audit.go`:

```go
package handler

import (
	"database/sql"
	"net/http"

	"sanjie/server/internal/httpx"
)

type AuditHandler struct {
	DB *sql.DB
}

func (h AuditHandler) Logs(w http.ResponseWriter, r *http.Request) {
	rows, err := h.DB.Query(`SELECT id, actor_id, action, target_type, target_id, before_json, after_json, note, created_at FROM audit_logs ORDER BY id DESC LIMIT 100`)
	if err != nil {
		httpx.Fail(w, http.StatusInternalServerError, 50005, err.Error())
		return
	}
	defer rows.Close()
	items := make([]map[string]interface{}, 0)
	for rows.Next() {
		var id, targetID int64
		var actorID sql.NullInt64
		var action, targetType, beforeJSON, afterJSON, note, createdAt string
		if err := rows.Scan(&id, &actorID, &action, &targetType, &targetID, &beforeJSON, &afterJSON, &note, &createdAt); err != nil {
			httpx.Fail(w, http.StatusInternalServerError, 50006, err.Error())
			return
		}
		items = append(items, map[string]interface{}{
			"id": id, "actorId": actorID.Int64, "action": action, "targetType": targetType,
			"targetId": targetID, "beforeJson": beforeJSON, "afterJson": afterJSON,
			"note": note, "createdAt": createdAt,
		})
	}
	httpx.OK(w, items)
}
```

- [ ] **Step 4: Register routes**

Modify `server/internal/handler/routes.go`:

```go
func Routes(db *sql.DB) http.Handler {
	r := chi.NewRouter()
	dashboard := DashboardHandler{DB: db}
	users := UsersHandler{DB: db}
	audit := AuditHandler{DB: db}

	r.Get("/api/health", func(w http.ResponseWriter, r *http.Request) {
		httpx.OK(w, map[string]string{"status": "ok"})
	})
	r.Get("/api/users/me", users.Me)
	r.Post("/api/users/switch", users.Switch)
	r.Get("/api/dashboard", dashboard.Dashboard)
	r.Get("/api/alerts", dashboard.Alerts)
	r.Post("/api/alerts/{id}/handle", dashboard.HandleAlert)
	r.Get("/api/audit-logs", audit.Logs)

	return r
}
```

- [ ] **Step 5: Run tests**

Run:

```powershell
cd server
go test ./...
```

Expected: PASS.

- [ ] **Step 6: Commit**

Run:

```powershell
git add server/internal
git commit -m "feat: add dashboard alerts users and audit api"
```

### Task 4: Capture Task Backend

**Files:**

- Create: `server/internal/repository/capture.go`
- Create: `server/internal/service/capture.go`
- Create: `server/internal/handler/capture.go`
- Modify: `server/internal/handler/routes.go`
- Test: `server/internal_test/capture_test.go`

- [ ] **Step 1: Write capture transition tests**

Create `server/internal_test/capture_test.go`:

```go
package internal_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"sanjie/server/internal/app"
)

func TestCaptureTaskCanStartAndComplete(t *testing.T) {
	database := NewTestDB(t)
	application := app.New(database)

	startReq := httptest.NewRequest(http.MethodPost, "/api/capture-tasks/1/start", nil)
	startRes := httptest.NewRecorder()
	application.Router.ServeHTTP(startRes, startReq)
	if startRes.Code != http.StatusOK {
		t.Fatalf("start status = %d body = %s", startRes.Code, startRes.Body.String())
	}

	completeReq := httptest.NewRequest(http.MethodPost, "/api/capture-tasks/1/complete", bytes.NewBufferString(`{"note":"魂魄已准时签收"}`))
	completeRes := httptest.NewRecorder()
	application.Router.ServeHTTP(completeRes, completeReq)
	if completeRes.Code != http.StatusOK {
		t.Fatalf("complete status = %d body = %s", completeRes.Code, completeRes.Body.String())
	}

	detailReq := httptest.NewRequest(http.MethodGet, "/api/capture-tasks/1", nil)
	detailRes := httptest.NewRecorder()
	application.Router.ServeHTTP(detailRes, detailReq)

	var body struct {
		Data struct {
			Status string `json:"status"`
		} `json:"data"`
	}
	_ = json.NewDecoder(detailRes.Body).Decode(&body)
	if body.Data.Status != "completed" {
		t.Fatalf("status = %s", body.Data.Status)
	}
}
```

- [ ] **Step 2: Run failing test**

Run:

```powershell
cd server
go test ./internal_test -run TestCaptureTaskCanStartAndComplete -v
```

Expected: FAIL with 404 for capture endpoints.

- [ ] **Step 3: Implement capture handler**

Create `server/internal/handler/capture.go`:

```go
package handler

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"

	"sanjie/server/internal/httpx"
)

type CaptureHandler struct {
	DB *sql.DB
}

func (h CaptureHandler) List(w http.ResponseWriter, r *http.Request) {
	status := r.URL.Query().Get("status")
	query := `SELECT t.id, t.soul_id, s.name, COALESCE(t.assignee_id, 0), COALESCE(u.nickname, ''), t.status, t.location, t.scheduled_time, t.actual_time, t.exception_type, t.exception_note, t.created_at, t.updated_at
		FROM capture_tasks t
		JOIN souls s ON s.id = t.soul_id
		LEFT JOIN users u ON u.id = t.assignee_id`
	args := []interface{}{}
	if status != "" {
		query += ` WHERE t.status = ?`
		args = append(args, status)
	}
	query += ` ORDER BY t.id DESC`
	rows, err := h.DB.Query(query, args...)
	if err != nil {
		httpx.Fail(w, http.StatusInternalServerError, 50011, err.Error())
		return
	}
	defer rows.Close()
	items := make([]map[string]interface{}, 0)
	for rows.Next() {
		item, err := scanCapture(rows)
		if err != nil {
			httpx.Fail(w, http.StatusInternalServerError, 50012, err.Error())
			return
		}
		items = append(items, item)
	}
	httpx.OK(w, items)
}

func (h CaptureHandler) Detail(w http.ResponseWriter, r *http.Request) {
	id, ok := parseID(w, r)
	if !ok {
		return
	}
	row := h.DB.QueryRow(`SELECT t.id, t.soul_id, s.name, COALESCE(t.assignee_id, 0), COALESCE(u.nickname, ''), t.status, t.location, t.scheduled_time, t.actual_time, t.exception_type, t.exception_note, t.created_at, t.updated_at
		FROM capture_tasks t
		JOIN souls s ON s.id = t.soul_id
		LEFT JOIN users u ON u.id = t.assignee_id
		WHERE t.id = ?`, id)
	item, err := scanCapture(row)
	if err != nil {
		httpx.Fail(w, http.StatusNotFound, 40411, "capture task not found")
		return
	}
	httpx.OK(w, item)
}

func (h CaptureHandler) Start(w http.ResponseWriter, r *http.Request) {
	id, ok := parseID(w, r)
	if !ok {
		return
	}
	now := time.Now().Format(time.RFC3339)
	res, err := h.DB.Exec(`UPDATE capture_tasks SET status = 'running', updated_at = ? WHERE id = ? AND status = 'pending'`, now, id)
	if err != nil {
		httpx.Fail(w, http.StatusInternalServerError, 50013, err.Error())
		return
	}
	if changed, _ := res.RowsAffected(); changed == 0 {
		httpx.Fail(w, http.StatusBadRequest, 40011, "only pending capture tasks can start")
		return
	}
	h.audit(r, "capture.start", "capture_task", id, "开始勾魂")
	httpx.OK(w, map[string]interface{}{"id": id, "status": "running"})
}

func (h CaptureHandler) Complete(w http.ResponseWriter, r *http.Request) {
	id, ok := parseID(w, r)
	if !ok {
		return
	}
	var input struct {
		Note string `json:"note"`
	}
	_ = json.NewDecoder(r.Body).Decode(&input)
	now := time.Now().Format(time.RFC3339)
	res, err := h.DB.Exec(`UPDATE capture_tasks SET status = 'completed', actual_time = ?, updated_at = ? WHERE id = ? AND status = 'running'`, now, now, id)
	if err != nil {
		httpx.Fail(w, http.StatusInternalServerError, 50014, err.Error())
		return
	}
	if changed, _ := res.RowsAffected(); changed == 0 {
		httpx.Fail(w, http.StatusBadRequest, 40012, "only running capture tasks can complete")
		return
	}
	h.audit(r, "capture.complete", "capture_task", id, input.Note)
	httpx.OK(w, map[string]interface{}{"id": id, "status": "completed"})
}

func (h CaptureHandler) Exception(w http.ResponseWriter, r *http.Request) {
	id, ok := parseID(w, r)
	if !ok {
		return
	}
	var input struct {
		Type string `json:"type"`
		Note string `json:"note"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil || input.Type == "" {
		httpx.Fail(w, http.StatusBadRequest, 40013, "exception type is required")
		return
	}
	now := time.Now().Format(time.RFC3339)
	if _, err := h.DB.Exec(`UPDATE capture_tasks SET status = 'exception', exception_type = ?, exception_note = ?, updated_at = ? WHERE id = ?`, input.Type, input.Note, now, id); err != nil {
		httpx.Fail(w, http.StatusInternalServerError, 50015, err.Error())
		return
	}
	h.audit(r, "capture.exception", "capture_task", id, input.Note)
	httpx.OK(w, map[string]interface{}{"id": id, "status": "exception"})
}

func (h CaptureHandler) Cancel(w http.ResponseWriter, r *http.Request) {
	id, ok := parseID(w, r)
	if !ok {
		return
	}
	now := time.Now().Format(time.RFC3339)
	if _, err := h.DB.Exec(`UPDATE capture_tasks SET status = 'canceled', updated_at = ? WHERE id = ?`, now, id); err != nil {
		httpx.Fail(w, http.StatusInternalServerError, 50016, err.Error())
		return
	}
	h.audit(r, "capture.cancel", "capture_task", id, "撤销勾魂任务")
	httpx.OK(w, map[string]interface{}{"id": id, "status": "canceled"})
}

type scanner interface {
	Scan(dest ...interface{}) error
}

func scanCapture(row scanner) (map[string]interface{}, error) {
	var id, soulID, assigneeID int64
	var soulName, assigneeName, status, location, scheduledTime, actualTime, exceptionType, exceptionNote, createdAt, updatedAt string
	if err := row.Scan(&id, &soulID, &soulName, &assigneeID, &assigneeName, &status, &location, &scheduledTime, &actualTime, &exceptionType, &exceptionNote, &createdAt, &updatedAt); err != nil {
		return nil, err
	}
	return map[string]interface{}{
		"id": id, "soulId": soulID, "soulName": soulName, "assigneeId": assigneeID, "assigneeName": assigneeName,
		"status": status, "location": location, "scheduledTime": scheduledTime, "actualTime": actualTime,
		"exceptionType": exceptionType, "exceptionNote": exceptionNote, "createdAt": createdAt, "updatedAt": updatedAt,
	}, nil
}

func parseID(w http.ResponseWriter, r *http.Request) (int64, bool) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil || id <= 0 {
		httpx.Fail(w, http.StatusBadRequest, 40010, "invalid id")
		return 0, false
	}
	return id, true
}

func (h CaptureHandler) audit(r *http.Request, action string, targetType string, targetID int64, note string) {
	_, _ = h.DB.Exec(`INSERT INTO audit_logs (actor_id, action, target_type, target_id, note, created_at) VALUES (?, ?, ?, ?, ?, ?)`,
		httpx.CurrentUserID(r), action, targetType, targetID, note, time.Now().Format(time.RFC3339))
}
```

- [ ] **Step 4: Register capture routes**

Modify `server/internal/handler/routes.go`:

```go
capture := CaptureHandler{DB: db}
r.Get("/api/capture-tasks", capture.List)
r.Get("/api/capture-tasks/{id}", capture.Detail)
r.Post("/api/capture-tasks/{id}/start", capture.Start)
r.Post("/api/capture-tasks/{id}/complete", capture.Complete)
r.Post("/api/capture-tasks/{id}/exception", capture.Exception)
r.Post("/api/capture-tasks/{id}/cancel", capture.Cancel)
```

- [ ] **Step 5: Run tests**

Run:

```powershell
cd server
go test ./...
```

Expected: PASS.

- [ ] **Step 6: Commit**

Run:

```powershell
git add server/internal server/internal_test/capture_test.go
git commit -m "feat: add soul capture workflow api"
```

### Task 5: Reincarnation Backend

**Files:**

- Create: `server/internal/handler/reincarnation.go`
- Modify: `server/internal/handler/routes.go`
- Test: `server/internal_test/reincarnation_test.go`

- [ ] **Step 1: Write reincarnation tests**

Create `server/internal_test/reincarnation_test.go`:

```go
package internal_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"sanjie/server/internal/app"
)

func TestReincarnationApproveAssignSoupAndComplete(t *testing.T) {
	database := NewTestDB(t)
	application := app.New(database)

	steps := []struct {
		method string
		path   string
		body   string
	}{
		{http.MethodPost, "/api/reincarnations/1/approve", `{"note":"功德够，准了"}`},
		{http.MethodPost, "/api/reincarnations/1/assign-quota", `{"quotaType":"普通胎"}`},
		{http.MethodPost, "/api/reincarnations/1/soup", `{"inventoryId":1,"dose":1,"memoryAfter":0}`},
		{http.MethodPost, "/api/reincarnations/1/complete", `{"note":"已入轮回"}`},
	}
	for _, step := range steps {
		req := httptest.NewRequest(step.method, step.path, bytes.NewBufferString(step.body))
		res := httptest.NewRecorder()
		application.Router.ServeHTTP(res, req)
		if res.Code != http.StatusOK {
			t.Fatalf("%s status = %d body = %s", step.path, res.Code, res.Body.String())
		}
	}

	req := httptest.NewRequest(http.MethodGet, "/api/reincarnations/1", nil)
	res := httptest.NewRecorder()
	application.Router.ServeHTTP(res, req)
	var body struct {
		Data struct {
			Status string `json:"status"`
			SoupStatus string `json:"soupStatus"`
		} `json:"data"`
	}
	_ = json.NewDecoder(res.Body).Decode(&body)
	if body.Data.Status != "completed" || body.Data.SoupStatus != "issued" {
		t.Fatalf("unexpected order state: %+v", body.Data)
	}
}
```

- [ ] **Step 2: Run failing test**

Run:

```powershell
cd server
go test ./internal_test -run TestReincarnationApproveAssignSoupAndComplete -v
```

Expected: FAIL with 404.

- [ ] **Step 3: Implement reincarnation handler**

Create `server/internal/handler/reincarnation.go` with list, detail, approve, reject, review, assign quota, soup, and complete handlers. Use these rules:

```text
approve: queued or reviewing -> approved
reject: queued or reviewing or approved -> rejected
review: queued or approved -> reviewing
assign-quota: approved or pending_soup -> pending_soup and set quota_type
soup: pending_soup -> pending_soup, soup_status = issued, insert soup_records, decrement soup_inventory.stock
complete: pending_soup with soup_status issued -> completed
```

Use the same JSON response pattern and audit insert style from Task 4. Return `40021` when a transition is invalid.

- [ ] **Step 4: Register routes**

Modify `server/internal/handler/routes.go`:

```go
reincarnation := ReincarnationHandler{DB: db}
r.Get("/api/reincarnations", reincarnation.List)
r.Get("/api/reincarnations/{id}", reincarnation.Detail)
r.Post("/api/reincarnations/{id}/approve", reincarnation.Approve)
r.Post("/api/reincarnations/{id}/reject", reincarnation.Reject)
r.Post("/api/reincarnations/{id}/review", reincarnation.Review)
r.Post("/api/reincarnations/{id}/assign-quota", reincarnation.AssignQuota)
r.Post("/api/reincarnations/{id}/soup", reincarnation.IssueSoup)
r.Post("/api/reincarnations/{id}/complete", reincarnation.Complete)
```

- [ ] **Step 5: Run tests**

Run:

```powershell
cd server
go test ./...
```

Expected: PASS.

- [ ] **Step 6: Commit**

Run:

```powershell
git add server/internal/handler/reincarnation.go server/internal/handler/routes.go server/internal_test/reincarnation_test.go
git commit -m "feat: add reincarnation workflow api"
```

### Task 6: Life Book and Soul Search Backend

**Files:**

- Create: `server/internal/handler/souls.go`
- Create: `server/internal/handler/life_book.go`
- Modify: `server/internal/handler/routes.go`
- Test: `server/internal_test/life_book_test.go`

- [ ] **Step 1: Write life book tests**

Create `server/internal_test/life_book_test.go`:

```go
package internal_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"sanjie/server/internal/app"
)

func TestLifeBookCanFreezeAndCreateApproval(t *testing.T) {
	database := NewTestDB(t)
	application := app.New(database)

	freezeReq := httptest.NewRequest(http.MethodPost, "/api/life-book/1/freeze", nil)
	freezeRes := httptest.NewRecorder()
	application.Router.ServeHTTP(freezeRes, freezeReq)
	if freezeRes.Code != http.StatusOK {
		t.Fatalf("freeze status = %d body = %s", freezeRes.Code, freezeRes.Body.String())
	}

	changeReq := httptest.NewRequest(http.MethodPost, "/api/life-book/1/change-lifespan", bytes.NewBufferString(`{"newLifespan":88,"reason":"救人一命，申请延寿"}`))
	changeRes := httptest.NewRecorder()
	application.Router.ServeHTTP(changeRes, changeReq)
	if changeRes.Code != http.StatusOK {
		t.Fatalf("change status = %d body = %s", changeRes.Code, changeRes.Body.String())
	}
}
```

- [ ] **Step 2: Run failing test**

Run:

```powershell
cd server
go test ./internal_test -run TestLifeBookCanFreezeAndCreateApproval -v
```

Expected: FAIL with 404.

- [ ] **Step 3: Implement soul list/detail**

Create `server/internal/handler/souls.go`. Implement:

```text
GET /api/souls?keyword=&status=
GET /api/souls/{id}
```

Query must join no extra tables and return core soul fields from `souls`.

- [ ] **Step 4: Implement life book list/detail/freeze/change**

Create `server/internal/handler/life_book.go`. Implement:

```text
GET /api/life-book?keyword=&riskFlag=
GET /api/life-book/{id}
POST /api/life-book/{id}/freeze
POST /api/life-book/{id}/change-lifespan
GET /api/life-book/{id}/audit-logs
```

`change-lifespan` must insert an approval row:

```sql
INSERT INTO approvals (type, target_id, title, status, applicant_id, reason, created_at, updated_at)
VALUES ('lifespan_change', ?, '寿命变更审批', 'pending', ?, ?, ?, ?)
```

Do not immediately update `actual_lifespan`; it changes only when the approval is approved in Task 7.

- [ ] **Step 5: Register routes**

Modify `server/internal/handler/routes.go`:

```go
souls := SoulsHandler{DB: db}
lifeBook := LifeBookHandler{DB: db}
r.Get("/api/souls", souls.List)
r.Get("/api/souls/{id}", souls.Detail)
r.Get("/api/life-book", lifeBook.List)
r.Get("/api/life-book/{id}", lifeBook.Detail)
r.Post("/api/life-book/{id}/change-lifespan", lifeBook.ChangeLifespan)
r.Post("/api/life-book/{id}/freeze", lifeBook.Freeze)
r.Get("/api/life-book/{id}/audit-logs", lifeBook.AuditLogs)
```

- [ ] **Step 6: Run tests**

Run:

```powershell
cd server
go test ./...
```

Expected: PASS.

- [ ] **Step 7: Commit**

Run:

```powershell
git add server/internal/handler server/internal_test/life_book_test.go
git commit -m "feat: add soul and life book api"
```

### Task 7: Approval Backend

**Files:**

- Create: `server/internal/handler/approvals.go`
- Modify: `server/internal/handler/routes.go`
- Test: `server/internal_test/approvals_test.go`

- [ ] **Step 1: Write approval tests**

Create `server/internal_test/approvals_test.go`:

```go
package internal_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"sanjie/server/internal/app"
)

func TestApprovalCanApproveRejectAndTransfer(t *testing.T) {
	database := NewTestDB(t)
	application := app.New(database)

	createReq := httptest.NewRequest(http.MethodPost, "/api/life-book/1/change-lifespan", bytes.NewBufferString(`{"newLifespan":88,"reason":"功德补偿"}`))
	createRes := httptest.NewRecorder()
	application.Router.ServeHTTP(createRes, createReq)
	if createRes.Code != http.StatusOK {
		t.Fatalf("create approval status = %d body = %s", createRes.Code, createRes.Body.String())
	}

	approveReq := httptest.NewRequest(http.MethodPost, "/api/approvals/1/approve", bytes.NewBufferString(`{"note":"判官复核通过"}`))
	approveRes := httptest.NewRecorder()
	application.Router.ServeHTTP(approveRes, approveReq)
	if approveRes.Code != http.StatusOK {
		t.Fatalf("approve status = %d body = %s", approveRes.Code, approveRes.Body.String())
	}
}
```

- [ ] **Step 2: Run failing test**

Run:

```powershell
cd server
go test ./internal_test -run TestApprovalCanApproveRejectAndTransfer -v
```

Expected: FAIL with 404 for approval endpoint.

- [ ] **Step 3: Implement approval handler**

Create `server/internal/handler/approvals.go`:

```text
GET list with ?status=&type=
GET detail by id
POST approve
POST reject
POST transfer
```

Approval rules:

```text
approve: pending or transferred -> approved
reject: pending or transferred -> rejected
transfer: pending -> transferred
```

When approving `lifespan_change`, parse the `reason` text only for display and do not infer lifespan from text. To support actual lifespan update, modify Task 6 `change-lifespan` to store JSON in `reason`:

```json
{"newLifespan":88,"reason":"救人一命，申请延寿"}
```

Then in approval approve:

```sql
UPDATE life_book_records SET actual_lifespan = ?, locked = 0, risk_flag = 'normal', updated_at = ? WHERE id = ?
```

Insert audit log for all approve/reject/transfer operations.

- [ ] **Step 4: Register approval routes**

Modify `server/internal/handler/routes.go`:

```go
approvals := ApprovalsHandler{DB: db}
r.Get("/api/approvals", approvals.List)
r.Get("/api/approvals/{id}", approvals.Detail)
r.Post("/api/approvals/{id}/approve", approvals.Approve)
r.Post("/api/approvals/{id}/reject", approvals.Reject)
r.Post("/api/approvals/{id}/transfer", approvals.Transfer)
```

- [ ] **Step 5: Run tests**

Run:

```powershell
cd server
go test ./...
```

Expected: PASS.

- [ ] **Step 6: Commit**

Run:

```powershell
git add server/internal/handler server/internal_test/approvals_test.go
git commit -m "feat: add approval center api"
```

### Task 8: Hell Backend

**Files:**

- Create: `server/internal/handler/hell.go`
- Modify: `server/internal/handler/routes.go`
- Test: `server/internal_test/hell_test.go`

- [ ] **Step 1: Write hell tests**

Create `server/internal_test/hell_test.go`:

```go
package internal_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"sanjie/server/internal/app"
)

func TestHellSentenceReviewAndFloorDispatch(t *testing.T) {
	database := NewTestDB(t)
	application := app.New(database)

	reviewReq := httptest.NewRequest(http.MethodPost, "/api/hell/sentences/1/review", bytes.NewBufferString(`{"note":"疑似判官手滑"}`))
	reviewRes := httptest.NewRecorder()
	application.Router.ServeHTTP(reviewRes, reviewReq)
	if reviewRes.Code != http.StatusOK {
		t.Fatalf("review status = %d body = %s", reviewRes.Code, reviewRes.Body.String())
	}

	dispatchReq := httptest.NewRequest(http.MethodPost, "/api/hell/floors/1/dispatch", bytes.NewBufferString(`{"targetFloorId":2,"amount":40}`))
	dispatchRes := httptest.NewRecorder()
	application.Router.ServeHTTP(dispatchRes, dispatchReq)
	if dispatchRes.Code != http.StatusOK {
		t.Fatalf("dispatch status = %d body = %s", dispatchRes.Code, dispatchRes.Body.String())
	}
}
```

- [ ] **Step 2: Run failing test**

Run:

```powershell
cd server
go test ./internal_test -run TestHellSentenceReviewAndFloorDispatch -v
```

Expected: FAIL with 404.

- [ ] **Step 3: Implement hell handler**

Create `server/internal/handler/hell.go`:

```text
GET /api/hell/floors
GET /api/hell/sentences?reviewStatus=
GET /api/hell/sentences/{id}
POST /api/hell/sentences/{id}/review
POST /api/hell/floors/{id}/dispatch
```

`review` must set `hell_sentences.review_status = 'reviewing'` and create an approval:

```sql
INSERT INTO approvals (type, target_id, title, status, applicant_id, reason, created_at, updated_at)
VALUES ('hell_review', ?, '地狱刑期复核', 'pending', ?, ?, ?, ?)
```

`dispatch` must run in a transaction:

```text
source.occupancy -= amount
target.occupancy += amount
load_level recomputed as normal < 75%, warning 75%-89%, critical >= 90%
```

- [ ] **Step 4: Register routes**

Modify `server/internal/handler/routes.go`:

```go
hell := HellHandler{DB: db}
r.Get("/api/hell/floors", hell.Floors)
r.Get("/api/hell/sentences", hell.Sentences)
r.Get("/api/hell/sentences/{id}", hell.SentenceDetail)
r.Post("/api/hell/sentences/{id}/review", hell.ReviewSentence)
r.Post("/api/hell/floors/{id}/dispatch", hell.DispatchFloor)
```

- [ ] **Step 5: Run tests**

Run:

```powershell
cd server
go test ./...
```

Expected: PASS.

- [ ] **Step 6: Commit**

Run:

```powershell
git add server/internal/handler/hell.go server/internal/handler/routes.go server/internal_test/hell_test.go
git commit -m "feat: add hell floor and sentence api"
```

### Task 9: Wishes and Soup Backend

**Files:**

- Create: `server/internal/handler/wishes.go`
- Create: `server/internal/handler/soup.go`
- Modify: `server/internal/handler/routes.go`
- Test: `server/internal_test/wishes_test.go`

- [ ] **Step 1: Write wish and soup tests**

Create `server/internal_test/wishes_test.go`:

```go
package internal_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"sanjie/server/internal/app"
)

func TestWishRouteResolveRejectAndSoupAdjust(t *testing.T) {
	database := NewTestDB(t)
	application := app.New(database)

	routeReq := httptest.NewRequest(http.MethodPost, "/api/wishes/1/route", bytes.NewBufferString(`{"assignedDeity":"财神","note":"可转财运组评估"}`))
	routeRes := httptest.NewRecorder()
	application.Router.ServeHTTP(routeRes, routeReq)
	if routeRes.Code != http.StatusOK {
		t.Fatalf("route status = %d body = %s", routeRes.Code, routeRes.Body.String())
	}

	resolveReq := httptest.NewRequest(http.MethodPost, "/api/wishes/1/resolve", bytes.NewBufferString(`{"note":"建议先上班"}`))
	resolveRes := httptest.NewRecorder()
	application.Router.ServeHTTP(resolveRes, resolveReq)
	if resolveRes.Code != http.StatusOK {
		t.Fatalf("resolve status = %d body = %s", resolveRes.Code, resolveRes.Body.String())
	}

	adjustReq := httptest.NewRequest(http.MethodPost, "/api/soup/inventory/1/adjust", bytes.NewBufferString(`{"delta":10,"note":"补充标准汤"}`))
	adjustRes := httptest.NewRecorder()
	application.Router.ServeHTTP(adjustRes, adjustReq)
	if adjustRes.Code != http.StatusOK {
		t.Fatalf("adjust status = %d body = %s", adjustRes.Code, adjustRes.Body.String())
	}
}
```

- [ ] **Step 2: Run failing test**

Run:

```powershell
cd server
go test ./internal_test -run TestWishRouteResolveRejectAndSoupAdjust -v
```

Expected: FAIL with 404.

- [ ] **Step 3: Implement wishes handler**

Create `server/internal/handler/wishes.go`:

```text
GET /api/wishes?status=&wishType=
GET /api/wishes/{id}
POST /api/wishes/{id}/route
POST /api/wishes/{id}/resolve
POST /api/wishes/{id}/reject
```

State rules:

```text
route: pending -> routed
resolve: pending or routed -> resolved
reject: pending or routed -> rejected
```

- [ ] **Step 4: Implement soup handler**

Create `server/internal/handler/soup.go`:

```text
GET /api/soup/inventory
POST /api/soup/inventory/{id}/adjust
GET /api/soup/records
```

`adjust` must reject any result where stock would become negative.

- [ ] **Step 5: Register routes**

Modify `server/internal/handler/routes.go`:

```go
wishes := WishesHandler{DB: db}
soup := SoupHandler{DB: db}
r.Get("/api/wishes", wishes.List)
r.Get("/api/wishes/{id}", wishes.Detail)
r.Post("/api/wishes/{id}/route", wishes.Route)
r.Post("/api/wishes/{id}/resolve", wishes.Resolve)
r.Post("/api/wishes/{id}/reject", wishes.Reject)
r.Get("/api/soup/inventory", soup.Inventory)
r.Post("/api/soup/inventory/{id}/adjust", soup.AdjustInventory)
r.Get("/api/soup/records", soup.Records)
```

- [ ] **Step 6: Run tests**

Run:

```powershell
cd server
go test ./...
```

Expected: PASS.

- [ ] **Step 7: Commit**

Run:

```powershell
git add server/internal/handler/wishes.go server/internal/handler/soup.go server/internal/handler/routes.go server/internal_test/wishes_test.go
git commit -m "feat: add wishes and mengpo soup api"
```

### Task 10: API Documentation and Backend Runbook

**Files:**

- Create: `docs/api.md`
- Create: `docs/runbook.md`

- [ ] **Step 1: Write API docs**

Create `docs/api.md` with all endpoints from Backend API Contract. For each endpoint include:

```text
method
path
query params
request body
response data shape
status transition notes
```

Include example request bodies for all POST endpoints.

- [ ] **Step 2: Write local runbook**

Create `docs/runbook.md`:

```markdown
# Sanjie Local Runbook

## Backend

```powershell
cd server
go mod tidy
go test ./...
go run .
```

Backend URL: `http://localhost:8080`

Health check:

```powershell
Invoke-RestMethod http://localhost:8080/api/health
```

## Mobile App

Open `app/` in HBuilderX or run the project with the repo's chosen uni-app CLI workflow.

Set API base URL in `app/api/http.js`:

```js
export const API_BASE_URL = 'http://localhost:8080'
```
```

- [ ] **Step 3: Run backend tests**

Run:

```powershell
cd server
go test ./...
```

Expected: PASS.

- [ ] **Step 4: Commit**

Run:

```powershell
git add docs/api.md docs/runbook.md
git commit -m "docs: add sanjie api and runbook"
```

### Task 11: uni-app Scaffold and uview-plus Setup

**Files:**

- Create: `app/package.json`
- Create: `app/pages.json`
- Create: `app/manifest.json`
- Create: `app/main.js`
- Create: `app/App.vue`
- Create: `app/uni.scss`
- Create: `app/api/http.js`
- Create: `app/stores/user.js`

- [ ] **Step 1: Create package**

Create `app/package.json`:

```json
{
  "name": "sanjie-mobile",
  "version": "0.1.0",
  "private": true,
  "scripts": {
    "dev:h5": "uni -p h5",
    "build:h5": "uni build -p h5"
  },
  "dependencies": {
    "@dcloudio/uni-app": "latest",
    "@dcloudio/uni-h5": "latest",
    "pinia": "latest",
    "uview-plus": "latest",
    "vue": "latest"
  },
  "devDependencies": {
    "@dcloudio/vite-plugin-uni": "latest",
    "vite": "latest"
  }
}
```

- [ ] **Step 2: Create app entry**

Create `app/main.js`:

```js
import { createSSRApp } from 'vue'
import { createPinia } from 'pinia'
import uviewPlus from 'uview-plus'
import App from './App.vue'
import 'uview-plus/index.scss'

export function createApp() {
  const app = createSSRApp(App)
  app.use(createPinia())
  app.use(uviewPlus)
  return { app }
}
```

Create `app/App.vue`:

```vue
<script>
export default {
  onLaunch() {}
}
</script>

<style lang="scss">
@import 'uview-plus/index.scss';
</style>
```

- [ ] **Step 3: Create pages config**

Create `app/pages.json` with all planned pages and tabBar entries:

```json
{
  "pages": [
    { "path": "pages/index/index", "style": { "navigationBarTitleText": "三界工作台" } },
    { "path": "pages/capture/index", "style": { "navigationBarTitleText": "勾魂任务" } },
    { "path": "pages/reincarnation/index", "style": { "navigationBarTitleText": "投胎队列" } },
    { "path": "pages/approval/index", "style": { "navigationBarTitleText": "审批中心" } },
    { "path": "pages/mine/index", "style": { "navigationBarTitleText": "我的" } },
    { "path": "pages/login/index", "style": { "navigationBarTitleText": "切换身份" } },
    { "path": "pages/capture/detail", "style": { "navigationBarTitleText": "勾魂详情" } },
    { "path": "pages/reincarnation/detail", "style": { "navigationBarTitleText": "投胎详情" } },
    { "path": "pages/approval/detail", "style": { "navigationBarTitleText": "审批详情" } },
    { "path": "pages/life-book/search", "style": { "navigationBarTitleText": "生死簿查询" } },
    { "path": "pages/life-book/detail", "style": { "navigationBarTitleText": "生死簿详情" } },
    { "path": "pages/hell/index", "style": { "navigationBarTitleText": "十八层地狱" } },
    { "path": "pages/hell/sentence-detail", "style": { "navigationBarTitleText": "刑期详情" } },
    { "path": "pages/wish/index", "style": { "navigationBarTitleText": "愿望工单" } },
    { "path": "pages/wish/detail", "style": { "navigationBarTitleText": "愿望详情" } },
    { "path": "pages/soup/index", "style": { "navigationBarTitleText": "孟婆汤库存" } },
    { "path": "pages/audit/index", "style": { "navigationBarTitleText": "审计日志" } }
  ],
  "tabBar": {
    "color": "#69707a",
    "selectedColor": "#a93226",
    "backgroundColor": "#fffaf0",
    "borderStyle": "black",
    "list": [
      { "pagePath": "pages/index/index", "text": "工作台" },
      { "pagePath": "pages/capture/index", "text": "勾魂" },
      { "pagePath": "pages/reincarnation/index", "text": "投胎" },
      { "pagePath": "pages/approval/index", "text": "审批" },
      { "pagePath": "pages/mine/index", "text": "我的" }
    ]
  }
}
```

- [ ] **Step 4: Create global styles**

Create `app/uni.scss`:

```scss
$sanjie-bg: #f7f1e3;
$sanjie-panel: #fffaf0;
$sanjie-ink: #25211b;
$sanjie-muted: #756f63;
$sanjie-red: #a93226;
$sanjie-gold: #b8860b;
$sanjie-border: #e5d8bd;

page {
  background: $sanjie-bg;
  color: $sanjie-ink;
}

.page {
  min-height: 100vh;
  padding: 24rpx;
  box-sizing: border-box;
}

.panel {
  background: $sanjie-panel;
  border: 1rpx solid $sanjie-border;
  border-radius: 8rpx;
  padding: 24rpx;
}
```

- [ ] **Step 5: Create HTTP wrapper**

Create `app/api/http.js`:

```js
export const API_BASE_URL = 'http://localhost:8080'

export function request(path, options = {}) {
  const userId = uni.getStorageSync('sanjie_user_id') || 1
  return new Promise((resolve, reject) => {
    uni.request({
      url: `${API_BASE_URL}${path}`,
      method: options.method || 'GET',
      data: options.data || {},
      header: {
        'Content-Type': 'application/json',
        'X-User-ID': String(userId),
        ...(options.header || {})
      },
      success(res) {
        const body = res.data
        if (res.statusCode >= 200 && res.statusCode < 300 && body && body.code === 0) {
          resolve(body.data)
          return
        }
        reject(new Error((body && body.message) || '请求失败'))
      },
      fail(err) {
        reject(new Error(err.errMsg || '网络异常'))
      }
    })
  })
}
```

- [ ] **Step 6: Create user store**

Create `app/stores/user.js`:

```js
import { defineStore } from 'pinia'
import { request } from '../api/http'

export const useUserStore = defineStore('user', {
  state: () => ({
    current: null
  }),
  actions: {
    async fetchMe() {
      this.current = await request('/api/users/me')
      return this.current
    },
    switchUser(userId) {
      uni.setStorageSync('sanjie_user_id', userId)
      return this.fetchMe()
    }
  }
})
```

- [ ] **Step 7: Install and build**

Run:

```powershell
cd app
npm install
npm run build:h5
```

Expected: build succeeds. If CLI support differs in the local uni-app setup, document the working command in `docs/runbook.md`.

- [ ] **Step 8: Commit**

Run:

```powershell
git add app docs/runbook.md
git commit -m "feat: scaffold uniapp mobile client"
```

### Task 12: Shared Mobile Components and API Modules

**Files:**

- Create all files under `app/api/`
- Create all files under `app/utils/`
- Create all files under `app/components/`

- [ ] **Step 1: Create status utility**

Create `app/utils/status.js`:

```js
export const statusText = {
  pending: '待处理',
  running: '执行中',
  completed: '已完成',
  exception: '异常',
  canceled: '已撤销',
  queued: '排队中',
  approved: '已通过',
  rejected: '已驳回',
  reviewing: '复核中',
  pending_soup: '待喝汤',
  issued: '已发放',
  failed: '失败',
  critical: '严重',
  warning: '预警',
  info: '提示'
}

export const statusType = {
  pending: 'warning',
  running: 'primary',
  completed: 'success',
  exception: 'error',
  canceled: 'info',
  queued: 'warning',
  approved: 'success',
  rejected: 'error',
  reviewing: 'primary',
  pending_soup: 'warning',
  issued: 'success',
  failed: 'error',
  critical: 'error',
  warning: 'warning',
  info: 'info'
}

export function labelOf(value) {
  return statusText[value] || value || '-'
}

export function typeOf(value) {
  return statusType[value] || 'info'
}
```

- [ ] **Step 2: Create API modules**

Create each API file with thin wrappers:

```js
// app/api/dashboard.js
import { request } from './http'
export const getDashboard = () => request('/api/dashboard')
export const getAlerts = () => request('/api/alerts')
export const handleAlert = id => request(`/api/alerts/${id}/handle`, { method: 'POST' })
```

Create these API modules with one exported function per endpoint:

```js
// app/api/capture.js
import { request } from './http'
export const getCaptureTasks = params => request(`/api/capture-tasks${params?.status ? `?status=${params.status}` : ''}`)
export const getCaptureTask = id => request(`/api/capture-tasks/${id}`)
export const startCaptureTask = id => request(`/api/capture-tasks/${id}/start`, { method: 'POST' })
export const completeCaptureTask = (id, data) => request(`/api/capture-tasks/${id}/complete`, { method: 'POST', data })
export const reportCaptureException = (id, data) => request(`/api/capture-tasks/${id}/exception`, { method: 'POST', data })
export const cancelCaptureTask = id => request(`/api/capture-tasks/${id}/cancel`, { method: 'POST' })
```

```js
// app/api/reincarnation.js
import { request } from './http'
export const getReincarnations = params => request(`/api/reincarnations${params?.status ? `?status=${params.status}` : ''}`)
export const getReincarnation = id => request(`/api/reincarnations/${id}`)
export const approveReincarnation = (id, data) => request(`/api/reincarnations/${id}/approve`, { method: 'POST', data })
export const rejectReincarnation = (id, data) => request(`/api/reincarnations/${id}/reject`, { method: 'POST', data })
export const reviewReincarnation = (id, data) => request(`/api/reincarnations/${id}/review`, { method: 'POST', data })
export const assignReincarnationQuota = (id, data) => request(`/api/reincarnations/${id}/assign-quota`, { method: 'POST', data })
export const issueReincarnationSoup = (id, data) => request(`/api/reincarnations/${id}/soup`, { method: 'POST', data })
export const completeReincarnation = (id, data) => request(`/api/reincarnations/${id}/complete`, { method: 'POST', data })
```

```js
// app/api/approvals.js
import { request } from './http'
export const getApprovals = params => request(`/api/approvals${params?.status ? `?status=${params.status}` : ''}`)
export const getApproval = id => request(`/api/approvals/${id}`)
export const approveApproval = (id, data) => request(`/api/approvals/${id}/approve`, { method: 'POST', data })
export const rejectApproval = (id, data) => request(`/api/approvals/${id}/reject`, { method: 'POST', data })
export const transferApproval = (id, data) => request(`/api/approvals/${id}/transfer`, { method: 'POST', data })
```

```js
// app/api/souls.js
import { request } from './http'
export const getSouls = query => request(`/api/souls${query ? `?${query}` : ''}`)
export const getSoul = id => request(`/api/souls/${id}`)
```

```js
// app/api/lifeBook.js
import { request } from './http'
export const getLifeBookRecords = query => request(`/api/life-book${query ? `?${query}` : ''}`)
export const getLifeBookRecord = id => request(`/api/life-book/${id}`)
export const changeLifespan = (id, data) => request(`/api/life-book/${id}/change-lifespan`, { method: 'POST', data })
export const freezeLifeBookRecord = id => request(`/api/life-book/${id}/freeze`, { method: 'POST' })
export const getLifeBookAuditLogs = id => request(`/api/life-book/${id}/audit-logs`)
```

```js
// app/api/hell.js
import { request } from './http'
export const getHellFloors = () => request('/api/hell/floors')
export const getHellSentences = query => request(`/api/hell/sentences${query ? `?${query}` : ''}`)
export const getHellSentence = id => request(`/api/hell/sentences/${id}`)
export const reviewHellSentence = (id, data) => request(`/api/hell/sentences/${id}/review`, { method: 'POST', data })
export const dispatchHellFloor = (id, data) => request(`/api/hell/floors/${id}/dispatch`, { method: 'POST', data })
```

```js
// app/api/wishes.js
import { request } from './http'
export const getWishes = query => request(`/api/wishes${query ? `?${query}` : ''}`)
export const getWish = id => request(`/api/wishes/${id}`)
export const routeWish = (id, data) => request(`/api/wishes/${id}/route`, { method: 'POST', data })
export const resolveWish = (id, data) => request(`/api/wishes/${id}/resolve`, { method: 'POST', data })
export const rejectWish = (id, data) => request(`/api/wishes/${id}/reject`, { method: 'POST', data })
```

```js
// app/api/soup.js
import { request } from './http'
export const getSoupInventory = () => request('/api/soup/inventory')
export const adjustSoupInventory = (id, data) => request(`/api/soup/inventory/${id}/adjust`, { method: 'POST', data })
export const getSoupRecords = () => request('/api/soup/records')
```

```js
// app/api/audit.js
import { request } from './http'
export const getAuditLogs = () => request('/api/audit-logs')
```

- [ ] **Step 3: Create StatusTag component**

Create `app/components/StatusTag.vue`:

```vue
<template>
  <u-tag :text="labelOf(value)" :type="typeOf(value)" size="mini" />
</template>

<script setup>
import { labelOf, typeOf } from '../utils/status'

defineProps({
  value: {
    type: String,
    default: ''
  }
})
</script>
```

- [ ] **Step 4: Create MetricCard component**

Create `app/components/MetricCard.vue`:

```vue
<template>
  <view class="metric-card">
    <text class="metric-card__label">{{ label }}</text>
    <text class="metric-card__value">{{ value }}</text>
  </view>
</template>

<script setup>
defineProps({
  label: { type: String, required: true },
  value: { type: [String, Number], required: true }
})
</script>

<style scoped lang="scss">
.metric-card {
  background: #fffaf0;
  border: 1rpx solid #e5d8bd;
  border-radius: 8rpx;
  padding: 22rpx;
}

.metric-card__label {
  display: block;
  color: #756f63;
  font-size: 24rpx;
}

.metric-card__value {
  display: block;
  margin-top: 10rpx;
  color: #25211b;
  font-size: 44rpx;
  font-weight: 700;
}
```

- [ ] **Step 5: Create EmptyState component**

Create `app/components/EmptyState.vue`:

```vue
<template>
  <view class="empty-state">
    <u-empty :text="text" mode="data" />
  </view>
</template>

<script setup>
defineProps({
  text: { type: String, default: '暂无数据' }
})
</script>

<style scoped>
.empty-state {
  padding: 80rpx 0;
}
</style>
```

- [ ] **Step 6: Create FixedActionBar component**

Create `app/components/FixedActionBar.vue`:

```vue
<template>
  <view class="fixed-action-bar">
    <slot />
  </view>
</template>

<style scoped>
.fixed-action-bar {
  position: fixed;
  left: 0;
  right: 0;
  bottom: 0;
  display: flex;
  gap: 16rpx;
  padding: 18rpx 24rpx calc(18rpx + env(safe-area-inset-bottom));
  background: #fffaf0;
  border-top: 1rpx solid #e5d8bd;
  box-sizing: border-box;
  z-index: 10;
}
</style>
```

- [ ] **Step 7: Build mobile app**

Run:

```powershell
cd app
npm run build:h5
```

Expected: PASS.

- [ ] **Step 8: Commit**

Run:

```powershell
git add app/api app/utils app/components
git commit -m "feat: add mobile api wrappers and shared components"
```

### Task 13: Workbench and Mine Pages

**Files:**

- Create: `app/pages/index/index.vue`
- Create: `app/pages/mine/index.vue`
- Create: `app/pages/login/index.vue`

- [ ] **Step 1: Implement workbench**

Create `app/pages/index/index.vue`. It must:

```text
fetch dashboard metrics on show
fetch alerts on show
show metric grid
show alert list
provide navigation to life book, hell, wishes, soup, audit
```

Use `MetricCard`, `u-grid`, `u-cell-group`, and `u-cell`.

- [ ] **Step 2: Implement mine**

Create `app/pages/mine/index.vue`. It must:

```text
show current user nickname and role
navigate to login switch page
navigate to audit logs
show API base URL
```

- [ ] **Step 3: Implement login switch**

Create `app/pages/login/index.vue`. It must list seed user IDs:

```text
1 玉帝
2 阎王
3 判官
4 黑无常
5 孟婆
```

On click, call `userStore.switchUser(id)` and navigate back.

- [ ] **Step 4: Build mobile app**

Run:

```powershell
cd app
npm run build:h5
```

Expected: PASS.

- [ ] **Step 5: Commit**

Run:

```powershell
git add app/pages/index app/pages/mine app/pages/login
git commit -m "feat: add workbench and user pages"
```

### Task 14: Capture Mobile Pages

**Files:**

- Create: `app/pages/capture/index.vue`
- Create: `app/pages/capture/detail.vue`

- [ ] **Step 1: Implement capture list**

Create `app/pages/capture/index.vue`. It must:

```text
show tabs for all, pending, running, exception, completed
fetch /api/capture-tasks with status filter
show soulName, location, scheduledTime, status
navigate to detail on click
show EmptyState when empty
```

- [ ] **Step 2: Implement capture detail**

Create `app/pages/capture/detail.vue`. It must:

```text
load id from route query
show soulName, location, assigneeName, scheduledTime, actualTime, exception fields
show Start button only when status is pending
show Complete and Exception buttons only when status is running
show Cancel button when status is pending or running
use u-modal or u-action-sheet for exception type selection
refresh detail after each operation
```

Exception options:

```text
emergency_rescue: 目标抢救中
wrong_soul: 勾错魂
escaped: 魂魄逃逸
not_found: 查无此人
ascended: 疑似飞升
```

- [ ] **Step 3: Build mobile app**

Run:

```powershell
cd app
npm run build:h5
```

Expected: PASS.

- [ ] **Step 4: Commit**

Run:

```powershell
git add app/pages/capture
git commit -m "feat: add capture task mobile pages"
```

### Task 15: Reincarnation Mobile Pages

**Files:**

- Create: `app/pages/reincarnation/index.vue`
- Create: `app/pages/reincarnation/detail.vue`

- [ ] **Step 1: Implement reincarnation list**

Create `app/pages/reincarnation/index.vue`. It must:

```text
show tabs for all, queued, reviewing, pending_soup, completed, rejected
fetch /api/reincarnations with status filter
show soulName, targetRealm, priority, soupStatus, quotaType, status
navigate to detail on click
```

- [ ] **Step 2: Implement reincarnation detail**

Create `app/pages/reincarnation/detail.vue`. It must:

```text
show all order fields
approve/reject/review buttons for queued/reviewing records
assign quota button for approved/pending_soup records
issue soup button for pending_soup records when soupStatus is pending
complete button when soupStatus is issued
use u-picker for quota type: 富贵胎, 普通胎, 清贫胎, 天才胎, 困难模式胎
use u-picker for soup inventory selection from /api/soup/inventory
```

- [ ] **Step 3: Build mobile app**

Run:

```powershell
cd app
npm run build:h5
```

Expected: PASS.

- [ ] **Step 4: Commit**

Run:

```powershell
git add app/pages/reincarnation
git commit -m "feat: add reincarnation mobile pages"
```

### Task 16: Approval Mobile Pages

**Files:**

- Create: `app/pages/approval/index.vue`
- Create: `app/pages/approval/detail.vue`

- [ ] **Step 1: Implement approval list**

Create `app/pages/approval/index.vue`. It must:

```text
show tabs for pending, transferred, approved, rejected
fetch /api/approvals with status filter
show type, title, reason summary, status
navigate to detail
```

- [ ] **Step 2: Implement approval detail**

Create `app/pages/approval/detail.vue`. It must:

```text
show type, title, targetId, applicantId, approverId, reason, resultNote
show approve/reject/transfer actions when status is pending or transferred
require note input for approve/reject/transfer
refresh after action
```

- [ ] **Step 3: Build mobile app**

Run:

```powershell
cd app
npm run build:h5
```

Expected: PASS.

- [ ] **Step 4: Commit**

Run:

```powershell
git add app/pages/approval
git commit -m "feat: add approval mobile pages"
```

### Task 17: Life Book Mobile Pages

**Files:**

- Create: `app/pages/life-book/search.vue`
- Create: `app/pages/life-book/detail.vue`

- [ ] **Step 1: Implement life book search**

Create `app/pages/life-book/search.vue`. It must:

```text
show u-search for keyword
show risk filter: all, normal, warning, critical
fetch /api/life-book
show soulName, actualLifespan, fateLevel, riskFlag, locked
navigate to detail
```

- [ ] **Step 2: Implement life book detail**

Create `app/pages/life-book/detail.vue`. It must:

```text
show life book fields
show Freeze button when not locked
show Change Lifespan button
change lifespan opens form with newLifespan and reason
submit creates approval
show audit logs for this record
```

- [ ] **Step 3: Build mobile app**

Run:

```powershell
cd app
npm run build:h5
```

Expected: PASS.

- [ ] **Step 4: Commit**

Run:

```powershell
git add app/pages/life-book
git commit -m "feat: add life book mobile pages"
```

### Task 18: Hell Mobile Pages

**Files:**

- Create: `app/pages/hell/index.vue`
- Create: `app/pages/hell/sentence-detail.vue`

- [ ] **Step 1: Implement hell floor board**

Create `app/pages/hell/index.vue`. It must:

```text
show floor cards from /api/hell/floors
show capacity, occupancy, load percentage, equipmentStatus, loadLevel
show sentence list from /api/hell/sentences
tap sentence to detail
provide dispatch action for critical floors
```

- [ ] **Step 2: Implement sentence detail**

Create `app/pages/hell/sentence-detail.vue`. It must:

```text
show soulName, floorName, crimeType, sentenceDays, painLevel, reviewStatus
show Review button when reviewStatus is none
submit review with note
```

- [ ] **Step 3: Build mobile app**

Run:

```powershell
cd app
npm run build:h5
```

Expected: PASS.

- [ ] **Step 4: Commit**

Run:

```powershell
git add app/pages/hell
git commit -m "feat: add hell mobile pages"
```

### Task 19: Wishes, Soup, and Audit Mobile Pages

**Files:**

- Create: `app/pages/wish/index.vue`
- Create: `app/pages/wish/detail.vue`
- Create: `app/pages/soup/index.vue`
- Create: `app/pages/audit/index.vue`

- [ ] **Step 1: Implement wish list and detail**

Create wish pages. They must:

```text
list by status filter
show requesterName, wishType, incenseAmount, meritScore, assignedDeity, status
route with assignedDeity and note
resolve with note
reject with note
```

- [ ] **Step 2: Implement soup page**

Create `app/pages/soup/index.vue`. It must:

```text
show inventory
show stock warning when stock <= warningStock
adjust inventory with positive or negative delta
show recent soup records
```

- [ ] **Step 3: Implement audit page**

Create `app/pages/audit/index.vue`. It must:

```text
fetch /api/audit-logs
show actorId, action, targetType, targetId, note, createdAt
```

- [ ] **Step 4: Build mobile app**

Run:

```powershell
cd app
npm run build:h5
```

Expected: PASS.

- [ ] **Step 5: Commit**

Run:

```powershell
git add app/pages/wish app/pages/soup app/pages/audit
git commit -m "feat: add wishes soup and audit mobile pages"
```

### Task 20: Full Integration and Manual Acceptance

**Files:**

- Modify: `docs/runbook.md`
- Create: `docs/acceptance-checklist.md`

- [ ] **Step 1: Create acceptance checklist**

Create `docs/acceptance-checklist.md`:

```markdown
# Sanjie Acceptance Checklist

- [ ] Backend starts on `http://localhost:8080`.
- [ ] `GET /api/health` returns code `0`.
- [ ] Mobile H5 starts and loads workbench metrics.
- [ ] User can switch between 玉帝, 阎王, 判官, 黑无常, 孟婆.
- [ ] Workbench shows alerts and metrics.
- [ ] Alert can be marked handled.
- [ ] Capture task can start, complete, cancel, and report exception.
- [ ] Reincarnation order can approve, reject, review, assign quota, issue soup, and complete.
- [ ] Life book can search, freeze, create lifespan approval, and show audit logs.
- [ ] Approval can approve, reject, and transfer.
- [ ] Hell page shows floors and sentences.
- [ ] Hell sentence review creates an approval.
- [ ] Hell floor dispatch changes source and target occupancy.
- [ ] Wish ticket can route, resolve, and reject.
- [ ] Soup inventory can adjust stock and rejects negative stock.
- [ ] Audit log records sensitive operations.
```

- [ ] **Step 2: Run backend tests**

Run:

```powershell
cd server
go test ./...
```

Expected: PASS.

- [ ] **Step 3: Build mobile app**

Run:

```powershell
cd app
npm run build:h5
```

Expected: PASS.

- [ ] **Step 4: Start backend**

Run:

```powershell
cd server
go run .
```

Expected: server logs `sanjie server listening on http://localhost:8080`.

- [ ] **Step 5: Start mobile H5**

Run in a second terminal:

```powershell
cd app
npm run dev:h5
```

Expected: H5 URL opens, usually `http://localhost:5173`.

- [ ] **Step 6: Execute checklist**

Open the mobile H5 URL and complete every item in `docs/acceptance-checklist.md`.

- [ ] **Step 7: Commit final docs**

Run:

```powershell
git add docs/acceptance-checklist.md docs/runbook.md
git commit -m "docs: add acceptance checklist"
```

---

## Self-Review

Spec coverage:

- Workbench, alerts, users, and role switching are covered by Tasks 3, 11, and 13.
- Soul capture is covered by Tasks 4 and 14.
- Reincarnation and Mengpo soup are covered by Tasks 5, 9, 15, and 19.
- Approval center is covered by Tasks 7 and 16.
- Life book and soul search are covered by Tasks 6 and 17.
- Eighteen hell floors and sentence review are covered by Tasks 8 and 18.
- Wish tickets are covered by Tasks 9 and 19.
- Audit logs are covered by Tasks 3, 4, 6, 7, 8, 9, 17, and 19.
- Documentation and acceptance are covered by Tasks 10 and 20.

Placeholder scan:

- No placeholder markers or intentional blanks remain.
- Tasks 5, 6, 7, 8, and 9 intentionally specify behavior contracts rather than full source listings because handlers are repetitive and should follow the complete pattern established in Task 4.

Type consistency:

- Backend constants and frontend status labels use the same status strings.
- Endpoint paths in backend, API contract, and mobile pages are aligned.
- SQLite table names match handler route names and seed data.
