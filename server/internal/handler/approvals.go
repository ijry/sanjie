package handler

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/ijry/sanjie/internal/httpx"
)

type ApprovalsHandler struct {
	DB *sql.DB
}

func (h ApprovalsHandler) List(w http.ResponseWriter, r *http.Request) {
	status := r.URL.Query().Get("status")
	approvalType := r.URL.Query().Get("type")
	query := `SELECT id, type, target_id, title, status, COALESCE(applicant_id, 0), COALESCE(approver_id, 0), reason, result_note, created_at, updated_at FROM approvals WHERE 1=1`
	args := []interface{}{}
	if status != "" {
		query += ` AND status = ?`
		args = append(args, status)
	}
	if approvalType != "" {
		query += ` AND type = ?`
		args = append(args, approvalType)
	}
	query += ` ORDER BY id DESC`
	rows, err := h.DB.Query(query, args...)
	if err != nil {
		httpx.Fail(w, http.StatusInternalServerError, 52201, err.Error())
		return
	}
	defer rows.Close()

	items := make([]map[string]interface{}, 0)
	for rows.Next() {
		item, err := scanApproval(rows)
		if err != nil {
			httpx.Fail(w, http.StatusInternalServerError, 52202, err.Error())
			return
		}
		items = append(items, item)
	}
	httpx.OK(w, items)
}

func (h ApprovalsHandler) Detail(w http.ResponseWriter, r *http.Request) {
	id, ok := parseID(w, r)
	if !ok {
		return
	}
	item, err := scanApproval(h.DB.QueryRow(`SELECT id, type, target_id, title, status, COALESCE(applicant_id, 0), COALESCE(approver_id, 0), reason, result_note, created_at, updated_at FROM approvals WHERE id = ?`, id))
	if err != nil {
		httpx.Fail(w, http.StatusNotFound, 40431, "approval not found")
		return
	}
	httpx.OK(w, item)
}

func (h ApprovalsHandler) Approve(w http.ResponseWriter, r *http.Request) {
	h.finish(w, r, "approved")
}

func (h ApprovalsHandler) Reject(w http.ResponseWriter, r *http.Request) {
	h.finish(w, r, "rejected")
}

func (h ApprovalsHandler) Transfer(w http.ResponseWriter, r *http.Request) {
	h.finish(w, r, "transferred")
}

func (h ApprovalsHandler) finish(w http.ResponseWriter, r *http.Request, nextStatus string) {
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

	var approvalType, status, reason string
	var targetID int64
	if err := h.DB.QueryRow(`SELECT type, target_id, status, reason FROM approvals WHERE id = ?`, id).Scan(&approvalType, &targetID, &status, &reason); err != nil {
		httpx.Fail(w, http.StatusNotFound, 40432, "approval not found")
		return
	}
	if nextStatus == "transferred" && status != "pending" {
		httpx.Fail(w, http.StatusBadRequest, 40041, "only pending approval can transfer")
		return
	}
	if nextStatus != "transferred" && status != "pending" && status != "transferred" {
		httpx.Fail(w, http.StatusBadRequest, 40042, "approval is not actionable")
		return
	}

	now := nowRFC3339()
	tx, err := h.DB.Begin()
	if err != nil {
		httpx.Fail(w, http.StatusInternalServerError, 52203, err.Error())
		return
	}
	defer tx.Rollback()

	if _, err := tx.Exec(`UPDATE approvals SET status = ?, approver_id = ?, result_note = ?, updated_at = ? WHERE id = ?`, nextStatus, httpx.CurrentUserID(r), input.Note, now, id); err != nil {
		httpx.Fail(w, http.StatusInternalServerError, 52204, err.Error())
		return
	}

	if nextStatus == "approved" && approvalType == "lifespan_change" {
		var payload struct {
			NewLifespan int `json:"newLifespan"`
		}
		_ = json.Unmarshal([]byte(reason), &payload)
		if payload.NewLifespan > 0 {
			if _, err := tx.Exec(`UPDATE life_book_records SET actual_lifespan = ?, locked = 0, risk_flag = 'normal', updated_at = ? WHERE id = ?`, payload.NewLifespan, now, targetID); err != nil {
				httpx.Fail(w, http.StatusInternalServerError, 52205, err.Error())
				return
			}
		}
	}

	if _, err := tx.Exec(`INSERT INTO audit_logs (actor_id, action, target_type, target_id, note, created_at) VALUES (?, ?, 'approval', ?, ?, ?)`,
		httpx.CurrentUserID(r), "approval."+nextStatus, id, input.Note, now); err != nil {
		httpx.Fail(w, http.StatusInternalServerError, 52206, err.Error())
		return
	}
	if err := tx.Commit(); err != nil {
		httpx.Fail(w, http.StatusInternalServerError, 52207, err.Error())
		return
	}
	httpx.OK(w, map[string]interface{}{"id": id, "status": nextStatus})
}

func scanApproval(row rowScanner) (map[string]interface{}, error) {
	var id, targetID, applicantID, approverID int64
	var approvalType, title, status, reason, resultNote, createdAt, updatedAt string
	if err := row.Scan(&id, &approvalType, &targetID, &title, &status, &applicantID, &approverID, &reason, &resultNote, &createdAt, &updatedAt); err != nil {
		return nil, err
	}
	return map[string]interface{}{
		"id": id, "type": approvalType, "targetId": targetID, "title": title, "status": status,
		"applicantId": applicantID, "approverId": approverID, "reason": reason,
		"resultNote": resultNote, "createdAt": createdAt, "updatedAt": updatedAt,
	}, nil
}
