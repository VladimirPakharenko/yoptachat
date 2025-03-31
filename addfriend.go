package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func addFriendHandler(c *gin.Context) {
	var request struct {
		FriendID int `json:"friend_id"`
	}
	if err := c.ShouldBindJSON(&request); err != nil {
		c.String(http.StatusBadRequest, "Неверный запрос")
		return
	}

	session, _ := store.Get(c.Request, "session")
	ID := session.Values["user_id"].(int)

	var exists bool
	err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM Friends WHERE user_id = ? AND friend_id = ?)", ID, request.FriendID).Scan(&exists)
	if err != nil {
		c.String(http.StatusInternalServerError, "Ошибка при проверке друзей")
		return
	}

	if exists {
		c.String(http.StatusConflict, "Пользователь уже в друзьях")
		return
	}

	_, err = db.Exec("INSERT INTO Friends (user_id, friend_id) VALUES (?, ?), (?, ?)", ID, request.FriendID, request.FriendID, ID)
	if err != nil {
		c.String(http.StatusInternalServerError, "Ошибка при добавлении друга")
		return
	}
	c.String(http.StatusOK, "Друг добавлен")
}
