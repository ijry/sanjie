package handler

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"

	"github.com/ijry/sanjie/internal/httpx"
)

func parseID(w http.ResponseWriter, r *http.Request) (int64, bool) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil || id <= 0 {
		httpx.Fail(w, http.StatusBadRequest, 40010, "invalid id")
		return 0, false
	}
	return id, true
}

func readJSON(w http.ResponseWriter, r *http.Request, dst interface{}) bool {
	if r.Body == nil {
		return true
	}
	if err := json.NewDecoder(r.Body).Decode(dst); err != nil {
		httpx.Fail(w, http.StatusBadRequest, 40011, "invalid json body")
		return false
	}
	return true
}

func nowRFC3339() string {
	return time.Now().Format(time.RFC3339)
}

func audit(db *sql.DB, r *http.Request, action string, targetType string, targetID int64, note string) {
	_, _ = db.Exec(
		`INSERT INTO audit_logs (actor_id, action, target_type, target_id, note, created_at) VALUES (?, ?, ?, ?, ?, ?)`,
		httpx.CurrentUserID(r), action, targetType, targetID, note, nowRFC3339(),
	)
}

func boolFromInt(v int) bool {
	return v == 1
}

func loadLevel(occupancy int, capacity int) string {
	if capacity <= 0 {
		return "normal"
	}
	percent := occupancy * 100 / capacity
	if percent >= 90 {
		return "critical"
	}
	if percent >= 75 {
		return "warning"
	}
	return "normal"
}
