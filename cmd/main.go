package main

import (
	"net/http"

	"yoptachat/pkg/auth"
	// "yoptachat/pkg/chat"
	"yoptachat/pkg/db"
)

func main() {
	db.InitDB()
	defer db.Close()

	fs := http.FileServer(http.Dir("../pkg/templates"))
	http.Handle("/", fs)

	// http.HandleFunc("/", auth.RedirectHandler)
    http.HandleFunc("/login", auth.LoginHandler)
    http.HandleFunc("/register", auth.RegisterHandler)

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(err)
	}
}
