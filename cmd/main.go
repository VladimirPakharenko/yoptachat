package main

import (
    "log"
    "net/http"
    "yoptachat/pkg/auth"
    "yoptachat/pkg/chat"
    "yoptachat/pkg/db"
)

func main() {
	db.InitDB()
    http.HandleFunc("/register", auth.RegisterHandler)
    http.HandleFunc("/login", auth.LoginHandler)
    http.HandleFunc("/index.html", chat.IndexHandler)

    log.Println("Сервер запущен на http://localhost:8080")
    log.Fatal(http.ListenAndServe(":8080", nil))
}