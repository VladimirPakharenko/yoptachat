package auth

import "net/http"

// Login аутентифицирует пользователя и создает сессию.
func (a *AuthService) Login(w http.ResponseWriter, r *http.Request) {
	login := r.FormValue("login")
	password := r.FormValue("password")
	rpassword := r.FormValue("rpassword")

	if password != rpassword {

	}
	user, err := a.getUserByLogin(login)
	if err != nil || !comparePasswords(user.Password, password) {
		http.Error(w, "Неверный логин или пароль", http.StatusUnauthorized)
		return
	}

	// Создание сессии
	session, _ := a.store.Get(r, "session-name")
	session.Values["user"] = user.ID
	session.Save(r, w)

	http.Redirect(w, r, "/index.html", http.StatusSeeOther)
}
