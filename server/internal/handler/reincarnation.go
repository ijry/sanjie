package handler

import (
	"database/sql"
	"net/http"

	"github.com/ijry/sanjie/internal/httpx"
)

type ReincarnationHandler struct {
	DB *sql.DB
}

func (h ReincarnationHandler) List(w http.ResponseWriter, r *http.Request) {
	status := r.URL.Query().Get("status")
	query := `SELECT o.id, o.soul_id, s.name, o.target_realm, o.priority, o.soup_status, o.quota_type, o.status, o.review_note, o.created_at, o.updated_at
		FROM reincarnation_orders o
		JOIN souls s ON s.id = o.soul_id`
	args := []interface{}{}
	if status != "" {
		query += ` WHERE o.status = ?`
		args = append(args, status)
	}
	query += ` ORDER BY o.priority DESC, o.id DESC`
	rows, err := h.DB.Query(query, args...)
	if err != nil {
		httpx.Fail(w, http.StatusInternalServerError, 53001, err.Error())
		return
	}
	defer rows.Close()
	items := make([]map[string]interface{}, 0)
	for rows.Next() {
		item, err := scanReincarnation(rows)
		if err != nil {
			httpx.Fail(w, http.StatusInternalServerError, 53002, err.Error())
			return
		}
		items = append(items, item)
	}
	httpx.OK(w, items)
}

func (h ReincarnationHandler) Detail(w http.ResponseWriter, r *http.Request) {
	id, ok := parseID(w, r)
	if !ok {
		return
	}
	item, err := scanReincarnation(h.DB.QueryRow(`SELECT o.id, o.soul_id, s.name, o.target_realm, o.priority, o.soup_status, o.quota_type, o.status, o.review_note, o.created_at, o.updated_at
		FROM reincarnation_orders o
		JOIN souls s ON s.id = o.soul_id
		WHERE o.id = ?`, id))
	if err != nil {
		httpx.Fail(w, http.StatusNotFound, 40441, "reincarnation order not found")
		return
	}
	httpx.OK(w, item)
}

func (h ReincarnationHandler) Approve(w http.ResponseWriter, r *http.Request) {
	h.simpleTransition(w, r, "approved", `status IN ('queued', 'reviewing')`, "reincarnation.approve")
}

func (h ReincarnationHandler) Reject(w http.ResponseWriter, r *http.Request) {
	h.simpleTransition(w, r, "rejected", `status IN ('queued', 'reviewing', 'approved')`, "reincarnation.reject")
}

func (h ReincarnationHandler) Review(w http.ResponseWriter, r *http.Request) {
	h.simpleTransition(w, r, "reviewing", `status IN ('queued', 'approved')`, "reincarnation.review")
}

func (h ReincarnationHandler) AssignQuota(w http.ResponseWriter, r *http.Request) {
	id, ok := parseID(w, r)
	if !ok {
		return
	}
	var input struct {
		QuotaType string `json:"quotaType"`
	}
	if !readJSON(w, r, &input) {
		return
	}
	if input.QuotaType == "" {
		httpx.Fail(w, http.StatusBadRequest, 40051, "quota type is required")
		return
	}
	now := nowRFC3339()
	res, err := h.DB.Exec(`UPDATE reincarnation_orders SET status = 'pending_soup', quota_type = ?, updated_at = ? WHERE id = ? AND status IN ('approved', 'pending_soup')`, input.QuotaType, now, id)
	if err != nil {
		httpx.Fail(w, http.StatusInternalServerError, 53003, err.Error())
		return
	}
	if changed, _ := res.RowsAffected(); changed == 0 {
		httpx.Fail(w, http.StatusBadRequest, 40052, "quota can only be assigned after approval")
		return
	}
	audit(h.DB, r, "reincarnation.assign_quota", "reincarnation_order", id, input.QuotaType)
	httpx.OK(w, map[string]interface{}{"id": id, "status": "pending_soup", "quotaType": input.QuotaType})
}

