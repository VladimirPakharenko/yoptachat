package db

// import (
// 	"database/sql"
// 	"fmt"
// 	"log"

// 	_ "github.com/go-sql-driver/mysql"
// )

// var db *sql.DB

// func initDB() (*sql.DB, error){
//     var err error
// 	// Подключаемся к базе данных "chat" с логином root и пустым паролем
// 	db, err := sql.Open("mysql", "root:@tcp(127.0.0.1)/chat")
// 	if err != nil {
// 		fmt.Println("Подключение не удачно")
// 		log.Fatalf("Не удалось подключиться к базе данных: %v", err)
// 	} else {
// 		fmt.Println("БД подклчена #1")
// 	}
// 	return db, err
// }