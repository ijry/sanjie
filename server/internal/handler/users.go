package handler

import (
	"database/sql"
	"net/http"

	"github.com/ijry/sanjie/internal/httpx"
)

type UsersHandler struct {
	DB *sql.DB
}

func (h UsersHandler) Me(w http.ResponseWriter, r *http.Request) {
	userID := httpx.CurrentUserID(r)
	var id int64
	var username, nickname, role, avatar, createdAt string
	err := h.DB.QueryRow(`SELECT id, username, nickname, role, avatar, created_at FROM users WHERE id = ?`, userID).
		Scan(&id, &username, &nickname, &role, &avatar, &createdAt)
	if err != nil {
		httpx.Fail(w, http.StatusNotFound, 40401, "user not found")
		return
	}
	httpx.OK(w, map[string]interface{}{
		"id": id, "username": username, "nickname": nickname,
		"role": role, "avatar": avatar, "createdAt": createdAt,
	})
}

func (h UsersHandler) Switch(w http.ResponseWriter, r *http.Request) {
	var input struct {
		UserID int64 `json:"userId"`
	}
	if !readJSON(w, r, &input) {
		return
	}
	if input.UserID <= 0 {
		httpx.Fail(w, http.StatusBadRequest, 40012, "invalid user id")
		return
	}
	httpx.OK(w, map[string]interface{}{"userId": input.UserID})
}
