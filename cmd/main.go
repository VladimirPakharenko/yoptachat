package main

// import (
// 	"fmt"
// 	"log"
// 	"net/http"

	// "yoptachat/pkg/auth"
	// "yoptachat/pkg/chat"
	// "yoptachat/pkg/db"
// )

// func main() {
	// Инициализация бд

	// Выполнение функций по запросам ниже
	// http.HandleFunc("/", auth.RedirectHandler)
	// http.HandleFunc("/login", auth.LoginHandler)
	// http.HandleFunc("/register", auth.RegisterHandler)
	// Открытие html файлов в браузере по пути ниже

// 	fs := http.FileServer(http.Dir("../pkg/templates"))
// 	http.Handle("/", fs)

// 	// Запуск сервера
// 	err := http.ListenAndServe(":8080", nil)
// 	if err == nil {
// 		fmt.Println("Запуск сервера на порту :8080")
// 	} else {
// 		fmt.Println("Проблема при запуске сервера")
// 		panic(err)
// 	}
// }


import (
    "database/sql"
    "log"
    "net/http"

    "github.com/gorilla/mux"
    _ "github.com/go-sql-driver/mysql"
    "yoptachat/pkg/auth"
    "yoptachat/pkg/chat"
    // "yoptachat/pkg/models"
)

func main() {
    // Настройка базы данных
    db, err := setupDatabase("user:password@tcp(127.0.0.1:3306)/dbname")
    if err != nil {
        log.Fatalf("Ошибка при подключении к базе данных: %v", err)
    }
    defer db.Close()

    // Инициализация сервисов
    authService := auth.NewAuthService(db)
    chatService := chat.NewChatService(db)

    // Настройка маршрутов
    router := mux.NewRouter()
    setupRoutes(router, authService, chatService)

    // Запуск сервера
    log.Println("Сервер запущен на :8080")
    if err := http.ListenAndServe(":8080", router); err != nil {
        log.Fatalf("Ошибка при запуске сервера: %v", err)
    }
}

// setupDatabase устанавливает соединение с базой данных.
func setupDatabase(dataSourceName string) (*sql.DB, error) {
    db, err := sql.Open("mysql", dataSourceName)
    if err != nil {
        return nil, err
    }

    if err = db.Ping(); err != nil {
        return nil, err
    }

    return db, nil
}

// setupRoutes настраивает маршруты для приложения.
func setupRoutes(router *mux.Router, authService *auth.AuthService, chatService *chat.ChatService) {
    // Маршруты для аутентификации
    router.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
        if r.Method == http.MethodPost {
            registerHandler(authService, w, r)
        } else {
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
}

// registerHandler обрабатывает регистрацию пользователя.
func registerHandler(authService *auth.AuthService, w http.ResponseWriter, r *http.Request) {
    login := r.FormValue("login")
    phone := r.FormValue("phone")
    password := r.FormValue("password")

    _, err := authService.Register(login, phone, password)
    if err != nil {
        http.Error(w, "Ошибка при регистрации", http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusCreated)
}







// http.HandleFunc("/pw", auth.Work)// Проверка связи