package handler

import (
	"database/sql"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/ijry/sanjie/internal/httpx"
)

func Routes(db *sql.DB) http.Handler {
	r := chi.NewRouter()
	dashboard := DashboardHandler{DB: db}
	users := UsersHandler{DB: db}
	audit := AuditHandler{DB: db}
	capture := CaptureHandler{DB: db}
	souls := SoulsHandler{DB: db}
	lifeBook := LifeBookHandler{DB: db}
	approvals := ApprovalsHandler{DB: db}
	reincarnation := ReincarnationHandler{DB: db}
	hell := HellHandler{DB: db}
	wishes := WishesHandler{DB: db}
	soup := SoupHandler{DB: db}

	r.Get("/api/health", func(w http.ResponseWriter, r *http.Request) {
		httpx.OK(w, map[string]string{"status": "ok"})
	})
	r.Get("/api/users/me", users.Me)
	r.Post("/api/users/switch", users.Switch)
	r.Get("/api/dashboard", dashboard.Dashboard)
	r.Get("/api/alerts", dashboard.Alerts)
	r.Post("/api/alerts/{id}/handle", dashboard.HandleAlert)
	r.Get("/api/audit-logs", audit.Logs)
	r.Get("/api/capture-tasks", capture.List)
	r.Get("/api/capture-tasks/{id}", capture.Detail)
	r.Post("/api/capture-tasks/{id}/start", capture.Start)
	r.Post("/api/capture-tasks/{id}/complete", capture.Complete)
	r.Post("/api/capture-tasks/{id}/exception", capture.Exception)
	r.Post("/api/capture-tasks/{id}/cancel", capture.Cancel)
	r.Get("/api/souls", souls.List)
	r.Get("/api/souls/{id}", souls.Detail)
	r.Get("/api/life-book", lifeBook.List)
	r.Get("/api/life-book/{id}", lifeBook.Detail)
	r.Post("/api/life-book/{id}/change-lifespan", lifeBook.ChangeLifespan)
	r.Post("/api/life-book/{id}/freeze", lifeBook.Freeze)
	r.Get("/api/life-book/{id}/audit-logs", lifeBook.AuditLogs)
	r.Get("/api/approvals", approvals.List)
	r.Get("/api/approvals/{id}", approvals.Detail)
	r.Post("/api/approvals/{id}/approve", approvals.Approve)
	r.Post("/api/approvals/{id}/reject", approvals.Reject)
	r.Post("/api/approvals/{id}/transfer", approvals.Transfer)
	r.Get("/api/reincarnations", reincarnation.List)
	r.Get("/api/reincarnations/{id}", reincarnation.Detail)
	r.Post("/api/reincarnations/{id}/approve", reincarnation.Approve)
	r.Post("/api/reincarnations/{id}/reject", reincarnation.Reject)
	r.Post("/api/reincarnations/{id}/review", reincarnation.Review)
	r.Post("/api/reincarnations/{id}/assign-quota", reincarnation.AssignQuota)
	r.Post("/api/reincarnations/{id}/soup", reincarnation.IssueSoup)
	r.Post("/api/reincarnations/{id}/complete", reincarnation.Complete)
	r.Get("/api/hell/floors", hell.Floors)
	r.Get("/api/hell/sentences", hell.Sentences)
	r.Get("/api/hell/sentences/{id}", hell.SentenceDetail)
	r.Post("/api/hell/sentences/{id}/review", hell.ReviewSentence)
	r.Post("/api/hell/floors/{id}/dispatch", hell.DispatchFloor)
	r.Get("/api/wishes", wishes.List)
	r.Get("/api/wishes/{id}", wishes.Detail)
	r.Post("/api/wishes/{id}/route", wishes.Route)
	r.Post("/api/wishes/{id}/resolve", wishes.Resolve)
	r.Post("/api/wishes/{id}/reject", wishes.Reject)
	r.Get("/api/soup/inventory", soup.Inventory)
	r.Post("/api/soup/inventory/{id}/adjust", soup.AdjustInventory)
	r.Get("/api/soup/records", soup.Records)

	return r
}
