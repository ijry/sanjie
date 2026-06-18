package handler

import (
	"database/sql"
	"net/http"

	"github.com/ijry/sanjie/internal/httpx"
)

type AuditHandler struct {
	DB *sql.DB
}

func (h AuditHandler) Logs(w http.ResponseWriter, r *http.Request) {
	rows, err := h.DB.Query(`SELECT id, COALESCE(actor_id, 0), action, target_type, target_id, before_json, after_json, note, created_at FROM audit_logs ORDER BY id DESC LIMIT 100`)
	if err != nil {
		httpx.Fail(w, http.StatusInternalServerError, 50005, err.Error())
		return
	}
	defer rows.Close()

	items := make([]map[string]interface{}, 0)
	for rows.Next() {
		var id, actorID, targetID int64
		var action, targetType, beforeJSON, afterJSON, note, createdAt string
		if err := rows.Scan(&id, &actorID, &action, &targetType, &targetID, &beforeJSON, &afterJSON, &note, &createdAt); err != nil {
			httpx.Fail(w, http.StatusInternalServerError, 50006, err.Error())
			return
		}
		items = append(items, map[string]interface{}{
			"id": id, "actorId": actorID, "action": action, "targetType": targetType,
			"targetId": targetID, "beforeJson": beforeJSON, "afterJson": afterJSON,
			"note": note, "createdAt": createdAt,
		})
	}
	httpx.OK(w, items)
}
