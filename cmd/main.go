package main

import (
	"fmt"
	"net/http"
	"yoptachat/pkg/db"
)

func main() {
	db.Connect()
	fmt.Println(":5050")
	http.ListenAndServe(":5050", nil)
}
