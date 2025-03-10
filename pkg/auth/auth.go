package auth
import (
	"fmt"
	"html/template"
	"net/http"
	"yoptachat/pkg/db"
	// "github.com/gorilla/sessions"
)

var (
// Создаем хранилище сессий
// store = sessions.NewCookieStore([]byte("secret-key"))
)

// Обработчик регистрации пользователя
func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		login := r.FormValue("login")
		phone := r.FormValue("phone")
		password := r.FormValue("password")
		rPassword := r.FormValue("rpassword")

		// Проверка соответствия паролей
		if password != rPassword {
			http.Error(w, "Пароли не совпадают", http.StatusBadRequest)
			return
		}

		// Сохранение пользователя в базе данных
		if err := saveUser(login, phone, password); err != nil {
			http.Error(w, "Ошибка при регистрации", http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/login", http.StatusSeeOther)
	} else {
		// Если метод не POST, показываем форму
		tmpl := template.Must(template.ParseFiles("../pkg/templates/regauth.html"))
		tmpl.Execute(w, nil)
	}
}

// Обработчик авторизации пользователя
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		login := r.FormValue("login")
		password := r.FormValue("password")
		rPassword := r.FormValue("rpassword")

		if password != rPassword {
			http.Error(w, "Пароли не совпадают", http.StatusBadRequest)
			return
		}
		// Проверка пользователя
		if err := checkUser(login, password); err != nil {
			http.Error(w, "Ошибка авторизации", http.StatusUnauthorized)
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:  "session",
			Value: login, // Храним имя пользователя в куке
			Path:  "/",
		})

		// Успешная авторизация
		http.Redirect(w, r, "/", http.StatusSeeOther)
	} else {
		// Если метод не POST, показываем форму
		tmpl := template.Must(template.ParseFiles("../pkg/templates/regauth.html"))
		tmpl.Execute(w, nil)
	}
}

// Функция для сохранения пользователя в базе данных
func saveUser(login, phone, password string) error {
	db := db.GetDB()
	_, err := db.Exec("INSERT INTO users (login, phone, password) VALUES (?, ?, ?)", login, phone, password)
	return err
}

// Функция для проверки пользователя
func checkUser(login, password string) error {
	db := db.GetDB()
	var storedPassword string
	err := db.QueryRow("SELECT password FROM users WHERE login = ?", login).Scan(&storedPassword)
	if err != nil {
		return err
	}

	if storedPassword != password {
		return fmt.Errorf("неверный пароль")
	}
	return nil
}

func CheckSessionMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("session")
		if err != nil || cookie.Value == "" {
			// Если сессия пустая, перенаправляем на страницу авторизации
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		next.ServeHTTP(w, r)
	}
}