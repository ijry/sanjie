package handler

import (
	"database/sql"
	"net/http"

	"github.com/ijry/sanjie/internal/httpx"
)

type SoulsHandler struct {
	DB *sql.DB
}

func (h SoulsHandler) List(w http.ResponseWriter, r *http.Request) {
	keyword := r.URL.Query().Get("keyword")
	status := r.URL.Query().Get("status")
	query := `SELECT id, name, birth_info, death_info, status, merit_score, karma_score, memory_residue, relationship_risk, created_at, updated_at FROM souls WHERE 1=1`
	args := []interface{}{}
	if keyword != "" {
		query += ` AND name LIKE ?`
		args = append(args, "%"+keyword+"%")
	}
	if status != "" {
		query += ` AND status = ?`
		args = append(args, status)
	}
	query += ` ORDER BY id DESC`

	rows, err := h.DB.Query(query, args...)
	if err != nil {
		httpx.Fail(w, http.StatusInternalServerError, 52001, err.Error())
		return
	}
	defer rows.Close()
	items := make([]map[string]interface{}, 0)
	for rows.Next() {
		item, err := scanSoul(rows)
		if err != nil {
			httpx.Fail(w, http.StatusInternalServerError, 52002, err.Error())
			return
		}
		items = append(items, item)
	}
	httpx.OK(w, items)
}

func (h SoulsHandler) Detail(w http.ResponseWriter, r *http.Request) {
	id, ok := parseID(w, r)
	if !ok {
		return
	}
	item, err := scanSoul(h.DB.QueryRow(`SELECT id, name, birth_info, death_info, status, merit_score, karma_score, memory_residue, relationship_risk, created_at, updated_at FROM souls WHERE id = ?`, id))
	if err != nil {
		httpx.Fail(w, http.StatusNotFound, 40421, "soul not found")
		return
	}
	httpx.OK(w, item)
}

func scanSoul(row rowScanner) (map[string]interface{}, error) {
	var id int64
	var meritScore, karmaScore, memoryResidue int
	var name, birthInfo, deathInfo, status, relationshipRisk, createdAt, updatedAt string
	if err := row.Scan(&id, &name, &birthInfo, &deathInfo, &status, &meritScore, &karmaScore, &memoryResidue, &relationshipRisk, &createdAt, &updatedAt); err != nil {
		return nil, err
	}
	return map[string]interface{}{
		"id": id, "name": name, "birthInfo": birthInfo, "deathInfo": deathInfo, "status": status,
		"meritScore": meritScore, "karmaScore": karmaScore, "memoryResidue": memoryResidue,
		"relationshipRisk": relationshipRisk, "createdAt": createdAt, "updatedAt": updatedAt,
	}, nil
}
