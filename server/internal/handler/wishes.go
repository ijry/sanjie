package handler

import (
	"database/sql"
	"net/http"

	"github.com/ijry/sanjie/internal/httpx"
)

type WishesHandler struct {
	DB *sql.DB
}

func (h WishesHandler) List(w http.ResponseWriter, r *http.Request) {
	status := r.URL.Query().Get("status")
	wishType := r.URL.Query().Get("wishType")
	query := `SELECT id, title, wish_type, requester_name, incense_amount, merit_score, status, assigned_deity, result_note, created_at, updated_at FROM wish_tickets WHERE 1=1`
	args := []interface{}{}
	if status != "" {
		query += ` AND status = ?`
		args = append(args, status)
	}
	if wishType != "" {
		query += ` AND wish_type = ?`
		args = append(args, wishType)
	}
	query += ` ORDER BY id DESC`
	rows, err := h.DB.Query(query, args...)
	if err != nil {
		httpx.Fail(w, http.StatusInternalServerError, 55001, err.Error())
		return
	}
	defer rows.Close()
	items := make([]map[string]interface{}, 0)
	for rows.Next() {
		item, err := scanWish(rows)
		if err != nil {
			httpx.Fail(w, http.StatusInternalServerError, 55002, err.Error())
			return
		}
		items = append(items, item)
	}
	httpx.OK(w, items)
}

func (h WishesHandler) Detail(w http.ResponseWriter, r *http.Request) {
	id, ok := parseID(w, r)
	if !ok {
		return
	}
	item, err := scanWish(h.DB.QueryRow(`SELECT id, title, wish_type, requester_name, incense_amount, merit_score, status, assigned_deity, result_note, created_at, updated_at FROM wish_tickets WHERE id = ?`, id))
	if err != nil {
		httpx.Fail(w, http.StatusNotFound, 40461, "wish ticket not found")
		return
	}
	httpx.OK(w, item)
}

func (h WishesHandler) Route(w http.ResponseWriter, r *http.Request) {
	id, ok := parseID(w, r)
	if !ok {
		return
	}
	var input struct {
		AssignedDeity string `json:"assignedDeity"`
		Note          string `json:"note"`
	}
	if !readJSON(w, r, &input) {
		return
	}
	if input.AssignedDeity == "" {
		httpx.Fail(w, http.StatusBadRequest, 40071, "assigned deity is required")
		return
	}
	now := nowRFC3339()
	res, err := h.DB.Exec(`UPDATE wish_tickets SET status = 'routed', assigned_deity = ?, result_note = ?, updated_at = ? WHERE id = ? AND status = 'pending'`, input.AssignedDeity, input.Note, now, id)
	if err != nil {
		httpx.Fail(w, http.StatusInternalServerError, 55003, err.Error())
		return
	}
	if changed, _ := res.RowsAffected(); changed == 0 {
		httpx.Fail(w, http.StatusBadRequest, 40072, "only pending wish can route")
		return
	}
	audit(h.DB, r, "wish.route", "wish_ticket", id, input.Note)
	httpx.OK(w, map[string]interface{}{"id": id, "status": "routed"})
}

func (h WishesHandler) Resolve(w http.ResponseWriter, r *http.Request) {
	h.finish(w, r, "resolved", "wish.resolve")
}

func (h WishesHandler) Reject(w http.ResponseWriter, r *http.Request) {
	h.finish(w, r, "rejected", "wish.reject")
}

func (h WishesHandler) finish(w http.ResponseWriter, r *http.Request, status string, action string) {
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
	res, err := h.DB.Exec(`UPDATE wish_tickets SET status = ?, result_note = ?, updated_at = ? WHERE id = ? AND status IN ('pending', 'routed')`, status, input.Note, now, id)
	if err != nil {
		httpx.Fail(w, http.StatusInternalServerError, 55004, err.Error())
		return
	}
	if changed, _ := res.RowsAffected(); changed == 0 {
		httpx.Fail(w, http.StatusBadRequest, 40073, "wish ticket is not actionable")
		return
	}
	audit(h.DB, r, action, "wish_ticket", id, input.Note)
	httpx.OK(w, map[string]interface{}{"id": id, "status": status})
}

func scanWish(row rowScanner) (map[string]interface{}, error) {
	var id int64
	var incenseAmount, meritScore int
	var title, wishType, requesterName, status, assignedDeity, resultNote, createdAt, updatedAt string
	if err := row.Scan(&id, &title, &wishType, &requesterName, &incenseAmount, &meritScore, &status, &assignedDeity, &resultNote, &createdAt, &updatedAt); err != nil {
		return nil, err
	}
	return map[string]interface{}{
		"id": id, "title": title, "wishType": wishType, "requesterName": requesterName,
		"incenseAmount": incenseAmount, "meritScore": meritScore, "status": status,
		"assignedDeity": assignedDeity, "resultNote": resultNote, "createdAt": createdAt, "updatedAt": updatedAt,
	}, nil
}
