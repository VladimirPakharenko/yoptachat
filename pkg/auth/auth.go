package auth

import (
	"html/template"
	"log"
	"net/http"
	"yoptachat/pkg/db"
)

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		login := r.FormValue("login")
		phone := r.FormValue("phone")
		password := r.FormValue("password")
		rPassword := r.FormValue("rpassword")

		if password != rPassword {
			log.Println("Пароли не совпадают")
			http.Error(w, "Пароли не совпадают", http.StatusBadRequest)
			return
		}

		_, err := db.DB.Exec("INSERT INTO users (login, phone, password) VALUES (?, ?, ?)", login, phone, password)
		if err != nil {
			log.Println("Ошибка при регистрации:", err)
			http.Error(w, "Ошибка при регистрации", http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/login", http.StatusSeeOther)
	} else {
		renderTemplate(w, "templates/regauth.html")
	}
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		login := r.FormValue("login")
		password := r.FormValue("password")
		rPassword := r.FormValue("rpassword")

		if password != rPassword {
			log.Println("Пароли не совпадают")
			http.Error(w, "Пароли не совпадают", http.StatusBadRequest)
			return
		}

		var storedPassword string
		err := db.DB.QueryRow("SELECT password FROM users WHERE login = ?", login).Scan(&storedPassword)
		if err != nil || storedPassword != password {
			log.Println("Ошибка авторизации:", err)
			http.Error(w, "Ошибка авторизации", http.StatusUnauthorized)
			return
		}

		// Создание сессии
		http.SetCookie(w, &http.Cookie{
			Name:  "session",
			Value: login,
			Path:  "/",
		})

		http.Redirect(w, r, "/index.html", http.StatusSeeOther)
	} else {
		renderTemplate(w, "templates/regauth.html")
	}
}

func renderTemplate(w http.ResponseWriter, tmpl string) {
	t, err := template.ParseFiles(tmpl)
	if err != nil {
		log.Println("Ошибка при рендеринге шаблона:", err)
		http.Error(w, "Ошибка сервера", http.StatusInternalServerError)
		return
	}
	t.Execute(w, nil)
}
