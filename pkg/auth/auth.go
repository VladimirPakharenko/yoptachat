package auth

import (
	"crypto/sha256"
	"database/sql"
	"fmt"
	"net/http"

	"github.com/gorilla/sessions"
	_ "github.com/go-sql-driver/mysql"
)

// Настройка сессий
var store = sessions.NewCookieStore([]byte("super-secret-key"))

// User структура для работы с пользователями
type User struct {
	ID       int
	Login    string
	Phone    string
	Password string
	IsAdmin  sql.NullBool
}

// Функция регистрации
func RegisterHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	if r.Method != "POST" {
		http.Redirect(w, r, "/regauth.html", http.StatusSeeOther)
		return
	}

	login := r.FormValue("login")
	phone := r.FormValue("phone")
	password := r.FormValue("password")
	rPassword := r.FormValue("rpassword")

	// Проверка совпадения паролей
	if password != rPassword {
		http.Error(w, "Пароли не совпадают", http.StatusBadRequest)
        fmt.Println("пароли при регистрации не совпадают")
		return
	} else {
        fmt.Println("пароли при регистрации совпали")
    }

	// Хэширование пароля
	hashedPassword := fmt.Sprintf("%x", sha256.Sum256([]byte(password)))

	// Вставка нового пользователя в базу данных
	_, err := db.Exec("INSERT INTO users (login, phone, password) VALUES (?, ?, ?)", login, phone, hashedPassword)
	if err != nil {
        http.Error(w, "Ошибка при регистрации", http.StatusInternalServerError)
        fmt.Println("ошибка при регистрации")
		return
	} else {
        fmt.Println("регистрация успешна")
    }

	http.Redirect(w, r, "/regauth.html", http.StatusSeeOther)
}

// Функция авторизации
func LoginHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	if r.Method != "POST" {
		http.Redirect(w, r, "/regauth.html", http.StatusSeeOther)
		return
	}

	loginOrPhone := r.FormValue("login")
	password := r.FormValue("password")
	rPassword := r.FormValue("rpassword")

    // Проверка совпадения паролей
	if password != rPassword {
		http.Error(w, "Пароли не совпадают", http.StatusBadRequest)
        fmt.Println("пароли при авторизации не совпадают")
		return
	} else {
        fmt.Println("пароли при авторизации совпали")
    }

	// Хэширование пароля для проверки
	hashedPassword := fmt.Sprintf("%x", sha256.Sum256([]byte(password)))

	// Проверка пользователя в базе данных
	var user User
	err := db.QueryRow("SELECT id, login, phone, password FROM users WHERE (login = ? OR phone = ?) AND password = ?", loginOrPhone, loginOrPhone, hashedPassword).Scan(&user.ID, &user.Login, &user.Phone, &user.Password)
	if err != nil {
		http.Error(w, "Неверный логин или пароль", http.StatusUnauthorized)
        fmt.Println("неверный логин или пароль")
		return
	} else {
        fmt.Println("успешный вход")
    }

	// Создание сессии
	session, _ := store.Get(r, "session")
	session.Values["user_id"] = user.ID
	session.Save(r, w)

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// Проверка сессии
func CheckSession(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session, _ := store.Get(r, "session")
		if session.Values["user_id"] == nil {
			http.Redirect(w, r, "/regauth.html", http.StatusSeeOther)
			return
		}
		next(w, r)
	}
}

func SearchUsersHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	query := r.URL.Query().Get("q")
	rows, err := db.Query("SELECT id, login, phone FROM users WHERE login LIKE ? OR phone LIKE ?", "%"+query+"%", "%"+query+"%")
	if err != nil {
		http.Error(w, "Ошибка поиска", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	// Формирование ответа
	var users []User
	for rows.Next() {
		var user User
		rows.Scan(&user.ID, &user.Login, &user.Phone)
		users = append(users, user)
	}
	// Отправка JSON-ответа
	// (реализация JSON-ответа опущена для краткости)
}