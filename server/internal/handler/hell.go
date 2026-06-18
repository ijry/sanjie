package handler

import (
	"database/sql"
	"net/http"

	"github.com/ijry/sanjie/internal/httpx"
)

type HellHandler struct {
	DB *sql.DB
}

func (h HellHandler) Floors(w http.ResponseWriter, r *http.Request) {
	rows, err := h.DB.Query(`SELECT id, floor_no, name, capacity, occupancy, equipment_status, load_level, updated_at FROM hell_floors ORDER BY floor_no ASC`)
	if err != nil {
		httpx.Fail(w, http.StatusInternalServerError, 54001, err.Error())
		return
	}
	defer rows.Close()
	items := make([]map[string]interface{}, 0)
	for rows.Next() {
		var id int64
		var floorNo, capacity, occupancy int
		var name, equipmentStatus, loadLevel, updatedAt string
		if err := rows.Scan(&id, &floorNo, &name, &capacity, &occupancy, &equipmentStatus, &loadLevel, &updatedAt); err != nil {
			httpx.Fail(w, http.StatusInternalServerError, 54002, err.Error())
			return
		}
		items = append(items, map[string]interface{}{"id": id, "floorNo": floorNo, "name": name, "capacity": capacity, "occupancy": occupancy, "equipmentStatus": equipmentStatus, "loadLevel": loadLevel, "updatedAt": updatedAt})
	}
	httpx.OK(w, items)
}

func (h HellHandler) Sentences(w http.ResponseWriter, r *http.Request) {
	reviewStatus := r.URL.Query().Get("reviewStatus")
	query := `SELECT hs.id, hs.soul_id, s.name, hs.floor_id, hf.name, hs.crime_type, hs.sentence_days, hs.pain_level, hs.review_status, hs.equipment_id, hs.created_at, hs.updated_at
		FROM hell_sentences hs
		JOIN souls s ON s.id = hs.soul_id
		JOIN hell_floors hf ON hf.id = hs.floor_id`
	args := []interface{}{}
	if reviewStatus != "" {
		query += ` WHERE hs.review_status = ?`
		args = append(args, reviewStatus)
	}
	query += ` ORDER BY hs.id DESC`
	rows, err := h.DB.Query(query, args...)
	if err != nil {
		httpx.Fail(w, http.StatusInternalServerError, 54003, err.Error())
		return
	}
	defer rows.Close()
	items := make([]map[string]interface{}, 0)
	for rows.Next() {
		item, err := scanHellSentence(rows)
		if err != nil {
			httpx.Fail(w, http.StatusInternalServerError, 54004, err.Error())
			return
		}
		items = append(items, item)
	}
	httpx.OK(w, items)
}

func (h HellHandler) SentenceDetail(w http.ResponseWriter, r *http.Request) {
	id, ok := parseID(w, r)
	if !ok {
		return
	}
	item, err := scanHellSentence(h.DB.QueryRow(`SELECT hs.id, hs.soul_id, s.name, hs.floor_id, hf.name, hs.crime_type, hs.sentence_days, hs.pain_level, hs.review_status, hs.equipment_id, hs.created_at, hs.updated_at
		FROM hell_sentences hs
		JOIN souls s ON s.id = hs.soul_id
		JOIN hell_floors hf ON hf.id = hs.floor_id
		WHERE hs.id = ?`, id))
	if err != nil {
		httpx.Fail(w, http.StatusNotFound, 40451, "hell sentence not found")
		return
	}
	httpx.OK(w, item)
}

func (h HellHandler) ReviewSentence(w http.ResponseWriter, r *http.Request) {
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
	tx, err := h.DB.Begin()
	if err != nil {
		httpx.Fail(w, http.StatusInternalServerError, 54005, err.Error())
		return
	}
	defer tx.Rollback()
	res, err := tx.Exec(`UPDATE hell_sentences SET review_status = 'reviewing', updated_at = ? WHERE id = ? AND review_status = 'none'`, now, id)
	if err != nil {
		httpx.Fail(w, http.StatusInternalServerError, 54006, err.Error())
		return
	}
	if changed, _ := res.RowsAffected(); changed == 0 {
		httpx.Fail(w, http.StatusBadRequest, 40061, "hell sentence is not reviewable")
		return
	}
	if _, err := tx.Exec(`INSERT INTO approvals (type, target_id, title, status, applicant_id, reason, created_at, updated_at) VALUES ('hell_review', ?, '地狱刑期复核', 'pending', ?, ?, ?, ?)`, id, httpx.CurrentUserID(r), input.Note, now, now); err != nil {
		httpx.Fail(w, http.StatusInternalServerError, 54007, err.Error())
		return
	}
	if err := tx.Commit(); err != nil {
		httpx.Fail(w, http.StatusInternalServerError, 54008, err.Error())
		return
	}
	audit(h.DB, r, "hell_sentence.review", "hell_sentence", id, input.Note)
	httpx.OK(w, map[string]interface{}{"id": id, "reviewStatus": "reviewing"})
}

