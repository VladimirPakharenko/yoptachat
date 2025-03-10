package db

import (
	"database/sql"
	"fmt"
	"log"
)

// Подключение к базе данных
func Connect() *sql.DB {
	dsn := "root:@tcp(127.0.0.1:3306)/chat"
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("Ошибка подключения к базе данных: %v", err)
        fmt.Println("ошибка подключения к бд")
	} else {
        fmt.Println("бд успешно подключена")
    }
	return db
}