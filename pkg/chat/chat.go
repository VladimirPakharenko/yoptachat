package chat

// import (
// 	"database/sql"
// 	"net/http"
// )

// Функция поиска пользователей
// func SearchUsersHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
// 	query := r.URL.Query().Get("q")
// 	rows, err := db.Query("SELECT id, login, phone FROM users WHERE login LIKE ? OR phone LIKE ?", "%"+query+"%", "%"+query+"%")
// 	if err != nil {
// 		http.Error(w, "Ошибка поиска", http.StatusInternalServerError)
// 		return
// 	}
// 	defer rows.Close()

// 	// Формирование ответа
// 	var users []User
// 	for rows.Next() {
// 		var user User
// 		rows.Scan(&user.ID, &user.Login, &user.Phone)
// 		users = append(users, user)
// 	}
// 	// Отправка JSON-ответа
// 	// (реализация JSON-ответа опущена для краткости)
// }

// Функция добавления друга
// func AddFriendHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
// 	userID := r.URL.Query().Get("user_id")
// 	friendID := r.URL.Query().Get("friend_id")

// 	_, err := db.Exec("INSERT INTO friends (userid1, userid2) VALUES (?, ?)", userID, friendID)
// 	if err != nil {
// 		http.Error(w, "Ошибка добавления друга", http.StatusInternalServerError)
// 		return
// 	}

// 	w.WriteHeader(http.StatusOK)
// }