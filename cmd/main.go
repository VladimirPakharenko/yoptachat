package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"yoptachat/pkg/auth"
	"yoptachat/pkg/chat"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
)

var (
	store = sessions.NewCookieStore([]byte("secret-key")) //определение переменной куда будут записыватсья сессии
)

func main() {
	// Настройка базы данных
	db, err := setupDatabase("root:@tcp(127.0.0.1)/chat") //выполнение функции setupDatabase для включения бд
	if err != nil {
		fmt.Println("ошибка при подключении к бд")
		log.Fatalf("Ошибка при подключении к базе данных: %v", err)
	} else {
		fmt.Println("БД подключена")
	}
	defer db.Close() //закрытие бд

	// Инициализация сервисов
	authService := auth.NewAuthService(db, store) //определение переменной для работы с функциями пакета auth
	chatService := chat.NewChatService(db)        //определение переменной для работы с функциями пакета chat

	// Настройка маршрутов
	router := mux.NewRouter()                     //создание переменной маршрутов для работы с запросами
	setupRoutes(router, authService, chatService) //определение функции для маршрутизации запросов.

	// Открытие html файлов в браузере по пути ниже
	fs := http.FileServer(http.Dir("../pkg/templates"))
	router.PathPrefix("/").Handler(fs)

	// Запуск сервера
	log.Println("Сервер запущен на :8080")
	if err := http.ListenAndServe(":8080", router); //запуск сервера на порту 8080
	err != nil {
		fmt.Println("сервак не замущен")
		log.Fatalf("Ошибка при запуске сервера: %v", err)
	}
}

// setupDatabase устанавливает соединение с базой данных.
func setupDatabase(dataSourceName string) (*sql.DB, error) { //Функция выполняет вход в бд путем принятия переменной определяющей путь и др данные
	db, err := sql.Open("mysql", dataSourceName)
	if err != nil {
		fmt.Println("ошибка 1")
		return nil, err
	}

	if err = db.Ping(); err != nil {
		fmt.Println("ошибка 1")
		return nil, err
	}

	return db, nil
}

// setupRourhtes настраивает маршруты для приложения.
func setupRoutes(router *mux.Router, authService *auth.AuthService, chatService *chat.ChatService) {
	// Маршруты для аутентификации
	router.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			fmt.Println("метод норм")
			registerHandler(authService, w, r)
		} else {
			fmt.Println("метод не разрешен")
			http.Error(w, "Метод не разрешен", http.StatusMethodNotAllowed)
		}
	}).Methods(http.MethodPost)

	router.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			authService.Login(w, r)
		} else {
			http.Error(w, "Метод не разрешен", http.StatusMethodNotAllowed)
		}
	}).Methods(http.MethodPost)

	// Маршруты для чата
	router.HandleFunc("/messages", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			chatService.SendMessage(w, r)
		} else {
			http.Error(w, "Метод не разрешен", http.StatusMethodNotAllowed)
		}
	}).Methods(http.MethodPost)

	router.HandleFunc("/messages", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			chatService.GetMessages(w, r)
		} else {
			http.Error(w, "Метод не разрешен", http.StatusMethodNotAllowed)
		}
	}).Methods(http.MethodGet)

	router.HandleFunc("/search", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			chatService.SearchUsers(w, r)
		} else {
			http.Error(w, "Метод не разрешен", http.StatusMethodNotAllowed)
		}
	}).Methods(http.MethodGet)

	// Маршруты для друзей
	router.HandleFunc("/friends", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			addFriendHandler(chatService, w, r)
		} else if r.Method == http.MethodGet {
			getFriendsHandler(chatService, w, r)
		} else {
			http.Error(w, "Метод не разрешен", http.StatusMethodNotAllowed)
		}
	}).Methods(http.MethodPost, http.MethodGet)

}

// registerHandler обрабатывает регистрацию пользователя.
func registerHandler(authService *auth.AuthService, w http.ResponseWriter, r *http.Request) {
	login := r.FormValue("login")
	phone := r.FormValue("phone")
	password := r.FormValue("password")
	rpassword := r.FormValue("rpassword")

	if rpassword != password {
		http.Error(w, "Хуйня давай по новой", http.StatusBadRequest)
		fmt.Println("Пароли несовпадают")
		return
	}

	fmt.Println("Попытка регистрации пользователя", login)
	_, err := authService.Register(login, phone, password)
	if err != nil {
		fmt.Printf("Ошибка при регистрации %s: %v\n", login, err)
		http.Error(w, "Ошибка при регистрации", http.StatusInternalServerError)
		return
	}

	// w.WriteHeader(http.StatusCreated)
	fmt.Println("переход на /login")
	http.Redirect(w, r, "/regauth.html", http.StatusSeeOther)
	fmt.Println("Редирект выполнен")
}

// addFriendHandler обрабатывает добавление друга.
func addFriendHandler(chatService *chat.ChatService, w http.ResponseWriter, r *http.Request) {
	userID := r.FormValue("user_id")     // ID текущего пользователя
	friendID := r.FormValue("friend_id") // ID друга

	// Преобразуем строки в int
	userIDInt, _ := strconv.Atoi(userID)
	friendIDInt, _ := strconv.Atoi(friendID)

	if err := chatService.AddFriend(userIDInt, friendIDInt); err != nil {
		http.Error(w, "Ошибка при добавлении друга", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// getFriendsHandler получает список друзей.
func getFriendsHandler(chatService *chat.ChatService, w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("user_id") // ID текущего пользователя

	userIDInt, _ := strconv.Atoi(userID)
	friends, err := chatService.GetFriends(userIDInt)
	if err != nil {
		http.Error(w, "Ошибка при получении списка друзей", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(friends)
}
