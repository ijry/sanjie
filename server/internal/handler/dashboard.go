package handler

import (
	"database/sql"
	"net/http"

	"github.com/ijry/sanjie/internal/httpx"
)

type DashboardHandler struct {
	DB *sql.DB
}

func (h DashboardHandler) Dashboard(w http.ResponseWriter, r *http.Request) {
	metrics := map[string]int{}
	queries := map[string]string{
		"pendingCaptureCount":       `SELECT COUNT(*) FROM capture_tasks WHERE status = 'pending'`,
		"pendingApprovalCount":      `SELECT COUNT(*) FROM approvals WHERE status = 'pending'`,
		"reincarnationQueueCount":   `SELECT COUNT(*) FROM reincarnation_orders WHERE status IN ('queued', 'reviewing', 'pending_soup')`,
		"criticalHellFloorCount":    `SELECT COUNT(*) FROM hell_floors WHERE load_level = 'critical'`,
		"lowSoupInventoryCount":     `SELECT COUNT(*) FROM soup_inventory WHERE stock <= warning_stock`,
		"unhandledAlertCount":       `SELECT COUNT(*) FROM alerts WHERE handled = 0`,
		"highRiskLifeBookCount":     `SELECT COUNT(*) FROM life_book_records WHERE risk_flag = 'critical'`,
		"pendingWishTicketCount":    `SELECT COUNT(*) FROM wish_tickets WHERE status = 'pending'`,
		"activeHellSentenceCount":   `SELECT COUNT(*) FROM hell_sentences WHERE review_status IN ('none', 'reviewing')`,
		"memoryResidueWarningCount": `SELECT COUNT(*) FROM souls WHERE memory_residue >= 10`,
	}
	for key, query := range queries {
		var count int
		if err := h.DB.QueryRow(query).Scan(&count); err != nil {
			httpx.Fail(w, http.StatusInternalServerError, 50001, err.Error())
			return
		}
		metrics[key] = count
	}
	httpx.OK(w, metrics)
}

func (h DashboardHandler) Alerts(w http.ResponseWriter, r *http.Request) {
	rows, err := h.DB.Query(`SELECT id, level, title, content, handled, created_at FROM alerts ORDER BY handled ASC, id DESC`)
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
			"handled": boolFromInt(handled), "createdAt": createdAt,
		})
	}
	httpx.OK(w, items)
}

func (h DashboardHandler) HandleAlert(w http.ResponseWriter, r *http.Request) {
	id, ok := parseID(w, r)
	if !ok {
		return
	}
	if _, err := h.DB.Exec(`UPDATE alerts SET handled = 1 WHERE id = ?`, id); err != nil {
		httpx.Fail(w, http.StatusInternalServerError, 50004, err.Error())
		return
	}
	audit(h.DB, r, "alert.handle", "alert", id, "告警已处理")
	httpx.OK(w, map[string]interface{}{"id": id, "handled": true})
}
