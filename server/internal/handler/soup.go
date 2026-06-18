package handler

import (
	"database/sql"
	"net/http"

	"github.com/ijry/sanjie/internal/httpx"
)

type SoupHandler struct {
	DB *sql.DB
}

func (h SoupHandler) Inventory(w http.ResponseWriter, r *http.Request) {
	rows, err := h.DB.Query(`SELECT id, name, stock, warning_stock, formula_note, updated_at FROM soup_inventory ORDER BY id ASC`)
	if err != nil {
		httpx.Fail(w, http.StatusInternalServerError, 56001, err.Error())
		return
	}
	defer rows.Close()
	items := make([]map[string]interface{}, 0)
	for rows.Next() {
		var id int64
		var stock, warningStock int
		var name, formulaNote, updatedAt string
		if err := rows.Scan(&id, &name, &stock, &warningStock, &formulaNote, &updatedAt); err != nil {
			httpx.Fail(w, http.StatusInternalServerError, 56002, err.Error())
			return
		}
		items = append(items, map[string]interface{}{"id": id, "name": name, "stock": stock, "warningStock": warningStock, "formulaNote": formulaNote, "updatedAt": updatedAt, "warning": stock <= warningStock})
	}
	httpx.OK(w, items)
}

func (h SoupHandler) AdjustInventory(w http.ResponseWriter, r *http.Request) {
	id, ok := parseID(w, r)
	if !ok {
		return
	}
	var input struct {
		Delta int    `json:"delta"`
		Note  string `json:"note"`
	}
	if !readJSON(w, r, &input) {
		return
	}
	var stock int
	if err := h.DB.QueryRow(`SELECT stock FROM soup_inventory WHERE id = ?`, id).Scan(&stock); err != nil {
		httpx.Fail(w, http.StatusNotFound, 40471, "soup inventory not found")
		return
	}
	next := stock + input.Delta
	if next < 0 {
		httpx.Fail(w, http.StatusBadRequest, 40081, "soup stock cannot be negative")
		return
	}
	now := nowRFC3339()
	if _, err := h.DB.Exec(`UPDATE soup_inventory SET stock = ?, updated_at = ? WHERE id = ?`, next, now, id); err != nil {
		httpx.Fail(w, http.StatusInternalServerError, 56003, err.Error())
		return
	}
	audit(h.DB, r, "soup.adjust", "soup_inventory", id, input.Note)
	httpx.OK(w, map[string]interface{}{"id": id, "stock": next})
}

func (h SoupHandler) Records(w http.ResponseWriter, r *http.Request) {
	rows, err := h.DB.Query(`SELECT sr.id, sr.soul_id, s.name, COALESCE(sr.order_id, 0), sr.inventory_id, si.name, sr.dose, sr.memory_after, COALESCE(sr.operator_id, 0), sr.created_at
		FROM soup_records sr
		JOIN souls s ON s.id = sr.soul_id
		JOIN soup_inventory si ON si.id = sr.inventory_id
		ORDER BY sr.id DESC LIMIT 100`)
	if err != nil {
		httpx.Fail(w, http.StatusInternalServerError, 56004, err.Error())
		return
	}
	defer rows.Close()
	items := make([]map[string]interface{}, 0)
	for rows.Next() {
		var id, soulID, orderID, inventoryID, operatorID int64
		var dose, memoryAfter int
		var soulName, inventoryName, createdAt string
		if err := rows.Scan(&id, &soulID, &soulName, &orderID, &inventoryID, &inventoryName, &dose, &memoryAfter, &operatorID, &createdAt); err != nil {
			httpx.Fail(w, http.StatusInternalServerError, 56005, err.Error())
			return
		}
		items = append(items, map[string]interface{}{"id": id, "soulId": soulID, "soulName": soulName, "orderId": orderID, "inventoryId": inventoryID, "inventoryName": inventoryName, "dose": dose, "memoryAfter": memoryAfter, "operatorId": operatorID, "createdAt": createdAt})
	}
	httpx.OK(w, items)
}
