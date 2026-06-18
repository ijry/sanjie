package handler

import (
	"database/sql"
	"net/http"

	"github.com/ijry/sanjie/internal/httpx"
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
		httpx.Fail(w, http.StatusInternalServerError, 51001, err.Error())
		return
	}
	defer rows.Close()

	items := make([]map[string]interface{}, 0)
	for rows.Next() {
		item, err := scanCapture(rows)
		if err != nil {
			httpx.Fail(w, http.StatusInternalServerError, 51002, err.Error())
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
	now := nowRFC3339()
	res, err := h.DB.Exec(`UPDATE capture_tasks SET status = 'running', updated_at = ? WHERE id = ? AND status = 'pending'`, now, id)
	if err != nil {
		httpx.Fail(w, http.StatusInternalServerError, 51003, err.Error())
		return
	}
	if changed, _ := res.RowsAffected(); changed == 0 {
		httpx.Fail(w, http.StatusBadRequest, 40021, "only pending capture tasks can start")
		return
	}
	audit(h.DB, r, "capture.start", "capture_task", id, "开始勾魂")
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
	if !readJSON(w, r, &input) {
		return
	}
	now := nowRFC3339()
	res, err := h.DB.Exec(`UPDATE capture_tasks SET status = 'completed', actual_time = ?, updated_at = ? WHERE id = ? AND status = 'running'`, now, now, id)
	if err != nil {
		httpx.Fail(w, http.StatusInternalServerError, 51004, err.Error())
		return
	}
	if changed, _ := res.RowsAffected(); changed == 0 {
		httpx.Fail(w, http.StatusBadRequest, 40022, "only running capture tasks can complete")
		return
	}
	audit(h.DB, r, "capture.complete", "capture_task", id, input.Note)
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
	if !readJSON(w, r, &input) {
		return
	}
	if input.Type == "" {
		httpx.Fail(w, http.StatusBadRequest, 40023, "exception type is required")
		return
	}
	now := nowRFC3339()
	if _, err := h.DB.Exec(`UPDATE capture_tasks SET status = 'exception', exception_type = ?, exception_note = ?, updated_at = ? WHERE id = ?`, input.Type, input.Note, now, id); err != nil {
		httpx.Fail(w, http.StatusInternalServerError, 51005, err.Error())
		return
	}
	audit(h.DB, r, "capture.exception", "capture_task", id, input.Note)
	httpx.OK(w, map[string]interface{}{"id": id, "status": "exception"})
}

func (h CaptureHandler) Cancel(w http.ResponseWriter, r *http.Request) {
	id, ok := parseID(w, r)
	if !ok {
		return
	}
	now := nowRFC3339()
	res, err := h.DB.Exec(`UPDATE capture_tasks SET status = 'canceled', updated_at = ? WHERE id = ? AND status IN ('pending', 'running', 'exception')`, now, id)
	if err != nil {
		httpx.Fail(w, http.StatusInternalServerError, 51006, err.Error())
		return
	}
	if changed, _ := res.RowsAffected(); changed == 0 {
		httpx.Fail(w, http.StatusBadRequest, 40024, "capture task cannot be canceled")
		return
	}
	audit(h.DB, r, "capture.cancel", "capture_task", id, "撤销勾魂任务")
	httpx.OK(w, map[string]interface{}{"id": id, "status": "canceled"})
}

type rowScanner interface {
	Scan(dest ...interface{}) error
}

func scanCapture(row rowScanner) (map[string]interface{}, error) {
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
