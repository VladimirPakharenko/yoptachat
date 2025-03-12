package auth

import (
	"fmt"
	"html/template"
	"net/http"
)

type PageData struct {
	Login string
}

// Login аутентифицирует пользователя и создает сессию.
func (a *AuthService) Login(w http.ResponseWriter, r *http.Request) {
	login := r.FormValue("login")
	password := r.FormValue("password")
	rpassword := r.FormValue("rpassword")

	if password != rpassword {
		http.Error(w, "Хуйня давай по новой", http.StatusBadRequest)
		fmt.Println("Пароли несовпадают")
		return
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

	// Рендеринг страницы с логином пользователя
	data := PageData{
		Login: user.Login,
	}

	fmt.Println(data)

	tmpl, err := template.ParseFiles("../pkg/templates/index.html")
	if err != nil {
		http.Error(w, "Ошибка при загрузке шаблона", http.StatusInternalServerError)
		return
	}
	fmt.Println("tmpl")
	err = tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, "Ошибка при рендеринге шаблона", http.StatusInternalServerError)
		return
	}
	// http.Redirect(w, r, "/", http.StatusSeeOther)

	// http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
	//     tmpl, _ := template.ParseFiles("templates/index.html")
	//     tmpl.Execute(w, data)
	// })
}