func (h ReincarnationHandler) IssueSoup(w http.ResponseWriter, r *http.Request) {
	id, ok := parseID(w, r)
	if !ok {
		return
	}
	var input struct {
		InventoryID int64 `json:"inventoryId"`
		Dose        int   `json:"dose"`
		MemoryAfter int   `json:"memoryAfter"`
	}
	if !readJSON(w, r, &input) {
		return
	}
	if input.InventoryID <= 0 || input.Dose <= 0 {
		httpx.Fail(w, http.StatusBadRequest, 40053, "inventory id and dose are required")
		return
	}

	now := nowRFC3339()
	tx, err := h.DB.Begin()
	if err != nil {
		httpx.Fail(w, http.StatusInternalServerError, 53004, err.Error())
		return
	}
	defer tx.Rollback()

	var soulID int64
	var status string
	if err := tx.QueryRow(`SELECT soul_id, status FROM reincarnation_orders WHERE id = ?`, id).Scan(&soulID, &status); err != nil {
		httpx.Fail(w, http.StatusNotFound, 40442, "reincarnation order not found")
		return
	}
	if status != "pending_soup" {
		httpx.Fail(w, http.StatusBadRequest, 40054, "soup can only be issued when pending soup")
		return
	}
	var stock int
	if err := tx.QueryRow(`SELECT stock FROM soup_inventory WHERE id = ?`, input.InventoryID).Scan(&stock); err != nil {
		httpx.Fail(w, http.StatusNotFound, 40443, "soup inventory not found")
		return
	}
	if stock < input.Dose {
		httpx.Fail(w, http.StatusBadRequest, 40055, "soup stock is insufficient")
		return
	}
	if _, err := tx.Exec(`UPDATE soup_inventory SET stock = stock - ?, updated_at = ? WHERE id = ?`, input.Dose, now, input.InventoryID); err != nil {
		httpx.Fail(w, http.StatusInternalServerError, 53005, err.Error())
		return
	}
	if _, err := tx.Exec(`INSERT INTO soup_records (soul_id, order_id, inventory_id, dose, memory_after, operator_id, created_at) VALUES (?, ?, ?, ?, ?, ?, ?)`, soulID, id, input.InventoryID, input.Dose, input.MemoryAfter, httpx.CurrentUserID(r), now); err != nil {
		httpx.Fail(w, http.StatusInternalServerError, 53006, err.Error())
		return
	}
	if _, err := tx.Exec(`UPDATE reincarnation_orders SET soup_status = 'issued', updated_at = ? WHERE id = ?`, now, id); err != nil {
		httpx.Fail(w, http.StatusInternalServerError, 53007, err.Error())
		return
	}
	if err := tx.Commit(); err != nil {
		httpx.Fail(w, http.StatusInternalServerError, 53008, err.Error())
		return
	}
	audit(h.DB, r, "reincarnation.issue_soup", "reincarnation_order", id, "孟婆汤已发放")
	httpx.OK(w, map[string]interface{}{"id": id, "soupStatus": "issued"})
}

func (h ReincarnationHandler) Complete(w http.ResponseWriter, r *http.Request) {
	id, ok := parseID(w, r)
	if !ok {
		return
	}
	var input struct {
		Note string `json:"note"`
	}
	if !readJSON(w, r, &input) {
		return
	}
	now := nowRFC3339()
	res, err := h.DB.Exec(`UPDATE reincarnation_orders SET status = 'completed', updated_at = ? WHERE id = ? AND status = 'pending_soup' AND soup_status = 'issued'`, now, id)
	if err != nil {
		httpx.Fail(w, http.StatusInternalServerError, 53009, err.Error())
		return
	}
	if changed, _ := res.RowsAffected(); changed == 0 {
		httpx.Fail(w, http.StatusBadRequest, 40056, "reincarnation requires issued soup before completion")
		return
	}
	audit(h.DB, r, "reincarnation.complete", "reincarnation_order", id, input.Note)
	httpx.OK(w, map[string]interface{}{"id": id, "status": "completed"})
}

func (h ReincarnationHandler) simpleTransition(w http.ResponseWriter, r *http.Request, nextStatus string, condition string, action string) {
	id, ok := parseID(w, r)
	if !ok {
		return
	}
	var input struct {
		Note string `json:"note"`
	}
	if !readJSON(w, r, &input) {
		return
	}
	now := nowRFC3339()
	res, err := h.DB.Exec(`UPDATE reincarnation_orders SET status = ?, review_note = ?, updated_at = ? WHERE id = ? AND `+condition, nextStatus, input.Note, now, id)
	if err != nil {
		httpx.Fail(w, http.StatusInternalServerError, 53010, err.Error())
		return
	}
	if changed, _ := res.RowsAffected(); changed == 0 {
		httpx.Fail(w, http.StatusBadRequest, 40057, "invalid reincarnation status transition")
		return
	}
	audit(h.DB, r, action, "reincarnation_order", id, input.Note)
	httpx.OK(w, map[string]interface{}{"id": id, "status": nextStatus})
}

func scanReincarnation(row rowScanner) (map[string]interface{}, error) {
	var id, soulID int64
	var priority int
	var soulName, targetRealm, soupStatus, quotaType, status, reviewNote, createdAt, updatedAt string
	if err := row.Scan(&id, &soulID, &soulName, &targetRealm, &priority, &soupStatus, &quotaType, &status, &reviewNote, &createdAt, &updatedAt); err != nil {
		return nil, err
	}
	return map[string]interface{}{
		"id": id, "soulId": soulID, "soulName": soulName, "targetRealm": targetRealm,
		"priority": priority, "soupStatus": soupStatus, "quotaType": quotaType,
		"status": status, "reviewNote": reviewNote, "createdAt": createdAt, "updatedAt": updatedAt,
	}, nil
}
