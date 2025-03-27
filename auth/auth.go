package auth

import (
	"database/sql"

	"github.com/gorilla/sessions"
)

// AuthService предоставляет методы для аутентификации пользователей.
type AuthService struct {
	db    *sql.DB
	store *sessions.CookieStore
}

// NewAuthService создает новый AuthService.
func NewAuthService(db *sql.DB, store *sessions.CookieStore) *AuthService {
	return &AuthService{db, store}
}