func (h HellHandler) DispatchFloor(w http.ResponseWriter, r *http.Request) {
	id, ok := parseID(w, r)
	if !ok {
		return
	}
	var input struct {
		TargetFloorID int64 `json:"targetFloorId"`
		Amount        int   `json:"amount"`
	}
	if !readJSON(w, r, &input) {
		return
	}
	if input.TargetFloorID <= 0 || input.Amount <= 0 || input.TargetFloorID == id {
		httpx.Fail(w, http.StatusBadRequest, 40062, "target floor and amount are invalid")
		return
	}
	now := nowRFC3339()
	tx, err := h.DB.Begin()
	if err != nil {
		httpx.Fail(w, http.StatusInternalServerError, 54009, err.Error())
		return
	}
	defer tx.Rollback()
	var sourceCap, sourceOcc, targetCap, targetOcc int
	if err := tx.QueryRow(`SELECT capacity, occupancy FROM hell_floors WHERE id = ?`, id).Scan(&sourceCap, &sourceOcc); err != nil {
		httpx.Fail(w, http.StatusNotFound, 40452, "source floor not found")
		return
	}
	if err := tx.QueryRow(`SELECT capacity, occupancy FROM hell_floors WHERE id = ?`, input.TargetFloorID).Scan(&targetCap, &targetOcc); err != nil {
		httpx.Fail(w, http.StatusNotFound, 40453, "target floor not found")
		return
	}
	if sourceOcc < input.Amount {
		httpx.Fail(w, http.StatusBadRequest, 40063, "source floor occupancy is insufficient")
		return
	}
	sourceOcc -= input.Amount
	targetOcc += input.Amount
	if _, err := tx.Exec(`UPDATE hell_floors SET occupancy = ?, load_level = ?, updated_at = ? WHERE id = ?`, sourceOcc, loadLevel(sourceOcc, sourceCap), now, id); err != nil {
		httpx.Fail(w, http.StatusInternalServerError, 54010, err.Error())
		return
	}
	if _, err := tx.Exec(`UPDATE hell_floors SET occupancy = ?, load_level = ?, updated_at = ? WHERE id = ?`, targetOcc, loadLevel(targetOcc, targetCap), now, input.TargetFloorID); err != nil {
		httpx.Fail(w, http.StatusInternalServerError, 54011, err.Error())
		return
	}
	if err := tx.Commit(); err != nil {
		httpx.Fail(w, http.StatusInternalServerError, 54012, err.Error())
		return
	}
	audit(h.DB, r, "hell_floor.dispatch", "hell_floor", id, "跨层分流")
	httpx.OK(w, map[string]interface{}{"sourceFloorId": id, "targetFloorId": input.TargetFloorID, "amount": input.Amount})
}

func scanHellSentence(row rowScanner) (map[string]interface{}, error) {
	var id, soulID, floorID int64
	var sentenceDays, painLevel int
	var soulName, floorName, crimeType, reviewStatus, equipmentID, createdAt, updatedAt string
	if err := row.Scan(&id, &soulID, &soulName, &floorID, &floorName, &crimeType, &sentenceDays, &painLevel, &reviewStatus, &equipmentID, &createdAt, &updatedAt); err != nil {
		return nil, err
	}
	return map[string]interface{}{
		"id": id, "soulId": soulID, "soulName": soulName, "floorId": floorID, "floorName": floorName,
		"crimeType": crimeType, "sentenceDays": sentenceDays, "painLevel": painLevel, "reviewStatus": reviewStatus,
		"equipmentId": equipmentID, "createdAt": createdAt, "updatedAt": updatedAt,
	}, nil
}
