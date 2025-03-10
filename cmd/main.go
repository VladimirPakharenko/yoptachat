package main

import (
	"net/http"
	"yoptachat/pkg/auth"
	// "yoptachat/pkg/chat"
	"yoptachat/pkg/db"
)

func main() {
	database := db.Connect()
	defer database.Close()

	http.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
		auth.RegisterHandler(w, r, database)
	})
	http.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		auth.LoginHandler(w, r, database)
	})
	// http.HandleFunc("/search", func(w http.ResponseWriter, r *http.Request) {
	// 	chat.SearchUsersHandler(w, r, database)
	// })
	// http.HandleFunc("/add_friend", func(w http.ResponseWriter, r *http.Request) {
	// 	chat.AddFriendHandler(w, r, database)
	// })

	// Главная страница с проверкой сессии
	http.HandleFunc("/", auth.CheckSession(func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./pkg/templates/index.html")
	}))

	// Страница регистрации/авторизации
	http.HandleFunc("/regauth.html", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./pkg/templates/regauth.html")
	})

	http.ListenAndServe(":5050", nil)
}