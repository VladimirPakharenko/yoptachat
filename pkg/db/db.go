package db

import (
    "database/sql"
    "log"

    _ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

func InitDB() {
    var err error
	// Подключаемся к базе данных "chat" с логином root и пустым паролем
	db, err = sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/chat")
	if err != nil {
		log.Fatalf("Не удалось подключиться к базе данных: %v", err)
	}

	// Проверяем подключение
	if err := db.Ping(); err != nil {
		log.Fatalf("Ошибка при проверке подключения: %v", err)
	}
}
func Close() {
    if err := db.Close(); err != nil {
        log.Fatal(err)
    }
}

// Дополнительные функции для работы с базой данных
func GetDB() *sql.DB {
    return db
}