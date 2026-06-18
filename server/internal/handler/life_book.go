package handler

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/ijry/sanjie/internal/httpx"
)

type LifeBookHandler struct {
	DB *sql.DB
}

func (h LifeBookHandler) List(w http.ResponseWriter, r *http.Request) {
	keyword := r.URL.Query().Get("keyword")
	riskFlag := r.URL.Query().Get("riskFlag")
	query := `SELECT l.id, l.soul_id, s.name, l.expected_lifespan, l.actual_lifespan, l.death_reason, l.fate_level, l.calamity_count, l.locked, l.risk_flag, l.updated_at
		FROM life_book_records l
		JOIN souls s ON s.id = l.soul_id
		WHERE 1=1`
	args := []interface{}{}
	if keyword != "" {
		query += ` AND s.name LIKE ?`
		args = append(args, "%"+keyword+"%")
	}
	if riskFlag != "" && riskFlag != "all" {
		query += ` AND l.risk_flag = ?`
		args = append(args, riskFlag)
	}
	query += ` ORDER BY l.id DESC`

	rows, err := h.DB.Query(query, args...)
	if err != nil {
		httpx.Fail(w, http.StatusInternalServerError, 52101, err.Error())
		return
	}
	defer rows.Close()
	items := make([]map[string]interface{}, 0)
	for rows.Next() {
		item, err := scanLifeBook(rows)
		if err != nil {
			httpx.Fail(w, http.StatusInternalServerError, 52102, err.Error())
			return
		}
		items = append(items, item)
	}
	httpx.OK(w, items)
}

func (h LifeBookHandler) Detail(w http.ResponseWriter, r *http.Request) {
	id, ok := parseID(w, r)
	if !ok {
		return
	}
	item, err := scanLifeBook(h.DB.QueryRow(`SELECT l.id, l.soul_id, s.name, l.expected_lifespan, l.actual_lifespan, l.death_reason, l.fate_level, l.calamity_count, l.locked, l.risk_flag, l.updated_at
		FROM life_book_records l
		JOIN souls s ON s.id = l.soul_id
		WHERE l.id = ?`, id))
	if err != nil {
		httpx.Fail(w, http.StatusNotFound, 40422, "life book record not found")
		return
	}
	httpx.OK(w, item)
}

func (h LifeBookHandler) Freeze(w http.ResponseWriter, r *http.Request) {
	id, ok := parseID(w, r)
	if !ok {
		return
	}
	now := nowRFC3339()
	if _, err := h.DB.Exec(`UPDATE life_book_records SET locked = 1, risk_flag = 'critical', updated_at = ? WHERE id = ?`, now, id); err != nil {
		httpx.Fail(w, http.StatusInternalServerError, 52103, err.Error())
		return
	}
	audit(h.DB, r, "life_book.freeze", "life_book_record", id, "冻结生死簿记录")
	httpx.OK(w, map[string]interface{}{"id": id, "locked": true})
}

func (h LifeBookHandler) ChangeLifespan(w http.ResponseWriter, r *http.Request) {
	id, ok := parseID(w, r)
	if !ok {
		return
	}
	var input struct {
		NewLifespan int    `json:"newLifespan"`
		Reason      string `json:"reason"`
	}
	if !readJSON(w, r, &input) {
		return
	}
	if input.NewLifespan <= 0 || input.Reason == "" {
		httpx.Fail(w, http.StatusBadRequest, 40031, "new lifespan and reason are required")
		return
	}
	raw, _ := json.Marshal(input)
	now := nowRFC3339()
	res, err := h.DB.Exec(`INSERT INTO approvals (type, target_id, title, status, applicant_id, reason, created_at, updated_at) VALUES ('lifespan_change', ?, '寿命变更审批', 'pending', ?, ?, ?, ?)`,
		id, httpx.CurrentUserID(r), string(raw), now, now)
	if err != nil {
		httpx.Fail(w, http.StatusInternalServerError, 52104, err.Error())
		return
	}
	approvalID, _ := res.LastInsertId()
	audit(h.DB, r, "life_book.change_lifespan.request", "life_book_record", id, input.Reason)
	httpx.OK(w, map[string]interface{}{"approvalId": approvalID})
}

func (h LifeBookHandler) AuditLogs(w http.ResponseWriter, r *http.Request) {
	id, ok := parseID(w, r)
	if !ok {
		return
	}
	rows, err := h.DB.Query(`SELECT id, COALESCE(actor_id, 0), action, target_type, target_id, note, created_at FROM audit_logs WHERE target_type = 'life_book_record' AND target_id = ? ORDER BY id DESC`, id)
	if err != nil {
		httpx.Fail(w, http.StatusInternalServerError, 52105, err.Error())
		return
	}
	defer rows.Close()
	items := make([]map[string]interface{}, 0)
	for rows.Next() {
		var logID, actorID, targetID int64
		var action, targetType, note, createdAt string
		if err := rows.Scan(&logID, &actorID, &action, &targetType, &targetID, &note, &createdAt); err != nil {
			httpx.Fail(w, http.StatusInternalServerError, 52106, err.Error())
			return
		}
		items = append(items, map[string]interface{}{"id": logID, "actorId": actorID, "action": action, "targetType": targetType, "targetId": targetID, "note": note, "createdAt": createdAt})
	}
	httpx.OK(w, items)
}

func scanLifeBook(row rowScanner) (map[string]interface{}, error) {
	var id, soulID int64
	var expectedLifespan, actualLifespan, calamityCount, locked int
	var soulName, deathReason, fateLevel, riskFlag, updatedAt string
	if err := row.Scan(&id, &soulID, &soulName, &expectedLifespan, &actualLifespan, &deathReason, &fateLevel, &calamityCount, &locked, &riskFlag, &updatedAt); err != nil {
		return nil, err
	}
	return map[string]interface{}{
		"id": id, "soulId": soulID, "soulName": soulName, "expectedLifespan": expectedLifespan,
		"actualLifespan": actualLifespan, "deathReason": deathReason, "fateLevel": fateLevel,
		"calamityCount": calamityCount, "locked": boolFromInt(locked), "riskFlag": riskFlag, "updatedAt": updatedAt,
	}, nil
}
